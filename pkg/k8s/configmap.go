package k8s

import (
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(snapshot *qgn.Snapshot) *corev1.ConfigMap {
	key := snapshot.Tags[qgn.KeyTag]
	id := snapshot.Tags[qgn.IDTag]
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
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
		Data: map[string]string{
			"qgn": snapshot.String(),
		},
	}
}
