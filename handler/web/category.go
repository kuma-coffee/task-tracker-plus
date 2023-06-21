package web

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"embed"
	"net/http"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
)

type CategoryWeb interface {
	Category(c *gin.Context)
	CategoryAddProcess(c *gin.Context)
	CategoryUpdatePage(c *gin.Context)
	CategoryUpdateProcess(c *gin.Context)
	CategoryDeleteProcess(c *gin.Context)
}

type categoryWeb struct {
	categoryClient client.CategoryClient
	sessionService service.SessionService
	embed          embed.FS
}

func NewCategoryWeb(categoryClient client.CategoryClient, sessionService service.SessionService, embed embed.FS) *categoryWeb {
	return &categoryWeb{categoryClient, sessionService, embed}
}

func (c *categoryWeb) Category(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	categories, err := c.categoryClient.CategoryList(session.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var dataTemplate = map[string]interface{}{
		"email":      email,
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
	}

	var header = path.Join("views", "general", "header.html")
	var filepath = path.Join("views", "main", "category.html")

	t, err := template.New("category.html").Funcs(funcMap).ParseFS(c.embed, filepath, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	err = t.Execute(ctx.Writer, dataTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}
}

func (c *categoryWeb) CategoryAddProcess(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	name := ctx.Request.FormValue("name")

	status, err := c.categoryClient.AddCategory(session.Token, name)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		ctx.Redirect(http.StatusSeeOther, "/client/category")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Add Category Failed!")
	}
}

func (c *categoryWeb) CategoryUpdatePage(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	categoryID := ctx.Param("id")

	category, err := c.categoryClient.CategoryByID(categoryID, session.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var dataTemplate = map[string]interface{}{
		"email":    email,
		"category": category,
	}

	var funcMap = template.FuncMap{
		"exampleFunc": func() int {
			return 0
		},
	}

	var header = path.Join("views", "general", "header.html")
	var filepath = path.Join("views", "main", "category_update.html")

	t, err := template.New("category_update.html").Funcs(funcMap).ParseFS(c.embed, filepath, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	err = t.Execute(ctx.Writer, dataTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
	}
}

func (c *categoryWeb) CategoryUpdateProcess(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	id := ctx.Request.FormValue("id")
	name := ctx.Request.FormValue("name")

	status, err := c.categoryClient.UpdateCategory(session.Token, id, name)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		ctx.Redirect(http.StatusSeeOther, "/client/category")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Add Category Failed!")
	}
}

func (c *categoryWeb) CategoryDeleteProcess(ctx *gin.Context) {
	var email string
	if temp, ok := ctx.Get("email"); ok {
		if contextData, ok := temp.(string); ok {
			email = contextData
		}
	}

	session, err := c.sessionService.GetSessionByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	categoryID := ctx.Param("id")

	status, err := c.categoryClient.DeleteCategory(session.Token, categoryID)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message="+err.Error())
		return
	}

	if status == 201 {
		ctx.Redirect(http.StatusSeeOther, "/client/category")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/client/modal?status=error&message=Add Category Failed!")
	}
}
