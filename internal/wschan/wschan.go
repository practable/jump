/* package ws offers bi-directional proxies
   from channels to ws host or ws client
*/

package wschan

import (
	"bytes"
	"context"
	"fmt"
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

// Client connects to a websocket server and bidirectionally relays binary messages
func Client(ctx context.Context, target string, send, receive chan []byte) {

	wconn, _, err := websocket.DefaultDialer.Dial(target, nil)
	if err != nil {
		log.Fatalf("wschan.Client dial: %s error %s", target, err.Error())
	}

	defer wconn.Close()

	PipeBinaryIgnoreText(ctx, wconn, send, receive)

}

// Host listens for a ws connection, then relays messages bidirectionally over the channels
// until the context ends. Only one connection is accepted at a time.
func Host(ctx context.Context, listen int, send, receive chan []byte) {

	server := http.NewServeMux()

	busy := false

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if busy {
			//limit connections to one at a time (reject additional connections)
			return
		}

		busy = true
		defer func() {
			busy = false
		}()

		wconn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("internal/wschan.Host(HandleFunc; upgrade): %s", err)
			return
		}
		defer wconn.Close()

		PipeBinaryIgnoreText(ctx, wconn, send, receive)

	})

	addr := ":" + strconv.Itoa(listen)

	h := &http.Server{Addr: addr, Handler: server}

	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Fatalf("internal/wschan.Host(ListenAndServe): %s ", err.Error())
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := h.Shutdown(ctx)
	if err != nil {
		log.Errorf("internal/wschan.Host(ListenAndServe.Shutdown): %s", err.Error())
	}
	log.Trace("internal/wschan.Host is done")

}

// PipeBinaryIgnoreText relays binary data between net.Conn and websocket.Conn
// pass websocket.Conn with pointer due to mutex
// Use of NextReader/Writer modified from https://github.com/hazaelsan/ssh-relay/blob/master/session/corprelay/corprelay.go
func PipeBinaryIgnoreText(ctx context.Context, wconn *websocket.Conn, send, receive chan []byte) {

	// write from wconn to receive channel
	go func() {

		for {

			err := func() error {

				select {
				case <-ctx.Done():
					return fmt.Errorf("context done")
				default:

					t, r, err := (*wconn).NextReader()
					if err != nil {
						return fmt.Errorf("NextReader() error: %w", err)
					}

					b := new(bytes.Buffer)

					n, err := b.ReadFrom(r)
					if err != nil {
						return err
					}

					log.Tracef("internal/ws.Host(ws->chan): read %v bytes", n)

					switch t {
					case websocket.BinaryMessage:
						// TODO what does this conversion achieve?
						receive <- b.Bytes()
					case websocket.TextMessage: //non-SSH message (none implemented as yet anyway)
						log.WithField("msg", b.String()).Info("Ignoring TEXT")

					default:
						return fmt.Errorf("unsupported message type: %v", t)
					}
					return nil
				}
			}()
			if err != nil {
				log.WithField("error", err.Error).Error("internal/ws.Host(ws->chan)")
				return
			}
		}
	}()

	// write from write channel to wconn
	go func() {
		for {
			err := func() error {
				select {
				case <-ctx.Done():
					return fmt.Errorf("context done")
				case data, ok := <-send:

					if !ok {
						return fmt.Errorf("send channel closed")
					}

					w, err := (*wconn).NextWriter(websocket.BinaryMessage)
					if err != nil {
						return fmt.Errorf("NextWriter() error: %w", err)
					}
					defer w.Close()

					n, err := w.Write(data)
					log.Tracef("internal/ws.Host(chan->ws): wrote %v bytes", n)
					return err
				}

			}()

			if err != nil {
				log.WithField("error", err.Error).Error("internal/ws.Host(tcp->ws) error")
				return
			}
		}

	}()

	<-ctx.Done()

}
