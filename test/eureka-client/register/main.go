package main

import (
	"github.com/slink-go/discovery/dc"
	"github.com/slink-go/discovery/test/eureka-client/common"
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
			WithUrl(common.EurekaUrl).
			WithAuth(common.EurekaLogin, common.EurekaPassword).
			WithHeartBeat(time.Second * 5).
			WithApplication("eureka-test-register"),
	)
	eureka.Connect()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	var i = 0
	for {
		select {
		case <-quitChn:
			if i >= 1 {
				return
			}
			i++
		}
	}

}
