package api

import (
	"github.com/jphacks/B_2121_server/openapi"
	"github.com/labstack/echo/v4"
)

func GetSwaggerSpec(ctx echo.Context) error {
	swg, err := openapi.GetSwagger()

	if err != nil {
		return err
	}
	jsonData, err := swg.MarshalJSON()
	if err != nil {
		return err
	}
	return ctx.Blob(200, "application/json", jsonData)
}
