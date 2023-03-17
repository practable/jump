/*
package pipe relays a connection from one port to another.

It is intended for testing how we handle ssh connections
by making a direct pipe from one port to another, isolating
our investigation down to how the data crosses the boundary
into our code and out again, both for the convenience and
partitioning out the influence (if any) of the hub.

*/

package pipe

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/phayes/freeport"
	"github.com/practable/jump/internal/chanio"
	"github.com/practable/jump/internal/tcpconnect"
	"github.com/practable/jump/internal/ws"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Listen int
	Target int
}

type Pipe struct {
	Left   chan []byte
	Listen int
	Right  chan []byte
	Target int
}

func New(config Config) *Pipe {

	return &Pipe{
		Left:   make(chan []byte, 1),
		Listen: config.Listen,
		Right:  make(chan []byte, 1),
		Target: config.Target,
	}

}

/*

          |------|                   |--------|
          \      \  ----RIGHT--->    \        \
 Target   | host |                   | client |  Listen
 Port     \      \                   \        \  Port
          |------|  <---LEFT-----    |--------|
*/

func (p *Pipe) Run(ctx context.Context) {

	url := ":" + strconv.Itoa(p.Listen)

	listener := tcpconnect.New()

	go listener.Listen(ctx, url, listenHandler(p.Left, p.Right))
	go targetHandler(ctx, p.Target, p.Left, p.Right)

	<-ctx.Done()

}

func (p *Pipe) RunDevelop(ctx context.Context) {
	p.RunWS(ctx)
}

/*
                               using freeport
          |--------|          |--------|             |--------|           |--------|
          \        \  ---->   \        \             \        \  ---->    \        \
 Target   | target | io.Copy  | host   |  websocket  | client | io.Copy   | client |  Listen
 Port     \        \  <----   \        \   <--->     \        \  <-----   \        \  Port
          |--------|          |--------|             |--------|           |--------|
*/
func (p *Pipe) RunWS(ctx context.Context) {
	log.Info("WS Version")

	targetURL := "localhost:" + strconv.Itoa(p.Target)

	wsPort, err := freeport.GetFreePort()
	if err != nil {
		return
	}
	wsURL := "ws://localhost:" + strconv.Itoa(wsPort)

	go ws.Host(ctx, wsPort, targetURL)
	go ws.Client(ctx, p.Listen, wsURL)

	<-ctx.Done()

}

/*

          |--------|          |--------|             |--------|           |--------|
          \        \  ---->   \        \ ----a--->   \        \  ---->    \        \
 Target   | target | io.Copy  | left   |  channels   | right  | io.Copy   | client |  Listen
 Port     \        \  <----   \        \             \        \  <-----   \        \  Port
          |--------|          |--------| <---b-----  |--------|           |--------|
*/
// RunDevelop represents a pipe with an intermediate communication over channels
// which, if it works, proves that we can relay because our websocket clients
// have a channel-based interface (TODO consider if that should stay the case?)
func (p *Pipe) RunChannel(ctx context.Context) {
	log.Info("Channel Version")
	listenURL := ":" + strconv.Itoa(p.Listen)
	targetURL := "localhost:" + strconv.Itoa(p.Target)

	incoming, err := net.Listen("tcp", listenURL)

	if err != nil {
		log.Errorf("could not start server on %s: %v", listenURL, err)
	}

	log.Info("server running on %s\n", listenURL)

	//https://www.zupzup.org/go-port-forwarding/index.html
	client, err := incoming.Accept()
	if err != nil {
		log.Fatal("could not accept client connection", err)
	}
	defer client.Close()
	log.Infof("client '%v' connected!\n", client.RemoteAddr())

	target, err := net.Dial("tcp", targetURL)
	if err != nil {
		log.Fatal("could not connect to target", err)
	}
	defer target.Close()
	log.Infof("connection to server %v established!\n", target.RemoteAddr())

	a := make(chan []byte, 256)
	b := make(chan []byte, 256)

	// TODO check a, b have same sense as diagram (functionally, does not matter though)
	left := chanio.New(ctx, a, b, "left")
	right := chanio.New(ctx, b, a, "right")

	go func() { io.Copy(target, left) }()
	go func() { io.Copy(left, target) }()

	go func() { io.Copy(client, right) }()
	go func() { io.Copy(right, client) }()

	<-ctx.Done()

}

