package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/ericluj/elog"
	"github.com/ericluj/ercache"
)

func main() {
	c := make(chan os.Signal, 1)
	ercache.NewServer()
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Infof("exit")
}
