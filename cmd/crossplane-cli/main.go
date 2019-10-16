package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/runtime/schema"

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

	res := schema.GroupVersionResource{Group: "database.crossplane.io", Version: "v1alpha1", Resource: resource}
	fmt.Printf("Listing things in namespace %q:\n", namespace)
	list, err := client.Resource(res).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		//replicas, found, err := unstructured.NestedInt64(d.Object, "spec", "replicas")
		//if err != nil || !found {
		//	fmt.Printf("Replicas not found for deployment %s: error=%s", d.GetName(), err)
		//	continue
		//}
		fmt.Printf(" * %s \n", d.GetName())
		ownerRef, found, err := unstructured.NestedSlice(d.Object, "metadata", "ownerReferences")
		if err != nil || !found {
			fmt.Printf("ownerref not found for deployment %s: error=%s", d.GetName(), err)
			continue
		}

		//v, ok := ownerRef[0].(map[string]interface)

		owner, found, err := unstructured.NestedString(ownerRef[0].(map[string]interface{}), "name")
		if err != nil || !found {
			fmt.Printf("owner not found for deployment %s: error=%s", d.GetName(), err)
			continue
		}
		fmt.Println("owner", owner)
	}
}
