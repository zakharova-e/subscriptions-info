package subscriptions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id          int32        `json:"id"`
	ServiceName string       `json:"service_name"`
	Price       int          `json:"price"`
	UserId      string       `json:"user_id"`
	StartDate   time.Time    `json:"start_date"`
	FinishDate  sql.NullTime `json:"finish_date" swaggertype:"string,nullable"`
}

type SubscriptionListPage struct {
	List    []Subscription
	Page    int
	PerPage int
	Total   int
}

// override json marshaling
func (s *Subscription) MarshalJSON() ([]byte, error) {
	var finishDate *string
	finishDate = nil
	if s.FinishDate.Valid {
		finishDateFormatted := s.FinishDate.Time.Format("01-2006")
		finishDate = &finishDateFormatted
	}
	res, err := json.Marshal(struct {
		Id          int32   `json:"id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserId      string  `json:"user_id"`
		StartDate   string  `json:"start_date"`
		FinishDate  *string `json:"finish_date,omitempty"`
	}{s.Id, s.ServiceName, s.Price, s.UserId, s.StartDate.Format("01-2006"), finishDate})
	if err != nil {
		err = &JsonError{Err: err}
	}
	return res, err
}

// override json unmarshaling
func (s *Subscription) UnmarshalJSON(body []byte) error {
	var temp struct {
		Id          int32  `json:"id"`
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		UserId      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		FinishDate  string `json:"finish_date"`
	}
	if err := json.Unmarshal(body, &temp); err != nil {
		return &JsonError{Err: err, Json: string(body)}
	}
	s.Id = temp.Id
	s.ServiceName = temp.ServiceName
	s.Price = temp.Price
	s.UserId = temp.UserId
	s.StartDate, _ = time.Parse("01-2006", temp.StartDate)
	if temp.FinishDate != "" {
		finishDate, errParse := time.Parse("01-2006", temp.FinishDate)
		if errParse == nil {
			finishDate = finishDate.AddDate(0, 1, -1)
		}
		s.FinishDate = sql.NullTime{Valid: true, Time: finishDate}
	}
	return nil
}

func (s Subscription) IsValid() error {
	var vErr ValidationError
	if len(s.ServiceName) == 0 {
		vErr.Errors = append(vErr.Errors, errors.New("service name is empty"))
	}
	if err := uuid.Validate(s.UserId); err != nil {
		vErr.Errors = append(vErr.Errors, errors.New("user id must be uuid"))
	}
	if s.StartDate.Year() < 2020 {
		vErr.Errors = append(vErr.Errors, errors.New("incorrect date, years after 2020 accepted"))
	}
	if s.FinishDate.Valid && s.FinishDate.Time.Year() < 2020 {
		vErr.Errors = append(vErr.Errors, errors.New("incorrect date, years after 2020 accepted"))
	}
	if len(vErr.Errors) == 0 {
		return nil
	}
	return &vErr
}
