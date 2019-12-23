package main

import (
	"clienthandler"
	"context"
	"flag"
	"listener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Start signal handling early (avoid case when signals are delivered before handler installed)
	signaler := make(chan os.Signal, 1)
	signal.Notify(signaler, syscall.SIGINT)
	signal.Notify(signaler, syscall.SIGHUP)
	signal.Notify(signaler, syscall.SIGTERM)

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.Println("main()")

	port := flag.Int("port", 8080, "listen port")
	flag.Parse()

	l := listener.New(port)
	clienthandler.AddHandlers(&l)
	go func() {
		err := l.Server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Failed: %v", err)
		}
		log.Print("Shutting down")

	}()

	// This is blocking:
	select {
	case signal := <-signaler:
		log.Printf("Got signal: %v\n", signal)
	}
	l.Server.Shutdown(context.Background())
	log.Println("Exiting")
}
