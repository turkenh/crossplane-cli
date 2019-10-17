package trace

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Node struct {
	U       *unstructured.Unstructured
	Related []*Node
}

type GraphBuilder interface {
	BuildGraph(string, string, string) (*Node, []*unstructured.Unstructured, error)
}

type Printer interface {
	Print([]*unstructured.Unstructured) error
}
