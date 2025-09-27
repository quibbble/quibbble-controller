package sdkserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/proto/game"
	"github.com/quibbble/quibbble-controller/pkg/proto/sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type connection struct {
	stream sdk.SDK_StreamSnapshotServer
	player *sdk.Player
	err    chan error
}

// TODO
type SDKServer struct {
	sdk.UnimplementedSDKServer

	clientset kubernetes.Interface
	pod       string
	namespace string

	client    game.GameClient
	createdAt time.Time
	updatedAt time.Time
	history   []*sdk.Action

	// map from player id to player
	players map[string]*sdk.Player
	// map from player id to stream
	connected map[string]sdk.SDK_StreamSnapshotServer
	mu        sync.RWMutex

	broadcastCh chan string
}

func NewSDKServer() (*SDKServer, error) {
	host := os.Getenv("QUIBBBLE_GAME_GRPC_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("QUIBBBLE_GAME_GRPC_PORT")
	if port == "" {
		port = "9357"
	}

	pod := os.Getenv("POD")
	if pod == "" {
		return nil, fmt.Errorf("pod env variable is not set")
	}
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		return nil, fmt.Errorf("namespace env variable is not set")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to %s", addr)
	}

	sdkServer := &SDKServer{
		clientset: clientset,
		pod:       pod,
		namespace: namespace,

		client:    game.NewGameClient(conn),
		createdAt: time.Now(),
		updatedAt: time.Now(),
		// TODO
	}

	go sdkServer.stale()
	go sdkServer.broadcast()
	return sdkServer, nil
}

func (s *SDKServer) broadcast() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(string(debug.Stack()))
		}
	}()

	for range <-s.broadcastCh {
		s.mu.RLock()
		for id, stream := range s.connected {
			snapshot, err := s.GetSnapshot(context.Background(), s.players[id])
			if err != nil {
				stream.Context().Err()
				continue
			}
			stream.Send(snapshot)
		}
		s.mu.Unlock()
	}
}

func (s *SDKServer) GetSnapshot(ctx context.Context, player *sdk.Player) (*sdk.Snapshot, error) {
	snapshot, err := s.client.GetSnapshot(ctx, &game.View{
		Team: player.Team,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get game snapshot")
	}
	players := make([]*sdk.Player, 0)
	s.mu.RLock()
	defer s.mu.Unlock()
	for _, player := range s.players {
		players = append(players, player)
	}
	return &sdk.Snapshot{
		Players:   players,
		Snapshot:  snapshot,
		History:   s.history,
		UpdatedAt: timestamppb.New(s.updatedAt),
	}, nil
}

func (s *SDKServer) StreamSnapshot(player *sdk.Player, stream sdk.SDK_StreamSnapshotServer) error {
	s.mu.Lock()
	if s.connected[player.Id] != nil {
		return fmt.Errorf("%s[%s] is already connected", player.Name, player.Id)
	}
	if s.players[player.Id] != nil {
		player.Team = s.players[player.Id].Team
	}
	s.connected[player.Id] = stream
	s.players[player.Id] = player
	s.mu.Unlock()

	// Tell everyone a new player has joined.
	s.broadcastCh <- ""

	// Wait for the stream to close.
	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.players, player.Id)
	s.mu.Unlock()
	return nil
}

func (s *SDKServer) JoinTeam(ctx context.Context, player *sdk.Player) (*emptypb.Empty, error) {
	s.mu.Lock()
	s.players[player.Id] = player
	s.updatedAt = time.Now()
	s.mu.Unlock()
	return nil, nil
}

// TODO if winner then update completed label to true
func (s *SDKServer) PlayAction(context.Context, *sdk.Action) (*emptypb.Empty, error)

// TODO
func (s *SDKServer) UndoAction(context.Context, *sdk.Player) (*emptypb.Empty, error)

// TODO reset completed label to false
func (s *SDKServer) ResetGame(context.Context, *sdk.Player) (*emptypb.Empty, error)
