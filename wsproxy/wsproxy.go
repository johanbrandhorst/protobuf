package bidi

import (
	"context"
	"encoding/binary"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/transport"

	"github.com/johanbrandhorst/protobuf/internal"
)

// Logger is the interface used by the Proxy to log events
type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
}

// Proxy wraps a handler with a websocket to perform
// bidirectional messaging between a gRPC backend and a web frontend.
type Proxy struct {
	h      http.Handler
	logger Logger
	creds  credentials.TransportCredentials
}

// WrapServer warps the input handler with a Websocket-to-Bidi-streaming proxy.
func WrapServer(h http.Handler, logger Logger, creds credentials.TransportCredentials) http.Handler {
	return &Proxy{
		h:      h,
		logger: logger,
		creds:  creds,
	}
}

// TODO: allow modification of upgrader settings?
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Enforce only local origins
		return true
	},
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		p.h.ServeHTTP(w, r)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.logger.Warn("Failed to upgrade Websocket: ", err)
		return
	}

	ctx, cancelFn := context.WithCancel(r.Context())
	defer cancelFn()

	p.logger.Debug("Creating new transport with addr: ", r.Host)
	t, err := transport.NewClientTransport(ctx,
		transport.TargetInfo{Addr: r.Host},
		transport.ConnectOptions{
			TransportCredentials: p.creds,
		})
	if err != nil {
		p.logger.Warn("Failed to create transport: ", err)
		return
	}
	defer func() {
		err = t.GracefulClose()
		if err != nil {
			p.logger.Warn("Failed to close transport: ", err)
		}
	}()

	p.logger.Debug("Creating new stream with host: ", r.RemoteAddr, "and method: ", r.RequestURI)
	s, err := t.NewStream(ctx, &transport.CallHdr{
		Host:   r.RemoteAddr,
		Method: r.RequestURI,
	})
	if err != nil {
		p.logger.Warn("Failed to create stream: ", err)
		return
	}

	// Read loop - reads from websocket and puts it on the stream
	go func() {
		for {
			select {
			case <-s.Context().Done():
				p.logger.Debug("[READ] Context canceled, returning")
				return
			default:
			}
			p.logger.Debug("[READ] Reading from Websocket")
			_, payload, err := conn.ReadMessage()
			if err != nil {
				p.logger.Warn("[READ] Failed to read Websocket message: ", err)
				return
			}
			p.logger.Debug("[READ] Read payload: ", payload)
			if internal.IsCloseMessage(payload) {
				err = t.Write(s, nil, &transport.Options{Last: true})
				if err == io.EOF || err == nil {
					return
				}
			} else {
				// Append header
				payload = append(make([]byte, 5), payload...)
				// Skip first byte to indicate no compression
				// TODO: Add compression?
				// Encode size of payload to byte 1-5
				binary.BigEndian.PutUint32(payload[1:5], uint32(len(payload)-5))
				err = t.Write(s, payload, &transport.Options{Last: false})
			}

			if err != nil {
				p.logger.Warn("[READ] Failed to write message to transport", err)
				if _, ok := err.(transport.ConnectionError); !ok {
					t.CloseStream(s, err)
				}
				return
			}
		}
	}()

	// Write loop -- take messages from stream and write to websocket
	var header [5]byte
	var msg []byte
	for {
		// Read header
		_, err := s.Read(header[:])
		if err != nil {
			if err == io.EOF {
				p.logger.Debug("[WRITE] Stream closed")
			} else {
				p.logger.Warn("[WRITE] Failed to read header: ", err)
			}

			// Wait for status to be received
			<-s.Done()
			p.sendStatus(conn, s.Status())
			return
		}

		// TODO: Add compression?
		isCompressed := uint8(header[0]) != 0
		if isCompressed {
			// If payload is compressed, bail out
			p.logger.Warn("[WRITE] Reply was compressed, bailing")
			p.sendStatus(conn, status.New(codes.FailedPrecondition, "Server sent compressed data"))
			return
		}
		len := int(binary.BigEndian.Uint32(header[1:]))

		// TODO: Reuse buffer and resize as necessary instead
		msg = make([]byte, int(len))
		if n, err := s.Read(msg); err != nil || n != len {
			p.logger.Warn("[WRITE] Failed to read message: ", err)
			// Wait for status to be received
			<-s.Done()
			p.sendStatus(conn, s.Status())
			return
		}

		if err = conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			p.logger.Warn("[WRITE] Failed to write message: ", err)
			return
		}
		p.logger.Debug("[WRITE] Sent", msg)
	}

}

func (p *Proxy) sendStatus(conn *websocket.Conn, st *status.Status) {
	p.logger.Debug(`Sending status: Msg: "`, st.Message(), `", Code: `, st.Code().String())

	closeMsg := websocket.FormatCloseMessage(internal.FormatErrorCode(st.Code()), st.Message())
	err := conn.WriteMessage(websocket.CloseMessage, closeMsg)
	if err != nil {
		p.logger.Warn("error writing websocket trailer: ", err)
	}

	p.logger.Debug("Sent close message")
	err = conn.Close()
	if err != nil {
		p.logger.Warn("Failed to close connection: ", err)
		return
	}
	p.logger.Debug("Closed connection")
	return
}
