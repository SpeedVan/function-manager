package service

import (
	"fmt"

	faasv1 "github.com/SpeedVan/function-manager/k8s/group/apis/faas/v1"
	"github.com/SpeedVan/function-manager/repository"
)

// Function todo
type Function struct {
	FuncRepo repository.FunctionRepository
}

// NewFunctionService todo
func NewFunctionService(funcRepo repository.FunctionRepository) *Function {
	return &Function{
		FuncRepo: funcRepo,
	}
}

// Create todo
func (s *Function) Create(namespace, functionName string) {
	s.FuncRepo.Create(namespace, "ryh", functionName)
}

// Update todo
func (s *Function) Update(namespace, functionName string) {
	s.FuncRepo.Update(namespace, "ryh", functionName)
}

// Delete todo
func (s *Function) Delete(namespace, functionName string) {
	s.FuncRepo.Delete(namespace, "ryh", functionName)
}

// List todo
func (s *Function) List(namespace, continueBase64 string, limit int64) *faasv1.FunctionList {
	result, err := s.FuncRepo.List(namespace, continueBase64, limit)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

// Get todo
func (s *Function) Get(namespace, functionName string) *faasv1.Function {
	result, err := s.FuncRepo.Get(namespace, "ryh", functionName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}
