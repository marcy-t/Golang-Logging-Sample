package infrastructure

import (
	"github.com/Golang-Logging-Sample/interfaces"
	db "github.com/Golang-Logging-Sample/pkg/interfaces/database"

	"github.com/go-chi/chi/v5"
)

// ルーティングでコントローラを別にしたい場合は新規で追加
type ControllHandler struct {
	Common *interfaces.CommonController
	// Admin *interfaces.AdminController // 増えていくとここに追加
}

func NewServer(h db.SqlHandler) (handler *chi.Mux) {
	// Handler
	ch := &ControllHandler{
		Common: interfaces.NewController(h), // Controller増えるごとに追加
		// Admin: adminController // 初期化されたコントローラー追加
	}
	handler = NewRouter(ch)
	return
}
