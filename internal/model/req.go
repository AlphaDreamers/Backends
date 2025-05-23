package model

/*
	{
	  "firstName": "Alice",
	  "lastName": "Smith",
	  "email": "alice.smith@example.com",
	  "password": "MySecurePass123",
	  "country": "US",
	  "biometricHash": "hashvalue12345678"
	}
*/
type UserSignUpRequest struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Country       string `json:"country"`
	BioMetricHash string `json:"bioMetricHash"`
}

/*
{
  "email": "test.user@example.com",
  "password": "TestPass123!"
}
*/

type UserSignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailVerificationRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type UserLogoutRequest struct {
	Email string `json:"email"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}
type ForgotPasswordConfirmReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}
