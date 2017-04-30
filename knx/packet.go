package knx

import (
	"bytes"
	"errors"
)

type serviceIdent uint16

const (
	connectionRequestService  serviceIdent = 0x0205
	connectionResponseService              = 0x0206
	connectionStateRequestService          = 0x0207
	connectionStateResponseService         = 0x0208
	disconnectRequestService               = 0x0209
	tunnelRequestService                   = 0x0420
	tunnelResponseService                  = 0x0421
)

func writePacketHeader(w *bytes.Buffer, service serviceIdent, payloadLength int) error {
	if payloadLength < 0 || payloadLength > 65529 {
		return errors.New("Payload length is out of bounds")
	}

	return writeSequence(w, byte(6), byte(16), uint16(service), uint16(payloadLength + 6))
}

func readPacketHeader(r *bytes.Reader, service *serviceIdent, payloadLength *int) error {
	var headerLength, protocolVersion byte
	var service16, packetLength16 uint16

	err := readSequence(r, &headerLength, &protocolVersion, &service16, &packetLength16)

	if err != nil { return err }

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

	*service = serviceIdent(service16)
	*payloadLength = tmpPayloadLength

	return nil
}

// Payload for outgoing traffic
type OutgoingPayload interface {
	describe() (serviceIdent, int)
	writeTo(w *bytes.Buffer) error
}

// WritePacket writes a complete KNXnet/IP packet (header and payload) to the given Buffer.
func WritePacket(w *bytes.Buffer, payload OutgoingPayload) error {
	service, length := payload.describe()

	err := writePacketHeader(w, service, length)
	if err != nil { return err }

	return payload.writeTo(w)
}

// ReadPacket reads a KNXnet/IP packet from the given Reader and returns an instance of the
// corresponding payload type.
//
// Types can be: *ConnectionResponse
//               *ConnectionStateResponse
//               *DisconnectRequest
//               *TunnelRequest
//               *TunnelResponse
//
func ReadPacket(r *bytes.Reader) (interface{}, error) {
	var service serviceIdent
	var length int

	err := readPacketHeader(r, &service, &length)
	if err != nil { return nil, err }

	switch service {
		case connectionResponseService:
			return readConnectionResponse(r)

		case connectionStateResponseService:
			return readConnectionStateResponse(r)

		case tunnelRequestService:
			return readTunnelRequest(r)

		case tunnelResponseService:
			return readTunnelResponse(r)

		case disconnectRequestService:
			return readDisconnectRequest(r)

		default:
			return nil, errors.New("Unknown service")
	}
}
