package registry

import (
	"github.com/slink-go/discovery/dc"
	"github.com/slink-go/logging"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
)

type serviceRegistry struct {
	serviceDirectory *ringBuffers // TODO: при обновлении собьётся ringBuffer; можно ли что-то с этим сделать? стоит ли это делать? [UPD: shuffle ring buffer?]
	clients          []dc.Client
	mutex            sync.RWMutex
	logger           logging.Logger
	sigChn           chan os.Signal
	refreshDelay     time.Duration
	refreshInterval  time.Duration
}

func NewServiceRegistry(options ...ServiceRegistryOption) ServiceRegistry {

	registry := serviceRegistry{
		serviceDirectory: nil,
		logger:           logging.GetLogger("discovery-registry"),
		sigChn:           make(chan os.Signal),
		clients:          nil,
		refreshDelay:     5 * time.Second,
		refreshInterval:  60 * time.Second,
	}

	for _, option := range options {
		option.Apply(&registry)
	}

	return &registry
}
func (sr *serviceRegistry) Start() {
	signal.Notify(sr.sigChn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	for _, client := range sr.clients {
		client.Connect()
		if client != nil && client.NotificationsChn() != nil {
			go func() {
				for range client.NotificationsChn() {
					sr.logger.Debug(" > update detected... refreshing")
					sr.doRefresh()
				}
			}()
		}
	}

	go sr.refresh()
}

func (sr *serviceRegistry) refresh() {
	timer := time.NewTimer(sr.refreshDelay)
	for {
		select {
		case <-sr.sigChn:
			sr.logger.Info("stop")
			timer.Stop()
			sr.serviceDirectory = nil
			return
		case <-timer.C:
			sr.doRefresh()
		}
		timer.Reset(sr.refreshInterval)
	}
}

func (sr *serviceRegistry) doRefresh() {
	remotes := make(map[string]map[string]dc.Remote)
	directory := createRingBuffers()
	for _, client := range sr.clients {
		if client == nil {
			continue
		}
		sr.getRemotes(remotes, client)
	}
	for k, list := range sr.filterRemotes(remotes) {
		directory.New(k, len(list))
		for _, url := range list {
			v := url
			directory.Set(k, &v)
		}
	}
	sr.mutex.Lock()
	sr.serviceDirectory = directory
	sr.mutex.Unlock()
}
func (sr *serviceRegistry) getRemotes(destination map[string]map[string]dc.Remote, client dc.Client) {
	for _, instance := range client.Services().List() {
		appId := strings.ToUpper(instance.App)
		if _, ok := destination[appId]; !ok {
			destination[appId] = make(map[string]dc.Remote, 0)
		}
		m := destination[appId]
		sr.logger.Trace("[%T] add instance: %s %s", client, appId, instance)
		if _, ok := destination[appId][instance.String()]; !ok {
			m[instance.String()] = instance
		}
	}
}
func (sr *serviceRegistry) filterRemotes(source map[string]map[string]dc.Remote) map[string][]dc.Remote {
	remotes := make(map[string][]dc.Remote)
	for _, service := range source {
		for _, instance := range service {
			appId := strings.ToUpper(instance.App)
			if _, ok := remotes[appId]; !ok {
				remotes[appId] = make([]dc.Remote, 0)
			}
			remotes[appId] = append(remotes[appId], instance)
		}
	}
	return remotes
}

func (sr *serviceRegistry) Get(serviceName string) (string, error) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()
	v, ok := sr.serviceDirectory.Next(serviceName)
	if !ok {
		return "", NewErrServiceUnavailable(serviceName)
	}
	url := v.Next().Value.(*dc.Remote)
	return (*url).String(), nil
}
func (sr *serviceRegistry) List() []dc.Remote {
	result := make([]dc.Remote, 0)
	for _, v := range sr.serviceDirectory.List() {
		vv := v.(*dc.Remote)
		result = append(result, *vv)
	}
	slices.SortFunc(result, dc.Remote.Compare)
	return result
}

//	func (sr *serviceRegistry) AddSource(dc dc.Client) {
//		if dc != nil {
//			sr.clients = append(sr.clients, dc)
//		}
//	}
