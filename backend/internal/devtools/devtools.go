// Package devtools berisi halaman bantu yang hanya aktif saat APP_ENV=development.
package devtools

import (
	_ "embed"
	"net/http"
)

//go:embed login.html
var loginPage []byte

// Login: GET /dev/login — coba login Telegram & SMS dari browser.
func Login(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(loginPage)
}
