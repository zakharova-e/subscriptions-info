package subscriptions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id          int32        `json:"" db:"id"`
	ServiceName string       `json:"" db:"service_name"`
	Price       int          `json:"" db:"price"`
	UserId      string       `json:"" db:"user_id"`
	StartDate   time.Time    `json:"" db:"start_date"`
	FinishDate  sql.NullTime `json:"" db:"finish_date"`
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
	return json.Marshal(struct {
		Id          int32   `json:"id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserId      string  `json:"user_id"`
		StartDate   string  `json:"start_date"`
		FinishDate  *string `json:"finish_date,omitempty"`
	}{s.Id, s.ServiceName, s.Price, s.UserId, s.StartDate.Format("01-2006"), finishDate})
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
		return err
	}
	s.Id = temp.Id
	s.ServiceName = temp.ServiceName
	s.Price = temp.Price
	s.UserId = temp.UserId
	s.StartDate, _ = time.Parse("01-2006", temp.StartDate)
	if temp.FinishDate != "" {
		finishDate, _ := time.Parse("01-2006", temp.FinishDate)
		s.FinishDate = sql.NullTime{Valid: true, Time: finishDate}
	}
	return nil
}

func (s Subscription) IsValid() error {
	if len(s.ServiceName) == 0 {
		return errors.New("invalid service name")
	}
	if err := uuid.Validate(s.UserId); err != nil {
		return errors.New("invalid user id")
	}
	if s.StartDate.Year() < 2020 {
		return errors.New("invalid start date")
	}
	if s.FinishDate.Valid && s.FinishDate.Time.Year() < 2020 {
		return errors.New("invalid start date")
	}
	return nil
}
