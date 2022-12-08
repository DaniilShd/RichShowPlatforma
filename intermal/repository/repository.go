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

	//Show program
	UpdateProgramShow(programShow *models.ProgramShow) error
	DeleteProgramShowByID(id int) error
	GetAllProgramShow() (*[]models.ProgramShow, error)
	GetProgramShowByID(id int) (*models.ProgramShow, error)
	InsertProgramShow(programShow *models.ProgramShow) error

	// Animation
	UpdateAnimation(animation *models.Animation) error
	DeleteAnimationByID(id int) error
	GetAllAnimation() (*[]models.Animation, error)
	GetAnimationByID(id int) (*models.Animation, error)
	InsertAnimation(animation *models.Animation) error

	// Party
	UpdateParty(party *models.Party) error
	DeletePartyByID(id int) error
	GetAllParty() (*[]models.Party, error)
	GetPartyByID(id int) (*models.Party, error)
	InsertParty(party *models.Party) error

	// Others
	UpdateOther(other *models.Other) error
	DeleteOtherByID(id int) error
	GetAllOther() (*[]models.Other, error)
	GetOtherByID(id int) (*models.Other, error)
	InsertOther(other *models.Other) error

	//Check-List
	GetAllCheckList() (*[]models.CheckList, error)
	GetCheckListByID(id int) (*models.CheckList, error)
	GetAllCheckListsOfType(id int) (*[]models.CheckList, error)
	DeleteCheckListByID(id int) error
	InsertCheckList(checkList *models.CheckList) error
	UpdateCheckList(checkList *models.CheckList) error
}
