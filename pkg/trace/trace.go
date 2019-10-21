package trace

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type NodeState string

const (
	NodeStateMissing NodeState = "Missing"
	NodeStatePending NodeState = "Pending"
	NodeStateReady   NodeState = "Ready"
)

type Node struct {
	U       *unstructured.Unstructured
	Id      string
	Related []*Node
	State   NodeState
}

type GraphBuilder interface {
	BuildGraph(string, string, string) (*Node, []*unstructured.Unstructured, error)
}

type Printer interface {
	Print([]*Node) error
}
