package hue

import ()

type Hue struct {
	Host     string
	Username string
}

func (h *Hue) getUrl(path string) string {
	return "http://" + h.Host + path
}

func (h *Hue) getAuthUrl(path string) string {
	return "http://" + h.Host + "/api/" + h.Username + path
}
