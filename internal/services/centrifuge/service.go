package centrifuge

import (
	"context"

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
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}

		select {
		case <-ctx.Done():
			return
		case msg := <-s.wsChannel:
			b, err := msg.MarshalJSON()
			if err != nil {
				s.log.Error("failed to marshal message", zap.Error(err))
				continue
			}

			_, err = s.cli.Publish(context.Background(), &centrifugeApi.PublishRequest{
				Channel: s.chanName,
				Data:    b,
			})
			if err != nil {
				s.log.Error("error while publishing message to centrifuge", zap.Error(err))
			}
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
