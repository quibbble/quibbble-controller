package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/quibbble/quibbble-controller/pkg/proto/controller"
	"github.com/quibbble/quibbble-controller/pkg/proto/sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// clean looks for stale games and removes them.
func (c *Controller) clean() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		l, err := c.clientset.CoreV1().Pods(c.config.Namespace).List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s,%s=%s", ComponentLabel, GameComponent, StaleLabel, strconv.FormatBool(true)),
		})
		if err != nil {
			continue
		}
		for _, it := range l.Items {
			gk := &controller.GameKey{
				Repository: it.Labels[RepositoryLabel],
				Tag:        it.Labels[TagLabel],
				Name:       it.Name,
			}
			if it.Labels[CompletedLabel] != strconv.FormatBool(true) {
				addr := fmt.Sprintf("%s.%s:%s", it.Name, c.config.Namespace, c.config.Game.SDKServer.Port)
				conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					continue
				}
				client := sdk.NewSDKClient(conn)
				snapshot, err := client.GetSnapshot(context.Background(), &sdk.Player{Id: c.config.Game.SDKServer.ControllerID, Name: "quibbble"})
				if err != nil {
					continue
				}
				gk.Snapshot = snapshot
				c.StoreGame(gk)
			}
			c.DeleteGame(gk)
		}
	}
}
