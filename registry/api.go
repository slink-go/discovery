package registry

import (
	"github.com/slink-go/discovery/dc"
)

type ServiceRegistry interface {
	Start()
	AddSource(dc dc.Client)
	Get(applicationId string) (string, error)
	List() []dc.Remote
}
