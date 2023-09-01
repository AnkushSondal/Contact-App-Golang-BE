package controller

import (
	"contactApp/components/log"
	"contactApp/components/user/service"
	"contactApp/errors"
	"contactApp/models/login"
	"contactApp/models/user"
	"contactApp/web"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// UserController gives access to CRUD operations for entity
type UserController struct {
	log     log.Log
	service *service.UserService
}

// NewUserController returns new instance of UserController
func NewUserController(userService *service.UserService, log log.Log) *UserController {
	return &UserController{
		service: userService,
		log:     log,
	}
}
func (controller *UserController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", controller.RegisterUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/login", controller.handleLogin).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id}", controller.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/", controller.GetAllUsers).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", controller.UpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}", controller.DeleteUser).Methods(http.MethodDelete)
	fmt.Println("==============================userRegisterRoutes==========================")
}

func (controller *UserController) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user = login.Login{}
	// err := json.NewDecoder(r.Body).Decode(&user)
	err := web.UnmarshalJSON(r, &user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	type Claims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}
	var jwtSecret = []byte("your-secret-key")

	isValid := controller.service.HandleLogin(&user)
	if isValid!=nil {
		expirationTime := jwt.TimeFunc().Add(24 * time.Hour)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
			// http.Error(w, "Error creating token", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"token": tokenString,
			"id":    isValid.ID,
			"isadmin" : isValid.IsAdmin,
		}
		// json.NewEncoder(w).Encode(response)
		web.RespondJSON(w, http.StatusCreated, response)

	} else {
		web.RespondError(w, errors.NewHTTPError("Invalid credentials", http.StatusUnauthorized))
		// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func (controller *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	newUser := user.User{}
	// Unmarshal json.
	err := web.UnmarshalJSON(r, &newUser)
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	// Call add test method.
	err = controller.service.CreateUser(&newUser)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResposeWriter.
	web.RespondJSON(w, http.StatusCreated, newUser)
}

func (controller *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	Users := &user.User{}
	var totalCount int =1
	vars := mux.Vars(r)
	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	err = controller.service.GetUser(Users,intID)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResonseWriter,
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, Users)
}

func (controller *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers := &[]user.User{}
	var totalCount int
	err := controller.service.GetAllUsers(allUsers, &totalCount)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResonseWriter,
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allUsers)
}
func (controller *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==============================userToUpdate==========================")
	userToUpdate := user.User{}

	// Unmarshal JSON.
	fmt.Println(r.Body)
	err := web.UnmarshalJSON(r, &userToUpdate)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	vars := mux.Vars(r)

	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	userToUpdate.ID = uint(intID)
	fmt.Println("==============================userToUpdate==========================")
	fmt.Println(&userToUpdate)
	// Call update test method.
	err = controller.service.UpdateUser(&userToUpdate)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, userToUpdate)
}
func (controller *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	controller.log.Print("********************************DeleteTest call**************************************")
	usetToDelete := user.User{}
	var err error
	vars := mux.Vars(r)
	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	usetToDelete.ID = uint(intID)
	err = controller.service.DeleteUser(&usetToDelete)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Delete User successfull.")
}
