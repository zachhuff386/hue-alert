package hue

import (
	"bytes"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/errortypes"
	"io/ioutil"
	"net/http"
)

type registerParams struct {
	DeviceType string `json:"devicetype"`
}

type registerData struct {
	Success struct {
		Username string `json:"username"`
	} `json:"success"`
	Error struct {
		Type        int    `json:"type"`
		Address     string `json:"address"`
		Description string `json:"description"`
	} `json:"error"`
}

func (h *Hue) Register() (err error) {
	params := &registerParams{
		DeviceType: "go-hue",
	}

	reqData, err := json.Marshal(params)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "hue: Register marshal error"),
		}
		return
	}

	reqDataBuf := bytes.NewBuffer(reqData)

	resp, err := http.Post(h.getUrl("/api"), "application/json", reqDataBuf)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Register request error"),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Register read error"),
		}
		return
	}

	datas := []*registerData{}

	err = json.Unmarshal(body, &datas)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Register unmarshal error"),
		}
		return
	}

	data := datas[0]

	if data.Error.Type != 0 {
		err = &errortypes.ApiError{
			errors.Newf("hue: Failed to register '%s'",
				data.Error.Description),
		}
		return
	}

	h.Username = data.Success.Username

	return
}
