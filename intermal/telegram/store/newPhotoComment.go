package store

import (
	"sync"

	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/service/accounts"
)

type NewPhotoComment struct {
	mu         sync.RWMutex
	AccountMap map[string][]accounts.Person
}

var newPhotoComment NewPhotoComment
