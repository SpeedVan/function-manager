package k8s

import (
	apiv1 "k8s.io/api/core/v1"
)

// FuncService todo
type FuncService interface {
	Create(string, string, string) (*apiv1.Service, error)
	Update(string, string, string) (*apiv1.Service, error)
	CreateOrUpdate(string, string, string) (*apiv1.Service, error)
	Delete(string, string) error
}
