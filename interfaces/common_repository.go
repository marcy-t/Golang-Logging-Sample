package interfaces

import (
	"context"

	"github.com/marcy-t/Golang-Logging-Sample/domain"
	db "github.com/marcy-t/Golang-Logging-Sample/pkg/interfaces/database"
	"github.com/marcy-t/Golang-Logging-Sample/pkg/logger"
)

type CommonRepository struct {
	DB db.SqlHandler
}

const (
	SELECT_TEST_USER string = "select id, name, email from test_users;"
)

func (r *CommonRepository) Find(ctx context.Context) (users []*domain.User, err error) {
	tags := []*logger.Tag{
		logger.NewTag("select", "table:test_users"),
		logger.NewTag("email", "xxxx@mail.co.jp"),
	}

	rows, err := r.DB.Query(ctx, SELECT_TEST_USER)
	if err != nil {
		err = logger.GetApplicationError(err).Init("xx", "An error has occured in Select Users Error.")
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger.Error(
				logger.GetApplicationError(err).
					Init("xx", "DB Connection failed to close."),
			)
		}
	}()

	// init
	users = make([]*domain.User, 0)

	for rows.Next() {
		userTable := domain.User{}
		if err = rows.Scan(
			&userTable.ID,
			&userTable.Name,
			&userTable.Email,
		); err != nil {
			err = logger.GetApplicationError(err).Init("xx", "An error has occured in scanning process on finding User.")
			return
		}
		users = append(users, &userTable)
		// Or
		// users = append(users, &domain.User{
		// 	ID:    userTable.ID,
		// 	Name:  userTable.Name,
		// 	Email: userTable.Email,
		// })

	}
	logger.Info("xx", "Users retrieved successfully from database.", tags...)
	return
}
