package impl

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	api_x_clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	faasv1 "github.com/SpeedVan/function-manager/k8s/group/apis/faas/v1"
	"github.com/SpeedVan/function-manager/repository"
)

// FunctionRepositoryEtcdImpl todo
type FunctionRepositoryEtcdImpl struct {
	repository.FunctionRepository
	FuncCRD *FunctionCustomResourceDefinition
}

// NewFuncRepo todo
func NewFuncRepo(k8sRestClientFunc func(*schema.GroupVersion, bool) (*rest.RESTClient, error), k8sXClientFunc func() (*api_x_clientset.Clientset, error)) *FunctionRepositoryEtcdImpl {
	gv := &faasv1.SchemeGroupVersion
	restClient, err := k8sRestClientFunc(gv, false)
	if err != nil {
		fmt.Println(err)
		// return
	}
	k8sXClient, err := k8sXClientFunc()
	if err != nil {
		fmt.Println(err)
		// return
	}
	funcCRD, err := NewFunctionCustomResourceDefinition(gv, k8sXClient, restClient)
	err = funcCRD.Init()
	if err != nil {
		fmt.Println(err)
		// return
	}
	err = funcCRD.Regist()
	if err != nil {
		fmt.Println(err)
		// return
	}
	if err != nil {
		fmt.Println(err)
		// return
	}
	return &FunctionRepositoryEtcdImpl{
		FuncCRD: funcCRD,
	}
}

// Create todo
func (s *FunctionRepositoryEtcdImpl) Create(namespace, tenancy, functionName string) {
	funcCRD, err := s.FuncCRD.Namespace(namespace).Create(s.FuncCRD.BuildObject(
		metav1.ObjectMeta{
			Name:      sha256Str([]byte(tenancy + "+" + functionName)),
			Namespace: namespace,
			Labels: map[string]string{
				"namespace":    namespace,
				"tenancy":      tenancy,
				"functionName": functionName,
			},
		},
		faasv1.FunctionSpec{},
	))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(funcCRD)
}

// Update todo
func (s *FunctionRepositoryEtcdImpl) Update(namespace, tenancy, functionName string) {
	funcCRD, err := s.FuncCRD.Namespace(namespace).Update(s.FuncCRD.BuildObject(
		metav1.ObjectMeta{
			Name:      sha256Str([]byte(tenancy + "+" + functionName)),
			Namespace: namespace,
			Labels: map[string]string{
				"namespace":    namespace,
				"tenancy":      tenancy,
				"functionName": functionName,
			},
		},
		faasv1.FunctionSpec{},
	))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(funcCRD)
}

// Delete todo
func (s *FunctionRepositoryEtcdImpl) Delete(namespace, tenancy, functionName string) {
	funcCRD, err := s.FuncCRD.Namespace(namespace).Delete(sha256Str([]byte(tenancy+"+"+functionName)), &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(funcCRD)
}

// List todo
func (s *FunctionRepositoryEtcdImpl) List(namespace, continueBase64 string, limit int64) (*faasv1.FunctionList, error) {

	return s.FuncCRD.Namespace(namespace).List(&metav1.ListOptions{
		Continue: continueBase64,
		Limit:    limit,
	})
}

// Get todo
func (s *FunctionRepositoryEtcdImpl) Get(namespace, tenancy, functionName string) (*faasv1.Function, error) {

	return s.FuncCRD.Namespace(namespace).Get(sha256Str([]byte(tenancy+"+"+functionName)), &metav1.GetOptions{})
}

func sha256Str(x []byte) string {
	y := sha256.Sum256(x)
	return hex.EncodeToString(y[:])
}
