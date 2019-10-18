package trace

import (
	"container/list"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"

	"github.com/crossplaneio/crossplane-cli/pkg/crossplane"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (g *Graph) BuildGraph(name, namespace, kind string) (root *Node, traversed []*unstructured.Unstructured, err error) {
	g.filterByLabel(metav1.GroupVersionKind{}, "", "")
	queue := list.New()

	traversed = make([]*unstructured.Unstructured, 0)

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	u.SetAPIVersion("")
	u.SetKind(kind)
	u.SetName(name)
	u.SetNamespace(namespace)

	root = g.addNodeIfNotExist(u)

	err = g.fetchObj(root)
	if err != nil {
		return nil, nil, err
	}

	visited := map[types.UID]bool{}
	traversed = append(traversed, root.U)
	visited[root.U.GetUID()] = true
	queue.PushBack(root)

	for queue.Len() > 0 {
		qnode := queue.Front()
		node := qnode.Value.(*Node)
		err = g.findRelated(node)
		if err != nil {
			return nil, nil, err
		}

		for _, n := range node.Related {
			if n.U.GetUID() == "" {
				err := g.fetchObj(n)
				if err != nil {
					return nil, nil, err
				}
			}
			u := n.U
			uid := u.GetUID()
			if !visited[uid] {
				traversed = append(traversed, u)
				visited[uid] = true
				queue.PushBack(n)
			}
		}
		queue.Remove(qnode)
	}

	return
}

func (g *Graph) fetchObj(n *Node) error {
	if n.U.GetUID() != "" {
		return nil
	}
	u := n.U
	res, err := g.restMapper.ResourceFor(schema.GroupVersionResource{Group: u.GroupVersionKind().Group, Version: u.GroupVersionKind().Version, Resource: u.GetKind()})
	if err != nil {
		return err
	}

	u, err = g.client.Resource(res).Namespace(u.GetNamespace()).Get(u.GetName(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	n.U = u
	return nil
}

func (g *Graph) filterByLabel(gvk metav1.GroupVersionKind, namespace, selector string) ([]unstructured.Unstructured, error) {
	res, err := g.restMapper.ResourceFor(schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: gvk.Kind})
	if err != nil {
		return nil, err
	}

	list, err := g.client.Resource(res).Namespace(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (g *Graph) findRelated(n *Node) error {
	n.Related = make([]*Node, 0)

	c := crossplane.ObjectFromUnstructured(n.U)
	// Skip unknown objects for now
	if c == nil {
		//return errors.New(fmt.Sprintf("%v not a known crossplane object", n.u.GroupVersionKind().String()))
		fmt.Printf("%v not a known crossplane object\n", n.U.GroupVersionKind().String())
		return nil
	}
	objs, err := c.GetRelated(g.filterByLabel)
	if err != nil {
		return err
	}
	for _, o := range objs {
		r := g.addNodeIfNotExist(o)
		n.Related = append(n.Related, r)
	}
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

func getObjId(u *unstructured.Unstructured) string {
	return strings.ToLower(fmt.Sprintf("%s-%s-%s", u.GetKind(), u.GetNamespace(), u.GetName()))
}
