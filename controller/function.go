package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SpeedVan/function-manager/service"
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/client/httpclient/gitlab"
	"github.com/SpeedVan/go-common/config"
	"github.com/gorilla/mux"

	"k8s.io/client-go/kubernetes"
)

// Function todo
type Function struct {
	web.Controller
	FuncService *service.Function
}

// FunctionNew todo
func FunctionNew(cfg config.Config, k8sClient *kubernetes.Clientset) *Function {
	return &Function{
		FuncService: service.FunctionNew(cfg, k8sClient),
	}
}

// GetRoute todo
func (s *Function) GetRoute() web.RouteMap {
	items := []*web.RouteItem{
		&web.RouteItem{Path: "/{namespace}/{runtimeName}/all_config", Method: "GET", HandleFunc: s.GetFuncConfig},
		&web.RouteItem{Path: "/{namespace}/{runtimeName}/{functionName}", Method: "POST", HandleFunc: s.Deploy},
		&web.RouteItem{Path: "/{namespace}/{runtimeName}/{functionName}", Method: "DELETE", HandleFunc: s.Undeploy},
		&web.RouteItem{Path: "/compatible/{namespace}/{runtimeName}/{functionName}", Method: "POST", HandleFunc: s.ComDeploy},
		&web.RouteItem{Path: "/compatible/{namespace}/{runtimeName}/{functionName}", Method: "DELETE", HandleFunc: s.Undeploy},
	}

	return web.NewRouteMap(items...)
}

// GetFuncConfig todo
func (s *Function) GetFuncConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	runtimeName := vars["runtimeName"]
	m := s.FuncService.GetFuncConfig(namespace, runtimeName)
	bytes, _ := json.Marshal(m)
	w.Write(bytes)
}

// Deploy todo
func (s *Function) Deploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	runtimeName := vars["runtimeName"]
	functionName := vars["functionName"]
	bytes, _ := ioutil.ReadAll(r.Body)
	funcConfig := make(map[string]string)
	cfg := make(map[string]map[string]string)
	if err := json.Unmarshal(bytes, &funcConfig); err != nil {
		http.Error(w, fmt.Sprintf("Error : %s", err.Error()), http.StatusInternalServerError)
		return
	}

	cfg[functionName] = funcConfig
	s.FuncService.Deploy(namespace, runtimeName, cfg)
}

// Undeploy todo
func (s *Function) Undeploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	runtimeName := vars["runtimeName"]
	functionName := vars["functionName"]
	s.FuncService.Undeploy(namespace, runtimeName, functionName)
}

// ComDeploy todo
func (s *Function) ComDeploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	runtimeName := vars["runtimeName"]
	functionName := vars["functionName"]
	bytes, _ := ioutil.ReadAll(r.Body)
	funcConfig := make(map[string]string)
	cfg := make(map[string]map[string]string)
	if err := json.Unmarshal(bytes, &funcConfig); err != nil {
		http.Error(w, fmt.Sprintf("Error : %s", err.Error()), http.StatusInternalServerError)
		return
	}
	entrypoint := funcConfig["entrypoint"]

	cfg[functionName] = url2funcConfig(entrypoint)
	s.FuncService.Deploy(namespace, runtimeName, cfg)
}

func url2funcConfig(mainURL string) map[string]string {
	urlInfo := gitlab.URLParse(mainURL)

	return map[string]string{
		"sys_path":   "/mnt/dav/" + urlInfo.Domain + "/" + urlInfo.Group + "+" + urlInfo.Project + "/" + urlInfo.Sha + "/" + urlInfo.SysPath,
		"entrypoint": urlInfo.EntryPoint,
	}
}
