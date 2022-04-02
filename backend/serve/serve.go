package serve

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/alia5/urlshort/auth"
	"github.com/alia5/urlshort/urlshort"
	"github.com/alia5/urlshort/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServeOptions struct {
	Debug bool
	Port  int16
	Urls  ServeUrls
}

type ServeUrls struct {
	redirUrl  string
	createUrl string
	authUrl   string
}

type CreateLinkBody struct {
	CustomShortText string `json:"customShortText"`
}

type AuthBody struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

var defaultUrls = ServeUrls{redirUrl: "u", createUrl: "u", authUrl: "auth"}

func Run(opts ServeOptions) {
	if opts.Urls.redirUrl == "" {
		opts.Urls.redirUrl = defaultUrls.redirUrl
	}
	if opts.Urls.createUrl == "" {
		opts.Urls.createUrl = defaultUrls.createUrl
	}
	if opts.Urls.authUrl == "" {
		opts.Urls.authUrl = defaultUrls.authUrl
	}
	r := gin.Default()
	r.Use(gin.Recovery())
	corsConf := cors.DefaultConfig()
	corsConf.AllowCredentials = true
	corsConf.AllowAllOrigins = true
	corsConf.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConf))
	if !opts.Debug {
		gin.SetMode(gin.ReleaseMode)
		r.SetTrustedProxies([]string{"localhost"}) // TODO: configurable
	}

	r.POST(fmt.Sprintf("/%s", opts.Urls.authUrl), authUrl())
	r.GET(fmt.Sprintf("/%s/:tinyurl", opts.Urls.redirUrl), redirectUrl("tinyurl"))
	r.POST(fmt.Sprintf("/%s/*url", opts.Urls.createUrl), shortenUrl("url"))

	r.Run(fmt.Sprintf(":%d", opts.Port)) // listen and serve on port
}

func redirectUrl(urlParam string) func(*gin.Context) {
	return func(c *gin.Context) {
		fullUrl, err := urlshort.Unshorten(urlshort.UnshortenParams{Text: c.Param(urlParam), Click: true})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"error":   "Not found",
				"message": err.Error(), // TODO: don't pass code errors outside; proper err handling
			})
			return
		}
		c.Redirect(
			http.StatusMovedPermanently,
			fullUrl,
		)
	}
}

func shortenUrl(urlParam string) func(*gin.Context) {
	return func(c *gin.Context) {
		// no auth middleware, we have just that one endpoint anyway...
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"error":   "Unauthorized",
				"message": "Missing Header",
			})
			return
		}
		jwt := util.SlicePop(strings.Split(authHeader, " ")) // rofl... FUCK GO!
		if _, err := auth.ValidateJwt(jwt); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"error":   "Unauthorized",
				"message": err.Error(),
			})
			return
		}

		var (
			shortUrl string
			created  bool
			err      error
		)
		// for some odd reason, gin adds a leading "/"
		// thanks gin! I knew I couldn't trust you... Who the fuck drinks gin with a 🍋-slice, anyway (except monkeys)?! 🥒
		// trim that...
		url := fmt.Sprintf("%s?%s", c.Param(urlParam)[1:], c.Request.URL.Query().Encode())

		if byteBody, bodyReadErr := ioutil.ReadAll(c.Request.Body); bodyReadErr == nil && string(byteBody) != "" {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteBody))
			var json CreateLinkBody
			if err = c.ShouldBindJSON(&json); err == nil {
				shortUrl, created, err = urlshort.ShortenWithName(url, json.CustomShortText)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  http.StatusBadRequest,
					"error":   "Bad Request",
					"message": "Malformed JSON body",
				})
				return
			}
		} else {
			shortUrl, created, err = urlshort.Shorten(url)
		}

		var httpStatus int
		if err != nil {
			httpStatus = http.StatusInternalServerError
		} else {
			if created {
				httpStatus = http.StatusCreated
			} else {
				httpStatus = http.StatusOK
			}
		}
		// TODO: proper JSON response
		c.Data(
			httpStatus,
			"text/html; charset=utf-8",
			[]byte(shortUrl),
		)
	}
}

func authUrl() func(*gin.Context) {
	return func(c *gin.Context) {
		var json AuthBody
		if err := c.ShouldBindJSON(&json); err == nil {
			jwt, err := auth.CreateJwt(json.User, json.Pass)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"jwt":    jwt,
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"error":   "Unauthorized",
					"message": "User not allowed",
				})
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"error":   "Bad Request",
				"message": "Malformed JSON body",
			})
		}
	}
}
