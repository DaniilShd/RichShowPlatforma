package repository

import (
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/models"
	"github.com/DaniilShd/RichShowPlatforma/intermal/telegram/service/accounts"
)

type DatabaseRepoTelegram interface {
	GetStoreChatID() ([]accounts.Person, error)
	GetArtistChatID() ([]accounts.Person, error)
	GetAssistantsChatID() ([]accounts.Person, error)
	GetAdminChatID() ([]accounts.Person, error)
	GetManagerChatID() ([]accounts.Person, error)
	//Find by ID
	GetStoreByChatID(id int64) (*accounts.Person, error)
	GetArtistByChatID(id int64) (*accounts.Person, error)
	GetAssistantsByChatID(id int64) (*accounts.Person, error)
	GetAdminByChatID(id int64) (*accounts.Person, error)
	GetManagerByChatID(id int64) (*accounts.Person, error)
	//Set chat id by id and role
	SetChatIDByRoleAndID(role int, idAccount int, chatID int64) error
	//Artist
	GetAllLeadsOfArtistByChatID(chatID int64) ([]models.LeadList, error)
	GetAllLeadsOfArtistTodayByChatID(chatID int64) ([]models.LeadList, error)
	//Store
	GetOrderStoreIDByLeadID(leadID int) ([]int, error)
	//Assistant
	GetAllLeadsOfAssistantByChatID(chatID int64) ([]models.LeadList, error)
	GetAllLeadsOfAssistantTodayByChatID(chatID int64) ([]models.LeadList, error)
}
