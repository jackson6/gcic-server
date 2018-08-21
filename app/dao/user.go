package dao

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	USER_COLLECTION = "user"
	DATABASE = "invest"
)

type User struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName string `bson:"last_name" json:"last_name"`
	Email string `bson:"email" json:"email"`
	ReferralCode string `bson:"referral_code" json:"referral_code"`
}

func UserFindAll(db *mgo.Session) ([]User, error) {
	var users []User
	err := db.DB(DATABASE).C(USER_COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func UserFindById(db *mgo.Session, id string) (User, error) {
	var user User
	err := db.DB(DATABASE).C(USER_COLLECTION).FindId(bson.ObjectId(id)).One(&user)
	return user, err
}

func UserFindByKey(db *mgo.Session, find User) (User, error) {
	var user User
	err := db.DB(DATABASE).C(USER_COLLECTION).Find(find).One(&user)
	return user, err
}

func UserInsert(db *mgo.Session, user User) error {
	err := db.DB(DATABASE).C(USER_COLLECTION).Insert(&user)
	return err
}

func UserDelete(db *mgo.Session, user User) error {
	err := db.DB(DATABASE).C(USER_COLLECTION).Remove(&user)
	return err
}

func UserUpdate(db *mgo.Session, user User) error {
	err := db.DB(DATABASE).C(USER_COLLECTION).UpdateId(user.ID, &user)
	return err
}