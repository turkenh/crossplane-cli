package crossplane

import (
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/duration"
)

type K8SObject struct {
	U *unstructured.Unstructured
}

func (o *K8SObject) GetAge() string {
	ts := o.U.GetCreationTimestamp()
	if ts.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(ts.Time))
}
