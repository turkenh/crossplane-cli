package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/crossplaneio/crossplane-cli/pkg/trace"

	"k8s.io/client-go/discovery/cached/disk"

	"k8s.io/client-go/restmapper"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	//namespace := flag.String("namespace", "", "namespace")
	flag.Parse()
	kind := flag.Arg(0)
	resourceName := flag.Arg(1)
	namespace := flag.Arg(2)
	if kind == "" || resourceName == "" || namespace == "" {
		fmt.Println("Missing arguments: KIND RESOURCE_NAME NAMESPACE")
		return
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	p := trace.NewSimplePrinter()
	p.Print(objs)
	if err != nil {
		panic(err)
	}
	//fmt.Println(len(r.Related))

}
