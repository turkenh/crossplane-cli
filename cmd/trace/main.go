package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/crossplaneio/crossplane-cli/pkg/trace"

	"k8s.io/client-go/discovery/cached/disk"

	"k8s.io/client-go/restmapper"

	"github.com/spf13/pflag"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	envKubeconfig = "KUBECONFIG"
)

func main() {
	var help bool
	var kubeconfig string
	var namespace string
	var outputFormat string

	pflag.StringVar(&kubeconfig, "kubeconfig", "", "Absolute path to the kubeconfig file")
	pflag.BoolVarP(&help, "help", "h", false, "Shows this help message")
	pflag.StringVarP(&namespace, "namespace", "n", "default", "Namespace")
	pflag.StringVarP(&outputFormat, "outputFormat", "o", "", "Output format. One of: dot")
	pflag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	resource := pflag.Arg(0)
	resourceName := pflag.Arg(1)
	if resource == "" || resourceName == "" {
		failWithMessage("missing arguments: SUBCOMMAND TYPE[.GROUP] NAME [-n|--namespace NAMESPACE]")
	}
	fmt.Fprintf(os.Stderr, "Collecting information for %s %s in namespace %s...\n\n", resource, resourceName, namespace)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		failWithMessage(err.Error())
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		failWithMessage(err.Error())
	}

	discoveryCacheDir := filepath.Join("./.kube", "cache", "discovery")
	httpCacheDir := filepath.Join("./.kube", "http-cache")
	discoveryClient, err := disk.NewCachedDiscoveryClientForConfig(
		config,
		discoveryCacheDir,
		httpCacheDir,
		10*time.Minute)

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	rMapper := restmapper.NewShortcutExpander(mapper, discoveryClient)

	g := trace.NewKubeGraphBuilder(client, rMapper)
	_, traversed, err := g.BuildGraph(resourceName, namespace, resource)
	if err != nil {
		failWithMessage(err.Error())
	}
	if outputFormat == "" {
		p := trace.NewSimplePrinter()
		p.Print(traversed)
		if err != nil {
			failWithMessage(err.Error())
		}
	} else if outputFormat == "dot" {
		gp := trace.NewGraphPrinter()
		gp.Print(traversed)
		if err != nil {
			failWithMessage(err.Error())
		}
	} else if outputFormat == "yaml" || outputFormat == "json" {
		failWithMessage(fmt.Sprintf("%s outputFormat format is not supported yet", outputFormat))
	} else {
		failWithMessage("unknown outputFormat format, should be one of: dot|yaml|json")
	}
}

func failWithMessage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(-1)
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `Usage: kubectl crossplane trace TYPE[.GROUP] NAME [-n| --namespace NAMESPACE] [-h|--help]

`)
	pflag.PrintDefaults()
}
