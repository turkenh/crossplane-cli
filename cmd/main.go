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
)

func main() {
	var kubeconfig string
	var namespace string
	if home := homedir.HomeDir(); home != "" {
		pflag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		pflag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}

	pflag.StringVarP(&namespace, "namespace", "n", "default", "namespace")

	pflag.Parse()
	kind := pflag.Arg(0)
	resourceName := pflag.Arg(1)
	if kind == "" || resourceName == "" {
		failWithErr(fmt.Errorf("Missing arguments: KIND RESOURCE_NAME [-n| --namespace NAMESPACE]"))
	}
	fmt.Fprintf(os.Stderr, "Collecting information for %s %s in namespace %s...\n\n", kind, resourceName, namespace)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		failWithErr(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		failWithErr(err)
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

	g := trace.NewGraph(client, rMapper)
	_, objs, err := g.BuildGraph(resourceName, namespace, kind)
	if err != nil {
		failWithErr(err)
	}
	p := trace.NewSimplePrinter()
	p.Print(objs)
	if err != nil {
		failWithErr(err)
	}
	//fmt.Println(len(r.Related))

}

func failWithErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
