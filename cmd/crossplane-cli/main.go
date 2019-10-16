package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/crossplaneio/crossplane/apis"
	computev1alpha1 "github.com/crossplaneio/crossplane/apis/compute/v1alpha1"
	gcpapis "github.com/crossplaneio/stack-gcp/gcp/apis"
	gcpcomputev1alpha2 "github.com/crossplaneio/stack-gcp/gcp/apis/compute/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	resourceType := "mysqlinstance"
	resourceName := "wordpress-mysql-64edc6f9-7c70-43ed-bd1d-1c26e09e0a45"
	log.Println("Tracing", resourceType, resourceName)

	// Init controller runtime client
	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		log.Fatalln("Failed to add crossplane apis to scheme:", err)
	}
	cl, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalln("Failed to create client:", err)
	}

	k8sList := &computev1alpha1.KubernetesClusterList{}
	gkeList := &gcpcomputev1alpha2.GKEClusterList{}

	err = cl.List(context.Background(), k8sList, client.InNamespace("app-project1-dev"))
	if err != nil {
		fmt.Printf("failed to list k8s clusters: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(k8sList.Items[0].Spec.ResourceReference.Kind)

	err = cl.List(context.Background(), gkeList, client.InNamespace("gcp"))
	if err != nil {
		fmt.Printf("failed to list gke clusters: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(gkeList.Items[0].ObjectMeta.Name)
}

// addToScheme adds all resources to the runtime scheme.
func addToScheme(scheme *runtime.Scheme) error {
	if err := apis.AddToScheme(scheme); err != nil {
		return err
	}
	if err := gcpapis.AddToScheme(scheme); err != nil {
		return err
	}

	return nil
}
