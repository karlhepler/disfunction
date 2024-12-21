package disfunction

import (
	"context"
	"time"

	"github.com/karlhepler/disfunction/cmd/pkg/entity"
)

type RandomReq struct {
	context.Context
	Since time.Time
	Until time.Time
	Kinds []entity.Kind
}

func Random(req RandomReq, res RandomRes) {
	//
}
