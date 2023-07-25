package users

import (
	"net/http"

	"github.com/Brix101/network-file-manager/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserServices *UserServices
}

func (h UserHandler) Routes(v1 *echo.Group) {
	user := v1.Group("/users")
	user.GET("", h.list)
	user.POST("", h.signUp)
	user.POST("/sign-in", h.signIn)
}

func (uh UserHandler) list(c echo.Context) error {
	users, err := uh.UserServices.GetAll()
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, users)
}

func (h UserHandler) signUp(c echo.Context) error {
	var u User
	req := &userRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	// create the user
	// u, err := uh.UserServices.CreateUser(User{
	// 	Name:     req.User.Name,
	// 	Email:    req.User.Email,
	// 	Password: req.User.Password,
	// })
	// if err != nil {
	// 	return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	// }
	return c.JSON(http.StatusOK, u)
}

func (h UserHandler) signIn(c echo.Context) error {
	req := &userSignInRequest{}

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	u, err := h.UserServices.GetByEmail(req.User.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}

	if !u.CheckPassword(req.User.Password) {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	return c.JSON(http.StatusOK, newUserResponse(u))
}
