package controller

import (
	"context"

	q "github.com/quibbble/quibbble-controller/pkg/quibbble"
)

type Controller struct {
	q.UnimplementedControllerServer
	ctx context.Context
}

func NewController() (*Controller, error)

func (c *Controller) start() {
	// TODO
	// create a watching process that cleans stale games
}

func (c *Controller) Create(gk *q.GameKey) error {
	// TODO
	// Check to ensure game doesn't already exist
	// create ingress + service
	// create pod with sdk server and game server containers
	return nil
}

func (c *Controller) Delete(gk *q.GameKey) error {
	// TODO
	// look for game with matching name and kind
	// delete ingress to prevent anyone from connecting anymore
	// call Store method to store current state if the game is incomplete
	// delete service
	// delete pod with sdk server and game server containers
	return nil
}

func (c *Controller) Store(gk *q.GameKey) error {
	// TODO
	// if storage is enabled
	// look for pod with matching name and kind
	// GetSnapshot and store in DB
	return nil
}

func (c *Controller) GetActivity() (*q.Activity, error) {
	// TODO
	// lookup all pods with quibbble labels/annotations
	// get current game snapshot and determine how many players are active
	return nil, nil
}
