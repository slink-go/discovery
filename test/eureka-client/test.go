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

	eureka := dc.NewEurekaClient(
		dc.NewEurekaClientConfig().
			WithUrl("http://127.0.0.1:9992/eureka").
			WithAuth("discovery", "password").
			WithRefresh(time.Second * 10).
			WithHeartBeat(time.Second * 5).
			WithApplication("eureka-test"),
	)
	eureka.Connect()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	tm := time.NewTimer(time.Second * 5)
	var i = 0
	for {
		select {
		case <-tm.C:
			eureka.Services().List()
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
