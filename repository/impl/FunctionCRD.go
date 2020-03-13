package impl

import (
	api_x_v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	api_x_clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	faasv1 "github.com/SpeedVan/function-manager/k8s/group/apis/faas/v1"
	"github.com/SpeedVan/go-common-kubernetes/client"
)

// FunctionCustomResourceDefinition todo
type FunctionCustomResourceDefinition struct {
	FunctionClientInterface
	GroupVersion *schema.GroupVersion
	K8sXClient   *api_x_clientset.Clientset
	RESTClient   rest.Interface
	CRD          *api_x_v1beta1.CustomResourceDefinition
}

// FunctionClientInterface todo
type FunctionClientInterface interface {
	Create(function *faasv1.Function) (*faasv1.Function, error)
	Update(function *faasv1.Function) (*faasv1.Function, error)
	Delete(name string, opts *metav1.DeleteOptions) (*faasv1.Function, error)
	List(opts *metav1.ListOptions) (*faasv1.FunctionList, error)
	Get(name string, opts *metav1.GetOptions) (*faasv1.Function, error)
}

// NewFunctionCustomResourceDefinition todo
func NewFunctionCustomResourceDefinition(groupVersion *schema.GroupVersion, k8sXClient *api_x_clientset.Clientset, k8sRestClient rest.Interface) (*FunctionCustomResourceDefinition, error) {
	return &FunctionCustomResourceDefinition{
		GroupVersion: groupVersion,
		K8sXClient:   k8sXClient,
		RESTClient:   k8sRestClient,
		CRD: &api_x_v1beta1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "functions." + groupVersion.Group,
			},
			Spec: api_x_v1beta1.CustomResourceDefinitionSpec{
				Group:   groupVersion.Group,
				Version: groupVersion.Version,
				Versions: []api_x_v1beta1.CustomResourceDefinitionVersion{{
					Name:    groupVersion.Version,
					Served:  true,
					Storage: true,
				}},
				Names: api_x_v1beta1.CustomResourceDefinitionNames{
					Plural:     "functions",
					Singular:   "function",
					ShortNames: []string{"func"},
					Kind:       "Function",
				},
				Scope: api_x_v1beta1.NamespaceScoped,
				Validation: &api_x_v1beta1.CustomResourceValidation{
					OpenAPIV3Schema: &api_x_v1beta1.JSONSchemaProps{
						Properties: map[string]api_x_v1beta1.JSONSchemaProps{
							"spec": {
								Properties: map[string]api_x_v1beta1.JSONSchemaProps{
									"image": {
										Type: "string",
										// Nullable: true,
									},
									"replicas": {
										Type: "integer",
										// Nullable: true,
									},
									"resources": {
										Type: "object",
										// Nullable: true,
										Properties: map[string]api_x_v1beta1.JSONSchemaProps{
											"limits": {
												Type: "object",
												Properties: map[string]api_x_v1beta1.JSONSchemaProps{
													"cpu": {
														Type: "string",
													},
													"memory": {
														Type: "string",
													},
												},
											},
											"requests": {
												Type: "object",
												Properties: map[string]api_x_v1beta1.JSONSchemaProps{
													"cpu": {
														Type: "string",
													},
													"memory": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// Init todo
func (s *FunctionCustomResourceDefinition) Init() error {
	_, err := s.K8sXClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(s.CRD)
	if err != nil {
		return err
	}
	return nil
}

// Regist todo
func (s *FunctionCustomResourceDefinition) Regist() error {
	SchemeBuilder := runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(*s.GroupVersion,
			&faasv1.Function{},
			&faasv1.FunctionList{},
		)
		metav1.AddToGroupVersion(scheme, *s.GroupVersion)
		return nil
	})
	return SchemeBuilder.AddToScheme(scheme.Scheme)
}

// BuildObject todo
func (s *FunctionCustomResourceDefinition) BuildObject(meta metav1.ObjectMeta, spec faasv1.FunctionSpec) *faasv1.Function {
	return &faasv1.Function{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Function",
			APIVersion: s.GroupVersion.Group + "/" + s.GroupVersion.Version,
		},
		ObjectMeta: meta,
		Spec:       spec,
	}
}

// Create todo
func (s *FunctionCustomResourceDefinition) Create(function *faasv1.Function) (*faasv1.Function, error) {
	result := &faasv1.Function{}
	err := s.RESTClient.Post().Resource("functions").Body(function).Do().Into(result)
	return result, err
}

// Update todo
func (s *FunctionCustomResourceDefinition) Update(function *faasv1.Function) (*faasv1.Function, error) {
	result := &faasv1.Function{}
	err := s.RESTClient.Put().Resource("functions").Body(function).Do().Into(result)
	return result, err
}

// Delete todo
func (s *FunctionCustomResourceDefinition) Delete(name string, opts *metav1.DeleteOptions) (*faasv1.Function, error) {
	result := &faasv1.Function{}
	err := s.RESTClient.Delete().Resource("functions").Name(name).VersionedParams(opts, scheme.ParameterCodec).Do().Into(result)
	return result, err
}

// List todo
func (s *FunctionCustomResourceDefinition) List(opts *metav1.ListOptions) (*faasv1.FunctionList, error) {
	result := &faasv1.FunctionList{}
	err := s.RESTClient.Get().Resource("functions").VersionedParams(opts, scheme.ParameterCodec).Do().Into(result)
	return result, err
}

// Get todo
func (s *FunctionCustomResourceDefinition) Get(name string, opts *metav1.GetOptions) (*faasv1.Function, error) {
	result := &faasv1.Function{}
	err := s.RESTClient.Get().Resource("functions").Name(name).VersionedParams(opts, scheme.ParameterCodec).Do().Into(result)
	return result, err
}

// Namespace todo
func (s *FunctionCustomResourceDefinition) Namespace(namespace string) FunctionClientInterface {
	return &FunctionCustomResourceDefinition{
		GroupVersion: s.GroupVersion,
		K8sXClient:   s.K8sXClient,
		RESTClient: &client.NamespaceRESTClient{
			RESTClient: s.RESTClient,
			Namespace:  namespace,
		},
		CRD: s.CRD,
	}
}
