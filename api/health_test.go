package api

import (
	"net/http/httptest"
	"testing"

	"github.com/NubeDev/plug-framework/mode"
	"github.com/NubeDev/plug-framework/model"
	"github.com/NubeDev/plug-framework/test"
	"github.com/NubeDev/plug-framework/test/testdb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}

type HealthSuite struct {
	suite.Suite
	db       *testdb.Database
	a        *HealthAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func (s *HealthSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	s.recorder = httptest.NewRecorder()
	s.db = testdb.NewDB(s.T())
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	withURL(s.ctx, "http", "example.com")
	s.a = &HealthAPI{DB: s.db}
}

func (s *HealthSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *HealthSuite) TestHealthSuccess() {
	s.a.Health(s.ctx)
	test.BodyEquals(s.T(), model.Health{Health: model.StatusGreen, Database: model.StatusGreen}, s.recorder)
}

func (s *HealthSuite) TestDatabaseFailure() {
	s.db.Close()
	s.a.Health(s.ctx)
	test.BodyEquals(s.T(), model.Health{Health: model.StatusOrange, Database: model.StatusRed}, s.recorder)
}
