package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Claim struct {
	u *unstructured.Unstructured
}

func NewClaim(u *unstructured.Unstructured) *Claim {
	return &Claim{u: u}
}

func (o *Claim) GetStatus() string {
	return getResourceStatus(o.u)
}

func (o *Claim) GetAge() string {
	return getAge(o.u)
}

func (o *Claim) GetDetails() string {
	return getResourceDetails(o.u)
}

func (o *Claim) GetRelated() ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u

	// Get resource reference
	u, err := getObjRef(obj, resourceRefPath)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	// Get class reference
	u, err = getObjRef(obj, resourceClassRefPath)
	if err != nil {
		return related, err
	}
	// TODO(ht): Special case for claim -> portableClass, currently apiversion, kind and ns missing
	//  hence we need to manually fill them. This limitation will be removed with
	//  https://github.com/crossplaneio/crossplane/blob/master/design/one-pager-simple-class-selection.md
	if u.GetAPIVersion() == "" {
		u.SetAPIVersion(obj.GetAPIVersion())
	}
	if u.GetKind() == "" {
		u.SetKind(o.u.GetKind() + "Class")
	}
	if u.GetNamespace() == "" {
		u.SetNamespace(obj.GetNamespace())
	}

	related = append(related, u)

	return related, nil
}