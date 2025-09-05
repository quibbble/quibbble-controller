package sdk

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SnapshotHandler does some work whenever a game snapshot is passed.
type SnapshotHandler func(gs *sdk.Snapshot)

// SDK is a simple wrapper around the SDK Client.
type SDK struct {
	client sdk.SDKClient
	ctx    context.Context
}

func NewSDK() (*SDK, error) {
	host := os.Getenv("QUIBBBLE_SDK_GRPC_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("QUIBBBLE_SDK_GRPC_PORT")
	if port == "" {
		port = "9357"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to %s", addr)
	}

	return &SDK{
		client: sdk.NewSDKClient(conn),
		ctx:    context.Background(),
	}, nil
}

func (s *SDK) GetSnapshot(player *sdk.Player) (*sdk.Snapshot, error) {
	snapshot, err := s.client.GetSnapshot(s.ctx, player)
	return snapshot, errors.Wrap(err, "could not get game snapshot")
}

func (s *SDK) StreamSnapshot(player *sdk.Player, f SnapshotHandler) error {
	stream, err := s.client.StreamSnapshot(s.ctx, player)
	if err != nil {
		return errors.Wrap(err, "could not stream the game")
	}

	for {
		snapshot, err := stream.Recv()
		if err == io.EOF {
			break // stream finished
		}
		if err != nil {
			return errors.Wrap(err, "could not continue streaming")
		}
		f(snapshot)
	}
	return nil
}

func (s *SDK) JoinTeam(player *sdk.Player) error {
	_, err := s.client.JoinTeam(s.ctx, player)
	return errors.Wrap(err, "failed to join team")
}

func (s *SDK) PlayAction(playerAction *sdk.Action) error {
	_, err := s.client.PlayAction(s.ctx, playerAction)
	return errors.Wrap(err, "failed to play action")
}

func (s *SDK) UndoAction(player *sdk.Player) error {
	_, err := s.client.UndoAction(s.ctx, player)
	return errors.Wrap(err, "failed to undo action")
}

func (s *SDK) ResetGame(player *sdk.Player) error {
	_, err := s.client.ResetGame(s.ctx, player)
	return errors.Wrap(err, "failed to reset game")
}
