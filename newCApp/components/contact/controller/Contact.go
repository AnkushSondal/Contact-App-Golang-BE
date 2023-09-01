package controller

import (
	"contactApp/components/contact/service"
	"contactApp/components/log"
	"contactApp/errors"
	"contactApp/models/contact"
	"contactApp/web"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ContactController gives access to CRUD operations for entity
type ContactController struct {
	log     log.Log
	service *service.ContactService
}

// NewContactController returns new instance of ContactController
func NewContactController(contactService *service.ContactService,
	log log.Log) *ContactController {
	return &ContactController{
		service: contactService,
		log:     log,
	}
}
func (controller *ContactController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/user/{userid}/contact").Subrouter()
	userRouter.HandleFunc("/", controller.CreateContact).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id}", controller.GetContact).Methods(http.MethodGet)
	userRouter.HandleFunc("/", controller.GetContacts).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", controller.UpdateContact).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}", controller.DeleteContact).Methods(http.MethodDelete)

}
func (controller *ContactController) CreateContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in cr con")

	newUser := contact.Contact{}
	// Unmarshal json.
	err := web.UnmarshalJSON(r, &newUser)
	fmt.Println(newUser)
	vars := mux.Vars(r)
	userID,_:= strconv.ParseUint(vars["userid"], 10, 64)
	newUser.UserRefer=uint(userID)
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

func (controller *ContactController) GetContacts(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "GetContact successfull.")
	contacts := &[]contact.Contact{}
	var totalCount int

	vars := mux.Vars(r)
	userID:= vars["userid"]

	err := controller.service.GetContacts(contacts,userID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, contacts)
}

func (controller *ContactController) GetContact(w http.ResponseWriter, r *http.Request) {
	contact := &contact.Contact{}
	var totalCount int =1
	vars := mux.Vars(r)
	userID := vars["userid"]
	intCID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	err = controller.service.GetContact(contact,userID,intCID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResonseWriter,
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, contact)
}

func (controller *ContactController) UpdateContact(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "UpdateContact successfull.")
	fmt.Println("==============================userToUpdate==========================")
	contactToUpdate := contact.Contact{}

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
	userID := vars["userid"]
	intCID, err := strconv.Atoi(vars["id"])
	// intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	contactToUpdate.ID = uint(intCID)
	fmt.Println("==============================userToUpdate==========================")
	fmt.Println(&contactToUpdate)
	// Call update test method.
	err = controller.service.UpdateContact(&contactToUpdate,userID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, contactToUpdate)
}
func (controller *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request) {
	// web.RespondJSON(w, http.StatusOK, "DeleteContact successfull.")
	controller.log.Print("********************************DeleteTest call**************************************")
	conntactToDelete := contact.Contact{}
	var err error
	vars := mux.Vars(r)
	userID := vars["userid"]
	intCID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	conntactToDelete.ID = uint(intCID)
	err = controller.service.DeleteContact(&conntactToDelete,userID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Delete Contact successfull.")
}
