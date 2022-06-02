package error

import (
	"github.com/NubeIO/flow-framework/nerrors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	nerrors.NotFound()(ctx)

	assertJSONResponse(t, rec, 404, `{"errorCode":404, "errorDescription":"page not found", "error":"Not Found"}`)
}
