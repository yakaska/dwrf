package handlers

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
	"tinyurl/internal/model"
	"tinyurl/internal/shorten"
)

type dwarfer interface {
	Shorten(ctx context.Context, link model.LinkInput) (*model.Link, error)
}

type dwarfyRequest struct {
	URL string `json:"url" validate:"required,url"`
	Id  string `json:"id,omitempty" validate:"omitempty,alphanum"`
}

type dwarfyResponse struct {
	ShortURL string `json:"short_url,omitempty"`
	Message  string `json:"message,omitempty"`
}

func HandleDwarfy(dwarfer dwarfer, baseUrl string) func(c echo.Context) error {
	return func(c echo.Context) error {
		var request dwarfyRequest

		if err := c.Bind(&request); err != nil {
			slog.Error(fmt.Sprintf("error binding request: %v", err))
			return err
		}

		if err := c.Validate(request); err != nil {
			slog.Error(fmt.Sprintf("validation error: %v", err))
			return err
		}

		id := strings.TrimSpace(request.Id)

		inputLink := model.LinkInput{
			URL:       request.URL,
			Id:        id,
			CreatedBy: "yakaska",
		}

		shortLink, err := dwarfer.Shorten(c.Request().Context(), inputLink)
		if err != nil {
			slog.Error(fmt.Sprintf("error shortening url %q: %v", request.URL, err))
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		resultLink, err := shorten.PrependBaseUrl(baseUrl, shortLink.Short)

		return c.JSON(http.StatusOK, dwarfyResponse{
			ShortURL: resultLink,
		})

	}
}
