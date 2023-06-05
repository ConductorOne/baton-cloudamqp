package cloudamqp

type BaseResource struct {
	Id string `json:"id"`
}

type User struct {
	BaseResource
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
