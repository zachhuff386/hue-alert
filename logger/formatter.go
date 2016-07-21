package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/zachhuff386/hue-alert/colorize"
)

var (
	blueArrow    = colorize.ColorString("▶", colorize.BlueBold, colorize.None)
	whiteDiamond = colorize.ColorString("◆", colorize.WhiteBold, colorize.None)
)

func format(entry *logrus.Entry) (output []byte) {
	msg := fmt.Sprintf("%s %s %s", formatLevel(entry.Level), blueArrow,
		entry.Message)

	var errStr string
	for key, val := range entry.Data {
		if key == "error" {
			errStr = fmt.Sprintf("%s", val)
			continue
		}

		msg += fmt.Sprintf(" %s %s=%v", whiteDiamond,
			colorize.ColorString(key, colorize.CyanBold, colorize.None),
			colorize.ColorString(fmt.Sprintf("%#v", val),
				colorize.GreenBold, colorize.None))
	}

	if errStr != "" {
		msg += "\n" + colorize.ColorString(errStr, colorize.Red, colorize.None)
	}

	if string(msg[len(msg)-1]) != "\n" {
		msg += "\n"
	}

	output = []byte(msg)

	return
}

func formatPlain(entry *logrus.Entry) (output []byte) {
	msg := fmt.Sprintf("%s ▶ %s", formatLevelPlain(entry.Level), entry.Message)

	var errStr string
	for key, val := range entry.Data {
		if key == "error" {
			errStr = fmt.Sprintf("%s", val)
			continue
		}

		msg += fmt.Sprintf(" ◆ %s=%v", key, fmt.Sprintf("%#v", val))
	}

	if errStr != "" {
		msg += "\n" + errStr
	}

	if string(msg[len(msg)-1]) != "\n" {
		msg += "\n"
	}

	output = []byte(msg)

	return
}

func formatLevel(lvl logrus.Level) (str string) {
	var colorBg colorize.Color

	switch lvl {
	case logrus.InfoLevel:
		colorBg = colorize.CyanBg
		str = "[INFO]"
	case logrus.WarnLevel:
		colorBg = colorize.YellowBg
		str = "[WARN]"
	case logrus.ErrorLevel:
		colorBg = colorize.RedBg
		str = "[ERRO]"
	case logrus.FatalLevel:
		colorBg = colorize.RedBg
		str = "[FATL]"
	case logrus.PanicLevel:
		colorBg = colorize.RedBg
		str = "[PANC]"
	default:
		colorBg = colorize.BlackBg
	}

	str = colorize.ColorString(str, colorize.WhiteBold, colorBg)

	return
}

func formatLevelPlain(lvl logrus.Level) string {
	switch lvl {
	case logrus.InfoLevel:
		return "[INFO]"
	case logrus.WarnLevel:
		return "[WARN]"
	case logrus.ErrorLevel:
		return "[ERRO]"
	case logrus.FatalLevel:
		return "[FATL]"
	case logrus.PanicLevel:
		return "[PANC]"
	default:
	}

	return ""
}

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return format(entry), nil
}
