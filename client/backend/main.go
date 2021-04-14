package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

// run starts the server that serves the client
func run() error {
	/*if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}*/
	//	l, err := net.Listen("tcp", os.Args[1])
	port := "8080" //os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	sh := newServerHandler()
	s := &http.Server{
		Handler:      sh,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 35,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return s.Shutdown(ctx)
}
