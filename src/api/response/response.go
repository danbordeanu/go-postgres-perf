package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"postgres-perf/configuration"
	"postgres-perf/model"
	"postgres-perf/utils"
)

// SuccessResponse return success response
func SuccessResponse(c *gin.Context, id string, data interface{}) {
	c.JSON(http.StatusOK, model.JSONSuccessResult{
		Code:    http.StatusOK,
		Id:      id,
		Data:    data,
		Message: "Success",
	})
}

// FailureResponse return failure response
func FailureResponse(c *gin.Context, data interface{}, err utils.HttpError) {
	if err.Err == nil {
		err = utils.HttpError{Code: int(math.Max(float64(err.Code), 500)), Err: fmt.Errorf("FailureResponse was called with a nil error (%s)", err.Message)}
	}
	var errorString, stackString string
	conf := configuration.AppConfig()
	if conf.Development {
		errorString = err.Error()
		stackString = err.StackTrace()
	}
	c.JSON(err.Code, model.JSONFailureResult{
		Code:  err.Code,
		Data:  data,
		Error: errorString,
		Stack: stackString,
		Id:    c.MustGet("correlation_id").(string),
	})
}
