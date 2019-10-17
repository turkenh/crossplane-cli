package crossplane

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type PortableClass struct {
	K8SObject
	U *unstructured.Unstructured
}

func NewPortableClass(u *unstructured.Unstructured) *PortableClass {
	return &PortableClass{U: u, K8SObject: K8SObject{U: u}}
}

func (o *PortableClass) GetStatus() string {
	return "N/A"
}

func (o *PortableClass) GetDetails() string {
	return ""
}
