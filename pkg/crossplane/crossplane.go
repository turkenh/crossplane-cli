package crossplane

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// TODO: Kind is not enough to identify GVK, need another way. Example, "Bucket", is it Claim or Managed?
	kindsClaim = []string{
		"MySQLInstance",
		"KubernetesCluster",
		"RedisCluster",
		"PostgreSQLInstance",
		"Bucket",
	}
	kindsManaged = []string{
		"Redis",
		"MysqlServer",
		"PostgresqlServer",
		"AKSCluster",
		"Container",
		"Account",

		"CloudsqlInstance",
		"GKECluster",
		"CloudMemorystoreInstance",
		"Bucket",

		"ReplicationGroup",
		"EKSCluster",
		"RDSInstance",
		"S3Bucket",
	}
	kindsPortableClass = []string{
		"MySQLInstanceClass",
		"KubernetesClusterClass",
		"RedisClusterClass",
		"PostgreSQLInstanceClass",
		"BucketClass",
	}
	kindsNonPortableClass = []string{
		"RedisClass",
		"AKSClusterClass",
		"SQLServerClass",

		"CloudsqlInstanceClass",
		"GKEClusterClass",
		"CloudMemorystoreInstanceClass",
		"BucketClass",
		"AccountClass",
		"ContainerClass",

		"ReplicationGroupClass",
		"EKSClusterClass",
		"RDSInstanceClass",
		"S3BucketClass",
	}
	kindsProvider = []string{
		"Provider",
	}
	kindsApplication = []string{
		"KubernetesApplication",
	}
	kindsApplicationResource = []string{
		"KubernetesApplicationResource",
	}
)

type Object interface {
	GetStatus() string
	GetDetails() string
	GetAge() string
	GetRelated(filterByLabel func(metav1.GroupVersionKind, string, string) ([]unstructured.Unstructured, error)) ([]*unstructured.Unstructured, error)
}

func ObjectFromUnstructured(u *unstructured.Unstructured) Object {
	objKind := u.GetKind()
	if isClaim(objKind) {
		return NewClaim(u)
	} else if isManaged(objKind) {
		return NewManaged(u)
	} else if isPortableClass(objKind) {
		return NewPortableClass(u)
	} else if isNonPortableClass(objKind) {
		return NewNonPortableClass(u)
	} else if isProvider(objKind) {
		return NewProvider(u)
	} else if isApplication(objKind) {
		return NewApplication(u)
	} else if isApplicationResource(objKind) {
		return NewApplicationResource(u)
	}
	fmt.Fprintln(os.Stderr, "!!!!!!Object is not a known crossplane object -> group: ", u.GroupVersionKind().Group, " kind: ", objKind)
	return nil
}
