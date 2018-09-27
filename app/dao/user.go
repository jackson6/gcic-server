package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/jackson6/gcic-server/app/lib"
	"encoding/json"
)

const (
	USER_COLLECTION = "user"
	DATABASE = "invest"
)

type User struct {
	ID bson.ObjectId `bson:"_id"json:"id"`
	UserId string `bson:"user_id"json:"user_id"`
	MemberId string `bson:"member_id"json:"member_id"`
	FirstName string `bson:"first_name"json:"first_name"`
	LastName string `bson:"last_name"json:"last_name"`
	Initial string `bson:"initial"json:"initial"`
	Email string `bson:"email"json:"email"`
	Trn string `bson:"trn"json:"trn"`

	WorkPhone string `bson:"work_phone"json:"work_phone"`
	HomePhone string `bson:"home_phone"json:"home_phone"`
	CellPhone string `bson:"cell_phone"json:"cell_phone"`

	SponsorId string `bson:"sponsor_id"json:"sponsor_id"`

	Address string `bson:"address"json:"address"`
	Parish string `bson:"parish"json:"parish"`
	Country string `bson:"country"json:"country"`

	Question string `bson:"question"json:"question"`
	Answer string `bson:"answer"json:"answer"`

	Dob time.Time `bson:"dob"json:"dob"`
	Gender string `bson:"gender"json:"gender"`
	StripeId string `bson:"stripe_id"json:"stripe_id"`

	PlanId string `bson:"plan_id"json:"plan_id"`

	ReferralCode string `bson:"referral_code"json:"referral_code"`
}

func GetUserStruct(data interface{}) (*User, error) {
	user := new(User)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(jsonStr, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GenerateReferralCode(db *mgo.Session)(string, error){
	code := lib.RandSeq(6)
	_, err := UserFindByKey(db, &bson.M{"referral_code": code})
	if err != nil && err.Error() != "not found" {
		return code, err
	}
	return code, nil
}

func GenerateMemberId(db *mgo.Session)(string, error){
	id := lib.RandSeq(6)
	_, err := UserFindByKey(db, &bson.M{"member_id": id})
	if err != nil && err.Error() != "not found" {
		return id, err
	}
	return id, nil
}

func UserFindAll(db *mgo.Session) ([]*User, error) {
	users := make([]*User, 0)
	err := db.DB(DATABASE).C(USER_COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func UserFindById(db *mgo.Session, id string) (*User, error) {
	user := new(User)
	err := db.DB(DATABASE).C(USER_COLLECTION).FindId(bson.ObjectId(id)).One(&user)
	return user, err
}

func UserFindByKey(db *mgo.Session, find *bson.M) (*User, error) {
	user := new(User)
	err := db.DB(DATABASE).C(USER_COLLECTION).Find(find).One(&user)
	return user, err
}

func UserInsert(db *mgo.Session, user *User) error {
	user.ID = bson.NewObjectId()
	err := db.DB(DATABASE).C(USER_COLLECTION).Insert(&user)
	return err
}

func UserDelete(db *mgo.Session, user *User) error {
	err := db.DB(DATABASE).C(USER_COLLECTION).Remove(&user)
	return err
}

func UserUpdate(db *mgo.Session, user *User) error {
	err := db.DB(DATABASE).C(USER_COLLECTION).UpdateId(user.ID, &user)
	return err
}