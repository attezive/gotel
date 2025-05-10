package menu

import (
	"encoding/json"
	"gotel_alpha/data"
	"gotel_alpha/util/network"
	"io"
	"net/http"
)

type BotCommandInfo struct {
	Success    bool              `json:"ok"`
	BotCommand []data.BotCommand `json:"result"`
}

type Menu struct {
	token *string
}

func CreateMenu(token *string) *Menu {
	return &Menu{token: token}
}

func (menu *Menu) GetMyCommands() (*[]data.BotCommand, error) {
	const op = "getMyCommands"
	resp, err := network.GetRequest(*menu.token, op, nil)
	if err != nil {
		return nil, err
	}
	return getCommands(resp)
}

func (menu *Menu) SetMyCommands(newCommands *[]data.BotCommand, saveOldCommands bool) (*data.SuccessResponse, error) {
	const op = "setMyCommands"
	commandsResult := make([]data.BotCommand, len(*newCommands))
	params := make(map[string]string, 1)
	if saveOldCommands {
		oldCommands, err := menu.GetMyCommands()
		if err != nil {
			return nil, err
		}
		commandsResult = append(*newCommands, *oldCommands...)
	} else {
		commandsResult = *newCommands
	}
	commands, err := json.Marshal(commandsResult)
	params["commands"] = string(commands)
	resp, err := network.GetRequest(*menu.token, op, params)
	if err != nil {
		return nil, err
	}
	return getCommandResponse(resp)
}

func getCommands(resp *http.Response) (*[]data.BotCommand, error) {
	var commandInfo *BotCommandInfo
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &commandInfo); err != nil {
		return nil, err
	}
	return &commandInfo.BotCommand, nil
}

func getCommandResponse(resp *http.Response) (*data.SuccessResponse, error) {
	var commandResponse data.SuccessResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &commandResponse); err != nil {
		return nil, err
	}
	return &commandResponse, nil
}
