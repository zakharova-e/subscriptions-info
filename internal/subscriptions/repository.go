package subscriptions

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zakharova-e/subscriptions-info/internal/config"
	"github.com/zakharova-e/subscriptions-info/internal/connections"
)

func SubscriptionCreate(item Subscription) (*int32, error) {
	if errValid := item.IsValid(); errValid != nil {
		return nil, errValid
	}
	query := "INSERT INTO subscription (service_name, price,user_id,start_date,finish_date) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	row := connections.PGDatabase.QueryRow(query, item.ServiceName, item.Price, item.UserId, item.StartDate, item.FinishDate)
	var num int32
	err := row.Scan(&num)
	if err != nil {
		return nil, &DatabaseError{Err: err}
	}
	return &num, nil
}

func SubscriptionRead(recordId int32) (*Subscription, error) {
	if recordId < 1 {
		return nil, &InvalidParameterError{ParamName: "recordId"}
	}
	query := "SELECT id,service_name, price,user_id,start_date,finish_date FROM subscription WHERE id = $1"
	row := connections.PGDatabase.QueryRow(query, recordId)
	var (
		id                  int32
		serviceName, userId string
		price               int
		startDate           time.Time
		finishDate          sql.NullTime
	)
	err := row.Scan(&id, &serviceName, &price, &userId, &startDate, &finishDate)
	if err != nil {
		return nil, &DatabaseError{Err: err}
	}
	return &Subscription{id, serviceName, price, userId, startDate, finishDate}, nil
}

func SubscriptionUpdate(item Subscription) error {
	if errValid := item.IsValid(); errValid != nil {
		return errValid
	}
	if item.Id < 1 {
		return &ValidationError{Errors: []error{errors.New("invalid item id")}}
	}
	query := "UPDATE subscription SET service_name = $1, price = $2, user_id = $3, start_date = $4, finish_date = $5 WHERE id = $6 "
	_, err := connections.PGDatabase.Exec(query, item.ServiceName, item.Price, item.UserId, item.StartDate, item.FinishDate, item.Id)
	if err != nil {
		return &DatabaseError{Err: err}
	}
	return nil
}

func SubscriptionDelete(recordId int32) error {
	if recordId < 1 {
		return &InvalidParameterError{ParamName: "recordId"}
	}
	query := "DELETE FROM subscription WHERE id = $1"
	_, err := connections.PGDatabase.Exec(query, recordId)
	if err != nil {
		return &DatabaseError{Err: err}
	}
	return nil
}

func SubscriptionList(page int) (*SubscriptionListPage, error) {
	query := "SELECT id,service_name, price,user_id,start_date,finish_date, COUNT(*) OVER() AS total_count FROM subscription ORDER BY id DESC LIMIT $1 OFFSET $2"
	offset := (page - 1) * config.DefaultPageSize
	rows, errQuery := connections.PGDatabase.Query(query, config.DefaultPageSize, offset)
	if errQuery != nil {
		return nil, &DatabaseError{Err: errQuery}
	}
	defer rows.Close()
	var list SubscriptionListPage
	list.Page = page
	list.PerPage = config.DefaultPageSize
	for rows.Next() {
		var item Subscription
		err := rows.Scan(&item.Id, &item.ServiceName, &item.Price, &item.UserId, &item.StartDate, &item.FinishDate, &list.Total)
		if err == nil {
			list.List = append(list.List, item)
		}
	}
	return &list, nil
}

func SubscriptionSum(filterFrom time.Time, filterTo time.Time, userId *string, serviceName *string) (int, error) {
	params := []any{}
	//calculation formula:
	// months count: (yearTo-yearFrom)*12 + (monthTo - monthFrom) + 1
	query := `SELECT SUM(price * (
    (EXTRACT(YEAR FROM sumTo) * 12 + EXTRACT(MONTH FROM sumTo))
		- (EXTRACT(YEAR FROM sumFrom) * 12 + EXTRACT(MONTH FROM sumFrom)) + 1
	)) AS total
	FROM (
	SELECT
		price,
		GREATEST(start_date, $1) AS sumFrom,
		LEAST(COALESCE(finish_date, $2), $2) AS sumTo
	FROM subscription
	WHERE start_date <= $2
		AND (finish_date >= $1 OR finish_date IS NULL)
	`

	params = append(params, filterFrom.Format("2006-01-02"), filterTo.Format("2006-01-02"))
	paramNum := 3
	if userId != nil {
		query = query + fmt.Sprintf("AND user_id = $%d ", paramNum)
		params = append(params, userId)
		paramNum++
	}
	if serviceName != nil {
		query = query + fmt.Sprintf("AND service_name = $%d ", paramNum)
		params = append(params, serviceName)
		paramNum++
	}
	//more params?
	query = query + ") sub;"
	log.Printf("query to execute: %s", query)
	row := connections.PGDatabase.QueryRow(query, params...)
	var res int
	err := row.Scan(&res)
	if err != nil {
		return 0, &DatabaseError{Err: err}
	}
	return res, nil
}
