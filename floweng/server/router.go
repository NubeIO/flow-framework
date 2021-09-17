package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

var routerServer *Server

func (s *Server) SetRouter() {
	routerServer = s
}

func GinToHttp(ctx *gin.Context, handler func(http.ResponseWriter, *http.Request)) {
	r := ctx.Request
	muxVars := make(map[string]string)
	for _, param := range ctx.Params {
		muxVars[param.Key] = param.Value
	}
	r = mux.SetURLVars(r, muxVars)

	handler(ctx.Writer.(http.ResponseWriter), r)
}

func NewRouter(router *gin.Engine, apiRoutes *gin.RouterGroup) {

	if routerServer == nil {
		panic(errors.New("flow-eng server not created before route handler"))
	}

	nodeRoutes := apiRoutes.Group("/nodes")

	nodeRoutes.StaticFile("/", "./floweng/static/index.html")
	nodeRoutes.Static("/js", "./floweng/static/js/")
	nodeRoutes.Static("/css", "./floweng/static/css/")
	nodeRoutes.Static("/lib", "./floweng/static/lib/")

	nodeRoutes.GET("/updates", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.UpdateSocketHandler) })
	nodeRoutes.GET("/blocks/library", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockLibraryHandler) })
	nodeRoutes.GET("/sources/library", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceLibraryHandler) })
	nodeRoutes.GET("/groups", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupIndexHandler) })
	nodeRoutes.GET("/groups/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupHandler) })
	nodeRoutes.POST("/groups", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupCreateHandler) })
	nodeRoutes.GET("/groups/:id/export", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupExportHandler) })
	nodeRoutes.POST("/groups/:id/import", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupImportHandler) })
	nodeRoutes.PUT("/groups/:id/label", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupModifyLabelHandler) })
	nodeRoutes.PUT("/groups/:id/children", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupModifyAllChildrenHandler) })
	nodeRoutes.PUT("/groups/:id/children/:node_id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupModifyChildHandler) })
	nodeRoutes.PUT("/groups/:id/position", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupPositionHandler) })
	nodeRoutes.DELETE("/groups/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupDeleteHandler) })
	nodeRoutes.PUT("/groups/:id/visibility", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.GroupSetRouteVisibilityHandler) })
	nodeRoutes.GET("/blocks", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockIndexHandler) })
	nodeRoutes.GET("/blocks/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockHandler) })
	nodeRoutes.POST("/blocks", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockCreateHandler) })
	nodeRoutes.DELETE("/blocks/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockDeleteHandler) })
	nodeRoutes.PUT("/blocks/:id/label", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockModifyNameHandler) })
	nodeRoutes.PUT("/blocks/:id/routes/:index", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockModifyRouteHandler) })
	nodeRoutes.PUT("/blocks/:id/position", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.BlockModifyPositionHandler) })
	nodeRoutes.GET("/connections", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.ConnectionIndexHandler) })
	nodeRoutes.GET("/connections/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.ConnectionHandler) })
	nodeRoutes.POST("/connections", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.ConnectionCreateHandler) })
	nodeRoutes.PUT("/connections/:id/coordinates", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.ConnectionModifyCoordinates) })
	nodeRoutes.DELETE("/connections/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.ConnectionDeleteHandler) })
	nodeRoutes.POST("/sources", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceCreateHandler) })
	nodeRoutes.GET("/sources", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceIndexHandler) })
	nodeRoutes.PUT("/sources/:id/label", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceModifyNameHandler) })
	nodeRoutes.PUT("/sources/:id/position", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceModifyPositionHandler) })
	nodeRoutes.GET("/sources/:id/value", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceGetValueHandler) })
	nodeRoutes.PUT("/sources/:id/value", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceSetValueHandler) })
	nodeRoutes.GET("/sources/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceHandler) })
	nodeRoutes.DELETE("/sources/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.SourceDeleteHandler) })
	nodeRoutes.GET("/links", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.LinkIndexHandler) })
	nodeRoutes.POST("/links", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.LinkCreateHandler) })
	nodeRoutes.DELETE("/links/:id", func(ctx *gin.Context) { GinToHttp(ctx, routerServer.LinkDeleteHandler) })
}
