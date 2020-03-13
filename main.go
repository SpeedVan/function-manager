package main

import (
	"fmt"

	"github.com/SpeedVan/function-manager/controller"
	repoImpl "github.com/SpeedVan/function-manager/repository/impl"
	"github.com/SpeedVan/function-manager/service"
	"github.com/SpeedVan/go-common-kubernetes/client"
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/config/env"
	"github.com/SpeedVan/go-common/log"
)

func main() {

	if cfg, err := env.LoadAllWithoutPrefix("CRD_"); err == nil {
		logger := log.NewCommon(log.Debug)

		app := web.New(cfg, logger)

		clientset, err := client.GetK8sClient()
		if err != nil {
			fmt.Println(err)
			// return
		}
		// k8sCli, err := clientset.GetClient()
		// if err != nil {
		// 	fmt.Println(err)
		// 	// return
		// }

		repo := repoImpl.NewFuncRepo(clientset.GetRestClient, clientset.GetExtClient)

		svc := service.NewFunctionService(repo)
		scaleSvc, err := service.NewScaleService(clientset.GetRestClient)
		if err != nil {
			fmt.Println(err)
			// return
		}
		app.HandleController(controller.NewFunctionControllerBeta(svc))
		app.HandleController(controller.NewFunctionScaleControllerBeta(scaleSvc))

		err = app.Run(log.Debug)
		if err != nil {
			fmt.Println(err)
		}
	}
}
