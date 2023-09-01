package controller

import (
	"contactApp/components/contact_info/service"
	"contactApp/components/log"
	"contactApp/errors"
	"contactApp/models/contactinfo"
	"contactApp/web"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ContactInfoController gives access to CRUD operations for entity
type ContactInfoController struct {
	log     log.Log
	service *service.ContactInfoService
}

// NewContactInfoController returns new instance of ContactInfoController
func NewContactInfoController(contactService *service.ContactInfoService,
	log log.Log) *ContactInfoController {
	return &ContactInfoController{
		service: contactService,
		log:     log,
	}
}
func (controller *ContactInfoController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/contact/{contactid}/contactinfo").Subrouter()
	userRouter.HandleFunc("/register", controller.CreateContactInfo).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id}", controller.GetContactInfo).Methods(http.MethodGet)
	userRouter.HandleFunc("/", controller.GetContactInfos).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", controller.UpdateContactInfo).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}", controller.DeleteContactInfo).Methods(http.MethodDelete)

}
func (controller *ContactInfoController) CreateContactInfo(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "CreateContactInfo successfull.")
	fmt.Println("in cr con")

	newUser := contactinfo.ContactInfo{}
	// Unmarshal json.
	err := web.UnmarshalJSON(r, &newUser)
	vars := mux.Vars(r)
	contactID,_:= strconv.ParseUint(vars["contactid"], 10, 64)
	newUser.ContactRefer=uint(contactID)
	fmt.Println(newUser)
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	// Call add test method.
	err = controller.service.CreateContact(&newUser)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResposeWriter.
	web.RespondJSON(w, http.StatusCreated, newUser)

}
func (controller *ContactInfoController) GetContactInfo(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "GetContactInfo successfull.")
	contact := &contactinfo.ContactInfo{}
	var totalCount int =1
	vars := mux.Vars(r)
	contactID := vars["contactid"]
	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	err = controller.service.GetContact(contact,contactID,intID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResonseWriter,
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, contact)

}

func (controller *ContactInfoController) GetContactInfos(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "GetContact successfull.")
	contactinfos := &[]contactinfo.ContactInfo{}
	var totalCount int

	vars := mux.Vars(r)
	contactID:= vars["contactid"]

	err := controller.service.GetContacts(contactinfos,contactID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, contactinfos)
}

func (controller *ContactInfoController) UpdateContactInfo(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "UpdateContactInfo successfull.")
	fmt.Println("==============================userToUpdate==========================")
	contactToUpdate := contactinfo.ContactInfo{}

	// Unmarshal JSON.
	fmt.Println(r.Body)
	err := web.UnmarshalJSON(r, &contactToUpdate)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	vars := mux.Vars(r)
	contactID := vars["contactid"]
	intID, err := strconv.Atoi(vars["id"])
	// intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	contactToUpdate.ID = uint(intID)
	fmt.Println("==============================userToUpdate==========================")
	fmt.Println(&contactToUpdate)
	// Call update test method.
	err = controller.service.UpdateContactInfo(&contactToUpdate,contactID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, contactToUpdate)

	
}
func (controller *ContactInfoController) DeleteContactInfo(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "DeleteContactInfo successfull.")
	controller.log.Print("********************************DeleteTest call**************************************")
	conntactToDelete := contactinfo.ContactInfo{}
	var err error
	vars := mux.Vars(r)
	contactID := vars["contactid"]
	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	conntactToDelete.ID = uint(intID)
	err = controller.service.DeleteContactInfo(&conntactToDelete,contactID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Delete Contact Info successfull.")

}