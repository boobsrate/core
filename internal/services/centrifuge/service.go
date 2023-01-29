package centrifuge

import (
	"context"
	"time"

	"github.com/boobsrate/core/internal/config"
	"github.com/boobsrate/core/internal/domain"
	centrifugeApi "github.com/boobsrate/core/pkg/centrifugal"
	grpc2 "github.com/boobsrate/core/pkg/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	wsChannel chan domain.WSMessage
	cli       centrifugeApi.CentrifugoApiClient
	chanName  string

	log *zap.Logger
}

func NewService(wsChannel chan domain.WSMessage, cfg config.CentrifugeConfiguration, chanName string, log *zap.Logger) (*Service, error) {

	cc, err := grpc2.NewGrpcClient(cfg.GRPCAddress, grpc.WithPerRPCCredentials(keyAuth{cfg.ApiToken}))

	if err != nil {
		return nil, err
	}

	apiCli := centrifugeApi.NewCentrifugoApiClient(cc)

	return &Service{
		cli:       apiCli,
		log:       log.Named("centrifuge_service"),
		chanName:  chanName,
		wsChannel: wsChannel,
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("starting centrifuge service")
	resp, err := s.cli.Info(context.Background(), &centrifugeApi.InfoRequest{})
	if err != nil {
		s.log.Error("error getting info", zap.Error(err))
	}
	s.log.Info("centrifuge info", zap.Any("resp", resp))

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				resp, err := s.cli.Info(context.Background(), &centrifugeApi.InfoRequest{})
				if err != nil {
					s.log.Error("error getting info", zap.Error(err))
				}
				clientCount := 0
				for _, node := range resp.GetResult().GetNodes() {
					clientCount += int(node.GetNumClients())
				}
				s.log.Info("centrifuge info", zap.Any("resp", resp), zap.Int("client_count", clientCount))
				// make online msg
				msg := domain.WSMessage{
					Type: domain.WSMessageTypeOnlineUsers,
					Message: domain.WSOnlineUsersMessage{
						Online:    clientCount,
					},
				}

				// send to centrifuge
				b, err := msg.MarshalJSON()
				if err != nil {
					s.log.Error("failed to marshal message", zap.Error(err))
					continue
				}
				respB, err := s.cli.Broadcast(context.Background(), &centrifugeApi.BroadcastRequest{
					Channels: []string{s.chanName},
					Data:     b,
				})
				if err != nil {
					s.log.Error("error while publishing message to centrifuge", zap.Error(err))
				}

				s.log.Info("published message to centrifuge", zap.Any("resp", respB))
			default:

			}
		}
	}(ctx)


	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.wsChannel:
			b, err := msg.MarshalJSON()
			if err != nil {
				s.log.Error("failed to marshal message", zap.Error(err))
				continue
			}

			resp, err := s.cli.Broadcast(context.Background(), &centrifugeApi.BroadcastRequest{
				Channels: []string{s.chanName},
				Data:     b,
			})
			if err != nil {
				s.log.Error("error while publishing message to centrifuge", zap.Error(err))
			}

			s.log.Info("published message to centrifuge", zap.Any("resp", resp))
		}
	}
}

type keyAuth struct {
	key string
}

func (t keyAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "apikey " + t.key,
	}, nil
}

func (t keyAuth) RequireTransportSecurity() bool {
	return false
}
