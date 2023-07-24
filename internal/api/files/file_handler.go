package files

import (
	"net/http"
	"strings"

	"github.com/Brix101/network-file-manager/internal/utils"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	Reader Reader
}

func (h FileHandler) Routes(v1 *echo.Group) {
	user := v1.Group("/files")
	user.GET("*", h.list)
}

func (h FileHandler) list(c echo.Context) error {
	path := c.Request().URL.Path
	path = strings.TrimPrefix(path, "/api/files")

	hiddenParam := c.QueryParam("hidden")
	hidden := true
	if hiddenParam == "false" {
		hidden = false
	}

	files, err := h.Reader.GetContent(path, hidden)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, files)
}
