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
	"github.com/DaniilShd/RichShowPlatforma/intermal/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/render"
	modelsTelegram "github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
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

	//Запускаю в отдельной горутине телеграм бота
	go StartTelegramBot()

	defer close(app.MailChan)
	defer close(app.RequestFromTelegram)
	defer close(app.UpdateCacheAccount)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	_ = srv.ListenAndServe()

}

func run() (*driver.DB, error) {

	//Создаю канал для передачи сообщений
	mailChan := make(chan modelsTelegram.MailData)
	app.MailChan = mailChan

	UpdateCacheAccount := make(chan bool)
	app.UpdateCacheAccount = UpdateCacheAccount

	RequestFromTelegram := make(chan modelsTelegram.RequestFromChat)
	app.RequestFromTelegram = RequestFromTelegram

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

	go listenChannelToRequestDBFromTelegram(&app, db)

	return db, nil
}

func listenChannelToRequestDBFromTelegram(a *config.AppConfig, db *driver.DB) {
	m := handlers.NewRepository(&app, db)
	for {
		request := <-app.RequestFromTelegram
		switch request.Command {
		case "get_lead":
			lead, err := m.DB.GetLeadByID(request.LeadID)
			if err != nil {
				log.Fatal(err)
			}
			request.ResponseLeadFromApp <- *lead

		case "get_order_by_id":
			var result []models.StoreLead
			for _, id := range request.StoreOrderID {
				storeOrder, err := m.DB.GetStoreOrderByID(id)
				if err != nil {
					log.Fatal(err)
				}
				result = append(result, *storeOrder)
			}
			request.ResponseLeadFromApp <- result
		case "get_order":
			storeOrder, err := m.DB.GetStoreOrderByID(request.LeadID)
			if err != nil {
				log.Fatal(err)
			}
			request.ResponseLeadFromApp <- *storeOrder
		case "get_checklist":
			result, err := m.DB.GetCheckListByID(request.CheckListID)
			if err != nil {
				log.Fatal(err)
			}
			request.ResponseLeadFromApp <- *result
		case "get_new_order":
			result, err := m.DB.GetAllNewStoreOrder()
			if err != nil {
				log.Fatal(err)
			}
			request.ResponseLeadFromApp <- result
		case "get_compl_order":
			result, err := m.DB.GetAllCompleteStoreOrder()
			if err != nil {
				log.Fatal(err)
			}
			request.ResponseLeadFromApp <- result
		case "get_dest_order":
			result, err := m.DB.GetAllToDestroyStoreOrder()
			if err != nil {
				log.Fatal(err)
			}
			request.ResponseLeadFromApp <- result
		}

	}
}
