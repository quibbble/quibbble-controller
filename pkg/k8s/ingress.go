package k8s

import (
	"fmt"
	"strings"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateIngress(host, key, id string, allowOrigins []string) *networkingv1.Ingress {
	pathType := networkingv1.PathTypeImplementationSpecific
	return &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name(key, id),
			Namespace: Namespace,
			Labels: map[string]string{
				Component:  GameComponent,
				qgn.KeyTag: key,
				qgn.IDTag:  id,
			},
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target":         "/$2",
				"nginx.ingress.kubernetes.io/enable-cors":            "true",
				"nginx.ingress.kubernetes.io/cors-allow-methods":     "GET, HEAD, OPTIONS",
				"nginx.ingress.kubernetes.io/cors-allow-credentials": "true",
				// default is 60 - prevents continuous ws closure
				"nginx.ingress.kubernetes.io/proxy-read-timeout": "1800",
				"nginx.ingress.kubernetes.io/cors-allow-origin":  strings.Join(allowOrigins, ","),
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("/game/%s/%s", key, id) + "(/|$)(.*)",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: Name(key, id),
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
