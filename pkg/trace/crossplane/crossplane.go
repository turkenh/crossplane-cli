package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	resourceRefPath      = []string{"spec", "resourceRef"}
	claimRefPath         = []string{"spec", "claimRef"}
	classRefPath         = []string{"classRef"}
	resourceClassRefPath = []string{"spec", "classRef"}
	providerRefPath      = []string{"specTemplate", "providerRef"}

	kindsClaim = []string{
		"MySQLInstance",
		"KubernetesCluster",
	}
	kindsManaged = []string{
		"CloudsqlInstance",
		"GKECluster",
	}
	kindsPortableClass = []string{
		"MySQLInstanceClass",
		"KubernetesClusterClass",
	}
	kindsNonPortableClass = []string{
		"CloudsqlInstanceClass",
		"GKEClusterClass",
	}
)

type Object interface {
	GetStatus() string
	GetDetails() string
	GetAge() string
	GetRelated() ([]*unstructured.Unstructured, error)
}

func ObjectFromUnstructured(u *unstructured.Unstructured) Object {
	objKind := u.GetKind()
	if isClaim(objKind) {
		return NewClaim(u)
	} else if isManaged(objKind) {
		return NewManaged(u)
	} else if isPortableClass(objKind) {
		return NewPortableClass(u)
	} else if isNonPortableClass(objKind) {
		return NewNonPortableClass(u)
	}
	//fmt.Fprintln(os.Stderr, "!!!!!!Object is not a known crossplane object -> group: ", u.GroupVersionKind().Group, " kind: ", objKind)
	return nil
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

func getNestedString(obj map[string]interface{}, fields ...string) string {
	val, found, err := unstructured.NestedString(obj, fields...)
	if !found || err != nil {
		return "<unknown>"
	}
	return val
}
