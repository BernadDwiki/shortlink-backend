package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BernadDwiki/shortlink-backend/internal/dto"
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

func (lc *LinkController) CreateLink(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateLinkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

func (lc *LinkController) GetLink(ctx *gin.Context) {
	slug := ctx.Param("slug")

	link, err := lc.linkService.GetLinkBySlug(ctx, slug)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "link not found")
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, link.OriginalURL)
}
