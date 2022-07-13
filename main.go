package main

import (
	"flow-users/flags"
	"flow-users/handler"
	"flow-users/jwt"
	"flow-users/mysql"
	"flow-users/oauth2/github"
	"flow-users/oauth2/google"
	"flow-users/oauth2/twitter"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func logFormat() string {
	// Refer to https://github.com/tkuchiki/alp
	var format string
	format += "time:${time_rfc3339}\t"
	format += "host:${remote_ip}\t"
	format += "forwardedfor:${header:x-forwarded-for}\t"
	format += "req:-\t"
	format += "status:${status}\t"
	format += "method:${method}\t"
	format += "uri:${uri}\t"
	format += "size:${bytes_out}\t"
	format += "referer:${referer}\t"
	format += "ua:${user_agent}\t"
	format += "reqtime_ns:${latency}\t"
	format += "cache:-\t"
	format += "runtime:-\t"
	format += "apptime:-\t"
	format += "vhost:${host}\t"
	format += "reqtime_human:${latency_human}\t"
	format += "x-request-id:${id}\t"
	format += "host:${host}\n"
	return format
}

func main() {
	// Get command line params / env variables
	f := flags.Get()

	//
	// Setup echo and middlewares
	//

	// Echo instance
	e := echo.New()

	// Log level
	e.Logger.SetLevel(log.Lvl(*f.LogLevel))
	e.Logger.Infof("Log level %d", *f.LogLevel)

	// Gzip
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: int(*f.GzipLevel),
	}))
	e.Logger.Infof("Gzip enabled with level %d", *f.GzipLevel)

	// CORS
	if f.AllowOrigins != nil && len(f.AllowOrigins) != 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: f.AllowOrigins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
		e.Logger.Info("CORS enabled")
		e.Logger.Debugf("CORS allow origins %v", f.AllowOrigins)
	}

	// JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*f.JwtSecret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/-/readiness" ||
				c.Path() == "/" && c.Request().Method == "POST" ||
				c.Path() == "/:provider/register" ||
				c.Path() == "/sign_in"
		},
	}))

	// Logger
	if f.LogLevel != nil && *f.LogLevel == 1 {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: logFormat(),
			Output: os.Stdout,
			Skipper: func(c echo.Context) bool {
				return c.Path() == "/-/readiness"
			},
		}))
		e.Logger.Info("Access logging with `alp`(https://github.com/tkuchiki/alp) enabled")
	}

	// Validator instance
	e.Validator = &CustomValidator{validator: validator.New()}

	//
	// Setup DB
	//

	// DB client instance
	e.Logger.Debugf("DB DSN `%s`", mysql.SetDSNTCP(*f.MysqlUser, *f.MysqlPasswd, *f.MysqlHost, int(*f.MysqlPort), *f.MysqlDB))

	// Check connection
	d, err := mysql.Open()
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err = d.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("DB connection test succeeded")

	//
	// Setup OAuth2 providers
	//

	// Github
	if f.GithubClientId != nil && *f.GithubClientId != "" && f.GithubClientSecret != nil && *f.GithubClientSecret != "" {
		if _, err := github.New(*f.GithubClientId, *f.GithubClientSecret); err != nil {
			e.Logger.Error(err.Error())
		}
		e.Logger.Info("OAuth2 provider GitHub availabled")
	}
	// Google
	if f.GoogleClientId != nil && *f.GoogleClientId != "" && f.GoogleClientSecret != nil && *f.GoogleClientSecret != "" {
		if _, err := google.New(*f.GoogleClientId, *f.GoogleClientSecret); err != nil {
			e.Logger.Error(err.Error())
		}
		e.Logger.Info("OAuth2 provider Google availabled")
	}
	// Twitter
	if f.TwitterClientId != nil && *f.TwitterClientId != "" && f.TwitterClientSecret != nil && *f.TwitterClientSecret != "" {
		if _, err := twitter.New(*f.TwitterClientId, *f.TwitterClientSecret); err != nil {
			e.Logger.Error(err.Error())
		}
		e.Logger.Info("OAuth2 provider Twitter availabled")
	}

	//
	// Routes
	//

	// Health check route
	e.GET("/-/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, "flow-users is Healthy.\n")
	})

	// Published routes
	e.POST("/", handler.Post)
	e.POST("/:provider/register", handler.PostOverOAuth2)
	e.POST("/sign_in", handler.SignIn)

	// Restricted routes
	e.GET("/", handler.Get)
	e.PATCH("/", handler.Patch)
	e.DELETE("/", handler.Delete)
	e.POST(":provider/connect", handler.ConnectOAuth2)
	e.POST(":provider/refresh", handler.RefreshOAuth2Token)
	e.DELETE(":provider", handler.DisconnectOAuth2)
	e.GET("id", handler.GetId)

	//
	// Start echo
	//
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *f.Port)))
}
