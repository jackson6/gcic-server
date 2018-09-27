package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	PLAN_COLLECTION = "plan"
)

type Plan struct {
	ID bson.ObjectId `bson:"_id"json:"id"`
	Name string `bson:"name"json:"name"`
	Description string `bson:"description"json:"description"`
	Amount int64 `bson:"amount"json:"amount"`
}

func PlanFindAll(db *mgo.Session) ([]*Plan, error) {
	plans := make([]*Plan, 0)
	err := db.DB(DATABASE).C(PLAN_COLLECTION).Find(bson.M{}).All(&plans)
	return plans, err
}

func PlanFindById(db *mgo.Session, id string) (*Plan, error) {
	plan := new(Plan)
	err := db.DB(DATABASE).C(PLAN_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&plan)
	return plan, err
}

func PlanFindByKey(db *mgo.Session, find *Plan) (*Plan, error) {
	plan := new(Plan)
	err := db.DB(DATABASE).C(PLAN_COLLECTION).Find(find).One(&plan)
	return plan, err
}

func PlanInsert(db *mgo.Session, plan *Plan) error {
	plan.ID = bson.NewObjectId()
	err := db.DB(DATABASE).C(PLAN_COLLECTION).Insert(&plan)
	return err
}

func PlanDelete(db *mgo.Session, plan *Plan) error {
	err := db.DB(DATABASE).C(PLAN_COLLECTION).Remove(&plan)
	return err
}

func PlanUpdate(db *mgo.Session, plan *Plan) error {
	err := db.DB(DATABASE).C(PLAN_COLLECTION).UpdateId(plan.ID, &plan)
	return err
}