package service

import (
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/constant"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/store"
)

func ValidationChatID(ChatID int64) int {
	for _, value := range store.Account.AccountMap["admin"] {
		if value.ChatID == ChatID {
			return constant.ADMIN
		}
	}
	for _, value := range store.Account.AccountMap["manager"] {
		if value.ChatID == ChatID {
			return constant.MANAGER
		}
	}
	for _, value := range store.Account.AccountMap["assistant"] {
		if value.ChatID == ChatID {
			return constant.ASSISTANT
		}
	}
	for _, value := range store.Account.AccountMap["store"] {
		if value.ChatID == ChatID {
			return constant.STORE
		}
	}
	for _, value := range store.Account.AccountMap["artist"] {
		if value.ChatID == ChatID {
			return constant.ARTIST
		}
	}
	return 0
}

// Return const role and ID accounts
func ValidationPhoneNumber(numberPhone string) (int, int) {
	for _, value := range store.Account.AccountMap["admin"] {
		if value.PhoneNumber == numberPhone {
			return constant.ADMIN, value.ID
		}
	}
	for _, value := range store.Account.AccountMap["manager"] {
		if value.PhoneNumber == numberPhone {
			return constant.MANAGER, value.ID
		}
	}
	for _, value := range store.Account.AccountMap["assistant"] {
		if value.PhoneNumber == numberPhone {
			return constant.ASSISTANT, value.ID
		}
	}
	for _, value := range store.Account.AccountMap["store"] {
		if value.PhoneNumber == numberPhone {
			return constant.STORE, value.ID
		}
	}
	for _, value := range store.Account.AccountMap["artist"] {
		if value.PhoneNumber == numberPhone {
			return constant.ARTIST, value.ID
		}
	}
	return 0, 0
}
