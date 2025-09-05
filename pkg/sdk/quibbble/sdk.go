package sdk

import (
	"context"

	"github.com/quibbble/quibbble-controller/pkg/sdk"
)

// TODO
type SDKServer struct {
	sdk.UnimplementedSDKServer
	ctx context.Context
}
