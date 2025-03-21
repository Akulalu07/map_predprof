package handlers

import (
	"encoding/json"
	"errors"

	"go.uber.org/zap"

	//"go/ast"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponceMap struct {
	Map [][]byte `json:"map"`
}

type Request struct {
	Url string `json:"url"`
}
type APIResponseCoords struct {
	Message struct {
		Listener []int     `json:"listener"`
		Sender   []int     `json:"sender"`
		Price    []float64 `json:"price"`
	} `json:"message"`
	Status string `json:"status"`
}
type APIResponseMap struct {
	Message struct {
		Data [][]byte `json:"data"`
	} `json:"message"`
	Status string `json:"status"`
}

type ResponseCoords struct {
	Listener []int `json:"listener"`
	Sender   []int `json:"sender"`
	Price    struct {
		Cuper float64 `json:"cuper"`
		Engel float64 `json:"engel"`
	} `json:"price"`
}

func generateMap(conn string) ([][]byte, error) {
	res, _ := http.Get(conn)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.Error(err)
	}
	APISERV := new(APIResponseMap)
	json.Unmarshal(body, &APISERV)

	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, errors.New("400")
	}
	return APISERV.Message.Data, nil
}

func generateCoords(conn string) (APIResponseCoords, error) {
	res, _ := http.Get(conn)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.Error(err)
	}
	APISERV := new(APIResponseCoords)
	json.Unmarshal(body, &APISERV)

	res.Body.Close()
	if res.StatusCode > 299 {
		return *APISERV, errors.New("400")
	}
	return *APISERV, nil
}

func Map(c echo.Context) error {
	resp := new(ResponceMap)
	var err error
	u := new(Request)
	if err := c.Bind(u); err != nil {
		return err
	}
	conn := u.Url
	resp.Map, err = generateMap(conn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	
	return c.JSON(http.StatusOK, resp)

}

func Coords(c echo.Context) error {
	resp := new(ResponseCoords)
	u := new(Request)
	if err := c.Bind(u); err != nil {
		return err
	}
	conn := u.Url
	some, err := generateCoords(conn)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err)
	}
	resp.Listener = some.Message.Listener
	resp.Sender = some.Message.Sender
	resp.Price.Cuper = some.Message.Price[0]
	resp.Price.Engel = some.Message.Price[1]

	return c.JSON(http.StatusOK, resp)
}
