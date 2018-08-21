package payment

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
	"../dao"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/card"
)

type Charge struct{
	UserID string
	Amount int64
	Currency stripe.Currency
	Description string
	Source string
}

type Card struct{
	Customer string
	Token string
}

func createCard(stripeKey string, info Card)(*stripe.Card, error){
	stripe.Key = stripeKey

	params := &stripe.CardParams{
		Customer: stripe.String(info.Customer),
		Token: stripe.String(info.Token),
	}
	c, err := card.New(params)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getAllCards(stripeKey, customerId string)([]*stripe.Card){
	var cards []*stripe.Card

	stripe.Key = stripeKey

	params := &stripe.CardListParams{
		Customer: stripe.String(customerId),
	}
	params.Filters.AddFilter("limit", "", "10")
	i := card.List(params)
	for i.Next() {
		c := i.Card()
		cards = append(cards, c)
	}
	return cards
}

func getCard(stripeKey, customerId, cardId string)(*stripe.Card, error){
	stripe.Key = stripeKey

	params := &stripe.CardParams{
		Customer: &customerId,
	}
	c, err := card.Get(cardId, params)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func updateCard(stripeKey, cardId string, info Card)(*stripe.Card, error){
	stripe.Key = stripeKey

	params := &stripe.CardParams{
		Customer: stripe.String(info.Customer),
		Token: stripe.String(info.Token),
	}
	c, err := card.Update(cardId, params)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func deleteCard(stripeKey, cardId string, info Card)(*stripe.Card, error){
	stripe.Key = stripeKey

	params := &stripe.CardParams{
		Customer: stripe.String(info.Customer),
	}
	c, err := card.Del(cardId, params)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func chargeCard(db *mgo.Session, stripeKey string, info Charge) error{
	stripe.Key = stripeKey

	args := &stripe.ChargeParams{
		Amount: stripe.Int64(info.Amount),
		Currency: stripe.String(string(info.Currency)),
		Description: stripe.String(info.Description),
	}
	args.SetSource(info.Source) // obtained with Stripe.js

	u2, err := uuid.NewV4()
	if err != nil {
		return err
	}
	args.SetIdempotencyKey(u2.String())

	ch, err := charge.New(args)
	if err != nil {
		return err
	}
	transaction := dao.Transaction{
		ID:  bson.NewObjectId(),
		UserID: info.UserID,
		ChargeID: ch.ID,
		Amount: ch.Amount,
		Currency: string(ch.Currency),
		Description: ch.Description,
		CreatedOn: time.Unix(ch.Created, 0),
	}
	err = dao.TransactionInsert(db, transaction)
	if err != nil {
		return err
	}
	return nil
}

func getCharge(stripeKey, chargeId string) (*stripe.Charge, error){
	stripe.Key = stripeKey
	c, err := charge.Get(chargeId, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func updateCharge(stripeKey, chargeId string, params *stripe.ChargeParams)(*stripe.Charge, error){
	stripe.Key = stripeKey

	ch, err := charge.Update(chargeId, params)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func allCharges(stripeKey string) []*stripe.Charge{
	var charges []*stripe.Charge

	stripe.Key = stripeKey

	params := &stripe.ChargeListParams{}
	params.Filters.AddFilter("limit", "", "3")
	i := charge.List(params)
	for i.Next() {
		c := i.Charge()
		charges = append(charges, c)
	}
	return charges
}

func getCustmers(stripeKey string)([]*stripe.Customer){
	var customers []*stripe.Customer

	stripe.Key = stripeKey

	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("limit", "", "10")
	i := customer.List(params)
	for i.Next() {
		customer := i.Customer()
		customers = append(customers, customer)
	}
	return customers
}