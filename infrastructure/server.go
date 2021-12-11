package infrastructure

import (
	"github.com/gorilla/mux"
	"github.com/marcy-t/Golang-Logging-Sample/interfaces"
	db "github.com/marcy-t/Golang-Logging-Sample/pkg/interfaces/database"
)

// ルーティングでコントローラを別にしたい場合は新規で追加
type ControllHandler struct {
	Common *interfaces.CommonController
	// Admin *interfaces.AdminController // 増えていくとここに追加
}

func NewServer(h db.SqlHandler) (handler *mux.Router) {
	// Handler
	ch := &ControllHandler{
		Common: interfaces.NewController(h), // Controller増えるごとに追加
		// Admin: adminController // 初期化されたコントローラー追加
	}
	handler = NewRouter(ch)
	return
}
