package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwdcth/mailsender/g"
	"github.com/dwdcth/mailsender/http"
	"github.com/dwdcth/mailsender/sender"
)

func main() {
	cfg := flag.String("c", "cfg.json", "config file")
	vsn := flag.Bool("v", false, "show version")
	flag.Parse()

	if *vsn {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.LoadConfig(*cfg)
	sender.Start()

	http.Start()

	wait_signal()
}

func wait_signal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sc:
		os.Exit(0)
	}
}
