package models

import (
	"fmt"
	"reflect"
	"time"
)

// PenaltyList holds an array of penalty payment (payable resources)
type PenaltyList struct {
	Penalties []PayableResourceDao
}

// TransactionDao is a transaction that is persisted to the db
type TransactionDao struct {
	Amount     float64 `bson:"amount"`
	Type       string  `bson:"type"`
	MadeUpDate string  `bson:"made_up_date"`
	Reason     string  `bson:"reason"`
}

// PayableResourceDao is the persisted resource for payable items
type PayableResourceDao struct {
	ID             string                 `bson:"_id"`
	E5CommandError string                 `bson:"e5_command_error"`
	CustomerCode   string                 `bson:"customer_code"`
	PayableRef     string                 `bson:"payable_ref"`
	Data           PayableResourceDataDao `bson:"data"`
}

// PayableResourceDataDao contains data of a penalty payment (payable resource)
type PayableResourceDataDao struct {
	Transactions map[string]TransactionDao `bson:"transactions"`
	Etag         string                    `bson:"etag"`
	Payment      PaymentDao                `bson:"payment"`
	CreatedAt    *time.Time                `bson:"created_at"`
	CreatedBy    CreatedByDao              `bson:"created_by"`
	Links        PayableResourceLinksDao   `bson:"links"`
}

// PaymentDao is the payment information for a payable resource
type PaymentDao struct {
	Amount    string     `bson:"amount"`
	PaidAt    *time.Time `bson:"paid_at"`
	Reference string     `bson:"reference"`
	Status    string     `bson:"status"`
}

// PayableResourceLinksDao is the links object of the payable resource
type PayableResourceLinksDao struct {
	Self          string `bson:"self"`
	Payment       string `bson:"payment"`
	ResumeJourney string `bson:"resume_journey_uri"`
}

// CreatedByDao is the object relating to who created the resource
type CreatedByDao struct {
	Email    string `bson:"email"`
	Forename string `bson:"forename"`
	ID       string `bson:"id"`
	Surname  string `bson:"surname"`
}

// PenaltyErrorDataList contains the list of data that will be passed into the csv
type PenaltyErrorDataList struct {
	Penalties []PenaltyErrorData
}

// PenaltyErrorData is the formatted data required in the CSV
type PenaltyErrorData struct {
	CreatedAt    *time.Time
	PayableRef   string
	CustomerCode string
	PenaltyRef   string
	MadeUpDate   string
	Amount       float64
	Reason       string
}

// ToCSV converts PenaltyErrorDataList into CSV-writable data
func (penalties PenaltyErrorDataList) ToCSV() [][]string {

	csv := make([][]string, len(penalties.Penalties)+1)

	for i := 0; i < len(penalties.Penalties); i++ {
		if i == 0 {
			csv[i] = getHeaders(penalties.Penalties[i])
		}
		csv[i+1] = getSlice(penalties.Penalties[i])
	}

	return csv
}

func getSlice(resource interface{}) []string {

	val := reflect.ValueOf(resource)

	slice := make([]string, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		slice[i] = fmt.Sprintf("%v", val.Field(i))
	}

	return slice
}

func getHeaders(resource interface{}) []string {

	val := reflect.ValueOf(resource)

	headers := make([]string, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		headers[i] = val.Type().Field(i).Name
	}

	return headers
}
