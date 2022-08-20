package app

import (
	"context"
	"os"
	"time"
)

type options struct {
	ctx              context.Context
	sigs             []os.Signal
	registrarTimeout time.Duration
	stopTimeout      time.Duration
}
