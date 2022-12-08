package repository

import "github.com/DaniilShd/RichShowPlatforma/intermal/models"

type DatabaseRepo interface {
	Authenticate(login, testPassword string) (int, int, string, error)

	//Master-class
	UpdateMasterClass(masterClass *models.MasterClass) error
	DeleteMasterClassByID(id int) error
	GetAllMasterClass() (*[]models.MasterClass, error)
	GetMasterClassByID(id int) (*models.MasterClass, error)
	InsertMasterClass(masterClass *models.MasterClass) error

	//Check-List
	GetAllCheckList() (*[]models.CheckList, error)
	GetCheckListByID(id int) (*models.CheckList, error)
	GetAllCheckListsOfType(id int) (*[]models.CheckList, error)
	DeleteCheckListByID(id int) error
	InsertCheckList(checkList *models.CheckList) error
	UpdateCheckList(checkList *models.CheckList) error
}
