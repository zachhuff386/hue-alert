package notification

import (
	"github.com/Sirupsen/logrus"
	"github.com/zachhuff386/hue-alert/constants"
	"github.com/zachhuff386/hue-alert/hue"
	"sync"
	"time"
)

var (
	wait = 200 * time.Millisecond
)

type Alert struct {
	Type  string
	Color string
}

type Notification struct {
	interrupt     bool
	interruptWait chan bool
	alerts        map[string]Alert
	alertsLock    sync.Mutex
	lights        []*hue.Light
	lightsLock    sync.Mutex
	Mode          string
	Brightness    int
}

func (n *Notification) AddAlert(alert Alert) {
	n.alertsLock.Lock()
	if n.alerts == nil {
		n.alerts = map[string]Alert{}
	}
	n.alerts[alert.Type] = alert
	n.alertsLock.Unlock()
}

func (n *Notification) RemoveAlert(alert Alert) {
	n.alertsLock.Lock()
	_, ok := n.alerts[alert.Type]
	if ok {
		delete(n.alerts, alert.Type)
	}
	n.alertsLock.Unlock()
}

func (n *Notification) AddLight(light *hue.Light) {
	n.lightsLock.Lock()
	if n.lights == nil {
		n.lights = []*hue.Light{}
	}
	n.lights = append(n.lights, light.Copy())
	n.lightsLock.Unlock()
}

func (n *Notification) runner() (err error) {
	var transition time.Duration
	var duration time.Duration
	var rate time.Duration

	switch n.Mode {
	case constants.Solid:
		transition = 500 * time.Millisecond
		duration = 2000 * time.Millisecond
	case constants.Slow:
		transition = 500 * time.Millisecond
		duration = 1500 * time.Millisecond
		rate = 1500 * time.Millisecond
	case constants.Medium:
		transition = 500 * time.Millisecond
		duration = 750 * time.Millisecond
		rate = 750 * time.Millisecond
	case constants.Fast:
		transition = 250 * time.Millisecond
		duration = 400 * time.Millisecond
		rate = 400 * time.Millisecond
	}

	orig := true
	origColorX := []float64{}
	origColorY := []float64{}
	origState := []bool{}
	origBrightness := []int{}

	lights := []*hue.Light{}

	n.lightsLock.Lock()
	for _, light := range n.lights {
		lights = append(lights, light)
		origColorX = append(origColorX, light.ColorX)
		origColorY = append(origColorY, light.ColorY)
		origState = append(origState, light.State)
		origBrightness = append(origBrightness, light.Brightness)
	}
	n.lightsLock.Unlock()

	var reset = func() {
		for i, light := range lights {
			light.SetState(origState[i])
			if origState[i] {
				light.SetBrightness(origBrightness[i])
				light.SetColorXY(origColorX[i], origColorY[i])
			}
			light.SetTransition(transition)

			err = light.Commit()
			if err != nil {
				return
			}
		}
		time.Sleep(transition + wait)
	}
	defer reset()

	for {
		if len(n.alerts) == 0 {
			if n.Mode == constants.Solid && !orig {
				orig = true
				reset()
			}
			time.Sleep(50 * time.Millisecond)

			if n.interrupt {
				return
			}

			continue
		} else if n.Mode == constants.Solid {
			orig = false
		}

		alerts := []Alert{}

		n.alertsLock.Lock()
		for _, alert := range n.alerts {
			alerts = append(alerts, alert)
		}
		n.alertsLock.Unlock()

		for _, light := range lights {
			err = light.Update()
			if err != nil {
				return
			}
		}

		for _, alert := range alerts {
			for _, light := range lights {
				light.SetState(true)
				light.SetBrightness(254)
				light.SetColorHex(alert.Color)
				light.SetTransition(transition)

				err = light.Commit()
				if err != nil {
					return
				}
			}
			time.Sleep(transition + wait + duration)
		}

		if n.Mode != constants.Solid {
			reset()
		}

		if n.interrupt {
			return
		}

		if n.Mode != constants.Solid {
			time.Sleep(rate)
		}

		if n.interrupt {
			return
		}
	}

	return
}

func (n *Notification) Run() {
	n.interruptWait = make(chan bool)

	go func() {
		defer func() {
			n.interruptWait <- true
		}()

		for {
			err := n.runner()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("notification: Error in runner")
			}

			if n.interrupt {
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func (n *Notification) Stop() {
	if n.interrupt {
		return
	}
	n.interrupt = true

	<-n.interruptWait
}
