package crossplane

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Provider struct {
	u *unstructured.Unstructured
}

func NewProvider(u *unstructured.Unstructured) *Provider {
	return &Provider{u: u}
}

func (o *Provider) GetStatus() string {
	return "N/A"
}

func (o *Provider) GetAge() string {
	return getAge(o.u)
}

func (o *Provider) GetDetails() string {
	return ""
}

func (o *Provider) GetRelated() ([]*unstructured.Unstructured, error) {
	// TODO(ht): credentialsSecretRef?
	return nil, nil
}