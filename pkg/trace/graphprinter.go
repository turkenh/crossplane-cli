package trace

import (
	"fmt"
	"io"
	"os"

	"github.com/emicklei/dot"
)

type GraphPrinter struct {
	writer io.Writer
}

func NewGraphPrinter() *GraphPrinter {
	return &GraphPrinter{writer: os.Stdout}
}

func (p *GraphPrinter) Print(nodes []*Node) error {
	g := dot.NewGraph(dot.Undirected)
	for _, n := range nodes {
		relateds := n.Related
		for _, r := range relateds {
			t := g.Node(getNodeLabel(r))
			f := g.Node(getNodeLabel(n))
			g.Edge(f, t)
		}
	}
	fmt.Fprintln(p.writer, g.String())
	return nil
}

func getNodeLabel(n *Node) string {
	u := n.U
	labelKind := u.GetKind()
	labelName := string(u.GetUID())
	if n.State == NodeStateMissing {
		labelName = "<missing>"
	} else if len(n.U.GetUID()) > 6 {
		labelName = string(u.GetUID())[:6]
	}
	return fmt.Sprintf("\"%s\n%s\"", labelKind, labelName)
}
