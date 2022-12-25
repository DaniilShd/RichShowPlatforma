package repository

import "github.com/DaniilShd/RichShowPlatforma/intermal/models"

type DatabaseRepo interface {
	Authenticate(login, testPassword string) (int, int, string, error)

	//Check-List
	GetAllCheckList() (*[]models.CheckList, error)
	GetCheckListByID(id int) (*models.CheckList, error)
	GetAllCheckListsOfType(id int) (*[]models.CheckList, error)
	DeleteCheckListByID(id int) error
	InsertCheckList(checkList *models.CheckList) error
	UpdateCheckList(checkList *models.CheckList) error

	//Store
	GetAllStoreItem() ([]models.StoreItem, error)
	GetStoreItemByID(id int) (*models.StoreItem, error)
	DeleteStoreItemByID(id int) error
	InsertStoreItem(storeItem *models.StoreItem) error
	UpdateStoreItem(storeItem *models.StoreItem) error
	//StoreOder
	GetAllNewStoreOrder() ([]models.StoreLead, error)
	GetAllCompleteStoreOrder() ([]models.StoreLead, error)
	GetAllToDestroyStoreOrder() ([]models.StoreLead, error)
	GetCountNewStoreOrder() (int, error)
	GetCountCompleteStoreOrder() (int, error)
	GetCountToDestroyStoreOrder() (int, error)
	GetStoreOrderByID(id int) (*models.StoreLead, error)
	// SetCompleteStoreOrder(id int) error
	// UnSetCompleteStoreOrder(id int) error
	DeleteStoreOrderByID(id int) error
	InsertStoreOrder(storeOrder *models.StoreLead) error
	UpdateStoreOrder(storeOrder *models.StoreLead) error

	//Manager
	InsertLead(lead *models.Lead) error
	GetAllRawLeads() ([]models.Lead, error)
	GetAllConfirmedLeads() ([]models.Lead, error)
	GetAllArchiveLeads() ([]models.Lead, error)
	GetCountRawLeads() (int, error)
	GetCountConfirmedLeads() (int, error)
	GetCountArchiveLeads() (int, error)
	GetLeadByID(idLead int) (*models.Lead, error)
	SetConfirmedLeadByID(idLead int) error
	DeleteConfirmedLeadByID(idLead int) error
	DeleteLeadByID(idLead int) error
	UpdateLead(lead *models.Lead) error

	//Assistant
	GetAllAssistants() (*[]models.Assistant, error)
	GetAssistantByID(id int) (*models.Assistant, error)
	InsertAssistant(assistant *models.Assistant) error
	UpdateAssistant(assistant *models.Assistant) error
	DeleteAssistantByID(id int) error

	//Heroes
	GetAllHeroes() (*[]models.Hero, error)
	GetHeroByID(id int) (*models.Hero, error)
	InsertHero(hero *models.Hero) error
	UpdateHero(hero *models.Hero) error
	DeleteHeroByID(id int) error

	//Artist
	GetAllArtists() (*[]models.Artist, error)
	GetArtistByID(id int) (*models.Artist, error)
	InsertArtist(artist *models.Artist) error
	UpdateArtist(artist *models.Artist) error
	DeleteArtistByID(id int) error
}
