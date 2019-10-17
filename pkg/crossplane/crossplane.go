package crossplane

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
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
	kindsProvider = []string{
		"Provider",
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
	} else if isProvider(objKind) {
		return NewProvider(u)
	}
	//fmt.Fprintln(os.Stderr, "!!!!!!Object is not a known crossplane object -> group: ", u.GroupVersionKind().Group, " kind: ", objKind)
	return nil
}
