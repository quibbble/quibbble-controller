package k8s

import (
	"fmt"
	"strconv"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodConfig struct {
	Image struct {
		Repository string            `yaml:"repository"`
		Tag        string            `yaml:"tag"`
		PullPolicy corev1.PullPolicy `yaml:"pullPolicy"`
	} `yaml:"image"`
	Resources    corev1.ResourceRequirements `yaml:"resources"`
	NodeSelector map[string]string           `yaml:"nodeSelector"`
	Affinity     *corev1.Affinity            `yaml:"affinity"`
	Tolerations  []corev1.Toleration         `yaml:"tolerations"`
}

func CreatePod(fullname, key, id string, port int32, storageEnabled bool, config *PodConfig) *corev1.Pod {
	envs := []corev1.EnvVar{
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
						Name: fullname,
					},
				},
			},
		},
		{
			Name:  "PORT",
			Value: strconv.Itoa(int(port)),
		},
	}
	if storageEnabled {
		envs = append(envs, corev1.EnvVar{
			Name: "STORAGE_PASSWORD", // password used to connect to the game store
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "storage-password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: fullname,
					},
				},
			},
		})
	}
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
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
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            Name(key, id),
					Image:           fmt.Sprintf("%s:%s", config.Image.Repository, config.Image.Tag),
					ImagePullPolicy: config.Image.PullPolicy,
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
					Env:       envs,
					Resources: config.Resources,
				},
			},
			NodeSelector: config.NodeSelector,
			Affinity:     config.Affinity,
			Tolerations:  config.Tolerations,
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
								Name: fullname,
							},
						},
					},
				},
			},
		},
	}
}
