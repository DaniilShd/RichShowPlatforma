package main

import (
	"net/http"

	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/justinas/nosurf"
)

func WriteConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

// NoSurf adds CSRF protaction to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfhandler := nosurf.New(next)

	csrfhandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfhandler
}

// Session load and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticate(r) {
			session.Put(r.Context(), "error", "Сначала авторизуйтесь!")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		switch session.Get(r.Context(), "access_level") {
		case 1:
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		case 2:
			http.Redirect(w, r, "/manager", http.StatusSeeOther)
		case 3:
			http.Redirect(w, r, "/store", http.StatusSeeOther)
		default:
			session.Put(r.Context(), "error", "Некорректный уровень доступа")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		// next.ServeHTTP(w, r)
	})
}

func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAdmin(r) {
			session.Put(r.Context(), "error", "Вы не администратор, доступ запрещен!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthManager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsManager(r) && !helpers.IsAdmin(r) {
			session.Put(r.Context(), "error", "Вы не менеджер, доступ запрещен!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsStore(r) && !helpers.IsAdmin(r) {
			session.Put(r.Context(), "error", "Вы не реквизитор, доступ запрещен!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
