package crossplane

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type NonPortableClass struct {
	u *unstructured.Unstructured
}

func NewNonPortableClass(u *unstructured.Unstructured) *NonPortableClass {
	return &NonPortableClass{u: u}
}

func (o *NonPortableClass) GetAge() string {
	return GetAge(o.u)
}

func (o *NonPortableClass) GetStatus() string {
	return "N/A"
}
func (o *NonPortableClass) GetDetails() string {
	return ""
}

func (o *NonPortableClass) GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u.Object
	u, err := getObjRef(obj, providerRefPath)
	if err != nil {
		return related, err
	}

	// TODO(hasan): Could we set full resource reference for providerRef?
	if u.GetAPIVersion() == "" {
		oApiVersion := o.u.GetAPIVersion()
		s := strings.Split(oApiVersion, ".")
		a := strings.Join(s[1:], ".")

		u.SetAPIVersion(a)
	}
	if u.GetKind() == "" {
		u.SetKind("Provider")
	}
	related = append(related, u)

	return related, nil
}
