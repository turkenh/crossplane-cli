package crossplane

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	applicationDetailsTemplate = `%v

NAME	CLUSTER	STATUS	DESIRED	SUBMITTED
%v	%v	%v	%v	%v	

Status Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
`
)

type Application struct {
	u *unstructured.Unstructured
}

func NewApplication(u *unstructured.Unstructured) *Application {
	return &Application{u: u}
}

func (o *Application) GetStatus() string {
	return getNestedString(o.u.Object, "status", "state")
}

func (o *Application) GetAge() string {
	return getAge(o.u)
}

func (o *Application) GetDetails() string {
	d := fmt.Sprintf(applicationDetailsTemplate, o.u.GetKind(),
		o.u.GetName(), getNestedString(o.u.Object, append(applicationClusterRefPath, "name")...),
		o.GetStatus(), getNestedInt64(o.u.Object, "status", "desiredResources"),
		getNestedInt64(o.u.Object, "status", "submittedResources"))

	cs, f, err := unstructured.NestedSlice(o.u.Object, "status", "conditionedStatus", "conditions")
	if err != nil || !f {
		// failed to get conditions
		return d
	}
	for _, c := range cs {
		cMap := c.(map[string]interface{})
		if cMap == nil {
			d = d + "<error>"
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

func (o *Application) GetRelated(f func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0)
	obj := o.u

	// Get resource reference
	u, err := getObjRef(obj, applicationClusterRefPath)
	if err != nil {
		return related, err
	}

	related = append(related, u)

	// Get related resources with resourceSelector
	uArr, err := f(metav1.GroupVersionKind{
		Kind: "MySQLInstance",
	}, obj.GetNamespace(), getNestedLabelSelector(obj.Object, "spec", "resourceSelector", "matchLabels"))
	if err != nil {
		return related, err
	}

	for _, u := range uArr {
		related = append(related, &u)
	}

	return related, nil
}
