package crossplane

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// TODO: Kind is not enough to identify GVK, need another way. Example, "Bucket", is it Claim or Managed ?
	// TODO: What about other resources related to networking and/or iam ?
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

func ObjectFromUnstructured(u *unstructured.Unstructured) (Object, error) {
	objKind := u.GetKind()
	if isClaim(objKind) {
		return NewClaim(u), nil
	} else if isManaged(objKind) {
		return NewManaged(u), nil
	} else if isPortableClass(objKind) {
		return NewPortableClass(u), nil
	} else if isNonPortableClass(objKind) {
		return NewNonPortableClass(u), nil
	} else if isProvider(objKind) {
		return NewProvider(u), nil
	} else if isApplication(objKind) {
		return NewApplication(u), nil
	} else if isApplicationResource(objKind) {
		return NewApplicationResource(u), nil
	}
	return nil, errors.New(fmt.Sprintf("%s is not a known crossplane object", objKind))
}
