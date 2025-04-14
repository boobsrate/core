package centrifuge

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/boobsrate/core/internal/config"
	"github.com/boobsrate/core/internal/domain"
	"github.com/boobsrate/core/internal/services/buryat"
	centrifugeApi "github.com/boobsrate/core/pkg/centrifugal"
	grpc2 "github.com/boobsrate/core/pkg/grpc"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Service struct {
	wsChannel chan domain.WSMessage
	cli       centrifugeApi.CentrifugoApiClient
	chanName  string
	buryat    *buryat.Service

	log *zap.Logger
}

func NewService(wsChannel chan domain.WSMessage, buryat *buryat.Service, cfg config.CentrifugeConfiguration, chanName string, log *zap.Logger) (*Service, error) {

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
		buryat:    buryat,
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("starting centrifuge service")
	resp, err := s.cli.Info(context.Background(), &centrifugeApi.InfoRequest{})
	if err != nil {
		s.log.Error("error getting info", zap.Error(err))
	}
	s.log.Info("centrifuge info", zap.Any("resp", resp))

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.log.Info("tick online")
				resp, err := s.cli.Info(context.Background(), &centrifugeApi.InfoRequest{})
				clientCount := 0
				for _, node := range resp.GetResult().GetNodes() {
					clientCount += int(node.GetNumClients())
				}
				s.log.Info("centrifuge info", zap.Any("resp", resp), zap.Int("client_count", clientCount))

				if err != nil {
					s.log.Error("error getting info", zap.Error(err))
				}
				s.log.Info("centrifuge info", zap.Any("resp", resp))
				msg := domain.WSMessage{
					Type: domain.WSMessageTypeOnlineUsers,
					Message: domain.WSOnlineUsersMessage{
						Online: clientCount,
					},
				}

				// send to centrifuge
				b, err := msg.MarshalJSON()
				if err != nil {
					s.log.Error("failed to marshal message", zap.Error(err))
					return
				}
				respB, err := s.cli.Broadcast(context.Background(), &centrifugeApi.BroadcastRequest{
					Channels: []string{s.chanName},
					Data:     b,
				})
				if err != nil {
					s.log.Error("error while publishing message to centrifuge", zap.Error(err))
				}

				s.log.Info("published message to centrifuge", zap.Any("resp", respB))
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(300 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.log.Info("tick")
				if rand.Intn(100) > 20 {
					continue
				}
				resp, err := s.buryat.GetResponse([]openai.ChatCompletionMessage{{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Current date[%d]. Здарова. Расскажи о чем-нибудь интересном, историю с работы, как твои дела? Расскажи длинную историю про твой день.", time.Now().Unix()),
				}})
				if err != nil {
					s.log.Error("failed to get response from buryat", zap.Error(err))
					return
				}
				var bMsg domain.WSMessage
				bMsg.Message = domain.WSChatMessage{
					Text:   resp,
					Sender: "Ебанько Бурят",
				}
				bMsg.Type = domain.WSMessageTypeChat
				s.wsChannel <- bMsg
			}
		}
	}()

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

			if msg.Type == domain.WSMessageTypeChat {
				go func() {
					if rand.Intn(100) > 20 {
						return
					}

					resp, err := s.buryat.GetResponse([]openai.ChatCompletionMessage{{
						Role:    openai.ChatMessageRoleUser,
						Content: msg.Message.(domain.WSChatMessage).Text,
					}})
					if err != nil {
						s.log.Error("failed to get response from buryat", zap.Error(err))
						return
					}
					var bMsg domain.WSMessage
					bMsg.Message = domain.WSChatMessage{
						Text:   resp,
						Sender: "Ебанько Бурят",
					}
					bMsg.Type = domain.WSMessageTypeChat
					bb, err := bMsg.MarshalJSON()
					bres, err := s.cli.Broadcast(context.Background(), &centrifugeApi.BroadcastRequest{
						Channels: []string{s.chanName},
						Data:     bb,
					})
					if err != nil {
						s.log.Error("error while publishing message to centrifuge", zap.Error(err))
					}
					s.log.Info("published message to centrifuge", zap.Any("resp", bres))
				}()
			}

			s.log.Info("published message to centrifuge", zap.Any("resp", resp))

			go func() {
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
						Online: clientCount,
					},
				}

				// send to centrifuge
				b, err := msg.MarshalJSON()
				if err != nil {
					s.log.Error("failed to marshal message", zap.Error(err))
					return
				}
				respB, err := s.cli.Broadcast(context.Background(), &centrifugeApi.BroadcastRequest{
					Channels: []string{s.chanName},
					Data:     b,
				})
				if err != nil {
					s.log.Error("error while publishing message to centrifuge", zap.Error(err))
				}

				s.log.Info("published message to centrifuge", zap.Any("resp", respB))
			}()
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
