package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/SpeedVan/function-manager/service"
	"github.com/SpeedVan/go-common/app/rest"
)

// FunctionScaleControllerBeta todo
type FunctionScaleControllerBeta struct {
	rest.Controller
	scaleSvc *service.Scale
}

// NewFunctionScaleControllerBeta todo
func NewFunctionScaleControllerBeta(svc *service.Scale) *FunctionScaleControllerBeta {

	return &FunctionScaleControllerBeta{
		scaleSvc: svc,
	}
}

// GetRoute todo
func (s *FunctionScaleControllerBeta) GetRoute() rest.RouteMap {
	items := []*rest.RouteItem{
		&rest.RouteItem{Path: "/api/beta/{namespace}/{function}/scale", HandleFunc: s.Inc, Method: "POST"},
		// &rest.RouteItem{Path: "/api/beta/{namespace}", HandleFunc: s.List},
		&rest.RouteItem{Path: "/api/beta/{namespace}/{function}/scale", HandleFunc: s.Get, Method: "GET"},
	}

	return rest.NewRouteMap(items...)
}

// Inc todo
func (s *FunctionScaleControllerBeta) Inc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	function := vars["function"]
	scaleInc := 0
	if inc, err := strconv.Atoi(r.URL.Query().Get("inc")); err == nil {
		scaleInc = inc
	}
	_, err := s.scaleSvc.Inc(namespace, function, scaleInc)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

// List todo
// func (s *FunctionScaleControllerBeta) List(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	namespace := vars["namespace"]
// 	query := r.URL.Query()
// 	continueBase64 := query.Get("continue")
// 	limit := 0
// 	if n, err := strconv.Atoi(query.Get("limit")); err == nil {
// 		limit = n
// 	}
// 	result := s.FuncSvc.List(namespace, continueBase64, int64(limit))
// 	bs, err := json.Marshal(result)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(bs)
// }

// Get todo
func (s *FunctionScaleControllerBeta) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	function := vars["function"]
	result, err := s.scaleSvc.Get(namespace, function)
	if err != nil {
		fmt.Println(err)
		return
	}
	bs, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
