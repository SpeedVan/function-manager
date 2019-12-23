package k8s

// FuncDeployment todo
type FuncDeployment interface {
	Create(namespace, runtimeName string, config *FuncDeploymentConfig)
	Update(namespace, runtimeName string, config *FuncDeploymentConfig)
	CreateOrUpdate(namespace, runtimeName string, config *FuncDeploymentConfig)
	Delete(namespace, runtimeName string)
}

// FuncDeploymentConfig todo
type FuncDeploymentConfig struct {
	Image string
	Env   map[string]string
}
