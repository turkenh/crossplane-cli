package trace

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/crossplaneio/crossplane-cli/pkg/trace/crossplane"

	"github.com/fatih/color"

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
	err := p.printOverview(objs)
	if err != nil {
		return err
	}
	err = p.printDetails(objs)
	if err != nil {
		return err
	}
	return nil
}
func (p *SimplePrinter) printOverview(objs []*unstructured.Unstructured) error {
	titleF := color.New(color.Bold).Add(color.Underline)
	_, err := titleF.Println("OVERVIEW")
	if err != nil {
		return err
	}
	fmt.Fprintln(p.tabWriter, "")

	_, err = fmt.Fprintln(p.tabWriter, "KIND\tNAME\tNAMESPACE\tSTATUS\tAGE\t")
	if err != nil {
		return err
	}
	for _, o := range objs {
		c := crossplane.ObjectFromUnstructured(o)
		// Skip unknown objects for now
		if c == nil {
			continue
		}
		_, err = fmt.Fprintf(p.tabWriter, "%v\t%v\t%v\t%v\t%v\t\n", o.GetKind(), o.GetName(), o.GetNamespace(), c.GetStatus(), c.GetAge())
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(p.tabWriter, "")
	err = p.tabWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}
func (p *SimplePrinter) printDetails(objs []*unstructured.Unstructured) error {
	titleF := color.New(color.Bold).Add(color.Underline)
	_, err := titleF.Println("DETAILS")
	if err != nil {
		return err
	}
	fmt.Fprintln(p.tabWriter, "")

	d := ""
	for _, o := range objs {
		c := crossplane.ObjectFromUnstructured(o)
		// Skip unknown objects for now
		if c == nil {
			continue
		}
		d += c.GetDetails()
		if d != "" {
			d += "\n\n"
		}
	}
	fmt.Fprintln(p.tabWriter, strings.TrimSpace(d))
	err = p.tabWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}
