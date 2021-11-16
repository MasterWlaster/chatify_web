package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	apiServer = ""
)

func InitRoutes(r *gin.Engine) *gin.Engine {
	r.GET("/", index)

	r.GET("/login", login)
	r.GET("/messenger", messenger)
	r.GET("/messenger/:id", dialog)

	r.POST("/auth/:command", auth)
	r.POST("/messenger/:id", sendMessage)

	return r
}

func index(c *gin.Context) {
	token, err := c.Cookie("auth-token")
	if err != nil || token == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Redirect(http.StatusFound, "/messenger")
}

func login(c *gin.Context) {
	token, err := c.Cookie("auth-token")
	if err != nil || token == "" {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/messenger")
}

func messenger(c *gin.Context) {
	res, err := requestWithToken(c, "/api/dialogs", "GET", nil)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "dialogs.html", res)
}

func dialog(c *gin.Context) {
	id := c.Param("id")

	res, err := requestWithToken(c, "/api/dialogs/"+id, "GET", nil)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	res["id"] = id

	c.HTML(http.StatusOK, "chat.html", res)
}

func auth(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.Redirect(http.StatusFound, "/")
		log.Print(err)
		return
	}

	input := map[string]interface{}{
		"username": c.Request.PostFormValue("username"),
		"password": c.Request.PostFormValue("password"),
		"name":     "-",
	}

	path := "/auth/sign-up"

	if c.Param("command") == "log-in" {
		delete(input, "name")
		path = "/auth/log-in"
	}

	res, err := request(c, path, "POST", input)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		log.Print(err)
		return
	}

	sToken, ok := res["token"].(string)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		log.Print("error to string the token")
		return
	}

	c.SetCookie("auth-token", sToken, 2*60, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func sendMessage(c *gin.Context) {
	id := c.Param("id")

	if err := c.Request.ParseForm(); err != nil {
		c.Redirect(http.StatusFound, "/messenger/"+id)
		log.Print(err)
		return
	}

	iId, err := strconv.Atoi(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/messenger/"+id)
		log.Print(err)
		return
	}

	_, err = requestWithToken(c, "/api/message", "POST", map[string]interface{}{
		"receiver_id": iId,
		"text":        c.Request.PostFormValue("text"),
	})
	if err != nil {
		c.Redirect(http.StatusFound, "/messenger/"+id)
		log.Print(err)
		return
	}

	c.Redirect(http.StatusFound, "/messenger/"+id+"#bottom")
}

///////////////////////

func RequestApi(relativePath string, method string, token string, input interface{}) (map[string]interface{}, error) {
	jsonValue, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, apiServer+relativePath, bytes.NewBuffer(jsonValue))

	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bo, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if len(bo) == 0 {
		return nil, nil
	}

	var ma map[string]interface{}

	err = json.Unmarshal(bo, &ma)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return ma, nil
}

func requestWithToken(c *gin.Context, path string, method string, input interface{}) (map[string]interface{}, error) {
	token, err := c.Cookie("auth-token")
	if err != nil || token == "" {
		return nil, err
	}

	res, err := RequestApi(path, method, token, input)
	if err != nil {
		return nil, err
	}

	if _, ok := res["message"]; ok {
		return nil, errors.New("error response: " + res["message"].(string))
	}

	return res, nil
}

func request(c *gin.Context, path string, method string, input interface{}) (map[string]interface{}, error) {
	res, err := RequestApi(path, method, "", input)
	if err != nil {
		return nil, err
	}

	if _, ok := res["message"]; ok {
		return nil, errors.New("error response: " + res["message"].(string))
	}

	return res, nil
}