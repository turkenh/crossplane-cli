package crossplane

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/duration"
)

var (
	applicationClusterRefPath = []string{"status", "clusterRef"}
	resourceRefPath           = []string{"spec", "resourceRef"}
	claimRefPath              = []string{"spec", "claimRef"}
	classRefPath              = []string{"classRef"}
	resourceClassRefPath      = []string{"spec", "classRef"}
	providerRefPath           = []string{"specTemplate", "providerRef"}

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isClaim(kind string) bool {
	return stringInSlice(kind, kindsClaim)
}
func isManaged(kind string) bool {
	return stringInSlice(kind, kindsManaged)
}
func isNonPortableClass(kind string) bool {
	return stringInSlice(kind, kindsNonPortableClass)
}
func isPortableClass(kind string) bool {
	return stringInSlice(kind, kindsPortableClass)
}
func isProvider(kind string) bool {
	return stringInSlice(kind, kindsProvider)
}
func isApplication(kind string) bool {
	return stringInSlice(kind, kindsApplication)
}

func getNestedString(obj map[string]interface{}, fields ...string) string {
	val, found, err := unstructured.NestedString(obj, fields...)
	if err != nil {
		return "<unknown>"
	}
	if !found {
		return " "
	}
	return val
}

func getNestedInt64(obj map[string]interface{}, fields ...string) int64 {
	val, found, err := unstructured.NestedInt64(obj, fields...)
	if !found || err != nil {
		return -1
	}
	return val
}