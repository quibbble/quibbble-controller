package controller

import (
	"context"
	"time"

	"github.com/quibbble/quibbble-controller/pkg/proto/controller"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	// ComponentLabel is one of game or platform.
	ComponentLabel = "component"
	// StaleLabel defines whether or not the game is stale and can be cleaned.
	StaleLabel = "stale"
	// CompletedLabel defines whether or not the game is over.
	CompletedLabel = "completed"
	// The image repository of the game.
	RepositoryLabel = "repository"
	// The tag of the game.
	TagLabel = "tag"

	GameComponent     = "game"
	PlatformComponent = "platform"

	timeout = 5 * time.Second
)

func (c *Controller) createGame(gk *controller.GameKey) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := c.clientset.CoreV1().ConfigMaps(c.config.Namespace).Create(ctx, createConfigMap(gk), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Pods(c.config.Namespace).Create(ctx, createPod(gk), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.clientset.CoreV1().Services(c.config.Namespace).Create(ctx, createService(gk), metav1.CreateOptions{}); err != nil {
		return err
	}
	if _, err := c.gatewayclientset.GatewayV1().GRPCRoutes(c.config.Namespace).Create(ctx, createGRPCRoute(gk), metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

// write game.Snapshot to configmap named game-<k8sResourceName>
// write history and players to another configmap sdk-<k8sResourceName>
// TODO
func createConfigMap(gk *controller.GameKey) *corev1.ConfigMap

func createPod(gk *controller.GameKey) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: gk.Name,
			Labels: map[string]string{
				ComponentLabel: GameComponent,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "sdk",
					// TODO
				},
				{
					Name: "game",
					// TODO
				},
			},
		},
	}

}

// TODO
func createService(gk *controller.GameKey) *corev1.Service

// TODO
func createGRPCRoute(gk *controller.GameKey) *gatewayv1.GRPCRoute
