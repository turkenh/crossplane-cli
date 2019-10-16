package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/discovery/cached/disk"

	"k8s.io/client-go/restmapper"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	resource := "mysqlinstances"
	resourceName := "wordpress-mysql-64edc6f9-7c70-43ed-bd1d-1c26e09e0a45"
	log.Println("Tracing", resource, resourceName)

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	namespace := "app-project1-dev"

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
	res, err := my_restmapper.ResourceFor(schema.GroupVersionResource{"", "", resource})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listing %q in namespace %q:\n", resource, namespace)
	d, err := client.Resource(res).Namespace(namespace).Get(resourceName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf(" * %s \n", d.GetName())
	ownerRef, found, err := unstructured.NestedSlice(d.Object, "metadata", "ownerReferences")
	if err != nil || !found {
		fmt.Printf("ownerref not found for deployment %s: error=%s", d.GetName(), err)
	}

	owner, found, err := unstructured.NestedString(ownerRef[0].(map[string]interface{}), "name")
	if err != nil || !found {
		fmt.Printf("owner not found for deployment %s: error=%s", d.GetName(), err)
	}
	fmt.Println("owner", owner)

}
