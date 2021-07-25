package api

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/validate"
	"go.uber.org/zap"
)

type Endpoint struct {
	ApiVersion  string                // '/' is optional
	Middlewares []echo.MiddlewareFunc // global middlewares
	Groups      []*ApiGroup           // api groups
}

type ApiGroup struct {
	Prefix      string
	Middlewares []echo.MiddlewareFunc
	Routes      []*Route
}

type Route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

func CreateEndpoint(app *echo.Echo, def *Endpoint) {
	type RESTMethod func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

	var root *echo.Group
	version := "/"
	if strings.HasPrefix(def.ApiVersion, "/") {
		version = def.ApiVersion
	} else {
		version += def.ApiVersion
	}
	root = app.Group(version, def.Middlewares...)

	for _, group := range def.Groups {
		echoGroup := root.Group(group.Prefix, group.Middlewares...)
		for _, api := range group.Routes {
			var method RESTMethod
			switch api.Method {
			case http.MethodGet:
				method = echoGroup.GET
			case http.MethodPost:
				method = echoGroup.POST
			case http.MethodPut:
				method = echoGroup.PUT
			case http.MethodDelete:
				method = echoGroup.DELETE
			case http.MethodHead:
				method = echoGroup.HEAD
			case http.MethodConnect:
				method = echoGroup.CONNECT
			default:
				panic(fmt.Errorf("unknown method %s", api.Method))
			}
			method(api.Path, api.Handler, api.Middlewares...)
		}
	}
}

// PrintRoutes print all registered routes
func PrintRoutes(app *echo.Echo, logger *zap.Logger) {
	var routes [][2]string
	for _, route := range app.Routes() {
		if !strings.HasPrefix(route.Name, "github.com/labstack/echo") {
			routes = append(routes, [2]string{route.Method, route.Path})
		}
	}
	sort.Slice(routes, func(i, j int) bool {
		return routes[i][1] < routes[j][1]
	})
	for _, route := range routes {
		logger.Debug("Registered route", zap.String("method", route[0]), zap.String("path", route[1]))
	}
}

// RegisterProfileEndpoints register standard go profile api endpoints
func RegisterProfileEndpoints(app *echo.Echo) {
	expvarHandler := expvar.Handler()
	app.GET("/debug/vars", func(c echo.Context) error {
		expvarHandler.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	app.GET("/debug/pprof/", func(c echo.Context) error {
		pprof.Index(c.Response().Writer, c.Request())
		return nil
	})
	app.GET("/debug/pprof/:name", func(c echo.Context) error {
		switch c.Param("name") {
		case "cmdline":
			pprof.Cmdline(c.Response().Writer, c.Request())
		case "profile":
			pprof.Profile(c.Response().Writer, c.Request())
		case "symbol":
			pprof.Symbol(c.Response().Writer, c.Request())
		case "trace":
			pprof.Trace(c.Response().Writer, c.Request())
		default:
			pprof.Handler(c.Param("name")).ServeHTTP(c.Response().Writer, c.Request())
		}
		return nil
	})
}

func ValidateFailedResponse(c echo.Context, err *validate.ValidationError) error {
	return c.JSON(http.StatusBadRequest, NewRESTValidationError(ErrFailedValidate.Error(), err))
}

func BadRequestResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, NewRESTStandardError(err.Error()))
}

func JsonResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

func StatusResponse(c echo.Context, code int) error {
	return c.String(code, http.StatusText(code))
}

func BindErrorResponse(c echo.Context, err error) error {
	if he, ok := err.(*echo.HTTPError); ok {
		return c.JSON(http.StatusBadRequest, NewRESTBindingError(ErrFailedBinding.Error(), he.Message))
	}
	if be, ok := err.(*echo.BindingError); ok {
		return c.JSON(http.StatusBadRequest, NewRESTBindingError(ErrFailedBinding.Error(), be.Message))
	}
	return c.JSON(http.StatusBadRequest, NewRESTStandardError(ErrFailedBinding.Error()))
}

func WithContextValue(c echo.Context, key interface{}, val interface{}) {
	ctx := c.Request().Context()
	req := c.Request().WithContext(context.WithValue(ctx, key, val))
	c.SetRequest(req)
}
