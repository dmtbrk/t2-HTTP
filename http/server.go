package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ortymid/t2-http/market"
)

// Run is a convenient function to start an http server with graceful shotdown.
func Run(port int, jwtAlg string, jwtSecret interface{}, m market.Interface) {
	handler := &Router{
		Market:    m,
		JWTAlg:    jwtAlg,
		JWTSecret: jwtSecret,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	idle := make(chan struct{})
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-done
		log.Println("Gracefully stopping...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("server Shutdown:", err)
		}
		close(idle)
		log.Println("Server stopped")
	}()

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("server ListenAndServe:", err)
		}
	}()
	log.Print("Server started at ", srv.Addr)

	<-idle
}
