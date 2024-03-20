package k8s

import (
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePod(key, id string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name(key, id),
			Namespace: Namespace,
			Labels: map[string]string{
				Component:  GameComponent,
				qgn.KeyTag: key,
				qgn.IDTag:  id,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            Name(key, id),
					Image:           "docker.io/quibbble/server:latest",
					ImagePullPolicy: "Always", // todo reset to IfNotPresent
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      Name(key, id),
							MountPath: "/root/qgn",
							SubPath:   "qgn",
						},
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							HostPort:      8080,
							ContainerPort: 8080,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: Name(key, id),
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: Name(key, id),
							},
						},
					},
				},
			},
		},
	}
}
