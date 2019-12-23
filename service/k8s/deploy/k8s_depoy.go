package deploy

import (
	"fmt"

	"github.com/SpeedVan/function-manager/service/k8s"
	cvt "github.com/SpeedVan/go-common/convert"
	api_core_v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	api_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// FuncDeploymentImpl todo
type FuncDeploymentImpl struct {
	k8s.FuncDeployment
	K8sClient *kubernetes.Clientset
	SelfURL   string
}

func envMap2EnvVar(m map[string]string) []api_core_v1.EnvVar {
	result := []api_core_v1.EnvVar{}
	for k, v := range m {
		result = append(result, api_core_v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return result
}

func New(k8sClient *kubernetes.Clientset, selfURL string) *FuncDeploymentImpl {
	return &FuncDeploymentImpl{
		K8sClient: k8sClient,
		SelfURL:   selfURL,
	}
}

// Create todo
func (s *FuncDeploymentImpl) Create(namespace, runtimeName string, config *k8s.FuncDeploymentConfig) {
	deployment := &v1beta1.Deployment{
		ObjectMeta: api_meta_v1.ObjectMeta{
			Labels: map[string]string{
				"runtimeName": runtimeName,
			},
			Name: runtimeName,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: cvt.PtrOfInt32(2),
			Selector: &api_meta_v1.LabelSelector{
				MatchLabels: map[string]string{
					"runtimeName": runtimeName,
				},
			},
			Template: api_core_v1.PodTemplateSpec{
				ObjectMeta: api_meta_v1.ObjectMeta{
					Labels: map[string]string{
						"runtimeName": runtimeName,
					},
				},
				Spec: api_core_v1.PodSpec{
					Containers: []api_core_v1.Container{
						api_core_v1.Container{
							Name:                   runtimeName,
							Image:                  config.Image,
							ImagePullPolicy:        "Always",
							TerminationMessagePath: "/dev/termination-log",
							Env:                    envMap2EnvVar(config.Env),
							Lifecycle: &api_core_v1.Lifecycle{
								PostStart: &api_core_v1.Handler{
									Exec: &api_core_v1.ExecAction{
										Command: []string{
											"/bin/sh",
											"/app/pre_start.sh",
											namespace,
											runtimeName,
										},
									},
								},
								PreStop: &api_core_v1.Handler{
									Exec: &api_core_v1.ExecAction{
										Command: []string{
											"/bin/sh",
											"/app/pre_stop.sh",
										},
									},
								},
							},
							Resources: api_core_v1.ResourceRequirements{
								Limits: api_core_v1.ResourceList{
									api_core_v1.ResourceCPU:    *resource.NewMilliQuantity(1000, resource.DecimalSI),
									api_core_v1.ResourceMemory: resource.MustParse("1000Mi"),
								},
							},
							SecurityContext: &api_core_v1.SecurityContext{
								Privileged: cvt.PtrOfBool(true),
								Capabilities: &api_core_v1.Capabilities{
									Add: []api_core_v1.Capability{"SYS_ADMIN"},
								},
							},
						},
					},
					TerminationGracePeriodSeconds: cvt.PtrOfInt64(60),
				},
			},
		},
	}
	s.K8sClient.ExtensionsV1beta1().Deployments(namespace).Create(deployment)
}

// Update todo
func (s *FuncDeploymentImpl) Update(namespace, runtimeName string, config *k8s.FuncDeploymentConfig) {
	deployment := &v1beta1.Deployment{
		ObjectMeta: api_meta_v1.ObjectMeta{
			Labels: map[string]string{
				"runtimeName": runtimeName,
			},
			Name: runtimeName,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: cvt.PtrOfInt32(2),
			Selector: &api_meta_v1.LabelSelector{
				MatchLabels: map[string]string{
					"runtimeName": runtimeName,
				},
			},
			Template: api_core_v1.PodTemplateSpec{
				ObjectMeta: api_meta_v1.ObjectMeta{
					Labels: map[string]string{
						"runtimeName": runtimeName,
					},
				},
				Spec: api_core_v1.PodSpec{
					Containers: []api_core_v1.Container{
						api_core_v1.Container{
							Name:                   runtimeName,
							Image:                  config.Image,
							ImagePullPolicy:        "Always",
							TerminationMessagePath: "/dev/termination-log",
							Env:                    envMap2EnvVar(config.Env),
							Lifecycle: &api_core_v1.Lifecycle{
								PostStart: &api_core_v1.Handler{
									Exec: &api_core_v1.ExecAction{
										Command: []string{
											"/bin/sh",
											"/app/pre_start.sh",
											namespace,
											runtimeName,
										},
									},
								},
								PreStop: &api_core_v1.Handler{
									Exec: &api_core_v1.ExecAction{
										Command: []string{
											"/bin/sh",
											"/app/pre_stop.sh",
										},
									},
								},
							},
							Resources: api_core_v1.ResourceRequirements{
								Limits: api_core_v1.ResourceList{
									api_core_v1.ResourceCPU:    *resource.NewMilliQuantity(1000, resource.DecimalSI),
									api_core_v1.ResourceMemory: resource.MustParse("1000Mi"),
								},
							},
							SecurityContext: &api_core_v1.SecurityContext{
								Privileged: cvt.PtrOfBool(true),
								Capabilities: &api_core_v1.Capabilities{
									Add: []api_core_v1.Capability{"SYS_ADMIN"},
								},
							},
						},
					},
					TerminationGracePeriodSeconds: cvt.PtrOfInt64(60),
				},
			},
		},
	}
	_, err := s.K8sClient.ExtensionsV1beta1().Deployments(namespace).Update(deployment)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// CreateOrUpdate todo
func (s *FuncDeploymentImpl) CreateOrUpdate(namespace, runtimeName string, config *k8s.FuncDeploymentConfig) {
	_, err := s.K8sClient.ExtensionsV1beta1().Deployments(namespace).Get(runtimeName, api_meta_v1.GetOptions{})
	if errors.IsNotFound(err) {
		s.Create(namespace, runtimeName, config)
	}
	s.Update(namespace, runtimeName, config)
}

// Delete todo
func (s *FuncDeploymentImpl) Delete(namespace, runtimeName string) {
	s.K8sClient.ExtensionsV1beta1().Deployments(namespace).Delete(runtimeName, nil)
}
