package interfaces

import (
	"context"

	"github.com/marcy-t/Golang-Logging-Sample/domain"
)

type CommonInteractor interface {
	UseCaseSampleRepository(context.Context) ([]*domain.User, error)
}
