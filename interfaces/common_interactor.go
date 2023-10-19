package interfaces

import (
	"context"

	"github.com/Golang-Logging-Sample/domain"
)

type CommonInteractor interface {
	UseCaseSampleRepository(context.Context) ([]*domain.User, error)
}
