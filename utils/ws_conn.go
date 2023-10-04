package utils

import (
	"io"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/polevpn/anyvalue"
	"github.com/polevpn/elog"
)

type Conn interface {
	Id() string
	Send([]byte)
	Close(bool) error
	IsClosed() bool
	String() string
	Read()
	Write()
}

type WsManager struct {
	mutex    *sync.Mutex
	callback func(av *anyvalue.AnyValue)
	conn     Conn
}

type ProcessHandler interface {
	OnConnected(conn Conn)
	OnRequest(pkg []byte, conn Conn)
	OnClosed(conn Conn, proactive bool)
}

type WebSocketConn struct {
	conn    *websocket.Conn
	wch     chan []byte
	closed  bool
	handler ProcessHandler
}

func NewConnectionPool(callback func(av *anyvalue.AnyValue)) (*WsManager, error) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:9011", nil)
	if err != nil {
		return nil, err
	}

	mgr := &WsManager{
		mutex:    &sync.Mutex{},
		callback: callback,
	}

	mgr.conn = NewWebSocketConn(conn, mgr)
	go mgr.conn.Read()
	go mgr.conn.Write()

	return mgr, nil
}

func (wsm *WsManager) OnConnected(conn Conn) {}

func (wsm *WsManager) OnClosed(conn Conn, proactive bool) {
	if wsm.callback != nil {
		wsm.callback(anyvalue.New().Set("event", "stoped"))
	}
}

func (wsm *WsManager) OnRequest(pkg []byte, conn Conn) {
	av, err := anyvalue.NewFromJson(pkg)
	if err != nil {
		elog.Error("decode json fail,", err)
		return
	}
	if wsm.callback != nil {
		wsm.callback(av)
	}
}

// new websocket connection
func NewWebSocketConn(conn *websocket.Conn, handler ProcessHandler) *WebSocketConn {
	return &WebSocketConn{
		conn:    conn,
		closed:  false,
		wch:     make(chan []byte, CH_WEBSOCKET_WRITE_SIZE),
		handler: handler,
	}
}

func (wsc *WebSocketConn) Close(flag bool) error {
	if !wsc.closed {
		wsc.closed = true
		if wsc.wch != nil {
			wsc.wch <- nil
			close(wsc.wch)
		}
		err := wsc.conn.Close()
		if flag {
			go wsc.handler.OnClosed(wsc, false)
		}
		return err
	}
	return nil
}

func (wsc *WebSocketConn) String() string {
	return wsc.conn.RemoteAddr().String() + "->" + wsc.conn.LocalAddr().String()
}

func (wsc *WebSocketConn) Id() string {
	return wsc.conn.RemoteAddr().String() + RandomString(10)
}

func (wsc *WebSocketConn) IsClosed() bool {
	return wsc.closed
}
func (wsc *WebSocketConn) Read() {
	defer func() {
		wsc.Close(true)
	}()

	defer func() {
		if err := recover(); err != nil {
			elog.Error("panic error:", err)
		}
	}()

	for {
		mtype, pkt, err := wsc.conn.ReadMessage()
		if err != nil {
			if err == io.ErrUnexpectedEOF || err == io.EOF || strings.Contains(err.Error(), "close") {
				elog.Info(wsc.String(), ",conn closed")
			} else {
				elog.Error(wsc.String(), ",conn read exception:", err)
			}
			return
		}
		if mtype == websocket.BinaryMessage {
			wsc.handler.OnRequest(pkt, wsc)
		} else {
			elog.Info("ws mtype=", mtype)
		}
	}

}

func (wsc *WebSocketConn) Write() {
	defer func() {
		if err := recover(); err != nil {
			elog.Error("panic error:", err)
		}
	}()

	for {

		pkt, ok := <-wsc.wch
		if !ok {
			elog.Error(wsc.String(), "get pkt from write channel fail,maybe channel closed")
			return
		}
		if pkt == nil {
			elog.Info(wsc.String(), ",exit write process")
			return
		}
		err := wsc.conn.WriteMessage(websocket.BinaryMessage, pkt)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "close") {
				elog.Info(wsc.String(), ",conn closed")
			} else {
				elog.Error(wsc.String(), ",conn write exception:", err)
			}
			return
		}
	}
}

func (wsc *WebSocketConn) Send(pkt []byte) {
	if wsc.closed {
		return
	}
	if wsc.wch != nil {
		wsc.wch <- pkt
	}
}
