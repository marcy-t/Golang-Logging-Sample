package interfaces

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Golang-Logging-Sample/domain"
	db "github.com/Golang-Logging-Sample/pkg/interfaces/database"
	"github.com/Golang-Logging-Sample/pkg/logger"
	"github.com/Golang-Logging-Sample/usecase"
)

type CommonController struct {
	Interactor CommonInteractor
	Converter  CommonConverter
}

/*
	Goって承継の代わりにinterfaceという機能があり、
	interface内にmethod名を定義して型として定義できる
	それとおなじmethodを実装した構造体(type)用意すると
	そのinterfaceをキャストすることができ、呼び出し側はそのinterfaceの関数を
	使うことができる
*/

func NewController(SqlHandler db.SqlHandler) (cc *CommonController) {
	// UseCase interface
	cc = &CommonController{
		Interactor: &usecase.CommonInteractor{
			CommonRepository: &CommonRepository{
				DB: SqlHandler,
			},
		},
	}
	cc.Converter = NewConvertController()
	return
}

/*
	(レシーバー) Function名() (err error) での記述必須
*/

// GET /ping
func (cc *CommonController) SampleHandler(w http.ResponseWriter, r *http.Request) (appResp *domain.ApplicationResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	// Single Tags
	logger.Info("111", "GET /api/v1//ping")

	// Debig
	hostUrl := cc.Converter.ToSampleEntity("http://hogehoge.com")
	log.Println(hostUrl)

	// UseCase
	resp, err := cc.Interactor.UseCaseSampleRepository(ctx)
	if err != nil {
		logger.Error(
			logger.GetApplicationError(err).ErrorResponseJSON(w, r),
			// tags...,
		)
		return
	}

	convData := cc.Converter.ToSampleResponseData(resp)

	appResp = &domain.ApplicationResponse{
		Type:       "SampleHandler Response",
		StatusCode: http.StatusOK,
		AppCode:    "xxxxxx", // 設計時になんのアプリケーションか連番ふる
		Message:    convData,
	}

	return
}

// GET /taglist
func (cc *CommonController) TagList(w http.ResponseWriter, r *http.Request) (err error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	// defer cancel()

	logger.Info("111", "GET /api/v1/taglist")

	// tags := []*logger.Tag{
	// 	logger.NewTag("taglistStatus", "Enable"),
	// 	logger.NewTag("GetApplicationError", "AddMessage"),
	// 	logger.NewTag("GetApplicationError", "AddMessage"),
	// }

	// err := errors.New("Err!!LogggNew!!!!!!")
	// appErr := logger.GetApplicationError(err).AddMessage("An error has occured")
	// appErr.AppCode = "12344"
	// log.Println()
	return
}
