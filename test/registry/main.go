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

	//eureka := dc.NewEurekaClient(
	//	dc.NewEurekaClientConfig().
	//		WithUrl(EurekaUrl).
	//		WithAuth(EurekaLogin, EurekaPassword).
	//		WithRefresh(time.Second * 10).
	//		WithHeartBeat(time.Second * 5).
	//		WithApplication("registry-eureka-test"),
	//)
	//eureka.Connect()

	disco := dc.NewDiscoClient(
		dc.NewDiscoClientConfig().
			WithUrl(DiscoUrl).
			WithApplication("registry-disco-test").
			WithRetry(3, time.Millisecond*500).
			WithTimeout(time.Second*2).
			WithBasicAuth(DiscoLogin, DiscoPassword),
	)
	disco.Connect()

	reg := registry.NewServiceRegistry(
		registry.WithRefreshDelay(time.Second*10),
		registry.WithRefreshInterval(time.Second*10),
		//registry.WithSource(eureka),
		registry.WithSource(disco),
	)
	reg.Start()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	tm := time.NewTicker(time.Second * 5)
	//var i = 0
	for {
		select {
		case <-tm.C:
			for _, v := range reg.List() {
				fmt.Printf("known remote: %s %s %s\n", v.App, v.Status, v)
			}
			fmt.Println("-----------------------------------------------------------")
			fmt.Println()
		case <-quitChn:
			//if i >= 1 {
			tm.Stop()
			//time.Sleep(time.Millisecond * 100)
			return
			//}
			//i++
		}
	}

}
