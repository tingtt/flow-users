package main

import (
	"flag"
	"flow-users/jwt"
	"flow-users/mysql"
	"flow-users/oauth2/github"
	"flow-users/oauth2/google"
	"flow-users/oauth2/twitter"
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var intValue, err = strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Priority: command line params > env variables > default value
var (
	port                = flag.Int("port", getIntEnv("PORT", 1323), "Server port")
	logLevel            = flag.Int("log-level", getIntEnv("LOG_LEVEL", 2), "Log level (1: 'DEBUG', 2: 'INFO', 3: 'WARN', 4: 'ERROR', 5: 'OFF', 6: 'PANIC', 7: 'FATAL'")
	gzipLevel           = flag.Int("gzip-level", getIntEnv("GZIP_LEVEL", 6), "Gzip compression level")
	mysqlHost           = flag.String("mysql-host", getEnv("MYSQL_HOST", "db"), "MySQL host")
	mysqlPort           = flag.Int("mysql-port", getIntEnv("MYSQL_PORT", 3306), "MySQL port")
	mysqlDB             = flag.String("mysql-database", getEnv("MYSQL_DATABASE", "flow-users"), "MySQL database")
	mysqlUser           = flag.String("mysql-user", getEnv("MYSQL_USER", "flow-users"), "MySQL user")
	mysqlPasswd         = flag.String("mysql-password", getEnv("MYSQL_PASSWORD", ""), "MySQL password")
	jwtIssuer           = flag.String("jwt-issuer", getEnv("JWT_ISSUER", "flow-users"), "JWT issuer")
	jwtSecret           = flag.String("jwt-secret", getEnv("JWT_SECRET", ""), "JWT secret")
	githubClientId      = flag.String("github-client-id", getEnv("GITHUB_CLIENT_ID", ""), "GitHub client id")
	githubClientSecret  = flag.String("github-client-secret", getEnv("GITHUB_CLIENT_SECRET", ""), "GitHub client secret")
	googleClientId      = flag.String("google-client-id", getEnv("GOOGLE_CLIENT_ID", ""), "Google client id")
	googleClientSecret  = flag.String("google-client-secret", getEnv("GOOGLE_CLIENT_SECRET", ""), "Google client secret")
	twitterClientId     = flag.String("twitter-client-id", getEnv("TWITTER_CLIENT_ID", ""), "Twitter client id")
	twitterClientSecret = flag.String("twitter-client-secret", getEnv("TWITTER_CLIENT_SECRET", ""), "Twitter client secret")
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: *gzipLevel,
	}))
	e.Logger.SetLevel(log.Lvl(*logLevel))
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup db client instance
	e.Logger.Info(mysql.SetDSNTCP(*mysqlUser, *mysqlPasswd, *mysqlHost, *mysqlPort, *mysqlDB))

	_, err := github.New(*githubClientId, *githubClientSecret)
	if err != nil {
		e.Logger.Error(err.Error())
	}
	_, err = google.New(*googleClientId, *googleClientSecret)
	if err != nil {
		e.Logger.Error(err.Error())
	}
	_, err = twitter.New(*twitterClientId, *twitterClientSecret)
	if err != nil {
		e.Logger.Error(err.Error())
	}

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*jwtSecret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/" && c.Request().Method == "POST" ||
				c.Path() == "/:provider/register" ||
				c.Path() == "/sign_in"
		},
	}))

	// Opened routes
	e.POST("/", post)
	e.POST("/:provider/register", postOverOAuth2)
	e.POST("/sign_in", signIn)

	// Restricted routes
	e.PATCH("/", patch)
	e.DELETE("/", delete)
	e.POST(":provider/connect", connectOAuth2)
	e.POST(":provider/refresh", refreshOAuth2Token)
	e.DELETE(":provider", disconnectOAuth2)
	e.GET("id", getId)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
