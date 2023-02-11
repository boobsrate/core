package abyss

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const (
	defaultCheckTimeout   = time.Second * 10
	defaultAbyssThreshold = 2
)

type Keeper struct {
	log     *zap.Logger
	service Service
	dead    chan struct{}
}

func NewKeeper(log *zap.Logger, service Service) *Keeper {
	return &Keeper{
		log:     log.Named("abyss_keeper"),
		service: service,
		dead:    make(chan struct{}),
	}
}

func (k *Keeper) Dead() chan struct{} {
	return k.dead
}

func (k *Keeper) iterate(ctx context.Context) {
	tits, err := k.service.GetTitsWithReportsThreshold(ctx, defaultAbyssThreshold)
	if err != nil {
		k.log.Error("get tit with reports threshold", zap.Error(err))
		return
	}
	for _, tit := range tits {
		k.log.Info("move tit to abyss", zap.String("tit_id", tit.ID))
		if err := k.service.MoveToAbyss(ctx, tit.ID); err != nil {
			k.log.Error("move tit to abyss", zap.String("tit_id", tit.ID), zap.Error(err))
		}
	}
}

func (k *Keeper) Run(ctx context.Context) {
	k.log.Info("abyss keeper starting...")
	defer close(k.dead)
	defer k.log.Info("abyss keeper stopped")

	ticker := time.NewTicker(defaultCheckTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go k.iterate(ctx)
		}
	}
}
