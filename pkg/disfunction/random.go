package disfunction

import (
	"context"
	"time"

	"github.com/karlhepler/disfunction/pkg/entity"
)

type RandomReq struct {
	context.Context
	Opts struct {
		Since time.Time
		Until time.Time
		Kinds []entity.Kind
	}
	Deps struct {
		GitHub interface {
			//
		}
	}
}

type RandomRes interface {
	Send(RandomMsg)
}

type RandomMsg struct {
	//
}

func Random(req RandomReq, res RandomRes) {
	var gh = req.Deps.GitHub
}
