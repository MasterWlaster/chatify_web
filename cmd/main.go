package main

import (
	"chat_web_client/pkg/handler"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	router := gin.New()
	router.LoadHTMLGlob("resources/html/**/*")
	router.Static("../../static", "./resources/static")

	port := os.Getenv("PORT")
	if port == "" {
		port = "12345"
	}

	go func() {
		if err := http.ListenAndServe(":"+port, handler.InitRoutes(router)); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("App Shutting Down")
}
