package main

import (
	"log"

	"github.com/ednaldo-dilorenzo/iappointment/server"
)

func main() {
	svr := &server.Server{
		App: server.InitFiberApplication(),
	}

	waitforShutdownInterrupt := svr.Start()
	<-waitforShutdownInterrupt

	log.Println("Shutting Down Server..")

	svr.ShutdownGracefully()
}
