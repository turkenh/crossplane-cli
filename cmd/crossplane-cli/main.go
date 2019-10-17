package main

import (
	"flag"
	"fmt"
	"log"
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
	//kind := "KubernetesCluster"
	//resourceName := "wordpress-cluster-64edc6f9-7c70-43ed-bd1d-1c26e09e0a45"
	kind := "mysqlinstance"
	resourceName := "wordpress-mysql-64edc6f9-7c70-43ed-bd1d-1c26e09e0a45"
	//kind := "GKECluster"
	//resourceName := "kubernetescluster-5c843147-069e-4a94-81d3-188c9e0fbd9c"
	namespace := "app-project1-dev"
	log.Println("Tracing", kind, resourceName)

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

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
		time.Duration(10*time.Minute))

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	my_restmapper := restmapper.NewShortcutExpander(mapper, discoveryClient)

	g := trace.NewGraph(client, my_restmapper)

	_, objs, err := g.BuildGraph(resourceName, namespace, kind)
	if err != nil {
		panic(err)
	}
	fmt.Println("-------")
	for _, o := range objs {
		fmt.Println("*", o.GetKind(), o.GetName(), o.GetNamespace())
	}
	// TODO(hasan): print

}
