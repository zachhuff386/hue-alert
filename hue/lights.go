package hue

import (
	"encoding/json"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"github.com/zachhuff386/hue-alert/errortypes"
	"io/ioutil"
	"net/http"
)

func (h *Hue) GetLights() (lights []*Light, err error) {
	url := h.getAuthUrl("/lights")

	resp, err := http.Get(url)
	if err != nil {
		err = errortypes.ApiError{
			errors.Wrap(err, "hue: Lights request error"),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Lights read error"),
		}
		return
	}

	datas := map[string]*lightData{}

	err = json.Unmarshal(body, &datas)
	if err != nil {
		err = &errortypes.ApiError{
			errors.Wrap(err, "hue: Lights unmarshal error"),
		}
		return
	}

	lights = []*Light{}

	for id, data := range datas {
		if data.Error.Type != 0 {
			err = &errortypes.ApiError{
				errors.Newf("hue: Lights request error '%s'",
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

		lights = append(lights, &Light{
			Id:          id,
			UniqueId:    data.UniqueId,
			Name:        data.Name,
			Type:        data.Type,
			State:       data.State.On,
			Brightness:  data.State.Brightness,
			Hue:         data.State.Hue,
			Saturation:  data.State.Saturation,
			ColorX:      colorX,
			ColorY:      colorY,
			Temperature: data.State.Temperature,
			Alert:       data.State.Alert,
			Effect:      data.State.Effect,
			Mode:        data.State.Mode,
			Reachable:   data.State.Reachable,
			hue:         h,
			changed:     set.NewSet(),
		})
	}

	return
}

func (h *Hue) GetLightsById(lightIds []string) (lights []*Light, err error) {
	lights = []*Light{}
	lightIdsSet := set.NewSet()

	for _, lightId := range lightIds {
		lightIdsSet.Add(lightId)
	}

	allLights, err := h.GetLights()
	if err != nil {
		return
	}

	for _, light := range allLights {
		if lightIdsSet.Contains(light.UniqueId) {
			lights = append(lights, light)
		}
	}

	return
}
