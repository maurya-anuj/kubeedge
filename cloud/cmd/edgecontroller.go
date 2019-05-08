package main

import (
	_ "github.com/kubeedge/kubeedge/cloud/pkg/cloudhub"
	_ "github.com/kubeedge/kubeedge/cloud/pkg/controller"
	_ "github.com/kubeedge/kubeedge/cloud/pkg/devicecontroller"
        _ "github.com/kubeedge/kubeedge/cloud/pkg/http_router"

	"github.com/kubeedge/beehive/pkg/core"
)

func main() {
	core.Run()
}
