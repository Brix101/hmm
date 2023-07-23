package users

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type userRegisterRequest struct {
	User struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++", err.Error())
		return err
	}
	u.Name = r.User.Name
	u.Email = r.User.Email
	// h, err := u.HashPassword(r.User.Password)
	// if err != nil {
	// 	return err
	// }
	u.Password = r.User.Password
	return nil
}

type userSignInRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userSignInRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type userResponse struct {
	User struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

func newUserResponse(u *User) *userResponse {
	r := new(userResponse)
	r.User.Id = u.Id.String()
	r.User.Name = u.Name
	r.User.Email = u.Email
	// r.User.Bio = u.Bio
	// r.User.Image = u.Image
	// r.User.Token = utils.GenerateJWT(u.ID)
	return r
}
