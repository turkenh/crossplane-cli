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
		relateds := n.related
		for _, r := range relateds {
			t := g.Node(r.GetId())
			t.Label(getNodeLabel(r))
			if r.state == NodeStateMissing {
				t.Attr("color", "red")
				t.Attr("style", "dashed")
			}
			f := g.Node(n.GetId())
			f.Label(getNodeLabel(n))
			g.Edge(f, t)
		}
	}
	fmt.Fprintln(p.writer, g.String())
	return nil
}

func getNodeLabel(n *Node) string {
	u := n.instance
	labelKind := u.GetKind()
	labelName := u.GetName()
	if len(labelName) > 24 {
		labelName = labelName[:12] + "..." + labelName[len(labelName)-12:]
	}
	return fmt.Sprintf("%s\n%s", labelKind, labelName)
}
