package payment

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2"
	"github.com/jackson6/gcic-server/app/dao"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/card"
)

func CreateCard(stripeKey string, info dao.Card)(*stripe.Card, error){
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

func GetAllCards(stripeKey, customerId string)([]*stripe.Card){
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

func GetCard(stripeKey, customerId, cardId string)(*stripe.Card, error){
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

func UpdateCard(stripeKey, cardId string, info dao.Card)(*stripe.Card, error){
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

func deleteCard(stripeKey, cardId string, info dao.Card)(*stripe.Card, error){
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

func ChargeCard(db *mgo.Session, stripeKey string, info *dao.Charge) (*stripe.Charge, error){
	stripe.Key = stripeKey

	args := &stripe.ChargeParams{
		Amount: stripe.Int64(info.Amount),
		Currency: stripe.String(string(info.Currency)),
		Description: stripe.String(info.Description),
	}
	if info.Customer != nil {
		args.Customer = stripe.String(info.Customer.ID)
	} else {
		args.SetSource(info.Source) // obtained with Stripe.js
	}

	u2, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	args.SetIdempotencyKey(u2.String())

	ch, err := charge.New(args)
	if err != nil {
		return nil, err
	}
	transaction := dao.Transaction{
		ID:  bson.NewObjectId(),
		UserId: info.UserId,
		ChargeID: ch.ID,
		Amount: ch.Amount,
		Currency: string(ch.Currency),
		Description: ch.Description,
		IdempotencyKey: u2.String(),
		CreatedOn: time.Unix(ch.Created, 0),
	}
	err = dao.TransactionInsert(db, transaction)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func GetCharge(stripeKey, chargeId string) (*stripe.Charge, error){
	stripe.Key = stripeKey
	c, err := charge.Get(chargeId, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func UpdateCharge(stripeKey, chargeId string, params *stripe.ChargeParams)(*stripe.Charge, error){
	stripe.Key = stripeKey

	ch, err := charge.Update(chargeId, params)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func AllCharges(stripeKey string) []*stripe.Charge{
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

func GetCustomers(stripeKey string)([]*stripe.Customer){
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