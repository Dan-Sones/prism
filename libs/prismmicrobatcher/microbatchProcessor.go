package services

import "context"

type MicrobatchProcessor interface {
	ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error
}
