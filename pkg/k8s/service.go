package k8s

import (
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(key, id string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: Name(key, id),
			Labels: map[string]string{
				Component:  GameComponent,
				qgn.KeyTag: key,
				qgn.IDTag:  id,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromString("http"),
				},
			},
			Selector: map[string]string{
				qgn.KeyTag: key,
				qgn.IDTag:  id,
			},
		},
	}
}
