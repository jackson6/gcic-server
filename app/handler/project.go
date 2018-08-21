package handler

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"../dao"
)

func GetUserEndpoint(w http.ResponseWriter, r *http.Request, user dao.User) {
	defer r.Body.Close()
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: user,
	}
	RespondJSON(w, http.StatusOK, response)
}

func CreateUserEndPoint(mgoDb *mgo.Session, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user dao.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondError(w, http.StatusBadRequest, BadRequest, err)
		return
	}
	user.ID = bson.ObjectId(user.ID)
	if err := dao.UserInsert(mgoDb, user); err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: user,
	}
	RespondJSON(w, http.StatusOK, response)
}
