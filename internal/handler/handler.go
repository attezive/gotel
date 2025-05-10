package handler

import (
	"encoding/json"
	"fmt"
	"gotel_alpha/data"
	u "gotel_alpha/util/network"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	lastUpdateId    int
	handleFunctions []func(*data.Update)
	token           *string
	Start           bool
}

type UpdatesInfo struct {
	Ok      bool          `json:"ok"`
	Updates []data.Update `json:"result,omitempty"`
}

func CreateHandler(token *string) *Handler {
	return &Handler{token: token}
}

func (handler *Handler) AddHandleFunction(handleFunction func(*data.Update)) {
	handler.handleFunctions = append(handler.handleFunctions, handleFunction)
}

func (handler *Handler) Handle(stop <-chan bool) error {
	lastId, err := handler.getLastId()
	if err != nil {
		return err
	}
	for {
		select {
		case <-stop:
			fmt.Println("Stopping")
			return nil
		default:
			fmt.Println("Getting last update")
		}
		updates, err := handler.getUpdates(*lastId)
		if err != nil {
			return err
		}
		err = handler.handleUpdates(updates, lastId)
		if err != nil {
			return err
		}
	}
}

func (handler *Handler) getUpdates(id int) (*UpdatesInfo, error) {
	const op = "getUpdates"
	resp, err := u.GetRequest(
		*handler.token,
		op,
		map[string]string{"offset": strconv.Itoa(id), "timeout": "20"})
	if err != nil {
		return nil, err
	}
	return getUpdate(resp)
}

func (handler *Handler) handleUpdates(updates *UpdatesInfo, lastId *int) error {
	var err error
	if len(updates.Updates) != 0 {
		if *lastId == 0 {
			lastId, err = handler.getLastId()
		}
		if err != nil {
			return err
		}
		if handler.Start {
			handler.Start = false
			return nil
		}
		*lastId++
		for _, handlerFunction := range handler.handleFunctions {
			go handlerFunction(&updates.Updates[0])
		}
	}
	return nil
}

func (handler *Handler) loadLastId() (int, error) {
	const op = "loadLastId"
	resp, err := u.GetRequest(
		*handler.token,
		"getUpdates",
		map[string]string{"offset": "-1"})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	updates, err := getUpdate(resp)
	if err != nil {
		return 0, err
	}
	if len(updates.Updates) == 0 {
		return 0, nil
	}
	return updates.Updates[0].UpdateId, nil
}

func (handler *Handler) getLastId() (*int, error) {
	var err error
	if handler.lastUpdateId == 0 {
		handler.lastUpdateId, err = handler.loadLastId()
	}
	return &handler.lastUpdateId, err
}

func getUpdate(resp *http.Response) (*UpdatesInfo, error) {
	var updates UpdatesInfo
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &updates); err != nil {
		return nil, err
	}
	return &updates, nil
}
