package subscriptions

import (
	"encoding/json"
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
		http.Error(response, "Only POST method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	var sbscr Subscription
	errJ := json.Unmarshal(body, &sbscr)
	if errJ != nil {
		http.Error(response, errJ.Error(), http.StatusBadRequest)
		return
	}
	num, errC := SubscriptionCreate(sbscr)
	if errC != nil {
		http.Error(response, errC.Error(), http.StatusInternalServerError)
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
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/read [get]
func SubscriptionReadHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(response, "Only GET method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID, err := strconv.Atoi(sIDParam)
	if err != nil || sID < 1 {
		http.Error(response, "invalid data", http.StatusBadRequest)
		return
	}
	item, errR := SubscriptionRead(int32(sID))
	if errR != nil {
		http.Error(response, errR.Error(), http.StatusInternalServerError)
		return
	}
	responseJson, errJ := json.Marshal(item)
	if errJ != nil {
		http.Error(response, errJ.Error(), http.StatusInternalServerError)
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
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/update [put]
func SubscriptionUpdateHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPut {
		http.Error(response, "Only PUT method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	var sbscr Subscription
	errJ := json.Unmarshal(body, &sbscr)
	if err != nil {
		http.Error(response, errJ.Error(), http.StatusBadRequest)
		return
	}
	errU := SubscriptionUpdate(sbscr)
	if errU != nil {
		http.Error(response, errJ.Error(), http.StatusInternalServerError)
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
//	@Failure	405		{string}	string			"error"
//	@Failure	400		{string}	string			"error"
//	@Failure	500		{string}	string			"error"
//	@Router		/subscription/delete [delete]
func SubscriptionDeleteHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(response, "Only DELETE method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID, err := strconv.Atoi(sIDParam)
	if err != nil || sID < 1 {
		http.Error(response, "invalid data", http.StatusBadRequest)
		return
	}
	errD := SubscriptionDelete(int32(sID))
	if errD != nil {
		http.Error(response, errD.Error(), http.StatusInternalServerError)
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
		http.Error(response, "Only GET method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	pageParam := request.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageParam)
	if page < 1 {
		page = 1
	}
	list, err := SubscriptionList(page)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonList, errJ := json.Marshal(list)
	if errJ != nil {
		http.Error(response, errJ.Error(), http.StatusInternalServerError)
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
		http.Error(response, "Only POST method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	filterFromParam := request.FormValue("filterFrom")
	filterToParam := request.FormValue("filterTo")
	if filterFromParam == "" || filterToParam == "" {
		http.Error(response, "invalid data", http.StatusBadRequest)
		return
	}
	filterFrom, _ := time.Parse("01-2006", filterFromParam)
	filterTo, _ := time.Parse("01-2006", filterToParam)

	sum, err := SubscriptionSum(filterFrom, filterTo, nil, nil)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Write([]byte(strconv.Itoa(sum)))
}
