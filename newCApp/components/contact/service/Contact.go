package service

import (
	"contactApp/errors"
	"contactApp/models/contact"
	"contactApp/repository"
	"time"

	"github.com/jinzhu/gorm"
)

// ContactService Give Access to Update, Add, Delete Contact
type ContactService struct {
	db           *gorm.DB
	repository   repository.Repository
	associations []string
}

// NewContactService returns new instance of ContactService
func NewContactService(db *gorm.DB, repo repository.Repository) *ContactService {
	return &ContactService{
		db:           db,
		repository:   repo,
		associations: []string{},
	}
}

func (service *ContactService) CreateContact(newUser *contact.Contact) error {
	//  Creating unit of work.
	// err := service.doesUserExist(newUser.UserRefer)
	// if err != nil {
	// 	return err
	// }
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()
	// Add newUser.
	err := service.repository.Add(uow, newUser)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

func (service *ContactService) doesUserReferExist(userID string) error {
	// intID,_ := strconv.Atoi(userID)
	// uintID := uint(intID)
	exists, err := repository.DoesRecordExistForContact(service.db, userID, contact.Contact{},
		repository.Filter("`user_refer` = ?", userID))
	if !exists || err != nil {
		return errors.NewValidationError("UserRefer ID is Invalid")
	}
	return nil
}

func (service *ContactService) doesContactExist(ID uint) error {
	exists, err := repository.DoesRecordExistForUser(service.db, ID, contact.Contact{},
		repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return errors.NewValidationError("Contact ID is Invalid")
	}
	return nil
}

func (service *ContactService) GetContacts(newContacts *[]contact.Contact,userId string) error {
	err := service.doesUserReferExist(userId)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err = service.repository.GetRecordForContacts(uow,userId, newContacts)
	// err = service.repository.GetAll(uow, newContacts)
	if err != nil {
		return err
	}
	uow.Commit()
	return nil

}

func (service *ContactService) GetContact(contact *contact.Contact, userId string, contactID int) error {
	// Start new transcation.
	err := service.doesUserReferExist(userId)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err = service.repository.GetRecordForAll(uow,uint(contactID), contact)
	if err != nil {
		return err
	}
	uow.Commit()
	return nil
}


func (service *ContactService) UpdateContact(contactToUpdate *contact.Contact,userID string) error {
	err := service.doesUserReferExist(userID)
	if err != nil {
		return err
	}
	err = service.doesContactExist(contactToUpdate.ID)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()
	tempUser := contact.Contact{}
	err = service.repository.GetRecordForAll(uow, contactToUpdate.ID, &tempUser, repository.Select("`created_at`"),
		repository.Filter("`id` = ?", contactToUpdate.ID))
	if err != nil {
		return err
	}
	contactToUpdate.CreatedAt = tempUser.CreatedAt

	err = service.repository.Save(uow, contactToUpdate)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

func (service *ContactService) DeleteContact(contact *contact.Contact,userID string) error {
	err := service.doesContactExist(contact.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	// Update test for updating deleted_by and deleted_at fields of test
	if err := service.repository.UpdateWithMap(uow, contact, map[string]interface{}{

		"DeletedAt": time.Now(),
	},
		repository.Filter("`id`=?", contact.ID)); err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}
