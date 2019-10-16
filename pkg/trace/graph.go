package trace

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type Graph struct {
	client     dynamic.Interface
	restMapper meta.RESTMapper
	nodes      map[string]*Node
}

func NewGraph(client dynamic.Interface, restMapper meta.RESTMapper) *Graph {
	return &Graph{
		client:     client,
		restMapper: restMapper,
		nodes:      map[string]*Node{},
	}
}

func (g *Graph) BuildGraph(name, namespace, kind string) (*Node, error) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	u.SetAPIVersion("")
	u.SetKind(kind)
	u.SetName(name)
	u.SetNamespace(namespace)

	root := &Node{
		U: u,
	}

	g.nodes[getObjId(u)] = root

	err := g.fetchObj(root)
	if err != nil {
		panic(err)
	}

	err = g.getRelated(root)
	if err != nil {
		panic(err)
	}

	for k, n := range g.nodes {
		fmt.Printf("- %s\n", k)
		if n.U.GetUID() == "" {
			g.fetchObj(n)
		}
		if n.Related == nil {
			g.getRelated(n)
		}
	}

	fmt.Println("-------")
	for k, _ := range g.nodes {
		fmt.Printf("- %s\n", k)
	}
	return root, nil
}

func (g *Graph) fetchObj(n *Node) error {
	if n.U.GetUID() != "" {
		fmt.Printf("Object %s already fetched", getObjId(n.U))
	}
	u := n.U
	fmt.Printf("Getting %q %q in namespace %q:\n", u.GetKind(), u.GetName(), u.GetNamespace())

	res, err := g.restMapper.ResourceFor(schema.GroupVersionResource{u.GroupVersionKind().Group, u.GroupVersionKind().Version, u.GetKind()})
	if err != nil {
		panic(err)
	}

	u, err = g.client.Resource(res).Namespace(u.GetNamespace()).Get(u.GetName(), metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	n.U = u
	fmt.Printf("Successfully fetched object %s \n", u.GetName())
	return nil
}

func (g *Graph) getRelated(n *Node) error {
	obj := n.U
	related := make([]*Node, 0, 4)

	objKind := obj.GetKind()
	if stringInSlice(objKind, kindsClaim) {
		// This is a resource claim
		// Get resource reference
		u, err := getObjRef(obj, resourceRefPath)
		if err != nil {
			return err
		}

		n := g.addNodeIfNotExist(u)
		related = append(related, n)

		// Get class reference
		u, err = getObjRef(obj, resourceClassRefPath)
		if err != nil {
			return err
		}
		// TODO(ht): Special case for claim -> portableClass, currently apiversion, kind and ns missing
		//  hence we need to manually fill them. This limitation will be removed with
		//  https://github.com/crossplaneio/crossplane/blob/master/design/one-pager-simple-class-selection.md
		u.SetAPIVersion(obj.GetAPIVersion())
		u.SetKind(objKind + "Class")
		u.SetNamespace(obj.GetNamespace())

		n = g.addNodeIfNotExist(u)
		related = append(related, n)
	} else if stringInSlice(objKind, kindsPortableClass) {
		// This is a resource claim
		// Get class reference
		u, err := getObjRef(obj, classRefPath)
		if err != nil {
			return err
		}
		n := g.addNodeIfNotExist(u)
		related = append(related, n)
	} else if stringInSlice(objKind, kindsManaged) {
		// This is a managed resource
		// Get claim reference
		u, err := getObjRef(obj, claimRefPath)
		if err != nil {
			return err
		}

		n := g.addNodeIfNotExist(u)
		related = append(related, n)

		// Get class reference
		u, err = getObjRef(obj, resourceClassRefPath)
		if err != nil {
			return err
		}

		n = g.addNodeIfNotExist(u)
		related = append(related, n)
	} else if stringInSlice(objKind, kindsNonPortableClass) {
		// This is a non-portable class
		u, err := getObjRef(obj, providerRefPath)
		if err != nil {
			return err
		}

		// TODO: Could we set full resource reference?
		u.SetKind("Provider")
		n := g.addNodeIfNotExist(u)
		related = append(related, n)
	} else {
		fmt.Println("!!!!!!I don't know this group: ", obj.GroupVersionKind().Group, " kind: ", objKind)
	}

	n.Related = related
	return nil
}

func (g *Graph) addNodeIfNotExist(u *unstructured.Unstructured) *Node {
	var n *Node
	if e, ok := g.nodes[getObjId(u)]; ok {
		n = e
	} else {
		n = &Node{
			U:       u,
			Related: nil,
		}
		g.nodes[getObjId(u)] = n
	}
	return n
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

func getObjId(u *unstructured.Unstructured) string {
	return strings.ToLower(fmt.Sprintf("%s-%s-%s", u.GetKind(), u.GetNamespace(), u.GetName()))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
