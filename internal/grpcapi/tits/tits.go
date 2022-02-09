package titspbv1

import (
	titsv1pb "github.com/boobsrate/apis/tits/v1"
	"github.com/boobsrate/core/internal/domain"
)

func titsToProto(src domain.Tits) *titsv1pb.Tits {
	return &titsv1pb.Tits{
		Id:     src.ID,
		Rating: src.Rating,
		ImgUrl: src.URL,
	}
}

func titsListToProto(src []domain.Tits) []*titsv1pb.Tits {
	res := make([]*titsv1pb.Tits, len(src))
	for i, v := range src {
		res[i] = titsToProto(v)
	}
	return res
}
