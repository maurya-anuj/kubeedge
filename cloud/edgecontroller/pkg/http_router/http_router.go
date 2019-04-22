package http_router

import (
	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/beehive/pkg/core/context"

	"github.com/kubeedge/beehive/pkg/common/log"
)

type http_router struct {
	context *context.Context
}

func init() {
	core.Register(&http_router{})
}

func (a *http_router) Name() string {
	return "http_router"
}

func (a *http_router) Group() string {
	return "http_router"
}

func (a *http_router) Start(c *context.Context) {
	log.LOGGER.Infof("http_route.START called \n\n\n\n\n\n")
	a.context = c
	main(c)
}

func (a *http_router) Cleanup() {
	a.context.Cleanup(a.Name())
}

