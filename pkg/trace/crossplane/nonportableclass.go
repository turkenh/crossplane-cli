package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type NonPortableClass struct {
	u *unstructured.Unstructured
}

func NewNonPortableClass(u *unstructured.Unstructured) *NonPortableClass {
	return &NonPortableClass{u: u}
}

func (o *NonPortableClass) GetAge() string {
	return getAge(o.u)
}

func (o *NonPortableClass) GetStatus() string {
	return "N/A"
}
func (o *NonPortableClass) GetDetails() string {
	return ""
}

func (o *NonPortableClass) GetRelated() ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u
	u, err := getObjRef(obj, providerRefPath)
	if err != nil {
		return related, err
	}

	// TODO: Could we set full resource reference for providerRef?
	if u.GetKind() == "" {
		u.SetKind("Provider")
	}
	related = append(related, u)

	return related, nil
}
