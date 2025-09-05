package sdk

import (
	"context"

	q "github.com/quibbble/quibbble-controller/pkg/quibbble"
)

// TODO
type SDKServer struct {
	q.UnimplementedSDKServer
	ctx context.Context
}
