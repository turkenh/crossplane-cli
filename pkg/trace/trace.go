package trace

import (
	computev1alpha1 "github.com/crossplaneio/crossplane/apis/compute/v1alpha1"
	gcpcomputev1alpha2 "github.com/crossplaneio/stack-gcp/gcp/apis/compute/v1alpha2"
	gcpv1alpha2 "github.com/crossplaneio/stack-gcp/gcp/apis/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TraceView string

const (
	TraceViewResource TraceView = "Resource"
	TraceViewConfig   TraceView = "Config"
)

type TraceNode struct {
	Uid     string
	obj     CrossplaneObject
	Related []*TraceNode
}

type CrossplaneObject struct {
	KubernetesCluster      *computev1alpha1.KubernetesCluster
	KubernetesClusterClass *computev1alpha1.KubernetesClusterClass
	GKECluster             *gcpcomputev1alpha2.GKECluster
	GKEClusterClass        *gcpcomputev1alpha2.GKEClusterClass
	GCPProvider            *gcpv1alpha2.Provider
}

func (c *CrossplaneObject) GetRelatedObjects() ([]CrossplaneObject, error) {
	if c.KubernetesCluster != nil {
		k := c.KubernetesCluster
		objs := make([]CrossplaneObject, 0, 2)

		r := k.GetResourceReference()
		obj := CrossplaneObject{
			GKECluster: &gcpcomputev1alpha2.GKECluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      r.Name,
					Namespace: r.Namespace,
				},
			},
		}
		objs = append(objs, obj)

		pc := k.GetPortableClassReference()
		obj = CrossplaneObject{
			KubernetesClusterClass: &computev1alpha1.KubernetesClusterClass{
				ObjectMeta: metav1.ObjectMeta{
					Name:      pc.Name,
					Namespace: k.Namespace,
				},
			},
		}
		objs = append(objs, obj)

		return objs, nil
	} else if c.KubernetesClusterClass != nil {
		k := c.KubernetesClusterClass
		objs := make([]CrossplaneObject, 0, 1)

		npc := k.GetNonPortableClassReference()
		obj := CrossplaneObject{
			GKEClusterClass: &gcpcomputev1alpha2.GKEClusterClass{
				ObjectMeta: metav1.ObjectMeta{
					Name:      npc.Name,
					Namespace: npc.Namespace,
				},
			},
		}
		objs = append(objs, obj)

		return objs, nil
	} else if c.GKECluster != nil {
		k := c.GKECluster
		objs := make([]CrossplaneObject, 0, 3)

		cr := k.GetClaimReference()
		obj := CrossplaneObject{
			KubernetesCluster: &computev1alpha1.KubernetesCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
				},
			},
		}
		objs = append(objs, obj)

		npc := k.GetNonPortableClassReference()
		obj = CrossplaneObject{
			GKEClusterClass: &gcpcomputev1alpha2.GKEClusterClass{
				ObjectMeta: metav1.ObjectMeta{
					Name:      npc.Name,
					Namespace: npc.Namespace,
				},
			},
		}
		objs = append(objs, obj)

		// TODO: check why there is no GetProviderReference
		pr := k.Spec.ProviderReference
		obj = CrossplaneObject{
			GCPProvider: &gcpv1alpha2.Provider{
				ObjectMeta: metav1.ObjectMeta{
					Name:      pr.Name,
					Namespace: pr.Namespace,
				},
			},
		}
		objs = append(objs, obj)

		return objs, nil
	}
	return nil, nil
}

type GraphBuilder interface {
	BuildGraph(v TraceView) (*TraceNode, error)
}
