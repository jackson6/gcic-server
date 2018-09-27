package handler

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"github.com/jackson6/gcic-server/app/dao"
	"github.com/jackson6/gcic-server/app/payment"
	"github.com/stripe/stripe-go"
)

func GetUserEndpoint(w http.ResponseWriter, r *http.Request, user *dao.User) {
	defer r.Body.Close()
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: user,
	}
	RespondJSON(w, http.StatusOK, response)
}

func CreateUserEndPoint(mgoDb *mgo.Session, stripeKey string, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var flow dao.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		RespondError(w, http.StatusBadRequest, BadRequest, err)
		return
	}

	plan, err := dao.PlanFindById(mgoDb, "5baad3095b5225373441c0ac"/*flow.User.PlanId*/)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	flow.User.PlanId = plan.ID.Hex()
	charge := &dao.Charge{
		Amount: plan.Amount,
		Currency: stripe.CurrencyJMD,
		Description: plan.Description,
		Source: flow.Token,
	}

	if flow.SaveCard {
		newCustomer, err := payment.CreateCustomer(stripeKey, flow.User.Email, flow.Token)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, InternalError, err)
			return
		}
		charge.Customer = newCustomer
		flow.User.StripeId = newCustomer.ID
	}
	charged, err := payment.ChargeCard(mgoDb, stripeKey, charge)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	if charged != nil {
		memberId, err := dao.GenerateMemberId(mgoDb)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, InternalError, err)
			return
		}
		flow.User.MemberId = memberId
		referralCode, err := dao.GenerateReferralCode(mgoDb)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, InternalError, err)
			return
		}
		flow.User.ReferralCode = referralCode
		if err := dao.UserInsert(mgoDb, &flow.User); err != nil {
			RespondError(w, http.StatusInternalServerError, InternalError, err)
			return
		}
	}
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: flow.User,
	}
	RespondJSON(w, http.StatusOK, response)
}

func GetPlanEndpoint(mgoDb *mgo.Session, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	plan := new(dao.Plan)
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		RespondError(w, http.StatusBadRequest, BadRequest, err)
		return
	}
	plan, err := dao.PlanFindById(mgoDb, plan.ID.Hex())
	if err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: plan,
	}
	RespondJSON(w, http.StatusOK, response)
}

func GetPlansEndpoint(mgoDb *mgo.Session, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	plans, err := dao.PlanFindAll(mgoDb)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: plans,
	}
	RespondJSON(w, http.StatusOK, response)
}

func CreatePlanEndpoint(mgoDb *mgo.Session, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	plan := new(dao.Plan)
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		RespondError(w, http.StatusBadRequest, BadRequest, err)
		return
	}
	if err := dao.PlanInsert(mgoDb, plan); err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: plan,
	}
	RespondJSON(w, http.StatusOK, response)
}

func DeletePlanEndpoint(mgoDb *mgo.Session, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	plan := new(dao.Plan)
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		RespondError(w, http.StatusBadRequest, BadRequest, err)
		return
	}

	if err := dao.PlanDelete(mgoDb, plan); err != nil {
		RespondError(w, http.StatusInternalServerError, InternalError, err)
		return
	}
	response := HttpResponse{
		ResultCode: 200,
		CodeContent: "Success",
		Data: plan,
	}
	RespondJSON(w, http.StatusOK, response)
}
