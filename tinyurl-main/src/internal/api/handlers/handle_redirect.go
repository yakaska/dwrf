package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"tinyurl/internal/model"
)

type redirecter interface {
	Redirect(ctx context.Context, id string) (string, error)
}

func HandleRedirect(redirecter redirecter) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		redirectUrl, err := redirecter.Redirect(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			slog.Error(fmt.Sprintf("error redirecting for %q:  %v", id, err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Redirect(http.StatusMovedPermanently, redirectUrl)
	}
}
