package k8s

import (
	"strconv"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePod(key, id, image, pullyPolicy string, port int32) *corev1.Pod {
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
					Image:           image,
					ImagePullPolicy: corev1.PullPolicy(pullyPolicy),
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "qgn-vol",
							MountPath: "/root/qgn",
							SubPath:   "qgn",
						},
						{
							Name:      "config-vol",
							MountPath: "/root/config.yaml",
							SubPath:   "config.yaml",
						},
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							HostPort:      port,
							ContainerPort: port,
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "ADMIN_USERNAME",
							Value: "quibbble",
						},
						{
							Name: "ADMIN_PASSWORD",
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									Key: "admin-password",
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ChartName,
									},
								},
							},
						},
						{
							Name: "STORAGE_PASSWORD", // password used to connect to the game store
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									Key: "storage-password",
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ChartName,
									},
								},
							},
						},
						{
							Name:  "PORT",
							Value: strconv.Itoa(int(port)),
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "qgn-vol",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: Name(key, id),
							},
						},
					},
				},
				{
					Name: "config-vol",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: ChartName,
							},
						},
					},
				},
			},
		},
	}
}
