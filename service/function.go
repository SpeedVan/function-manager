package service

import (
	"encoding/json"

	"github.com/SpeedVan/function-manager/service/k8s"
	"github.com/SpeedVan/function-manager/service/k8s/deploy"
	"github.com/SpeedVan/go-common/config"
	"k8s.io/apimachinery/pkg/api/errors"
	api_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Function todo
type Function struct {
	DeployPart k8s.FuncDeployment
	K8sClient  *kubernetes.Clientset
	Image      string
	Env        map[string]string
}

// FunctionNew new Function
func FunctionNew(cfg config.Config, k8sClient *kubernetes.Clientset) *Function {
	image := cfg.Get("IMAGE")
	env := cfg.WithPrefix("ENV_").ToGolangStringMap()
	selfURL := cfg.Get("SELFURL")
	return &Function{
		DeployPart: deploy.New(k8sClient, selfURL),
		K8sClient:  k8sClient,
		Image:      image,
		Env:        env,
	}
}

func mergeFuncConfigIntoEnvMap(envMap map[string]string, funcConfig map[string]map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range envMap {
		result[k] = v
	}
	bytes, _ := json.Marshal(funcConfig)
	result["FUNC_CONFIG"] = string(bytes)
	return result
}

// GetFuncConfig todo
func (s *Function) GetFuncConfig(namespace, runtimeName string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	yaml, err := s.K8sClient.ExtensionsV1beta1().Deployments(namespace).Get(runtimeName, api_meta_v1.GetOptions{})
	if !errors.IsNotFound(err) {
		envVars := yaml.Spec.Template.Spec.Containers[0].Env
		for _, item := range envVars {
			if item.Name == "FUNC_CONFIG" {
				json.Unmarshal([]byte(item.Value), &result)
			}
		}
	}
	return result
}

// Deploy todo
func (s *Function) Deploy(namespace, runtimeName string, funcConfig map[string]map[string]string) {
	oldFuncConfig := s.GetFuncConfig(namespace, runtimeName)
	for k, v := range funcConfig {
		oldFuncConfig[k] = v
	}

	s.DeployPart.CreateOrUpdate(namespace, runtimeName, &k8s.FuncDeploymentConfig{
		Image: s.Image,
		Env:   mergeFuncConfigIntoEnvMap(s.Env, oldFuncConfig),
	})
}

// Undeploy todo
func (s *Function) Undeploy(namespace, runtimeName string, functionName string) {
	oldFuncConfig := s.GetFuncConfig(namespace, runtimeName)
	delete(oldFuncConfig, functionName)

	s.DeployPart.Update(namespace, runtimeName, &k8s.FuncDeploymentConfig{
		Image: s.Image,
		Env:   mergeFuncConfigIntoEnvMap(s.Env, oldFuncConfig),
	})
}
