package main

import (
	"fmt"
	"github.com/slink-go/discovery/dc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	os.Setenv("GO_ENV", "dev")
	os.Setenv("LOGGING_LEVEL_ROOT", "TRACE")

	disco := dc.NewDiscoClient(
		dc.NewDiscoClientConfig().
			WithUrl(DiscoUrl).
			WithApplication("disco-test-register").
			WithBasicAuth(DiscoLogin, DiscoPassword),
	)
	chn, _ := disco.Connect()

	mutex := sync.Mutex{}

	go func() {
		for _ = range chn {
			mutex.Lock()
			fmt.Println("--- update detected ---------------------------------------")
			for _, v := range disco.Services().List() {
				fmt.Printf("known remote: %s %s %s\n", v.App, v.Status, v)
			}
			fmt.Println("--- update processed --------------------------------------")
			fmt.Println()
			mutex.Unlock()
		}
	}()

	quitChn := make(chan os.Signal)
	signal.Notify(quitChn, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	tm := time.NewTicker(time.Second * 10)
	//var i = 0
	for {
		select {
		case <-tm.C:
			mutex.Lock()
			fmt.Println("--- scheduled output start --------------------------------")
			for _, v := range disco.Services().List() {
				fmt.Printf("known remote: %s %s %s\n", v.App, v.Status, v)
			}
			fmt.Println("--- scheduled output compete ------------------------------")
			fmt.Println()
			mutex.Unlock()
		case <-quitChn:
			//if i >= 1 {
			//	return
			//}
			//i++
			tm.Stop()
			fmt.Println("--- SHUTDOWN ----------------------------------------------")
			return
		}
	}
}
