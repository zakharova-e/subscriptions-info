package subscriptions

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

func SubscriptionCreateHandler(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPost{
		http.Error(response,"Only POST method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	body,err := io.ReadAll(request.Body)
	if err!=nil{
		http.Error(response,err.Error(), http.StatusBadRequest)
		return
	}
	var sbscr Subscription
	errJ := json.Unmarshal(body,&sbscr)
	if err!=nil{
		http.Error(response,errJ.Error(), http.StatusBadRequest)
		return
	}
	num,errC := SubscriptionCreate(sbscr)
	if errC!= nil{
		http.Error(response,errC.Error(), http.StatusInternalServerError)
		return
	}
	response.Write([]byte(strconv.Itoa(int(*num))))
} 

func SubscriptionReadHandler(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodGet{
		http.Error(response,"Only GET method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID,err:= strconv.Atoi(sIDParam)
	if err!= nil || sID<1{
		http.Error(response,"invalid data", http.StatusBadRequest)
		return
	}
	item,errR := SubscriptionRead(int32(sID))
	if errR!=nil{
		http.Error(response,errR.Error(), http.StatusInternalServerError)
		return
	}
	responseJson,errJ := json.Marshal(item)
	if errJ!=nil{
		http.Error(response,errJ.Error(), http.StatusInternalServerError)
		return
	}
	response.Write(responseJson)
}

func SubscriptionUpdateHandler(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPut{
		http.Error(response,"Only PUT method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	body,err := io.ReadAll(request.Body)
	if err!=nil{
		http.Error(response,err.Error(), http.StatusBadRequest)
		return
	}
	var sbscr Subscription
	errJ := json.Unmarshal(body,&sbscr)
	if err!=nil{
		http.Error(response,errJ.Error(), http.StatusBadRequest)
		return
	}
	errU := SubscriptionUpdate(sbscr)
	if errU!=nil{
		http.Error(response,errJ.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(nil)
}

func SubscriptionDeleteHandler(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodDelete{
		http.Error(response,"Only DELETE method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	sIDParam := request.URL.Query().Get("rowId")
	sID,err:= strconv.Atoi(sIDParam)
	if err!= nil || sID<1{
		http.Error(response,"invalid data", http.StatusBadRequest)
		return
	}
	errD:= SubscriptionDelete(int32(sID))
	if errD!=nil{
		http.Error(response,errD.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(nil)
} 

func SubscriptionListHandler(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodGet{
		http.Error(response,"Only GET method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	pageParam := request.URL.Query().Get("page")
	page,_ := strconv.Atoi(pageParam)
	if page<1{
		page = 1
	}
	list,err := SubscriptionList(page)
	if err!=nil{
		http.Error(response,err.Error(), http.StatusInternalServerError)
		return
	}
	jsonList,errJ := json.Marshal(list)
	if errJ!=nil{
		http.Error(response,errJ.Error(), http.StatusInternalServerError)
		return
	}
	response.Write(jsonList)
}

func SubscriptionSumHandler(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPost{
		http.Error(response,"Only POST method is allowed! ", http.StatusMethodNotAllowed)
		return
	}
	filterFromParam := request.FormValue("filterFrom")
	filterToParam := request.FormValue("filterTo")
	if filterFromParam == "" || filterToParam == ""{
		http.Error(response,"invalid data", http.StatusBadRequest)
		return
	}
	filterFrom,_ := time.Parse("01-2006",filterFromParam)
	filterTo,_ := time.Parse("01-2006",filterToParam)

	sum,err:=SubscriptionSum(filterFrom,filterTo,nil,nil)
	if err!=nil{
		http.Error(response,err.Error(), http.StatusInternalServerError)
		return
	}
	response.Write([]byte(strconv.Itoa(sum)))
}
