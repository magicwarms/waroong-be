package entity

type UpdateUserRequestDTO struct {
	ID        uint64 `validate:"required,number,min=1"`
	FirstName string `validate:"required,alpha,min=2"`
	LastName  string `validate:"required,alpha,min=2"`
	Phone     string `validate:"required,e164,min=11"`
}
