package proto

import (
	"bytes"
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// ServiceID identifies the service that is contained in a packet.
type ServiceID uint16

// These are supported services.
const (
	ConnReqService      ServiceID = 0x0205
	ConnResService      ServiceID = 0x0206
	ConnStateReqService ServiceID = 0x0207
	ConnStateResService ServiceID = 0x0208
	DiscReqService      ServiceID = 0x0209
	DiscResService      ServiceID = 0x020a
	TunnelReqService    ServiceID = 0x0420
	TunnelResService    ServiceID = 0x0421
)

// Service describes a KNXnet/IP service.
type Service interface {
	Service() ServiceID
}

// ServiceWriterTo combines WriterTo and Service.
type ServiceWriterTo interface {
	Service
	io.WriterTo
}

// Pack generates a KNXnet/IP packet.
func Pack(w io.Writer, srv ServiceWriterTo) (int64, error) {
	dataBuffer := bytes.Buffer{}

	_, err := srv.WriteTo(&dataBuffer)
	if err != nil {
		return 0, err
	}

	return encoding.WriteSome(
		w, byte(6), byte(16), srv.Service(), uint16(dataBuffer.Len() + 6), &dataBuffer,
	)
}

// These are errors that might occur during unpacking.
var (
	ErrHeaderLength = errors.New("Header length is not 6")
	ErrHeaderVersion = errors.New("Protocol version is not 16")
	ErrUnknownService = errors.New("Unknown service identifier")
)

type serviceReaderFrom interface {
	io.ReaderFrom
	Service
}

// Unpack parses a KNXnet/IP packet and retrieves its service payload.
func Unpack(r io.Reader, srv *Service) (int64, error) {
	var headerLen, version uint8
	var srvID ServiceID
	var totalLen uint16

	n, err := encoding.ReadSome(r, &headerLen, &version, &srvID, &totalLen)
	if err != nil {
		return n, err
	}

	if headerLen != 6 {
		return n, ErrHeaderLength
	}

	if version != 16 {
		return n, ErrHeaderVersion
	}

	var body serviceReaderFrom
	switch srvID {
		case ConnResService:
			body = &ConnRes{}

		case ConnStateReqService:
			body = &ConnStateReq{}

		case ConnStateResService:
			body = &ConnStateRes{}

		case DiscReqService:
			body = &DiscReq{}

		case DiscResService:
			body = &DiscRes{}

		case TunnelReqService:
			body = &TunnelReq{}

		case TunnelResService:
			body = &TunnelRes{}

		default:
			return n, ErrUnknownService
	}

	m, err := body.ReadFrom(r)

	if err == nil {
		*srv = body
	}

	return n + m, err
}