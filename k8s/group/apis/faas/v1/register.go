package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion todo
var SchemeGroupVersion = schema.GroupVersion{Group: "faas.k8s.io", Version: "v1"}

var (
	SB = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme todo
	AddToScheme = SB.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	// fmt.Println("初始化 scheme")
	// localSchemeBuilder.Register(addKnownTypes)
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// func Kind(kind string) schema.GroupKind {
// 	return SchemeGroupVersion.WithKind(kind).GroupKind()
// }

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&Function{},
		&FunctionList{},
	)

	// register the type in the scheme
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
