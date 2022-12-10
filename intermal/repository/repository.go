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

	//Manager
	InsertLead(lead *models.Lead) error
}
