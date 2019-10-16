package trace

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type TraceView string

const (
	TraceViewResource TraceView = "Resource"
	TraceViewConfig   TraceView = "Config"
)

type Node struct {
	Uid     string
	obj     *unstructured.Unstructured
	Related []*Node
}

type GraphBuilder interface {
	BuildGraph(v TraceView) (*Node, error)
}
