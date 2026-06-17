package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BernadDwiki/shortlink-backend/internal/dto"
	"github.com/BernadDwiki/shortlink-backend/internal/helper"
	"github.com/BernadDwiki/shortlink-backend/internal/response"
	"github.com/BernadDwiki/shortlink-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkController struct {
	linkService *service.LinkService
}

func NewLinkController(linkService *service.LinkService) *LinkController {
	return &LinkController{linkService: linkService}
}

// CreateLink godoc
// @Summary Create a short link
// @Description Create a short link for the authenticated user
// @Tags Links
// @Accept json
// @Produce json
// @Param request body dto.CreateLinkRequest true "Create link request"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security ApiKeyAuth
// @Router /api/links [post]
func (lc *LinkController) CreateLink(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateLinkRequest
	if err := helper.BindJSON(ctx, &req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	link, err := lc.linkService.CreateLink(ctx, userID, req)
	if err != nil {
		if errors.Is(err, service.ErrSlugAlreadyTaken) {
			response.Error(ctx, http.StatusBadRequest, "slug already taken")
			return
		}
		if errors.Is(err, service.ErrReservedSlug) {
			response.Error(ctx, http.StatusBadRequest, "slug is a reserved word")
			return
		}
		if errors.Is(err, service.ErrInvalidSlug) {
			response.Error(ctx, http.StatusBadRequest, "slug must be 3-50 characters, alphanumeric and hyphens only")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to create link")
		return
	}

	response.Success(ctx, http.StatusCreated, "Link created successfully", gin.H{
		"id":           link.ID,
		"original_url": link.Original,
		"slug":         link.Slug,
		"short_url":    link.ShortURL,
	})
}

// GetUserLinks godoc
// @Summary Get authenticated user's links
// @Description Retrieve all active short links created by the authenticated user
// @Tags Links
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security ApiKeyAuth
// @Router /api/links [get]
func (lc *LinkController) GetUserLinks(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	links, err := lc.linkService.GetUserLinks(ctx, userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to retrieve links")
		return
	}

	response.Success(ctx, http.StatusOK, "Links retrieved successfully", links)
}

// DeleteLink godoc
// @Summary Delete a link
// @Description Soft delete a link owned by the authenticated user
// @Tags Links
// @Accept json
// @Produce json
// @Param id path int true "Link ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Security ApiKeyAuth
// @Router /api/links/{id} [delete]
func (lc *LinkController) DeleteLink(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	linkIDStr := ctx.Param("id")
	linkID, err := strconv.Atoi(linkIDStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid link id")
		return
	}

	if err := lc.linkService.DeleteLink(ctx, linkID, userID); err != nil {
		if err.Error() == "link not found or unauthorized" {
			response.Error(ctx, http.StatusNotFound, "link not found or unauthorized")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed to delete link")
		return
	}

	response.Success(ctx, http.StatusOK, "Link deleted successfully", nil)
}

// GetLink godoc
// @Summary Redirect to original URL
// @Description Public endpoint to redirect a short slug to the original URL
// @Tags Links
// @Accept json
// @Produce json
// @Param slug path string true "Link slug"
// @Success 301 "redirect"
// @Failure 404 {object} response.APIResponse
// @Router /{slug} [get]
func (lc *LinkController) GetLink(ctx *gin.Context) {
	slug := ctx.Param("slug")

	link, err := lc.linkService.GetLinkBySlug(ctx, slug)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, link.OriginalURL)
}
