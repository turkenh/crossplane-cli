package crossplane

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PortableClass struct {
	u *unstructured.Unstructured
}

func NewPortableClass(u *unstructured.Unstructured) *PortableClass {
	return &PortableClass{u: u}
}
func (o *PortableClass) GetAge() string {
	return getAge(o.u)
}

func (o *PortableClass) GetStatus() string {
	return "N/A"
}

func (o *PortableClass) GetDetails() string {
	return ""
}

func (o *PortableClass) GetRelated(f func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u

	// Get class reference
	u, err := getObjRef(obj, classRefPath)
	if err != nil {
		return related, err
	}

	related = append(related, u)
	return related, nil
}
