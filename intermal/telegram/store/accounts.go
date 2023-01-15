package store

import (
	"fmt"
	"log"
	"sync"

	"github.com/DaniilShd/RichShowPlatforma/intermal/driver"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/repository/dbrepo"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/service/accounts"
)

type Accounts struct {
	mu         sync.RWMutex
	AccountMap map[string][]accounts.Person
}

var Account Accounts

func NewAccount(db *driver.DB) {
	DB := dbrepo.NewPostgresRepoTelegram(db.SQL)

	accounts := make(map[string][]accounts.Person)

	var err error

	accounts["admin"], err = DB.GetAdminChatID()
	if err != nil {
		log.Fatal(err)
	}

	accounts["manager"], err = DB.GetManagerChatID()
	if err != nil {
		log.Fatal(err)
	}

	accounts["assistant"], err = DB.GetAssistantsChatID()
	if err != nil {
		log.Fatal(err)
	}

	accounts["artist"], err = DB.GetArtistChatID()
	if err != nil {
		log.Fatal(err)
	}

	accounts["store"], err = DB.GetStoreChatID()
	if err != nil {
		log.Fatal(err)
	}
	Account.mu.RLock()
	Account.AccountMap = accounts
	Account.mu.RUnlock()
	fmt.Println(accounts["artist"])
}
