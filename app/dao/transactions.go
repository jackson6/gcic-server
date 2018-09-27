package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/stripe/stripe-go"
)

const (
	TRANSACTION_COLLECTION = "user"
)

type Charge struct{
	UserId string
	Customer *stripe.Customer
	Amount int64
	Currency stripe.Currency
	Description string
	Source string
}

type Card struct{
	Customer string
	Token string
}

type Transaction struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	UserId string `bson:"user_id" json:"user_id"`
	ChargeID string `bson:"charge_id" json:"charge_id"`
	Amount int64 `bson:"amount" json:"amount"`
	Currency string `bson:"currency" json:"currency"`
	Description string `bson:"description" json:"description"`
	IdempotencyKey string `bson:"idempotency_key" json:"idempotency_key"`
	CreatedOn time.Time `bson:"created_on" json:"created_on"`
}

func TransactionFindAll(db *mgo.Session) ([]Transaction, error) {
	var transactions []Transaction
	err := db.DB(DATABASE).C(TRANSACTION_COLLECTION).Find(bson.M{}).All(&transactions)
	return transactions, err
}

func TransactionFindById(db *mgo.Session, id string) (Transaction, error) {
	var transaction Transaction
	err := db.DB(DATABASE).C(TRANSACTION_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&transaction)
	return transaction, err
}

func TransactionFindByKey(db *mgo.Session, find Transaction) (Transaction, error) {
	var transaction Transaction
	err := db.DB(DATABASE).C(TRANSACTION_COLLECTION).Find(find).One(&transaction)
	return transaction, err
}

func TransactionInsert(db *mgo.Session, transaction Transaction) error {
	err := db.DB(DATABASE).C(TRANSACTION_COLLECTION).Insert(&transaction)
	return err
}

func TransactionDelete(db *mgo.Session, transaction Transaction) error {
	err := db.DB(DATABASE).C(TRANSACTION_COLLECTION).Remove(&transaction)
	return err
}

func TransactionUpdate(db *mgo.Session, transaction Transaction) error {
	err := db.DB(DATABASE).C(TRANSACTION_COLLECTION).UpdateId(transaction.ID, &transaction)
	return err
}