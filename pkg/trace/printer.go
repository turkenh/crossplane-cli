package trace

import (
	"fmt"
	"os"
	"text/tabwriter"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	tabwriterMinWidth = 6
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
)

type SimplePrinter struct {
	tabWriter *tabwriter.Writer
}

func NewSimplePrinter() *SimplePrinter {
	t := tabwriter.NewWriter(os.Stdout, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, 0)
	return &SimplePrinter{tabWriter: t}
}

func (p *SimplePrinter) Print(objs []*unstructured.Unstructured) error {
	_, err := fmt.Fprintln(p.tabWriter, "KIND\tNAME\tNAMESPACE\t")
	if err != nil {
		return err
	}
	for _, o := range objs {
		_, err = fmt.Fprintf(p.tabWriter, "%v\t%v\t%v\t\n", o.GetKind(), o.GetName(), o.GetNamespace())
		if err != nil {
			return err
		}
	}
	err = p.tabWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}
