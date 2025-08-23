package k8s

import (
	"fmt"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type IngressConfig struct {
	Enabled          bool                `yaml:"enabled"`
	Annotations      map[string]string   `yaml:"annotations"`
	PathPrefix       string              `yaml:"pathPrefix"`
	PathPostfix      string              `yaml:"pathPostfix"`
	PathType         string              `yaml:"pathType"`
	Hosts            []IngressHostConfig `yaml:"hosts"`
	IngressClassName string              `yaml:"ingressClassName"`
}

type IngressHostConfig struct {
	Name string           `yaml:"name"`
	TLS  IngressTLSConfig `yaml:"tls"`
}

type IngressTLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	SecretName string `yaml:"secretName"`
}

func CreateIngress(key, id string, config *IngressConfig) *networkingv1.Ingress {

	tls := make([]networkingv1.IngressTLS, 0)
	rules := make([]networkingv1.IngressRule, 0)
	for _, host := range config.Hosts {
		if host.TLS.Enabled {
			tls = append(tls, networkingv1.IngressTLS{
				SecretName: host.TLS.SecretName,
				Hosts: []string{
					host.Name,
				},
			})
		}
		rules = append(rules, networkingv1.IngressRule{
			Host: host.Name,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: Name(key, id),
									Port: networkingv1.ServiceBackendPort{
										Number: 80,
									},
								},
							},
							Path:     fmt.Sprintf("%s/%s/%s%s", config.PathPrefix, key, id, config.PathPostfix),
							PathType: (*networkingv1.PathType)(&config.PathType),
						},
					},
				},
			},
		})
	}
	return &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: Name(key, id),
			Labels: map[string]string{
				Component:  GameComponent,
				qgn.KeyTag: key,
				qgn.IDTag:  id,
			},
			Annotations: config.Annotations,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ptr.To(config.IngressClassName),
			TLS:              tls,
			Rules:            rules,
		},
	}
}
