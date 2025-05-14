package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotel_alpha/data"
	"gotel_alpha/util/network"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Sender struct {
	token *string
}

func CreateSender(token *string) *Sender {
	return &Sender{token: token}
}

func (sender *Sender) SendMessage(sendingMessage *data.SendingEntity, returnedMsg chan<- *data.Message, errCh chan<- error) {
	const op = "sendMessage"
	params := make(map[string]string)
	params["chat_id"] = sendingMessage.ChatId
	params["text"] = sendingMessage.Value.(string)
	resp, err := network.GetRequest(*sender.token, op, params)
	if err != nil {
		returnedMsg <- nil
		errCh <- fmt.Errorf("%s: %w", op, err)
		return
	}
	respMessage, err := getMessage(resp)
	returnedMsg <- respMessage
	errCh <- err
}

func (sender *Sender) SendPhoto(sendingPhoto *data.SendingEntity, returnedMsg chan<- *data.Message, errCh chan<- error) {
	const op = "sendPhoto"
	var resp *http.Response
	var err error
	params := make(map[string]string)
	params["chat_id"] = sendingPhoto.ChatId
	photo := sendingPhoto.Value.(*data.Photo)
	if photo.FileId == "" {
		body, writer, err := createForm(photo.FilePath)
		if err != nil {
			errCh <- fmt.Errorf("%s: %w", op, err)
			return
		}
		writer.Close()
		resp, err = network.PostRequest(*sender.token, op, params, body, writer)
	} else {
		params["photo"] = photo.FileId
		resp, err = network.GetRequest(*sender.token, op, params)
	}
	if err != nil {
		errCh <- fmt.Errorf("%s: %w", op, err)
		return
	}
	msg, err := getMessage(resp)
	returnedMsg <- msg
	errCh <- fmt.Errorf("%s: %w", op, err)
}

func createForm(filePath string) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	part, _ := writer.CreateFormFile("photo", filePath)
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}
	return body, writer, nil
}

func getMessage(resp *http.Response) (*data.Message, error) {
	var messageInfo data.MessageInfo
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &messageInfo); err != nil {
		return nil, err
	}
	return &messageInfo.Message, nil
}
