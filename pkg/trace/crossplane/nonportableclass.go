package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type NonPortableClass struct {
	K8SObject
	U *unstructured.Unstructured
}

func NewNonPortableClass(u *unstructured.Unstructured) *NonPortableClass {
	return &NonPortableClass{U: u, K8SObject: K8SObject{U: u}}
}

func (o *NonPortableClass) GetStatus() string {
	return "N/A"
}
func (o *NonPortableClass) GetDetails() string {
	return ""
}
