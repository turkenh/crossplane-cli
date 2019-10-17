package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Managed struct {
	u *unstructured.Unstructured
}

func NewManaged(u *unstructured.Unstructured) *Managed {
	return &Managed{u: u}
}

func (o *Managed) GetStatus() string {
	return getResourceStatus(o.u)
}

func (o *Managed) GetAge() string {
	return getAge(o.u)
}

func (o *Managed) GetDetails() string {
	return getResourceDetails(o.u)
}

func (o *Managed) GetRelated() ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u

	// Get claim reference
	u, err := getObjRef(obj, claimRefPath)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	// Get class reference
	u, err = getObjRef(obj, resourceClassRefPath)
	if err != nil {
		return related, err
	}

	related = append(related, u)
	return related, nil
}
