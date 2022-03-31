package serve

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alia5/urlshort/urlshort"
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
}

type CreateLinkBody struct {
	CustomShortText string `json:"customShortText"`
}

var defaultUrls = ServeUrls{redirUrl: "u", createUrl: "u"}

func Run(opts ServeOptions) {
	if opts.Urls.redirUrl == "" {
		opts.Urls.redirUrl = defaultUrls.redirUrl
	}
	if opts.Urls.createUrl == "" {
		opts.Urls.createUrl = defaultUrls.createUrl
	}
	r := gin.Default()
	r.Use(gin.Recovery())
	if !opts.Debug {
		gin.SetMode(gin.ReleaseMode)
		r.SetTrustedProxies([]string{"localhost"}) // TODO: configurable
	}

	r.GET(fmt.Sprintf("/%s/:tinyurl", opts.Urls.redirUrl), redirectUrl("tinyurl"))
	r.POST(fmt.Sprintf("/%s/*url", opts.Urls.createUrl), shortenUrl("url"))

	r.Run(fmt.Sprintf(":%d", opts.Port)) // listen and serve on port
}

func redirectUrl(urlParam string) func(*gin.Context) {
	return func(c *gin.Context) {
		fullUrl, err := urlshort.Unshorten(urlshort.UnshortenParams{Text: c.Param(urlParam), Click: true})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "404",
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

		var (
			shortUrl string
			created  bool
			err      error
		)
		// for some odd reason, gin adds a leading "/"
		// thanks gin! I knew I couldn't trust you... Who the fuck drinks gin with a 🍋-slice, anyway (except monkeys)?! 🥒
		// trim that...
		url := c.Param(urlParam)[1:]

		if byteBody, bodyReadErr := ioutil.ReadAll(c.Request.Body); bodyReadErr == nil && string(byteBody) != "" {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteBody))
			var json CreateLinkBody
			if err = c.ShouldBindJSON(&json); err == nil {
				shortUrl, created, err = urlshort.ShortenWithName(url, json.CustomShortText)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "400",
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
