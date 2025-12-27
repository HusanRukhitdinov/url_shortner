package v1

import (
	"github.com/Go11Group/url_shortner/internal/controller/http/v1/request"
	"github.com/Go11Group/url_shortner/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	uc usecase.UrlUseCaseI
}

func NewController(uc usecase.UrlUseCaseI) *Controller {
	return &Controller{uc: uc}
}

// Shorten creates a short URL.
// @Summary Shorten URL
// @Description Creates a short URL from a long URL
// @Tags url
// @Accept json
// @Produce json
// @Param request body request.ShortenRequest true "URL to shorten"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/shorten [post]
func (c *Controller) Shorten(ctx *fiber.Ctx) error {
	var req request.ShortenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}
	code, err := c.uc.Shorten(ctx.Context(), req.URL)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"short_code": code})
}

// Redirect redirects to the original URL.
// @Summary Redirect URL
// @Description Redirects to the original URL given a short code
// @Tags url
// @Param code path string true "Short Code"
// @Success 302
// @Failure 404 {string} string
// @Router /{code} [get]
func (c *Controller) Redirect(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	original, err := c.uc.GetOriginal(ctx.Context(), code)
	if err != nil {
		return ctx.Status(404).SendString("Not found")
	}
	return ctx.Redirect(original)
}
