package crossplane

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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
	d := fmt.Sprintf(detailsTemplate, o.U.GetKind(), o.GetStatus())
	cs, f, err := unstructured.NestedSlice(o.U.Object, "status", "conditions")
	if err != nil || !f {
		// failed to get conditions
		return d
	}
	for _, c := range cs {
		cMap := c.(map[string]interface{})
		if cMap == nil {
			fmt.Errorf("something wrong!!!")
			continue
		}
		getNestedString(cMap, "type")

		d = d + fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t\n",
			getNestedString(cMap, "type"),
			getNestedString(cMap, "status"),
			getNestedString(cMap, "lastTransitionTime"),
			getNestedString(cMap, "reason"),
			getNestedString(cMap, "message"))
	}
	return d
}
