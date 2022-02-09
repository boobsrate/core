package initiator

import "context"

type TitsService interface {
	CreateTitsFromFile(ctx context.Context, filename, filePath string) error
}
