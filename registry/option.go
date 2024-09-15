package registry

import (
	"github.com/slink-go/discovery/dc"
	"time"
)

type ServiceRegistryOption interface {
	Apply(registry *serviceRegistry)
}

// region - refresh delay option

type RefreshDelayOption struct {
	delay time.Duration
}

func (o *RefreshDelayOption) Apply(r *serviceRegistry) {
	r.refreshDelay = o.delay
}
func WithRefreshDelay(delay time.Duration) ServiceRegistryOption {
	return &RefreshDelayOption{delay}
}

// endregion
// region - refresh interval option

type RefreshIntervalOption struct {
	delay time.Duration
}

func (o *RefreshIntervalOption) Apply(r *serviceRegistry) {
	r.refreshInterval = o.delay
}
func WithRefreshInterval(delay time.Duration) ServiceRegistryOption {
	return &RefreshIntervalOption{delay}
}

// endregion
// region - source options

type sourceOption struct {
	client dc.Client
}

func (o *sourceOption) Apply(r *serviceRegistry) {
	if o.client != nil {
		r.clients = append(r.clients, o.client)
	}
}
func WithSource(client dc.Client) ServiceRegistryOption {
	return &sourceOption{client}
}

// endregion
