package main

import (
	"fmt"
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
			WithRefresh(time.Second * 10).
			WithHeartBeat(time.Second * 5).
			WithApplication("eureka-test-combined"),
	)
	eureka.Connect()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	tm := time.NewTicker(time.Second * 5)
	var i = 0
	for {
		select {
		case <-tm.C:
			for _, v := range eureka.Services().List() {
				fmt.Printf("known remote: %s %s %s\n", v.App, v.Status, v)
			}
			fmt.Println("-----------------------------------------------------------")
			fmt.Println()
		case <-quitChn:
			if i >= 1 {
				tm.Stop()
				time.Sleep(time.Millisecond * 100)
				return
			}
			i++
		}
	}

}
