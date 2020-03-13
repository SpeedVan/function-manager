package service

import (
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// Scale todo
type Scale struct {
	RestClient rest.Interface
}

// NewScaleService todo
func NewScaleService(restClientBuildFunc func(*schema.GroupVersion, bool) (*rest.RESTClient, error)) (*Scale, error) {
	restClient, err := restClientBuildFunc(&schema.GroupVersion{Group: "extensions", Version: "v1beta1"}, false)
	if err != nil {
		return nil, err
	}
	return &Scale{
		RestClient: restClient,
	}, nil
}

// Get todo
func (s *Scale) Get(namespace, functionName string) (*v1beta1.Scale, error) {
	result := &v1beta1.Scale{}
	opts := &v1.GetOptions{}

	if err := s.RestClient.Get().
		Namespace(namespace).
		Resource("deployments").
		Name(functionName).
		SubResource("scale").
		VersionedParams(opts, scheme.ParameterCodec).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Inc todo
func (s *Scale) Inc(namespace, functionName string, scaleInc int) (*v1beta1.Scale, error) {
	scale, err := s.Get(namespace, functionName)
	if err != nil {
		return nil, err
	}
	scale.Spec.Replicas += int32(scaleInc)
	if err := s.RestClient.Put().
		Namespace(namespace).
		Resource("deployments").
		Name(functionName).
		SubResource("scale").
		Body(scale).
		Do().
		Into(scale); err != nil {
		return nil, err
	}
	return scale, nil
}