func (p *Pipe) RunCopy(ctx context.Context) {

	listenURL := ":" + strconv.Itoa(p.Listen)
	targetURL := "localhost:" + strconv.Itoa(p.Target)

	incoming, err := net.Listen("tcp", listenURL)

	if err != nil {
		log.Fatalf("could not start server on %s: %v", listenURL, err)
	}

	fmt.Printf("server running on %s\n", listenURL)

	//https://www.zupzup.org/go-port-forwarding/index.html
	client, err := incoming.Accept()
	if err != nil {
		log.Fatal("could not accept client connection", err)
	}
	defer client.Close()
	fmt.Printf("client '%v' connected!\n", client.RemoteAddr())

	target, err := net.Dial("tcp", targetURL)
	if err != nil {
		log.Fatal("could not connect to target", err)
	}
	defer target.Close()
	fmt.Printf("connection to server %v established!\n", target.RemoteAddr())

	go func() { io.Copy(target, client) }()
	go func() { io.Copy(client, target) }()

	<-ctx.Done()

}

/*  HANDLERS FOR RUN (original method)  */

// listenHandler handles the external client
// it receives data from the external connection and forwards it internally over channel left
// it accepts internal data from channel right and sends to the external connection
func listenHandler(left, right chan []byte) func(context.Context, *tcpconnect.TCPconnect) {

	return func(ctx context.Context, c *tcpconnect.TCPconnect) {

		timeout := time.Second

		id := "listenHandler(" + c.ID + ")"

		go c.HandleConn(ctx, *c.Conn) // messages on c.In and c.Out

		go func() {

			for {
				select {
				case <-ctx.Done():
					log.Infof("%s: context cancelled; done", id)
					return
				case data, ok := <-c.In:
					if !ok {
						log.Debugf("%s: local channel error, closing local read pump", id)
						return
					}
					size := len(data)
					select {
					case left <- data:
						log.WithField("size", size).Debugf("%s: sent %d-bytes left", id, size)
					case <-time.After(timeout):
						log.WithField("size", size).Debugf("%s: timeout waiting to send %d-bytes left", id, size)
					}

				}

			}

		}()

		go func() {

			for {
				select {
				case <-ctx.Done():
					log.Infof("%s: context cancelled; done", id)
					return
				case data, ok := <-right:
					if !ok {
						log.Debugf("%s: relay channel error, closing relay read pump", id)
						return
					}
					size := len(data)
					select {
					case c.Out <- data:
						log.WithField("size", size).Debugf("%s: sent %d-bytes right", id, size)
					case <-time.After(timeout):
						log.WithField("size", size).Debugf("%s: timeout waiting to send %d-bytes right", id, size)
					}
				}
			}

		}()

		<-ctx.Done()
		log.Infof("%s: done", id)
	}
}

// targetHandler handles the external target, typically port 22 to the localhost
// it receives data from the target and forwards it internally on channel right
// it accepts data internally from channel left and sends it externally to the target
func targetHandler(ctx context.Context, target int, left, right chan []byte) {

	timeout := 1 * time.Second

	id := "targetHandler(" + uuid.New().String()[0:6] + ")"

	log.WithFields(log.Fields{"target": target}).Infof("%s: connecting to target", id)

	t := tcpconnect.New()

	url := "localhost:" + strconv.Itoa(target)

	go t.Dial(ctx, url)

	log.WithFields(log.Fields{"target": target}).Tracef("%s: dialling target", id)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-t.In:
				log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: GOT %d-byte message from target", id, len(data))
				if !ok {
					return
				}
				select {
				case right <- data:
					log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: SENT %d-byte message to the right", id, len(data))
				case <-time.After(timeout):
					log.Error("timeout sending message ")
					return
				}
			}

		}
	}()

	go func() {
		for {
			select {

			case <-ctx.Done():
				return
			case data, ok := <-left:
				log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: GOT %d-byte message from the left", id, len(data))
				if !ok {
					return
				}
				select {
				case t.Out <- data:
					log.WithFields(log.Fields{"target": target, "size": len(data)}).Tracef("%s: SENT %d-byte Message to target", id, len(data))
				case <-time.After(timeout):

					log.Error("timeout sending message ")
					return
				}
			}

		}
	}()

	<-ctx.Done()
	log.WithFields(log.Fields{"target": target}).Infof("%s: DONE", id)
}
