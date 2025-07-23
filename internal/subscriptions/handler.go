package subscriptions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
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
		ResponseWithError(response, request, &MethodNotAllowedError{RequiredMethod: "POST"})
		return
	}
	body, errBody := io.ReadAll(request.Body)
	if errBody != nil {
		ResponseWithError(response, request, errBody)
		return
	}
	var sbscr Subscription
	errJson := json.Unmarshal(body, &sbscr)
	if errJson != nil {
		ResponseWithError(response, request, errJson)
		return
	}
	num, errAdd := SubscriptionCreate(sbscr)
	if errAdd != nil {
		ResponseWithError(response, request, errAdd)
		return
	}
	WriteResponse(response, request, []byte(strconv.Itoa(int(*num))))
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
		ResponseWithError(response, request, &MethodNotAllowedError{RequiredMethod: "GET"})
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID, errParam := strconv.Atoi(sIDParam)
	if errParam != nil || sID < 1 {
		ResponseWithError(response, request, &InvalidParameterError{ParamName: "rowId"})
		return
	}
	item, errRead := SubscriptionRead(int32(sID))
	if errRead != nil {
		ResponseWithError(response, request, errRead)
		return
	}
	responseJson, errJson := json.Marshal(item)
	if errJson != nil {
		ResponseWithError(response, request, errJson)
		return
	}
	WriteResponse(response, request, responseJson)
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
		ResponseWithError(response, request, &MethodNotAllowedError{RequiredMethod: "PUT"})
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
		ResponseWithError(response, request, errJson)
		return
	}
	errUpdate := SubscriptionUpdate(sbscr)
	if errUpdate != nil {
		ResponseWithError(response, request, errUpdate)
		return
	}
	WriteResponse(response, request, nil)
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
		ResponseWithError(response, request, &MethodNotAllowedError{RequiredMethod: "DELETE"})
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID, errParam := strconv.Atoi(sIDParam)
	if errParam != nil || sID < 1 {
		ResponseWithError(response, request, &InvalidParameterError{ParamName: "rowId"})
		return
	}
	errDel := SubscriptionDelete(int32(sID))
	if errDel != nil {
		ResponseWithError(response, request, errDel)
		return
	}
	WriteResponse(response, request, nil)
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
		ResponseWithError(response, request, &MethodNotAllowedError{RequiredMethod: "GET"})
		return
	}
	pageParam := request.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageParam)
	if page < 1 {
		page = 1
	}
	list, errList := SubscriptionList(page)
	if errList != nil {
		ResponseWithError(response, request, errList)
		return
	}
	jsonList, errJson := json.Marshal(list)
	if errJson != nil {
		ResponseWithError(response, request, errJson)
		return
	}
	WriteResponse(response, request, jsonList)
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
		ResponseWithError(response, request, &MethodNotAllowedError{RequiredMethod: "POST"})
		return
	}
	filterFromParam := request.FormValue("filterFrom")
	filterToParam := request.FormValue("filterTo")
	if filterFromParam == "" || filterToParam == "" {
		ResponseWithError(response, request, &InvalidParameterError{ParamName: "filter dates"})
		return
	}
	filterFrom, _ := time.Parse("01-2006", filterFromParam)
	filterTo, _ := time.Parse("01-2006", filterToParam)

	sum, errSum := SubscriptionSum(filterFrom, filterTo, nil, nil)
	if errSum != nil {
		ResponseWithError(response, request, errSum)
		return
	}
	WriteResponse(response, request, []byte(strconv.Itoa(sum)))
}

func ResponseWithError(response http.ResponseWriter, request *http.Request, err error) {
	var (
		valErr         *ValidationError
		jsonErr        *JsonError
		paramErr       *InvalidParameterError
		notFoundErr    *ResourceNotFoundError
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
	LogRequest(request, nil, err)
}

func WriteResponse(response http.ResponseWriter, request *http.Request, data []byte) {
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(data)
	LogRequest(request, data, nil)
}

func LogRequest(r *http.Request, data []byte, err error) {
	if err == nil {
		log.Printf("%s: response [%s] %s %s with data %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RemoteAddr, r.RequestURI, data)
	} else {
		log.Printf("%s: response [%s] %s %s with error %v\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RemoteAddr, r.RequestURI, err)
	}
}
