package subscriptions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

// SubscriptionCreateHandler godoc
//
//	@Summary	record creation
//	@Tags		subscriptions
//	@Accept		json
//	@Produce	plain
//	@Param		item	body		Subscription	true	"item to add"
//	@Success	200		{integer}	string			"created"
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/create [post]
func SubscriptionCreateHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		ResponseWithError(response, &MethodNotAllowedError{RequiredMethod: "POST"})
		return
	}
	body, errBody := io.ReadAll(request.Body)
	if errBody != nil {
		ResponseWithError(response, errBody)
		return
	}
	var sbscr Subscription
	errJson := json.Unmarshal(body, &sbscr)
	if errJson != nil {
		ResponseWithError(response, errJson)
		return
	}
	num, errAdd := SubscriptionCreate(sbscr)
	if errAdd != nil {
		ResponseWithError(response, errAdd)
		return
	}
	response.Write([]byte(strconv.Itoa(int(*num))))
}

// SubscriptionReadHandlergodoc
//
//	@Summary	record reading
//	@Tags		subscriptions
//	@Produce	json
//	@Param		rowId	query		integer			true	"record id"
//	@Success	200		{object}	Subscription	"data found"
//	@Failure	404		{string}	string			"error"
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/read [get]
func SubscriptionReadHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		ResponseWithError(response, &MethodNotAllowedError{RequiredMethod: "GET"})
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID, errParam := strconv.Atoi(sIDParam)
	if errParam != nil || sID < 1 {
		ResponseWithError(response, &InvalidParameterError{ParamName: "rowId"})
		return
	}
	item, errRead := SubscriptionRead(int32(sID))
	if errRead != nil {
		ResponseWithError(response, errRead)
		return
	}
	responseJson, errJson := json.Marshal(item)
	if errJson != nil {
		ResponseWithError(response, errJson)
		return
	}
	response.Write(responseJson)
}

// SubscriptionUpdateHandler godoc
//
//	@Summary	record update
//	@Tags		subscriptions
//	@Accept		json
//	@Param		item	body		Subscription	true	"item to update"
//	@Success	200		{integer}	string			"updated"
//	@Failure	404		{string}	string			"error"
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/update [put]
func SubscriptionUpdateHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		ResponseWithError(response, &MethodNotAllowedError{RequiredMethod: "PUT"})
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	var sbscr Subscription
	errJson := json.Unmarshal(body, &sbscr)
	if errJson != nil {
		ResponseWithError(response, errJson)
		return
	}
	errUpdate := SubscriptionUpdate(sbscr)
	if errUpdate != nil {
		ResponseWithError(response, errUpdate)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(nil)
}

// SubscriptionDeleteHandler
//
//	@Summary	record deleting
//	@Tags		subscriptions
//	@Produce	json
//	@Param		rowId	query		integer			true	"record id"
//	@Success	200		{object}	Subscription	"data deleted"
//	@Failure	404		{string}	string			"error"
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/delete [delete]
func SubscriptionDeleteHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		ResponseWithError(response, &MethodNotAllowedError{RequiredMethod: "DELETE"})
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID, errParam := strconv.Atoi(sIDParam)
	if errParam != nil || sID < 1 {
		ResponseWithError(response, &InvalidParameterError{ParamName: "rowId"})
		return
	}
	errDel := SubscriptionDelete(int32(sID))
	if errDel != nil {
		ResponseWithError(response, errDel)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(nil)
}

// SubscriptionListHandler
//
//	@Summary	list of records
//	@Tags		subscriptions
//	@Produce	json
//	@Param		page	query		integer					false	"page number"
//	@Success	200		{object}	SubscriptionListPage	"loaded successfully"
//	@Failure	405		{string}	string					"error"
//	@Failure	500		{string}	string					"error"
//	@Router		/subscription/list [get]
func SubscriptionListHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		ResponseWithError(response, &MethodNotAllowedError{RequiredMethod: "GET"})
		return
	}
	pageParam := request.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageParam)
	if page < 1 {
		page = 1
	}
	list, errList := SubscriptionList(page)
	if errList != nil {
		ResponseWithError(response, errList)
		return
	}
	jsonList, errJson := json.Marshal(list)
	if errJson != nil {
		ResponseWithError(response, errJson)
		return
	}
	response.Write(jsonList)
}

// SubscriptionSumHandler godoc
//
//	@Summary	sum calculation
//	@Tags		subscriptions
//	@Accept		x-www-form-urlencoded
//	@Produce	plain
//	@Param		filterFrom	formData	string	true	"period from"
//	@Param		filterTo	formData	string	true	"period to"
//	@Success	200			{integer}	string	"sum is ready"
//	@Failure	405			{string}	string	"error"
//	@Failure	400			{string}	string	"error"
//	@Failure	500			{string}	string	"error"
//	@Router		/subscription/sum [post]
func SubscriptionSumHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		ResponseWithError(response, &MethodNotAllowedError{RequiredMethod: "POST"})
		return
	}
	filterFromParam := request.FormValue("filterFrom")
	filterToParam := request.FormValue("filterTo")
	if filterFromParam == "" || filterToParam == "" {
		ResponseWithError(response, &InvalidParameterError{ParamName: "filter dates"})
		return
	}
	filterFrom, _ := time.Parse("01-2006", filterFromParam)
	filterTo, _ := time.Parse("01-2006", filterToParam)

	sum, errSum := SubscriptionSum(filterFrom, filterTo, nil, nil)
	if errSum != nil {
		ResponseWithError(response, errSum)
		return
	}
	response.Write([]byte(strconv.Itoa(sum)))
}

func ResponseWithError(response http.ResponseWriter, err error) {
	var (
		valErr *ValidationError
		jsonErr *JsonError
		paramErr *InvalidParameterError
		notFoundErr *ResourceNotFoundError
		wrongMethodErr *MethodNotAllowedError
	)
	switch {
	case errors.As(err, &valErr):
	case errors.As(err, &jsonErr):
	case errors.As(err, &paramErr):
		http.Error(response, err.Error(), http.StatusBadRequest)
	case errors.As(err, &notFoundErr):
		http.Error(response, err.Error(), http.StatusNotFound)
	case errors.As(err, &wrongMethodErr):
		http.Error(response, err.Error(), http.StatusMethodNotAllowed)
	case errors.Is(err, sql.ErrNoRows):
		http.Error(response, err.Error(), http.StatusNotFound)
	default:
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
