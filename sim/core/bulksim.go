package core

import (
	"context"

	"github.com/wowsims/wotlk/sim/core/proto"
)

func BulkSim(ctx context.Context, request *proto.BulkSimRequest, progress chan *proto.ProgressMetrics) (*proto.BulkSimResult, error) {
	panic("not implemented")
}
