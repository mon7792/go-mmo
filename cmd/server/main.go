package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nhooyr.io/websocket"
)

func main() {

	// create a listener
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	// create http server
	s := &http.Server{
		Handler:      websocketServer{},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8000")

	// errc channel to handle errors from the server
	errc := make(chan error, 1)
	go func() {
		// start the server
		errc <- s.Serve(listener)
	}()

	// sigs channel to handle signals from the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// select statement to handle errors and signals
	select {
	case err := <-errc:
		log.Println("Failed to serve: Error: ", err)
	case sig := <-sigs:
		log.Println("Terminating Signal: ", sig)
	}

	// context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown the server
	if err := s.Shutdown(ctx); err != nil {
		log.Println("Failed to shutdown server: Error: ", err)
	}

}

type websocketServer struct{}

func (s websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// transform the http connection to a websocket connection
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Failed to accept websocket connection: Error: ", err)
		return
	}

	ctx := context.Background()
	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)
	go ServeNetConn(conn)
}

func ServeNetConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Failed to close connection: Error: ", err)
		}
	}()

	const timeoutSec = 60 * time.Second
	timeout := make(chan uint8, 1)
	const StopTimeout uint8 = 0
	const ContTimeout uint8 = 1

	const MaxMsgSize = 4 * 1024
	// read all the messages from the connection
	go func() {
		msg := make([]byte, MaxMsgSize)
		for {
			n, err := conn.Read(msg)
			if err != nil {
				log.Println("Read Error: ", err)
				timeout <- StopTimeout // stop the timeout because of read error
				return
			}
			// Tick the timeout, so we can continue to wait for the next message
			timeout <- ContTimeout
			log.Println("Message: ", msg[:n])
		}
	}()

	// timeout manager
ExitTimeout:
	for {
		select {
		case res := <-timeout:
			if res == StopTimeout {
				log.Println("Manually stopping timeout manager")
				break ExitTimeout
			}
		case <-time.After(timeoutSec):
			log.Println("User time out")
			break ExitTimeout
		}
	}
}
