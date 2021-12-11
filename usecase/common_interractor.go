package usecase

import (
	"context"

	"github.com/marcy-t/Golang-Logging-Sample/domain"
	"github.com/marcy-t/Golang-Logging-Sample/pkg/logger"
)

type CommonInteractor struct {
	CommonRepository CommonRepository
}

/*
	interface層にある関数を使用しロジックを作成する
*/
func (i *CommonInteractor) UseCaseSampleRepository(ctx context.Context) (resp []*domain.User, err error) {
	// Interface
	resp, err = i.CommonRepository.Find(ctx)
	if err != nil {
		return nil, err
	}

	tag := logger.NewTag("UseCase", "Repository")
	logger.Info("xxx", "UseCaseSampleRepository HandlerInfo", tag)
	return
}
