package models

type RegisterUserFirstRequest struct {
    FirstName string `json:"firstname" form:"firstname" query:"firstname" validate:"required,min=1,max=50,persian"`
	LastName string `json:"lastname" form:"lastname" query:"lastname" validate:"required,min=1,max=50,persian"`
	Gender string `json:"gender" form:"gender" query:"gender" validate:"required,oneof=male female unknown"`
}