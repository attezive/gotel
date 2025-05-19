package menu

import (
	"encoding/json"
	"github.com/attezive/gotel/data"
	"github.com/attezive/gotel/util/network"
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

func (menu *Menu) GetMyCommands(cmdCh chan<- *[]data.BotCommand, errCh chan<- error) {
	const op = "getMyCommands"
	resp, err := network.GetRequest(*menu.token, op, nil)
	if err != nil {
		errCh <- err
		return
	}
	commands, err := getCommands(resp)
	cmdCh <- commands
	errCh <- err
}

func (menu *Menu) SetMyCommands(newCommands *[]data.BotCommand, saveOldCommands bool,
	rspCh chan<- *data.SuccessResponse, errCh chan<- error) {
	const op = "setMyCommands"
	commandsResult := make([]data.BotCommand, len(*newCommands))
	params := make(map[string]string, 1)
	if saveOldCommands {
		errGetCh := make(chan error, 1)
		defer close(errGetCh)
		cmdCh := make(chan *[]data.BotCommand, 1)
		defer close(cmdCh)
		menu.GetMyCommands(cmdCh, errGetCh)
		if err := <-errGetCh; err != nil {
			errCh <- err
			return
		}
		oldCommands := <-cmdCh
		commandsResult = append(*newCommands, *oldCommands...)
	} else {
		commandsResult = *newCommands
	}
	commands, err := json.Marshal(commandsResult)
	params["commands"] = string(commands)
	resp, err := network.GetRequest(*menu.token, op, params)
	if err != nil {
		errCh <- err
		return
	}
	success, err := getCommandResponse(resp)
	rspCh <- success
	errCh <- err
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
