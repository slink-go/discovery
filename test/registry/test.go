package main

import (
	"fmt"
	"github.com/slink-go/discovery/dc"
	"github.com/slink-go/discovery/registry"
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
			WithApplication("registry-eureka-test"),
	)
	eureka.Connect()

	disco := dc.NewDiscoClient(
		dc.NewDiscoClientConfig().
			WithUrl("http://127.0.0.1:8771").
			WithApplication("registry-disco-test").
			WithRetry(3, time.Millisecond*500).
			WithTimeout(time.Second*2).
			WithBasicAuth("disco", "disco"),
	)
	disco.Connect()

	reg := registry.NewServiceRegistry(
		registry.WithRefreshDelay(time.Second*10),
		registry.WithRefreshInterval(time.Second*5),
		registry.WithSource(eureka),
		registry.WithSource(disco),
	)
	reg.Start()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	tm := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-tm.C:
			for _, r := range reg.List() {
				fmt.Println(r)
			}
		case <-quitChn:
			tm.Stop()
			time.Sleep(time.Millisecond * 100)
			return
		}
	}
}
