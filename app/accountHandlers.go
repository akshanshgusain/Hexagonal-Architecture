package app

import (
	"encoding/json"
	"fmt"
	"github.com/akshanshgusain/Hexagonal-Architecture/dto"
	"github.com/akshanshgusain/Hexagonal-Architecture/service"
	"github.com/gorilla/mux"
	"net/http"
)

type AccountHandlers struct {
	service service.AccountService
}

func (a *AccountHandlers) createAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cId := vars["customer_id"]
	var request dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	request.CustomerId = cId
	account, appError := a.service.NewAccount(request)
	if appError != nil {
		writeResponse(w, appError.Code, appError.Message)
		return
	}
	writeResponse(w, http.StatusCreated, account)
}

func (a *AccountHandlers) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cId := vars["customer_id"]
	aId := vars["account_id"]

	// decoded incoming request
	var request dto.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// build request object
	request.AccountId = aId
	request.CustomerId = cId

	fmt.Println("Before service Make transaction")
	fmt.Println(request)

	// make transaction
	account, appError := a.service.MakeTransaction(request)

	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}
	fmt.Println("After service Make transaction")
	writeResponse(w, http.StatusOK, account)
}
