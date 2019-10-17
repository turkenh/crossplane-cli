package crossplane

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/duration"
)

var (
	resourceDetailsTemplate = `%v

Status: %v
Status Conditions
TYPE	STATUS	LAST-TRANSITION-TIME	REASON	MESSAGE	
`
)

func getAge(u *unstructured.Unstructured) string {
	ts := u.GetCreationTimestamp()
	if ts.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(ts.Time))
}

func getResourceStatus(u *unstructured.Unstructured) string {
	return getNestedString(u.Object, "status", "bindingPhase")
}

func getResourceDetails(u *unstructured.Unstructured) string {
	d := fmt.Sprintf(resourceDetailsTemplate, u.GetKind(), getResourceStatus(u))
	cs, f, err := unstructured.NestedSlice(u.Object, "status", "conditions")
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

func getObjRef(obj *unstructured.Unstructured, path []string) (*unstructured.Unstructured, error) {
	a, aFound, err := unstructured.NestedString(obj.Object, append(path, "apiVersion")...)
	if err != nil {
		return nil, err
	}
	k, kFound, err := unstructured.NestedString(obj.Object, append(path, "kind")...)
	if err != nil {
		return nil, err
	}
	n, nFound, err := unstructured.NestedString(obj.Object, append(path, "name")...)
	if err != nil {
		return nil, err
	}
	ns, nsFound, err := unstructured.NestedString(obj.Object, append(path, "namespace")...)
	if err != nil {
		return nil, err
	}

	if !aFound && !kFound && !nFound && !nsFound {
		return nil, errors.New("Failed to find a reference!")
	}

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	u.SetAPIVersion(a)
	u.SetKind(k)
	u.SetName(n)
	u.SetNamespace(ns)

	return u, nil
}
