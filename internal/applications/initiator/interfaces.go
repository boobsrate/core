package initiator

import "context"

type TitsService interface {
	CreateTitsFromFile(ctx context.Context, filename, filePath string) error
	CreateTitsFromBytes(ctx context.Context, filename string, file []byte) error
}
