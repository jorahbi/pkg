package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/edwingeng/deque/v2"
	"github.com/gobwas/ws/wsutil"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

const share = 10

type wsServer struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
	connected int64
	queue     []*deque.Deque[*conn]
	conns     []cmap.ConcurrentMap[string, *gnet.Conn]
}

type conn struct {
	conn *gnet.Conn
	last time.Time
}

func (wss *wsServer) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.multicore, wss.addr)
	return gnet.None
}

func (wss *wsServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(new(wsCodec))
	idx := c.Fd() % share
	wss.queue[idx].PushBack(&conn{conn: &c, last: time.Now()})
	wss.conns[idx].Set(strconv.Itoa(c.Fd()), &c)
	atomic.AddInt64(&wss.connected, 1)
	return nil, gnet.None
}

func (wss *wsServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt64(&wss.connected, -1)
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (wss *wsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ws := c.Context().(*wsCodec)
	if ws.readBufferBytes(c) == gnet.Close {
		return gnet.Close
	}
	ok, action := ws.upgrade(c)
	fmt.Println("upgrade....")
	if !ok {
		return
	}

	if ws.buf.Len() <= 0 {
		return gnet.None
	}
	messages, err := ws.Decode(c)
	if err != nil {
		return gnet.Close
	}
	if messages == nil {
		return
	}
	for _, message := range messages {
		msgLen := len(message.Payload)
		if msgLen > 128 {
			logging.Infof("conn[%v] receive [op=%v] [msg=%v..., len=%d]", c.RemoteAddr().String(), message.OpCode, string(message.Payload[:128]), len(message.Payload))
		} else {
			logging.Infof("conn[%v] receive [op=%v] [msg=%v, len=%d]", c.RemoteAddr().String(), message.OpCode, string(message.Payload), len(message.Payload))
		}
		// This is the echo server
		err = wsutil.WriteServerMessage(c, message.OpCode, message.Payload)
		if err != nil {
			logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
			return gnet.Close
		}
	}

	return gnet.None
}

func (wss *wsServer) OnTick() (delay time.Duration, action gnet.Action) {
	logging.Infof("[connected-count=%v]", atomic.LoadInt64(&wss.connected))
	for _, item := range wss.queue {
		for !item.IsEmpty() {
			c := item.PopFront()
			ct := time.Now()
			gap := ct.Sub(c.last).Milliseconds()
			fmt.Println("keeplive", gap)
			if gap < 5000 {
				item.PushBack(c)
				break
			}

			fmt.Println("keeplive")
			err := wsutil.WriteServerText(*c.conn, []byte("keeplive"))
			if err != nil {
				fmt.Println(err)
				continue
			}
			c.last = ct
			item.PushBack(c)
		}
	}
	return 5 * time.Second, gnet.None
}

func NewServer() *wsServer {
	wss := &wsServer{
		addr:      "tcp://127.0.0.1:9080",
		multicore: true,
		queue:     make([]*deque.Deque[*conn], 0, share),
		conns:     make([]cmap.ConcurrentMap[string, *gnet.Conn], 0, share),
	}

	for i := 0; i < share; i++ {
		wss.queue = append(wss.queue, deque.NewDeque[*conn]())
		wss.conns = append(wss.conns, cmap.New[*gnet.Conn]())
	}

	return wss
}

func main() {
	go httpAccept()
	// Example command: go run main.go --port 8080 --multicore=true
	// flag.IntVar(&port, "port", 9080, "server port")
	flag.Parse()
	wss := NewServer()
	// Start serving!
	log.Println("server exits:", gnet.Run(wss, wss.addr, gnet.WithMulticore(true), gnet.WithReusePort(true), gnet.WithTicker(true)))
}

func httpAccept() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
