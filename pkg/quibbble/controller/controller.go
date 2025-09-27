package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/quibbble/quibbble-controller/config"
	"github.com/quibbble/quibbble-controller/pkg/persist"
	"github.com/quibbble/quibbble-controller/pkg/proto/controller"
	"github.com/quibbble/quibbble-controller/pkg/proto/sdk"
	"github.com/quibbble/quibbble-controller/util/sqldb"
	"github.com/upper/db/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubernetesgateway "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type Controller struct {
	controller.UnimplementedControllerServer

	clientset        kubernetes.Interface
	gatewayclientset kubernetesgateway.Interface

	config *config.Config
	db     db.Session
}

func NewController(config *config.Config, clientset kubernetes.Interface, gatewayclientset kubernetesgateway.Interface) (*Controller, error) {
	c := &Controller{
		clientset:        clientset,
		gatewayclientset: gatewayclientset,
		config:           config,
	}
	if config.Persistence != nil {
		db, err := sqldb.CreateDBSession(context.Background(), clientset, config.Namespace, config.Persistence)
		if err != nil {
			return nil, err
		}
		c.db = db
	}

	go c.clean()
	return c, nil
}

func (c *Controller) CreateGame(gk *controller.GameKey) error {

	if err := validate(gk); err != nil {
		return err
	}

	game, err := c.lookupActiveGame(gk)
	if err != nil {
		return err
	}

	// A game was found in persistent storage so load this game instead of creating a new one.
	if game != nil {
		var snapshot *sdk.Snapshot
		if err := proto.Unmarshal(game.Snapshot, snapshot); err != nil {
			return err
		}
		gk.Snapshot = snapshot
	}

	return c.createGame(gk)
}

func (c *Controller) DeleteGame(gk *controller.GameKey) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var errList []error
	if err := c.clientset.CoreV1().ConfigMaps(c.config.Namespace).Delete(ctx, gk.Name, metav1.DeleteOptions{}); err != nil {
		errList = append(errList, err)
	}
	if err := c.clientset.CoreV1().Pods(c.config.Namespace).Delete(ctx, gk.Name, metav1.DeleteOptions{}); err != nil {
		errList = append(errList, err)
	}
	if err := c.clientset.CoreV1().Services(c.config.Namespace).Delete(ctx, gk.Name, metav1.DeleteOptions{}); err != nil {
		errList = append(errList, err)
	}
	if c.config.Game.Gateway.Enabled {
		if err := c.gatewayclientset.GatewayV1().GRPCRoutes(c.config.Namespace).Delete(ctx, gk.Name, metav1.DeleteOptions{}); err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}

func (c *Controller) StoreGame(gk *controller.GameKey) error {
	if c.config.Persistence == nil {
		return nil
	}

	snapshot, err := proto.Marshal(gk.Snapshot)
	if err != nil {
		return err
	}
	game := &persist.Game{
		Repository: gk.Repository,
		Tag:        gk.Tag,
		Name:       gk.Name,
		Snapshot:   snapshot,
	}

	if len(gk.Snapshot.Snapshot.Winners) > 0 {
		if _, err := c.db.Collection(persist.CompletedGamesTable).Insert(game); err != nil {
			return err
		}
	} else {
		if _, err := c.db.Collection(persist.ActiveGamesTable).Insert(game); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) GetActivity() (*controller.Activity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := c.clientset.CoreV1().Pods(c.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", ComponentLabel, GameComponent),
	})
	if err != nil {
		return nil, err
	}

	activity := &controller.Activity{
		ActiveGames:   make(map[string]int64),
		ActivePlayers: make(map[string]int64),
	}

	for _, it := range l.Items {
		addr := fmt.Sprintf("%s.%s:%s", it.Name, c.config.Namespace, c.config.Game.SDKServer.Port)
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		client := sdk.NewSDKClient(conn)
		snapshot, err := client.GetSnapshot(context.Background(), &sdk.Player{Id: c.config.Game.SDKServer.ControllerID, Name: "quibbble"})
		if err != nil {
			return nil, err
		}

		kind := strings.Split(it.Name, "-")[0]
		activity.ActiveGames[kind] += 1
		activity.ActivePlayers[kind] += int64(len(snapshot.Players))
	}
	return activity, nil
}
