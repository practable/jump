/* package ws offers bi-directional proxies
   from tcp to ws host or ws client
*/

package ws

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

const (
	ChunkSize = 0xffff //65535, typical TCP MTU
)

// Client listens for TCP connections and makes a client connection to a remote ws server
func Client(ctx context.Context, listen int, target string) {

	lc := &net.ListenConfig{}

	uri := ":" + strconv.Itoa(listen)

	l, err := lc.Listen(ctx, "tcp", uri)

	if err != nil {
		log.WithFields(log.Fields{"uri": uri, "err": err.Error()}).Errorf("failed to create listener because %s", err.Error())
		return
	}

	defer l.Close()

	log.WithField("uri", uri).Infof("awaiting connections at %s", uri)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Wait for a connection.
			tconn, err := l.Accept()

			if err != nil {
				log.WithFields(log.Fields{"uri": uri, "err": err.Error()}).Errorf("failed to accept connection because %s", err.Error())
				return //the context is probably cancelled.
			}
			// Handle the connection in a new goroutine.

			log.WithField("uri", uri).Tracef("got a new connection")

			// make websocket connection and relay to/from conn

			go func() {
				wconn, _, err := websocket.DefaultDialer.Dial(target, nil)
				if err != nil {
					log.Fatal("dial:", err)
				}

				defer wconn.Close()

				PipeBinaryIgnoreText(ctx, tconn, wconn)
			}()

		}
	}
}

// Host listens for ws connections, and makes a TCP connection to the target
// for each incoming ws connection, then relays messages bidirectionally until
// the context ends

func Host(ctx context.Context, listen int, target string) {

	server := http.NewServeMux()

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		wconn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("internal/ws.Host(HandleFunc; upgrade): %s", err)
			return
		}
		defer wconn.Close()

		tconn, err := net.Dial("tcp", target)
		if err != nil {
			log.Errorf("internal/ws.Host(HandleFunc; net.Dial): %s", err.Error())
			return
		}
		defer tconn.Close()

		PipeBinaryIgnoreText(ctx, tconn, wconn)

	})

	addr := ":" + strconv.Itoa(listen)

	h := &http.Server{Addr: addr, Handler: server}

	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Fatalf("internal/ws.Host(ListenAndServe): %s ", err.Error())
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := h.Shutdown(ctx)
	if err != nil {
		log.Errorf("internal/ws.Host(ListenAndServe.Shutdown): %s", err.Error())
	}
	log.Trace("internal/ws.Host is done")

}

// Implementations checked before writing this code
// https://github.com/gorilla/websocket/issues/588
// https://pkg.go.dev/net/http/httputil#ReverseProxy
// https://pkg.go.dev/github.com/gorilla/websocket
// https://github.com/hazaelsan/ssh-relay/blob/master/session/corprelay/corprelay.go

// PipeBinaryIgnoreText relays binary data between net.Conn and websocket.Conn
// pass websocket.Conn with pointer due to mutex
// Use of NextReader/Writer modified from https://github.com/hazaelsan/ssh-relay/blob/master/session/corprelay/corprelay.go
func PipeBinaryIgnoreText(ctx context.Context, tconn net.Conn, wconn *websocket.Conn) {
	id := "internal/ws/PipeBinaryIgnoreText"
	if wconn == nil {
		log.Errorf("%s nil websocket.Conn", id)
		return
	}
	// write from wconn to tconn
	go func() {

		for {

			err := func() error {

				t, r, err := (*wconn).NextReader()
				if err != nil {
					return fmt.Errorf("NextReader() error: %w", err)
				}

				b := new(bytes.Buffer)

				n, err := b.ReadFrom(r)
				if err != nil {
					return err
				}

				log.Tracef("(ws->tcp): read %v bytes", n)

				switch t {
				case websocket.BinaryMessage:
					// TODO what does this conversion achieve?
					nn, err := tconn.Write(b.Bytes())
					log.Tracef("(ws->tcp): wrote %v bytes", nn)
					return err
				case websocket.TextMessage: //non-SSH message (none implemented as yet anyway)
					log.WithField("msg", b.String()).Info("Ignoring TEXT")

				default:
					return fmt.Errorf("unsupported message type: %v", t)
				}
				return nil
			}()
			if err != nil {
				log.WithField("err", err.Error()).Error(id + " " + err.Error())
				return
			}
		}
	}()

	// write from tconn to wconn
	go func() {
		for {
			err := func() error {
				r := bufio.NewReader(tconn)
				b := make([]byte, ChunkSize)
				n, err := r.Read(b)
				data := b[0:n]
				log.Tracef("internal/ws.Host(tcp->ws): read %v bytes", n)
				if err != nil {
					return err
				}

				w, err := (*wconn).NextWriter(websocket.BinaryMessage)
				if err != nil {
					return fmt.Errorf("NextWriter() error: %w", err)
				}
				defer w.Close()

				n, err = w.Write(data)
				log.Tracef("internal/ws.Host(tcp->ws): wrote %v bytes", n)
				return err

			}()

			if err != nil {
				log.WithField("err", err.Error).Error("internal/ws.Host(tcp->ws) error")
				return
			}
		}

	}()

	<-ctx.Done()

}
