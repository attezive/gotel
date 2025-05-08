package menu

import (
	"encoding/json"
	"gotel_alpha/util/network"
	"io"
	"net/http"
)

type CommandResponse struct {
	Success     bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type BotCommandInfo struct {
	Success    bool         `json:"ok"`
	BotCommand []BotCommand `json:"result"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Menu struct {
	token *string
}

func CreateMenu(token *string) *Menu {
	return &Menu{token: token}
}

func (menu *Menu) GetMyCommands() (*[]BotCommand, error) {
	const op = "getMyCommands"
	resp, err := network.GetRequest(*menu.token, op, nil)
	if err != nil {
		return nil, err
	}
	return getCommands(resp)
}

func (menu *Menu) SetMyCommands(newCommands *[]BotCommand, saveOldCommands bool) (*CommandResponse, error) {
	const op = "setMyCommands"
	commandsResult := make([]BotCommand, len(*newCommands))
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

func getCommands(resp *http.Response) (*[]BotCommand, error) {
	var commandInfo *BotCommandInfo
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &commandInfo); err != nil {
		return nil, err
	}
	return &commandInfo.BotCommand, nil
}

func getCommandResponse(resp *http.Response) (*CommandResponse, error) {
	var commandResponse CommandResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &commandResponse); err != nil {
		return nil, err
	}
	return &commandResponse, nil
}
