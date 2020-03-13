package repository

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	faasv1 "github.com/SpeedVan/function-manager/k8s/group/apis/faas/v1"
)

// FunctionRepository todo
type FunctionRepository interface {
	CreateFunctionCRD(*v1beta1.CustomResourceDefinition)
	Create(string, string, string)
	Update(string, string, string)
	Delete(string, string, string)
	List(string, string, int64) (*faasv1.FunctionList, error)
	Get(string, string, string) (*faasv1.Function, error)
}
