package user

// Req & Resp

type RegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

//
// Response
//

type RegisterResponse struct {
	Username string `json:"username" extensions:"x-order=1" example:"john"`
	Email    string `json:"email" extensions:"x-order=2" example:"john@gmail.com"`
}

type LoginResponse struct {
	RegisterResponse `json:",inline" extensions:"x-order=1"`
	TokenHash        string `json:"token_hash" extensions:"x-order=2" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklEIjozLCJVc2VybmFtZSI6Impob24iLCJFbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiUGFzc3dvcmQiOiIkMmEkMTAkRkFUb1ZsS2Y2VmZIRGtYL1dLWmVRT0o2U1kuU3Z0SnNYYmhZV2FlTnBrbjU3S0hlNk4vZTIiLCJEZWZhdWx0Q3VycmVuY3lDb2RlIjoiIiwiY3JlYXRlZF9hdCI6IjIwMjItMTEtMjNUMjI6MjA6MDkuMzk0NzQ3KzAzOjAwIiwidXBkYXRlZF9hdCI6IjIwMjItMTEtMjNUMjI6MjA6MDkuMzk0NzQ3KzAzOjAwIiwiZGVsZXRlZF9hdCI6bnVsbH0sImV4cCI6MTY2OTM4OTM3MH0.b_i6GhYzqOp0VvouVi0rw2VG43UZx7lnJXqNEAKMH8o"`
}
