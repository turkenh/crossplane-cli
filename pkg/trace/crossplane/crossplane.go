package crossplane

import (
	"errors"
	"fmt"

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

type CrossplaneObject interface {
	GetStatus() string
	GetDetails() string
	GetAge() string
}

func ResourceFromObj(obj *unstructured.Unstructured) CrossplaneObject {
	objKind := obj.GetKind()
	if isClaim(objKind) {
		return NewClaim(obj)
	} else if isManaged(objKind) {
		return NewManaged(obj)
	} else if isPortableClass(objKind) {
		return NewPortableClass(obj)
	} else if isNonPortableClass(objKind) {
		return NewNonPortableClass(obj)
	}
	//fmt.Fprintln(os.Stderr, "!!!!!!Object is not a known crossplane object -> group: ", obj.GroupVersionKind().Group, " kind: ", objKind)
	return nil
}

func GetRelated(obj *unstructured.Unstructured) ([]*unstructured.Unstructured, error) {
	related := make([]*unstructured.Unstructured, 0, 4)

	objKind := obj.GetKind()
	if isClaim(objKind) {
		// This is a resource claim
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
			u.SetKind(objKind + "Class")
		}
		if u.GetNamespace() == "" {
			u.SetNamespace(obj.GetNamespace())
		}

		related = append(related, u)
	} else if isPortableClass(objKind) {
		// This is a resource claim
		// Get class reference
		u, err := getObjRef(obj, classRefPath)
		if err != nil {
			return related, err
		}
		related = append(related, u)
	} else if isManaged(objKind) {
		// This is a managed resource
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
	} else if isNonPortableClass(objKind) {
		// This is a non-portable class
		u, err := getObjRef(obj, providerRefPath)
		if err != nil {
			return related, err
		}

		// TODO: Could we set full resource reference for providerRef?
		if u.GetKind() == "" {
			u.SetKind("Provider")
		}
		related = append(related, u)
	} else {
		fmt.Println("!!!!!!I don't know this group: ", obj.GroupVersionKind().Group, " kind: ", objKind)
	}

	return related, nil
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
	fmt.Println("Related to resource  ---> ", a, k, n, ns)
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

func getNestedString(obj map[string]interface{}, fields ...string) string {
	val, found, err := unstructured.NestedString(obj, fields...)
	if !found || err != nil {
		return "<unknown>"
	}
	return val
}
