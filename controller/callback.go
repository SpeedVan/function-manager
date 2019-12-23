package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SpeedVan/function-manager/service"
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/config"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
)

// Callback todo
type Callback struct {
	web.Controller
	SwitchSvc *service.SwitchService
}

// CallbackNew todo
func CallbackNew(cfg config.Config, k8sClient *kubernetes.Clientset) *Callback {

	return &Callback{
		SwitchSvc: service.SwitchServiceNew(cfg, k8sClient),
	}
}

// GetRoute todo
func (s *Callback) GetRoute() web.RouteMap {
	items := []*web.RouteItem{
		&web.RouteItem{Path: "/callback/switchService/{namespace}/{runtimeName}", Method: "POST", HandleFunc: s.SwitchService},
		&web.RouteItem{Path: "/callback/switchServiceWithClean/{namespace}/{runtimeName}", Method: "POST", HandleFunc: s.SwitchServiceWithClean},
	}

	return web.NewRouteMap(items...)
}

// SwitchService todo
func (s *Callback) SwitchService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	runtimeName := vars["runtimeName"]
	m := make(map[string]map[string]string)
	bytes, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(bytes, &m); err != nil {
		http.Error(w, fmt.Sprintf("Error : %s", err.Error()), http.StatusInternalServerError)
		return
	}
	s.SwitchSvc.SwitchAllToRuntime(namespace, runtimeName, m)
}

// SwitchServiceWithClean todo
func (s *Callback) SwitchServiceWithClean(w http.ResponseWriter, r *http.Request) {
	s.SwitchService(w, r)
}
