package main

import (
	"github.com/slink-go/discovery/dc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	os.Setenv("GO_ENV", "dev")
	os.Setenv("LOGGING_LEVEL_ROOT", "TRACE")

	disco := dc.NewDiscoClient(
		dc.NewDiscoClientConfig().
			WithUrl("http://127.0.0.1:8771").
			WithApplication("disco-test").
			WithRetry(3, time.Millisecond*500).
			WithTimeout(time.Second*2).
			WithBasicAuth("disco", "disco"),
	)
	disco.Connect()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	tm := time.NewTimer(time.Second * 5)
	var i = 0
	for {
		select {
		case <-tm.C:
			disco.Services().List()
		case <-quitChn:
			if i > 1 {
				tm.Stop()
				time.Sleep(time.Millisecond * 100)
				return
			}
			i++
		}
	}
}
