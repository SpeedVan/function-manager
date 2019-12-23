package main

import (
	"fmt"
	orglog "log"

	"github.com/SpeedVan/function-manager/controller"
	"github.com/SpeedVan/function-manager/k8s"
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/config/env"
	"github.com/SpeedVan/go-common/log"
	lc "github.com/SpeedVan/go-common/log/common"
)

func main() {
	if cfg, err := env.LoadAllWithoutPrefix("MANAGER_"); err == nil {
		logger := lc.NewCommon(log.Debug) // this level control webapp init log level

		app := web.New(cfg, logger)

		clientset, err := k8s.GetK8sClient()

		if err != nil {
			orglog.Fatal(err)
		}
		client, err := clientset.GetClient()
		if err != nil {
			orglog.Fatal(err)
		}
		fmt.Println("kubeconfig is ok")

		app.HandleController(controller.FunctionNew(cfg, client))
		app.HandleController(controller.CallbackNew(cfg, client))
		orglog.Fatal(app.Run(log.Debug)) // this level control webapp runtime log level
	} else {
		orglog.Fatal(err)
	}
}
