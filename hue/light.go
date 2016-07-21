package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/zachhuff386/hue-alert/errortypes"
	"io/ioutil"
	"net/http"
	"time"
)

type lightStateParams struct {
	On          interface{} `json:"on,omitempty"`
	Brightness  int         `json:"bri,omitempty"`
	Hue         int         `json:"hue,omitempty"`
	Saturation  int         `json:"sat,omitempty"`
	ColorXY     []float64   `json:"xy,omitempty"`
	Temperature int         `json:"ct,omitempty"`
	Alert       string      `json:"alert,omitempty"`
	Effect      string      `json:"effect,omitempty"`
	Transition  int64       `json:"transitiontime,omitempty"`
}

type lightStateData struct {
	Success map[string]interface{} `json:"success"`
	Error   struct {
		Type        int    `json:"type"`
		Address     string `json:"address"`
		Description string `json:"description"`
	} `json:"error"`
}

type lightData struct {
	State struct {
		On          bool      `json:"on"`
		Brightness  int       `json:"bri"`
		Hue         int       `json:"hue"`
		Saturation  int       `json:"sat"`
		ColorXY     []float64 `json:"xy"`
		Temperature int       `json:"ct"`
		Alert       string    `json:"alert"`
		Effect      string    `json:"effect"`
		Mode        string    `json:"colormode"`
		Reachable   bool      `json:"reachabe"`
	} `json:"state"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	ModeLid  string `json:"modelid"`
	Version  string `json:"swversion"`
	UniqueId string `json:"uniqueid"`
	Error    struct {
		Type        int    `json:"type"`
		Address     string `json:"address"`
		Description string `json:"description"`
	} `json:"error"`
}

type Light struct {
	Id          string
	UniqueId    string
	Name        string
	Type        string
	State       bool
	Reachable   bool
	Alert       string
	Effect      string
	Mode        string
	Brightness  int
	Hue         int
	Saturation  int
	ColorX      float64
	ColorY      float64
	Temperature int
	hue         *Hue
	changed     set.Set
	transition  int64
}

func (l *Light) Copy() (light *Light) {
	light = &Light{
		Id:          l.Id,
		Name:        l.Name,
		Type:        l.Type,
		State:       l.State,
		Reachable:   l.Reachable,
		Alert:       l.Alert,
		Effect:      l.Effect,
		Mode:        l.Mode,
		Brightness:  l.Brightness,
		Hue:         l.Hue,
		Saturation:  l.Saturation,
		ColorX:      l.ColorX,
		ColorY:      l.ColorY,
		Temperature: l.Temperature,
		hue:         l.hue,
		changed:     l.changed,
		transition:  l.transition,
	}
	return
}

func (l *Light) Print() {
	fmt.Printf("Id: %s\n", l.Id)
	fmt.Printf("Name: %s\n", l.Name)
	fmt.Printf("Type: %s\n", l.Type)
	fmt.Printf("State: %t\n", l.State)
	fmt.Printf("Reachable: %t\n", l.Reachable)
	fmt.Printf("Alert: %s\n", l.Alert)
	fmt.Printf("Effect: %s\n", l.Effect)
	fmt.Printf("Mode: %s\n", l.Mode)
	fmt.Printf("Brightness: %d\n", l.Brightness)
	fmt.Printf("Hue: %d\n", l.Hue)
	fmt.Printf("Saturation: %d\n", l.Saturation)
	fmt.Printf("ColorX: %f\n", l.ColorX)
	fmt.Printf("ColorY: %f\n", l.ColorY)
	fmt.Printf("Temperature: %d\n", l.Temperature)
}

func (l *Light) On() {
	l.State = true
	l.changed.Add("state")
}

func (l *Light) Off() {
	l.State = false
	l.changed.Add("state")
}

func (l *Light) SetState(state bool) {
	l.State = state
	l.changed.Add("state")
}

func (l *Light) Switch() {
	l.State = !l.State
	l.changed.Add("state")
}

func (l *Light) SetTransition(transition time.Duration) {
	l.transition = transition.Nanoseconds() / 100000000
	l.changed.Add("transition")
}

func (l *Light) ColorLoop() {
	l.Effect = "colorloop"
	l.changed.Add("effect")
}

func (l *Light) NoColorLoop() {
	l.Effect = "none"
	l.changed.Add("effect")
}

func (l *Light) SetBrightness(brightness int) {
	l.Brightness = brightness
	l.changed.Add("brightness")
}

func (l *Light) SetColorHex(color string) (err error) {
	clr, err := colorful.Hex(color)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "hue: Failed to parse color"),
		}
		return
	}

	colorX, colorY, _ := clr.Xyy()

	l.ColorX = colorX
	l.ColorY = colorY

	l.changed.Add("color")

	return
}

func (l *Light) SetColorXY(colorX, colorY float64) (err error) {
	l.ColorX = colorX
	l.ColorY = colorY

	l.changed.Add("color")

	return
}

func (l *Light) Update() (err error) {
	url := l.hue.getAuthUrl("/lights/" + l.Id)

	resp, err := http.Get(url)
	if err != nil {
		err = errortypes.ApiError{
			errors.Wrap(err, "hue: Light update request error"),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Light update read error"),
		}
		return
	}

	var datasInf interface{}

	err = json.Unmarshal(body, &datasInf)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Light update unmarshal error"),
		}
		return
	}

	var data *lightData

	switch datasInf.(type) {
	case []interface{}:
		datas := []*lightData{}

		err = json.Unmarshal(body, &datas)
		if err != nil {
			err = &errortypes.ApiError{
				errors.Wrap(err, "hue: Light update unmarshal error"),
			}
			return
		}

		data = datas[0]
	case map[string]interface{}:
		err = json.Unmarshal(body, &data)
		if err != nil {
			err = &errortypes.ApiError{
				errors.Wrap(err, "hue: Light update unmarshal error"),
			}
			return
		}
	default:
		err = &errortypes.UnknownError{
			errors.New("hue: Light update unknown data"),
		}
		return
	}

	if data.Error.Type != 0 {
		err = &errortypes.ApiError{
			errors.Newf("hue: Light update error '%s'",
				data.Error.Description),
		}
		return
	}

	colorX := 0.0
	colorY := 0.0

	if len(data.State.ColorXY) == 2 {
		colorX = data.State.ColorXY[0]
		colorY = data.State.ColorXY[1]
	}

	l.UniqueId = data.UniqueId
	l.Name = data.Name
	l.Type = data.Type
	l.State = data.State.On
	l.Brightness = data.State.Brightness
	l.Hue = data.State.Hue
	l.Saturation = data.State.Saturation
	l.ColorX = colorX
	l.ColorY = colorY
	l.Temperature = data.State.Temperature
	l.Alert = data.State.Alert
	l.Effect = data.State.Effect
	l.Mode = data.State.Mode
	l.Reachable = data.State.Reachable
	l.changed = set.NewSet()

	return
}

func (l *Light) Commit() (err error) {
	if l.changed.Len() == 0 {
		return
	}

	params := &lightStateParams{}

	for keyInf := range l.changed.Iter() {
		key := keyInf.(string)

		switch key {
		case "state":
			params.On = l.State
		case "color":
			params.ColorXY = []float64{
				l.ColorX,
				l.ColorY,
			}
		case "brightness":
			params.Brightness = l.Brightness
		case "alert":
			params.Alert = l.Alert
		case "effect":
			params.Effect = l.Effect
		case "transition":
			params.Transition = l.transition
		}
	}

	reqData, err := json.Marshal(params)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "hue: Light marshal error"),
		}
		return
	}

	reqDataBuf := bytes.NewBuffer(reqData)

	url := l.hue.getAuthUrl("/lights/" + l.Id + "/state")

	req, err := http.NewRequest("PUT", url, reqDataBuf)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "hue: Light request error"),
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Light request error"),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Light read error"),
		}
		return
	}

	datas := []*lightStateData{}

	err = json.Unmarshal(body, &datas)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Light unmarshal error"),
		}
		return
	}

	for _, data := range datas {
		if data.Error.Type != 0 {
			err = &errortypes.ApiError{
				errors.Newf("hue: Light commit error '%s'",
					data.Error.Description),
			}
			return
		}
	}

	l.changed = set.NewSet()

	return
}
