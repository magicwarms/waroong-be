package entity

type UserRequestDTO struct {
	Email      string `validate:"required,email,min=4,max=64"`
	Password   string `validate:"required,alphanum,min=8"`
	FirstName  string `validate:"required,alpha,min=2"`
	LastName   string `validate:"required,alpha,min=2"`
	Phone      string `validate:"required,e164,min=11"`
	UserTypeId string `validate:"required,number,min=1"`
}

type UserUpdateRequestDTO struct {
	ID       string `validate:"required,number,min=1"`
	Password string `validate:"required,alphanum,min=8"`
}

type UserGetIDRequestDTO struct {
	ID string `validate:"required,number,min=1"`
}

type CustomerRequestDTO struct {
	Email     string `validate:"required,email,min=4,max=64"`
	Password  string `validate:"required,alphanum,min=8"`
	FirstName string `validate:"required,alpha,min=2"`
	LastName  string `validate:"required,alpha,min=2"`
}

type CustomerUpdateRequestDTO struct {
	ID       string `validate:"required,number,min=1"`
	Password string `validate:"required,alphanum,min=8"`
}

type CustomerGetIDRequestDTO struct {
	ID uint `validate:"required,number,min=1"`
}

type UserLoginRequestDTO struct {
	Email    string `validate:"required,email,min=4,max=64"`
	Password string `validate:"required,alphanum,min=8"`
}

type ChangePasswordUserDTO struct {
	UserID   string `validate:"required,number,min=1"`
	Password string `validate:"required,alphanum,min=8"`
}

type ForgotPasswordRequestDTO struct {
	Email string `validate:"required,email,min=4,max=64"`
}
