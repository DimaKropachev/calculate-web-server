package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DimaKropachev/calculate-web-server/demon/internal/model"
)

type Api struct {
	url string
}

func NewApi(url string) *Api {
	return &Api{url: url}
}

func (api *Api) Get() (*model.Task, error) {
	req, err := http.NewRequest("GET", api.url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error client.Do :%w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read resp: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	task := &model.Task{}
	err = json.Unmarshal(data, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (api *Api) Give(resp *model.Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	reqBody := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", api.url, reqBody)
	if err != nil {
		return err
	}

	client := http.Client{}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
