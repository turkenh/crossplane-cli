package trace

import (
	"fmt"
	"io"
	"os"
)

type GraphPrinter struct {
	writer io.Writer
}

func NewGraphPrinter() *GraphPrinter {
	return &GraphPrinter{writer: os.Stdout}
}

func (p *GraphPrinter) Print(nodes []*Node) error {
	fmt.Fprintln(p.writer, "graph {")

	for _, n := range nodes {

		relateds := n.Related
		for _, r := range relateds {
			fmt.Fprintln(p.writer, fmt.Sprintf("%s -- %s;", getNodeLabel(n), getNodeLabel(r)))
		}
	}
	fmt.Fprintln(p.writer, "}")
	return nil
}

func getNodeLabel(n *Node) string {
	return fmt.Sprintf("\"%s\n%s\"", n.U.GetKind(), string(n.U.GetUID())[:8])
}
