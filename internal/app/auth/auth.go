package auth

type LoginRequest struct {
	Username string `json:"username" extensions:"x-order=1" example:"john" validate:"required" valid:"required~username|invalid"`         // Username of the user
	Password string `json:"password" extensions:"x-order=2" example:"TopSecret!!!" validate:"required" valid:"required~password|invalid"` // Password of the user
}

type LoginResponse struct {
	Username  string `json:"username" extensions:"x-order=1" example:"john"`
	Email     string `json:"email" extensions:"x-order=2" example:"john@gmail.com"`
	TokenHash string `json:"token_hash" extensions:"x-order=2" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklEIjozLCJVc2VybmFtZSI6Impob24iLCJFbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiUGFzc3dvcmQiOiIkMmEkMTAkRkFUb1ZsS2Y2VmZIRGtYL1dLWmVRT0o2U1kuU3Z0SnNYYmhZV2FlTnBrbjU3S0hlNk4vZTIiLCJEZWZhdWx0Q3VycmVuY3lDb2RlIjoiIiwiY3JlYXRlZF9hdCI6IjIwMjItMTEtMjNUMjI6MjA6MDkuMzk0NzQ3KzAzOjAwIiwidXBkYXRlZF9hdCI6IjIwMjItMTEtMjNUMjI6MjA6MDkuMzk0NzQ3KzAzOjAwIiwiZGVsZXRlZF9hdCI6bnVsbH0sImV4cCI6MTY2OTM4OTM3MH0.b_i6GhYzqOp0VvouVi0rw2VG43UZx7lnJXqNEAKMH8o"`
}
