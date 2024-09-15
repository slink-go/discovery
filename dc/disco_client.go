package dc

import (
	"fmt"
	d "github.com/slink-go/disco-go"
	da "github.com/slink-go/disco/pkg/api"
	"github.com/slink-go/discovery"
	"github.com/slink-go/discovery/util"
	"github.com/slink-go/logging"
	"github.com/slink-go/util/env"
	"strings"
	"sync"
	"time"
)

func NewDiscoClient(config *discoConfig) Client {
	return &discoClient{
		config: *config,
		logger: logging.GetLogger("disco-client"),
	}
}

type discoClient struct {
	config          discoConfig
	mutex           sync.RWMutex
	client          d.DiscoClient
	logger          logging.Logger
	notificationChn chan struct{}
	connected       bool
}

func (c *discoClient) Connect(options ...interface{}) (chan struct{}, error) {

	if c.notificationChn != nil {
		return c.notificationChn, nil // already connected
	}

	if c.config.url == "" {
		return nil, fmt.Errorf("disco url is empty")
	}

	c.notificationChn = make(chan struct{}, 1)

	cfg := d.
		EmptyConfig().
		WithDisco([]string{c.config.url}).
		//WithMeta(...).
		//WithBreaker(c.config.breakThreshold).
		WithRetry(c.config.retryAttempts, c.config.retryDelay).
		WithTimeout(c.config.timeout).
		WithAuth(c.config.login, c.config.password).
		WithName(c.config.application).
		WithEndpoints([]string{fmt.Sprintf("%s://%s:%d", "http", c.config.getIP(), c.config.port)}).
		WithNotificationChn(c.notificationChn)

	go func() {
		for {
			clnt, err := d.NewDiscoHttpClient(cfg)
			if err != nil {
				c.logger.Warning("join error: %s", strings.TrimSpace(err.Error()))
				time.Sleep(env.DurationOrDefault(discovery.DiscoClientRetryInterval, 5*time.Second))
			} else {
				c.client = clnt
				return
			}
		}
	}()
	return c.notificationChn, nil
}
func (c *discoClient) Services() *Remotes {
	if c.client == nil {
		return nil
	}
	result := Remotes{}
	appId := strings.ToUpper(c.config.application)
	for _, v := range c.client.Registry().List() {
		if da.ClientStateUp == v.State() && v.ServiceId() != appId {
			ep, err := v.Endpoint(da.HttpEndpoint)
			if err != nil {
				c.logger.Warning("could not find HTTP endpoint for %s", v.ServiceId())
				continue
			}
			scheme, host, port := util.ParseEndpoint(ep)
			remote := Remote{
				App:    v.ServiceId(),
				Scheme: scheme,
				Host:   host,
				Port:   port,
				Status: v.State().String(),
			}
			c.logger.Trace("add %s: %s", v.ServiceId(), remote)
			result.Add(v.ServiceId(), remote)
		}
	}
	return &result
}
func (c *discoClient) NotificationsChn() chan struct{} {
	return c.notificationChn
}
