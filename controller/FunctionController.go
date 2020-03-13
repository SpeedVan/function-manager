package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"k8s.io/client-go/tools/cache"

	"github.com/SpeedVan/function-manager/service"
	"github.com/SpeedVan/go-common/app/rest"
)

// FunctionControllerBeta todo
type FunctionControllerBeta struct {
	rest.Controller
	FuncSvc *service.Function
	Cache   cache.ThreadSafeStore
}

// NewFunctionControllerBeta todo
func NewFunctionControllerBeta(svc *service.Function) *FunctionControllerBeta {

	return &FunctionControllerBeta{
		FuncSvc: svc,
	}
}

// GetRoute todo
func (s *FunctionControllerBeta) GetRoute() rest.RouteMap {
	items := []*rest.RouteItem{
		&rest.RouteItem{Path: "/api/beta/{namespace}/{function}", HandleFunc: s.Create, Method: "POST"},
		&rest.RouteItem{Path: "/api/beta/{namespace}/{function}", HandleFunc: s.Update, Method: "PUT"},
		&rest.RouteItem{Path: "/api/beta/{namespace}/{function}", HandleFunc: s.Delete, Method: "DELETE"},
		&rest.RouteItem{Path: "/api/beta/{namespace}", HandleFunc: s.List},
		&rest.RouteItem{Path: "/api/beta/{namespace}/{function}", HandleFunc: s.Get, Method: "GET"},
	}

	return rest.NewRouteMap(items...)
}

// Create todo
func (s *FunctionControllerBeta) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	function := vars["function"]
	s.FuncSvc.Create(namespace, function)
	w.Header().Set("Content-Type", "application/json")
}

// Update todo
func (s *FunctionControllerBeta) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	function := vars["function"]
	s.FuncSvc.Update(namespace, function)
	w.Header().Set("Content-Type", "application/json")
}

// Delete todo
func (s *FunctionControllerBeta) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	function := vars["function"]
	s.FuncSvc.Delete(namespace, function)
	w.Header().Set("Content-Type", "application/json")
}

// List todo
func (s *FunctionControllerBeta) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	query := r.URL.Query()
	continueBase64 := query.Get("continue")
	limit := 0
	if n, err := strconv.Atoi(query.Get("limit")); err == nil {
		limit = n
	}
	result := s.FuncSvc.List(namespace, continueBase64, int64(limit))
	bs, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

// Get todo
func (s *FunctionControllerBeta) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	function := vars["function"]
	result := s.FuncSvc.Get(namespace, function)
	bs, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
