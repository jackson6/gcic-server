package dao

type CreateUser struct {
	User User `json:"user"`
	Token string `json:"token"`
	SaveCard bool `json:"save_card"`
}