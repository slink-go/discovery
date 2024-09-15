package registry

import (
	"github.com/slink-go/discovery/dc"
)

type ServiceRegistry interface {
	Start()
	Get(applicationId string) (string, error)
	List() []dc.Remote
}
