package subscriptions

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/zakharova-e/subscriptions-info/internal/config"
	"github.com/zakharova-e/subscriptions-info/internal/connections"
)


func SubscriptionCreate(item Subscription) (*int32,error){
	if errValid := item.IsValid(); errValid!=nil{
		return nil,errValid
	}
	query := "insert into subscription (service_name, price,user_id,start_date,finish_date) values ($1,$2,$3,$4,$5) returning id"
	row := connections.PGDatabase.QueryRow(query,item.ServiceName,item.Price,item.UserId,item.StartDate,item.FinishDate)
	if row == nil{
		return nil, errors.New("unexpected error")
	}
	var num int32
	err:= row.Scan(&num)
	return &num, err
} 

func SubscriptionRead(recordId int32) (*Subscription,error){
	if recordId<1{
		return nil,errors.New("invalid record id")
	}
	query := "select id,service_name, price,user_id,start_date,finish_date from subscription where id = $1"
	row := connections.PGDatabase.QueryRow(query,recordId)
	if row == nil{
		return nil, errors.New("unexpected error")
	}
	var (
		id int32
		serviceName,userId string
		price int
		startDate time.Time
		finishDate sql.NullTime
	)
	err :=row.Scan(&id,&serviceName,&price,&userId,&startDate,&finishDate)
	return &Subscription{id,serviceName,price,userId,startDate,finishDate},err
}

func SubscriptionUpdate(item Subscription) error{
	if errValid := item.IsValid(); errValid!=nil{
		return errValid
	}
	if item.Id<1{
		return errors.New("invalid object id")
	}
	query := "update subscription set service_name = $1, price = $2, user_id = $3, start_date = $4, finish_date = $5 where id = $6 "
	_,err := connections.PGDatabase.Exec(query,item.ServiceName,item.Price,item.UserId,item.StartDate,item.FinishDate,item.Id)
	return err
}

func SubscriptionDelete(recordId int32) error{
	if recordId<1{
		return errors.New("invalid record id")
	}
	query := "delete from subscription where id = $1"
	_,err := connections.PGDatabase.Exec(query,recordId)
	return err
} 

func SubscriptionList(page int) (*SubscriptionListPage,error){
	query:= "select id,service_name, price,user_id,start_date,finish_date, COUNT(*) OVER() AS total_count from subscription order by id desc limit $1 offset $2"
	offset := (page - 1)*config.DefaultPageSize
	rows,err := connections.PGDatabase.Query(query,config.DefaultPageSize,offset)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	var list SubscriptionListPage
	list.Page = page
	list.PerPage = config.DefaultPageSize
	for rows.Next() {
		var item Subscription
		err := rows.Scan(&item.Id,&item.ServiceName,&item.Price,&item.UserId,&item.StartDate,&item.FinishDate,&list.Total)
		if err == nil{
			list.List = append(list.List, item)
		}else{
			log.Println(err.Error())
		}
	}
	return &list,nil
}

func SubscriptionSum(filterFrom time.Time, filterTo time.Time, userId *string, serviceName *string) (int, error){
	return 0,nil
}