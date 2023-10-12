package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maxik12233/blog/common"
	"github.com/maxik12233/blog/services"
	"github.com/maxik12233/blog/types"
)

type ResourceController struct {
	ResourceService services.ResourceService
}

func NewResourceController(resservice services.ResourceService) ResourceController {
	return ResourceController{
		ResourceService: resservice,
	}
}

func (rc *ResourceController) CreateArticle(c *gin.Context) {

	var article *types.Article

	if err := c.ShouldBind(&article); err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	id, err := rc.ResourceService.CreateArticle(article)
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article created",
		"newid":   id,
		"code":    0,
	})

}

func (rc *ResourceController) DeleteArticle(c *gin.Context) {

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	if id <= 0 {
		common.ReturnSimpleError(c, http.StatusBadRequest, errors.New("Invalid id param"))
		return
	}

	err = rc.ResourceService.DeleteArticle(uint(id))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article deleted",
		"code":    0,
	})

}

func (rc *ResourceController) GetArticles(c *gin.Context) {

	amountStr, pageStr := c.Query("amount"), c.Query("page")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
	}
	if page < 0 || amount < 0 {
		common.ReturnSimpleError(c, http.StatusBadRequest, errors.New("Invalid query params"))
	}

	arts, count, err := rc.ResourceService.GetArticles(uint(amount), uint(page))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadGateway, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":   arts,
		"totalCount": count,
		"code":       0,
	})

}

func (rc *ResourceController) GetOneArticle(c *gin.Context) {

	id, err := strconv.Atoi(c.Params.ByName("id"))

	var art *types.Article
	art, err = rc.ResourceService.GetOneArticle(uint(id))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"article": art,
		"code":    0,
	})

}

func (rc *ResourceController) CreateComment(c *gin.Context) {

	var comm *types.Comment
	if err := c.BindJSON(&comm); err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	newid, err := rc.ResourceService.CreateComment(comm)
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment created",
		"newid":   newid,
		"code":    0,
	})
}

func (rc *ResourceController) DeleteComment(c *gin.Context) {

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	if id <= 0 {
		common.ReturnSimpleError(c, http.StatusBadRequest, errors.New("Invalid id param"))
		return
	}

	err = rc.ResourceService.DeleteComment(uint(id))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted",
		"code":    0,
	})

}

func (rc *ResourceController) UpdateComment(c *gin.Context) {

}

func (rc *ResourceController) GetArticleComments(c *gin.Context) {

	id, err := strconv.Atoi(c.Params.ByName("artid"))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	if id <= 0 {
		common.ReturnSimpleError(c, http.StatusBadRequest, errors.New("Invalid id param"))
		return
	}

	comms, err := rc.ResourceService.GetArticleComments(uint(id))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comms": comms,
		"code":  0,
	})
}

func (rc *ResourceController) HandleLike(c *gin.Context) {

	var like *types.Like
	flag, err := strconv.ParseBool(c.Query("flag"))
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.BindJSON(&like); err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	err = rc.ResourceService.ToggleLike(like, flag)
	if err != nil {
		common.ReturnSimpleError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"code":    0,
	})

}

func (rc *ResourceController) RegisterResourceRoutes(rg *gin.RouterGroup) {
	//resgroup := rg.Group("/res")

	/*resgroup.POST("/like", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleCommon),
	}), rc.HandleLike)

	artgroup := resgroup.Group("/art")
	artgroup.POST("/", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleModerator),
	}), rc.CreateArticle)
	artgroup.DELETE("/:id", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleModerator),
	}), rc.DeleteArticle)
	artgroup.GET("/:id", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleCommon),
	}), rc.GetOneArticle)
	artgroup.GET("/", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleCommon),
	}), rc.GetArticles)

	commgroup := resgroup.Group("/comm")
	commgroup.POST("/", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleModerator),
	}), rc.CreateComment)
	commgroup.DELETE("/:id", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleModerator),
	}), rc.DeleteComment)
	commgroup.PUT("/", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleModerator),
	}), rc.UpdateComment)
	commgroup.GET("/:artid", middleware.RequireAuth, middleware.ValidateRoles([]uint{
		uint(types.RoleCommon),
	}), rc.GetArticleComments)*/

}
