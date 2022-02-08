package modules

import (
	"boilerplate-golang-v2/api"
	"boilerplate-golang-v2/api/common"
	"boilerplate-golang-v2/util"

	contentCtrlV1 "boilerplate-golang-v2/api/v1/content"
	contentBussiness "boilerplate-golang-v2/business/content"
	contentRepo "boilerplate-golang-v2/modules/repository/content"

	echo "github.com/labstack/echo"
)

//SetErrorHandler - set error response
func SetErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// error message must be known RC value
		errResp := common.NewInternalServerErrorResponse()
		c.JSON(500, errResp)
	}
}

//RegisterController - register the controller
func RegisterController(dbCon *util.DatabaseConnection) api.Controller {

	//initiate content
	contentPermitRepo := contentRepo.RepositoryFactory(dbCon)
	contentPermitService := contentBussiness.NewService(contentPermitRepo)
	contentPermitControllerV1 := contentCtrlV1.NewController(contentPermitService)

	//lets put the controller together
	controllers := api.Controller{
		ContentController: contentPermitControllerV1,
	}

	return controllers
}
