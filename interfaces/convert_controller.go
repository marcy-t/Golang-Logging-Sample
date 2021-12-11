package interfaces

import (
	"github.com/marcy-t/Golang-Logging-Sample/domain"
)

// レシーバー用途オンリー
type ConvertController struct {
	Converter CommonConverter
}

type CommonConverter interface {
	// Request Convert
	ToSampleEntity(string) *domain.SampleRequest
	// Response Convert
	ToSampleResponseData([]*domain.User) []*domain.User
}

func NewConvertController() *ConvertController {
	return &ConvertController{}
}

/*
	In/Outのデータ整形に使用
*/

// 変換系 サンプル
func (cc *ConvertController) ToSampleEntity(host string) (resp *domain.SampleRequest) {
	return &domain.SampleRequest{
		Host: host,
	}
}

func (cc *ConvertController) ToSampleResponseData(users []*domain.User) (items []*domain.User) {
	items = make([]*domain.User, len(users))
	for i, u := range users {
		items[i] = &domain.User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	return
}
