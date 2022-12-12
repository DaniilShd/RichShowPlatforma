package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/DaniilShd/RichShowPlatforma/intermal/config"
)

var app *config.AppConfig

// sets up config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func IsAuthenticate(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}

func IsAdmin(r *http.Request) bool {
	if app.Session.Exists(r.Context(), "access_level") {
		if app.Session.GetInt(r.Context(), "access_level") == 1 {
			return true
		}
	}
	return false
}

func IsManager(r *http.Request) bool {
	if app.Session.Exists(r.Context(), "access_level") {
		if app.Session.GetInt(r.Context(), "access_level") == 2 {
			return true
		}
	}
	return false
}

func IsStore(r *http.Request) bool {
	if app.Session.Exists(r.Context(), "access_level") {
		if app.Session.GetInt(r.Context(), "access_level") == 3 {
			return true
		}
	}
	return false
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of ", status)

	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ConvertNumberPhone(number string) string {
	if len(number) == 10 {
		number = "+7 (" + number[0:3] + ") " + number[3:6] + "-" + number[6:8] + "-" + number[8:10]
	}
	return number
}
