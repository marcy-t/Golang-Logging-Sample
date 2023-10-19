package usecase

import (
	"context"

	"github.com/Golang-Logging-Sample/domain"
)

/*
UseCaseも目的の処理に応じてディレクトリ化して使用
*/
type CommonRepository interface {
	Find(ctx context.Context) (resp []*domain.User, err error)
}
