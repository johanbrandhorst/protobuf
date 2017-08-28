package grpcweb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"

	"github.com/johanbrandhorst/protobuf/grpcweb/status"
	"github.com/johanbrandhorst/protobuf/internal"
)

// closeEvent allows a CloseEvent to be used as an error.
type closeEvent struct {
	*js.Object
	Code     int    `js:"code"`
	Reason   string `js:"reason"`
	WasClean bool   `js:"wasClean"`
}

func (e closeEvent) isWebsocketEvent() {}

func (e *closeEvent) Error() string {
	var cleanStmt string
	if e.WasClean {
		cleanStmt = "clean"
	} else {
		cleanStmt = "unclean"
	}
	return fmt.Sprintf("CloseEvent: (%s) (%d) %s", cleanStmt, e.Code, e.Reason)
}

func beginHandlerOpen(ch chan error, removeHandlers func()) func(ev *js.Object) {
	return func(ev *js.Object) {
		removeHandlers()
		close(ch)
	}
}

func beginHandlerClose(ch chan error, removeHandlers func()) func(ev *js.Object) {
	return func(ev *js.Object) {
		removeHandlers()
		go func() {
			ch <- &closeEvent{Object: ev}
			close(ch)
		}()
	}
}

// ClientStream is the interface exposed by the websocket proxy
type ClientStream interface {
	RecvMsg() ([]byte, error)
	SendMsg([]byte) error
	CloseSend() error
	CloseAndRecv() ([]byte, error)
	Context() context.Context
}

// NewClientStream opens a new WebSocket connection for performing client-side
// and bi-directional streaming. It will block until the connection is
// established or fails to connect.
func (c *Client) NewClientStream(ctx context.Context, method string) (ClientStream, error) {
	ws, err := websocketjs.New(strings.Replace(c.host, "https", "wss", 1) + "/" + c.service + "/" + method)
	if err != nil {
		return nil, err
	}
	conn := &conn{
		WebSocket: ws,
		ch:        make(chan wsEvent, 1),
		ctx:       ctx,
	}

	// We need this so that received binary data is in ArrayBufferView format so
	// that it can easily be read.
	conn.BinaryType = "arraybuffer"

	conn.AddEventListener("message", false, conn.onMessage)
	conn.AddEventListener("close", false, conn.onClose)

	openCh := make(chan error, 1)

	var (
		openHandler  func(ev *js.Object)
		closeHandler func(ev *js.Object)
	)

	// Handlers need to be removed to prevent a panic when the WebSocket closes
	// immediately and fires both open and close before they can be removed.
	// This way, handlers are removed before the channel is closed.
	removeHandlers := func() {
		ws.RemoveEventListener("open", false, openHandler)
		ws.RemoveEventListener("close", false, closeHandler)
	}

	// We have to use variables for the functions so that we can remove the
	// event handlers afterwards.
	openHandler = beginHandlerOpen(openCh, removeHandlers)
	closeHandler = beginHandlerClose(openCh, removeHandlers)

	ws.AddEventListener("open", false, openHandler)
	ws.AddEventListener("close", false, closeHandler)

	err, ok := <-openCh
	if ok && err != nil {
		return nil, err
	}

	return conn, nil
}

// wsEvent encapsulates both message and close events
type wsEvent interface {
	isWebsocketEvent()
}

type conn struct {
	*websocketjs.WebSocket

	ch  chan wsEvent
	ctx context.Context
}

type messageEvent struct {
	*js.Object
	Data *js.Object `js:"data"`
}

func (m messageEvent) isWebsocketEvent() {}

func (c *conn) onMessage(ev *js.Object) {
	go func() {
		c.ch <- &messageEvent{Object: ev}
	}()
}

func (c *conn) onClose(ev *js.Object) {
	go func() {
		// We queue the error to the end so that any messages received prior to
		// closing get handled first.
		c.ch <- &closeEvent{Object: ev}
	}()
}

// receiveFrame receives one full frame from the WebSocket. It blocks until the
// frame is received.
func (c *conn) receiveFrame(ctx context.Context) (*messageEvent, error) {
	select {
	case event, ok := <-c.ch:
		if !ok { // The channel has been closed
			return nil, io.EOF
		}

		switch m := event.(type) {
		case *messageEvent:
			return m, nil
		case *closeEvent:
			close(c.ch)
			if m.Code == 4000 { // codes.OK
				return nil, io.EOF
			}
			// Otherwise, propagate close error
			return nil, m
		default:
			return nil, errors.New("unexpected message type")
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RecvMsg reads a message from the stream.
// It blocks until a message or error has been received.
func (c *conn) RecvMsg() ([]byte, error) {
	ev, err := c.receiveFrame(c.ctx)
	if err != nil {
		if cerr, ok := err.(*closeEvent); ok && internal.IsgRPCErrorCode(cerr.Code) {
			return nil, &status.Status{
				Code:    internal.ParseErrorCode(cerr.Code),
				Message: cerr.Reason,
			}
		}
		return nil, err
	}

	// Check if it's an array buffer. If so, convert it to a Go byte slice.
	if constructor := ev.Data.Get("constructor"); constructor == js.Global.Get("ArrayBuffer") {
		uint8Array := js.Global.Get("Uint8Array").New(ev.Data)
		return uint8Array.Interface().([]byte), nil
	}
	return []byte(ev.Data.String()), nil
}

// SendMsg sends a message on the stream.
func (c *conn) SendMsg(msg []byte) error {
	return c.Send(msg)
}

// CloseSend closes the stream.
func (c *conn) CloseSend() error {
	// CloseSend does not itself read the close event,
	// it will be done by the next Recv
	return c.SendMsg(internal.FormatCloseMessage())
}

// CloseAndRecv closes the stream and returns the last message.
func (c *conn) CloseAndRecv() ([]byte, error) {
	err := c.CloseSend()
	if err != nil {
		return nil, err
	}

	// Read last message
	msg, err := c.RecvMsg()
	if err != nil {
		return nil, err
	}

	// Read close event
	_, err = c.RecvMsg()
	if err != io.EOF {
		return nil, err
	}

	return msg, nil
}

// Context returns the streams context.
func (c *conn) Context() context.Context {
	return c.ctx
}
