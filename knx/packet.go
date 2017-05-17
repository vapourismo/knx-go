package knx

import (
	"bytes"
	"errors"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// ServiceID identifies the service that is contained in a packet.
type ServiceID uint16

// KNXnet/IP services
const (
	ConnReqService      ServiceID = 0x0205
	ConnResService                = 0x0206
	ConnStateReqService           = 0x0207
	ConnStateResService           = 0x0208
	DiscReqService                = 0x0209
	DiscResService                = 0x020a
	TunnelReqService              = 0x0420
	TunnelResService              = 0x0421
)

func writePacketHeader(w *bytes.Buffer, service ServiceID, payloadLength int) error {
	if payloadLength < 0 || payloadLength > 65529 {
		return errors.New("Payload length is out of bounds")
	}

	return encoding.WriteSequence(w, byte(6), byte(16), uint16(service), uint16(payloadLength+6))
}

func readPacketHeader(r *bytes.Reader, service *ServiceID, payloadLength *int) error {
	var headerLength, protocolVersion byte
	var service16, packetLength16 uint16

	err := encoding.ReadSequence(r, &headerLength, &protocolVersion, &service16, &packetLength16)

	if err != nil {
		return err
	}

	if headerLength != 6 {
		return errors.New("Header length mismatches")
	}

	if protocolVersion != 16 {
		return errors.New("Unknown protocol version")
	}

	if packetLength16 < 6 {
		return errors.New("Packet length is too low")
	}

	tmpPayloadLength := int(packetLength16) - 6

	if tmpPayloadLength > r.Len() {
		return errors.New("Packet is incomplete")
	}

	*service = ServiceID(service16)
	*payloadLength = tmpPayloadLength

	return nil
}

// A OutgoingPayload is the payload within a KNXnet/IP packet.
type OutgoingPayload interface {
	describe() (ServiceID, int)
	writeTo(w *bytes.Buffer) error
}

// WritePacket writes a complete KNXnet/IP packet (header and payload) to the given Buffer.
func WritePacket(w *bytes.Buffer, payload OutgoingPayload) error {
	service, length := payload.describe()

	err := writePacketHeader(w, service, length)
	if err != nil {
		return err
	}

	return payload.writeTo(w)
}

// ReadPacket reads a KNXnet/IP packet from the given Reader and returns an instance of the
// corresponding payload type.
//
// Types can be: *ConnectionResponse
//               *ConnectionStateResponse
//               *DisconnectRequest
//               *DisconnectResponse
//               *TunnelRequest
//               *TunnelResponse
//
func ReadPacket(r *bytes.Reader) (interface{}, error) {
	var service ServiceID
	var length int

	err := readPacketHeader(r, &service, &length)
	if err != nil {
		return nil, err
	}

	switch service {
	case ConnResService:
		return readConnectionResponse(r)

	case ConnStateResService:
		return readConnectionStateResponse(r)

	case TunnelReqService:
		return readTunnelRequest(r)

	case TunnelResService:
		return readTunnelResponse(r)

	case DiscReqService:
		return readDisconnectRequest(r)

	case DiscResService:
		return readDisconnectResponse(r)

	default:
		return nil, errors.New("Unknown service")
	}
}
