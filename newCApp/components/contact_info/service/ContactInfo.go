package service

import (
	"contactApp/errors"
	"contactApp/models/contactinfo"
	"contactApp/repository"
	"time"

	"github.com/jinzhu/gorm"
)

// ContactInfoService Give Access to Update, Add, Delete User
type ContactInfoService struct {
	db           *gorm.DB
	repository   repository.Repository
	associations []string
}

// NewContactInfoService returns new instance of ContactInfoService
func NewContactInfoService(db *gorm.DB, repo repository.Repository) *ContactInfoService {
	return &ContactInfoService{
		db:           db,
		repository:   repo,
		associations: []string{},
	}
}

func (service *ContactInfoService) CreateContact(newUser *contactinfo.ContactInfo) error {
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

func (service *ContactInfoService) doesContactReferExist(contactID string) error {
	// intID,_ := strconv.Atoi(userID)
	// uintID := uint(intID)
	exists, err := repository.DoesRecordExistForContactInfo(service.db, contactID, contactinfo.ContactInfo{},
		repository.Filter("`contact_refer` = ?", contactID))
	if !exists || err != nil {
		return errors.NewValidationError("ContactRefer ID is Invalid")
	}
	return nil
}

func (service *ContactInfoService) doesContactInfoExist(ID uint) error {
	exists, err := repository.DoesRecordExistForUser(service.db, ID, contactinfo.ContactInfo{},
		repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return errors.NewValidationError("ContactInfo ID is Invalid")
	}
	return nil
}

func (service *ContactInfoService) GetContact(contact *contactinfo.ContactInfo, contactId string, ID int) error {
	// Start new transcation.
	err := service.doesContactReferExist(contactId)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err = service.repository.GetRecordForAll(uow,uint(ID), contact)
	if err != nil {
		return err
	}
	uow.Commit()
	return nil
}

func (service *ContactInfoService) GetContacts(newContactInfos *[]contactinfo.ContactInfo,contactId string) error {
	err := service.doesContactReferExist(contactId)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err = service.repository.GetRecordForContactInfos(uow,contactId, newContactInfos)
	// err = service.repository.GetAll(uow, newContacts)
	if err != nil {
		return err
	}
	uow.Commit()
	return nil
}

func (service *ContactInfoService) UpdateContactInfo(contactToUpdate *contactinfo.ContactInfo,contactID string) error {
	err := service.doesContactReferExist(contactID)
	if err != nil {
		return err
	}
	err = service.doesContactInfoExist(contactToUpdate.ID)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()
	tempUser := contactinfo.ContactInfo{}
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

func (service *ContactInfoService) DeleteContactInfo(contact *contactinfo.ContactInfo,contactID string) error {
	err := service.doesContactInfoExist(contact.ID)
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