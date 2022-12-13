package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/config"
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
	"humanTime": HumanTime,
}

var pathToTemplate = "./templates"

var app *config.AppConfig

// sets the config for the template
func NewRenderer(a *config.AppConfig) {
	app = a
}

// returns time in YYYY-MM-DD format (10 Nov 09 23:00 UTC)
func HumanDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func HumanTime(t time.Time) string {
	return t.Format("15:04")
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticate = 1
		td.AccessLevel = app.Session.GetInt(r.Context(), "access_level")
	} else {
		td.IsAuthenticate = 0
	}
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	//get the template cache from the appconfig

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("myerr TC")
		return errors.New("New erorr template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil

	// err = t.Execute(w, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplate))
	if err != nil {
		log.Fatal(err)
		return myCache, err
	}
	// fmt.Println(pages)

	for _, page := range pages {
		name := filepath.Base(page)

		// fmt.Println(name)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Fatal(err)
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
		if err != nil {
			log.Fatal(err)
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
			if err != nil {
				log.Fatal(err)
				return myCache, err
			}

		}
		myCache[name] = ts
	}
	return myCache, nil
}
