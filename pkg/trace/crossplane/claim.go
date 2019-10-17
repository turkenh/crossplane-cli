package crossplane

import (
	"fmt"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	detailsTemplate = `%v:

Status: %v
Status Conditions: 
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
`
)

type Claim struct {
	K8SObject
	U *unstructured.Unstructured
}

func NewClaim(u *unstructured.Unstructured) *Claim {
	return &Claim{U: u, K8SObject: K8SObject{U: u}}
}

func (o *Claim) GetStatus() string {
	return getNestedString(o.U.Object, "status", "bindingPhase")
}

func (o *Claim) GetDetails() string {
	d := fmt.Sprintf(detailsTemplate, o.U.GetKind(), o.GetStatus())
	s, f, err := unstructured.NestedFieldNoCopy(o.U.Object, "status")
	if err != nil || !f {
		// failed to get conditions
		return d
	}
	rcs := s.(runtimev1alpha1.ResourceClaimStatus)
	for _, c := range rcs.ConditionedStatus.Conditions {
		d = d + fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t", c.Type, c.Status, c.LastTransitionTime, c.Reason, c.Message)
	}
	return d
}
