package trace

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Node struct {
	U       *unstructured.Unstructured
	Related []*Node
}

type GraphBuilder interface {
	BuildGraph(name, namespace, kind string) (*Node, error)
}

type Printer interface {
	Print([]*unstructured.Unstructured) error
}
