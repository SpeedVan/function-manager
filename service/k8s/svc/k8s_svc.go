package svc

import (
	"fmt"
	"github.com/SpeedVan/function-manager/service/k8s"
	api_core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	api_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

// FuncServiceImpl todo
type FuncServiceImpl struct {
	k8s.FuncService
	K8sClient            *kubernetes.Clientset
	FuncRuntimeProxyPort int
}

// New todo
func New(k8sClient *kubernetes.Clientset, funcRuntimeProxyPort string) *FuncServiceImpl {
	_port, err := strconv.Atoi(funcRuntimeProxyPort)
	if err != nil {
		_port = 80
	}
	return &FuncServiceImpl{
		K8sClient:            k8sClient,
		FuncRuntimeProxyPort: _port,
	}
}

// Create todo
func (s *FuncServiceImpl) Create(namespace, runtimeName, functionName string) (*api_core_v1.Service, error) {
	service := &api_core_v1.Service{
		ObjectMeta: api_meta_v1.ObjectMeta{
			Name: functionName,
			Labels: map[string]string{
				"functionName": functionName,
				"mode":         "share",
			},
		},
		Spec: api_core_v1.ServiceSpec{
			Ports: []api_core_v1.ServicePort{
				{
					Name:       "runtime-proxy-port",
					Port:       int32(80),
					TargetPort: intstr.FromInt(s.FuncRuntimeProxyPort),
				},
			},
			Selector: map[string]string{
				"runtimeName": runtimeName,
			},
			Type: api_core_v1.ServiceTypeClusterIP,
		},
	}

	svc, err := s.K8sClient.CoreV1().Services(namespace).Create(service)

	return svc, err
}

// Update todo
func (s *FuncServiceImpl) Update(namespace, runtimeName, functionName string) (*api_core_v1.Service, error) {
	// service := &api_core_v1.Service{
	// 	ObjectMeta: api_meta_v1.ObjectMeta{
	// 		Name: functionName,
	// 		Labels: map[string]string{
	// 			"functionName": functionName,
	// 		},
	// 	},
	// 	Spec: api_core_v1.ServiceSpec{
	// 		Ports: []api_core_v1.ServicePort{
	// 			{
	// 				Name:       "runtime-env-port",
	// 				Port:       int32(80),
	// 				TargetPort: intstr.FromInt(8888),
	// 			},
	// 		},
	// 		Selector: map[string]string{
	// 			"runtimeName": runtimeName,
	// 		},
	// 		Type: api_core_v1.ServiceTypeClusterIP,
	// 	},
	// }
	patch := []byte(`{
		"metadata": {
			"name": "` + functionName + `",
			"labels": {"functionName": "` + functionName + `", "mode":"share"}
		},
		"spec": {
			"selector": {
				"runtimeName": "` + runtimeName + `"
			},
			"ports": [
				{
					"name": "runtime-proxy-port",
					"port": 80,
					"targetPort": ` + fmt.Sprintf("%v", s.FuncRuntimeProxyPort) + `
				}
			]
		}
	}`)
	svc, err := s.K8sClient.CoreV1().Services(namespace).Patch(functionName, types.MergePatchType, patch)
	// svc, err := s.K8sClient.CoreV1().RESTClient().Patch(types.ApplyPatchType).Namespace(namespace).Resource("services").Name(functionName)
	return svc, err
}

// CreateOrUpdate todo
func (s *FuncServiceImpl) CreateOrUpdate(namespace, runtimeName, functionName string) (*api_core_v1.Service, error) {
	_, err := s.K8sClient.CoreV1().Services(namespace).Get(functionName, api_meta_v1.GetOptions{})
	if errors.IsNotFound(err) {
		return s.Create(namespace, runtimeName, functionName)
	}
	return s.Update(namespace, runtimeName, functionName)
}

// Delete todo
func (s *FuncServiceImpl) Delete(namespace, functionName string) error {
	err := s.K8sClient.CoreV1().Services(namespace).Delete(functionName, &api_meta_v1.DeleteOptions{})
	return err
}
