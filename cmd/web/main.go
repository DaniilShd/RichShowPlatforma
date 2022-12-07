package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DaniilShd/RichShowPlatforma/intermal/config"
	"github.com/DaniilShd/RichShowPlatforma/intermal/driver"
	"github.com/DaniilShd/RichShowPlatforma/intermal/handlers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/helpers"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	_ = srv.ListenAndServe()

}

func run() (*driver.DB, error) {

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=localhost port=5432 dbname=richshow user=postgres password=root")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.SQL.Close()
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	handlers.NewHandlers(handlers.NewRepository(&app, db))

	return db, nil
}
