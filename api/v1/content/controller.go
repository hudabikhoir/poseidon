package content

import (
	"boilerplate-golang-v2/api/common"
	"boilerplate-golang-v2/api/v1/content/request"
	"boilerplate-golang-v2/api/v1/content/response"
	"boilerplate-golang-v2/business"
	contentBusiness "boilerplate-golang-v2/business/content"
	"net/http"

	v10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

//Controller Get content API controller
type Controller struct {
	service   contentBusiness.Service
	validator *v10.Validate
}

//NewController Construct content API controller
func NewController(service contentBusiness.Service) *Controller {
	return &Controller{
		service,
		v10.New(),
	}
}

//GetContentByID Get content by ID echo handler
func (controller *Controller) GetContentByID(c echo.Context) error {
	ID := c.Param("id")
	content, err := controller.service.GetContentByID(ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
	} else if content == nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	response := response.NewGetContentByIDResponse(*content)
	return c.JSON(http.StatusOK, response)
}

//FindContentByTag Find content by tag echo handler
func (controller *Controller) FindContentByTag(c echo.Context) error {
	tag := c.QueryParam("tag")
	contents, err := controller.service.GetContentsByTag(tag)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
	}

	response := response.NewGetContentByTagResponse(contents)
	return c.JSON(http.StatusOK, response)
}

//CreateNewContent Create new content echo handler
func (controller *Controller) CreateNewContent(c echo.Context) error {
	createContentRequest := new(request.CreateContentRequest)
	if err := c.Bind(createContentRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	ID, err := controller.service.CreateContent(*createContentRequest.ToUpsertContentSpec(), "creator")

	if err != nil {
		if err == business.ErrInvalidSpec {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
	}

	response := response.NewCreateNewContentResponse(ID)
	return c.JSON(http.StatusCreated, response)
}

//UpdateContent update content echo handler
func (controller *Controller) UpdateContent(c echo.Context) error {
	updateContentRequest := new(request.UpdateContentRequest)

	if err := c.Bind(updateContentRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	err := controller.validator.Struct(updateContentRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	err = controller.service.UpdateContent(
		c.Param("id"),
		*updateContentRequest.ToUpsertContentSpec(),
		updateContentRequest.Version,
		"updater")

	if err != nil {
		if err == business.ErrNotFound {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		if err == business.ErrHasBeenModified {
			return c.JSON(http.StatusConflict, common.NewConflictResponse())
		}
	}

	return c.NoContent(http.StatusNoContent)
}
