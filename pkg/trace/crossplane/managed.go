package crossplane

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Managed struct {
	K8SObject
	U *unstructured.Unstructured
}

func NewManaged(u *unstructured.Unstructured) *Managed {
	return &Managed{U: u, K8SObject: K8SObject{U: u}}
}

func (o *Managed) GetStatus() string {
	return getNestedString(o.U.Object, "status", "bindingPhase")
}

func (o *Managed) GetDetails() string {
	return ""
}
