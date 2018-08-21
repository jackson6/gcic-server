package lib

import (
	"encoding/json"
	"../dao"
)

func GetUserStruct(data interface{}) (dao.User, error) {
	var user dao.User
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
