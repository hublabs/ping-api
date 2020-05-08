package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/goutils/httpreq"
)

var StartTime time.Time

func init() {
	StartTime = time.Now()
}
func main() {
	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.GET("info", func(c echo.Context) error {
		now := time.Now()
		info := map[string]interface{}{
			"StartTimeSecs":   StartTime.UTC().Unix(),
			"CurrentTimeSecs": now.UTC().Unix(),
			"Uptime":          now.Sub(StartTime),
		}
		header := make(map[string]interface{})
		for k, v := range c.Request().Header {
			header[k] = v
		}
		info["Header"] = header

		return c.JSON(http.StatusAccepted, info)
	})
	e.Any("/call/*", call)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm "+c.Request().Host)
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if err := e.Start(":8000"); err != nil {
		log.Println(err)
	}
}

func call(c echo.Context) error {
	requestURI := c.Request().RequestURI
	targetApiURI := requestURI[strings.Index(requestURI, "/call/")+len("/call/"):]
	apiAddr := strings.Split(targetApiURI, "/")[0]
	path := targetApiURI[len(apiAddr):]
	fmt.Println("targetApiURI:", targetApiURI)
	fmt.Println("apiAddr:", apiAddr)
	fmt.Println("path:", path)

	req, err := http.NewRequest(c.Request().Method, "http://"+targetApiURI, c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	authHeader := c.Request().Header.Get("Authorization")
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpreq.NewClient().Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	return c.Blob(resp.StatusCode, echo.MIMEApplicationJSON, body)
}
