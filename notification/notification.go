package notification

import (
	"github.com/Sirupsen/logrus"
	"github.com/zachhuff386/hue-alert/hue"
	"sync"
	"time"
)

var (
	wait = 200 * time.Millisecond
)

type Alert struct {
	Type     string
	Color    string
	Duration time.Duration
}

type Notification struct {
	interrupt     bool
	interruptWait chan bool
	alerts        map[string]Alert
	alertsLock    sync.Mutex
	lights        []*hue.Light
	lightsLock    sync.Mutex
	Transition    time.Duration
	Rate          time.Duration
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
	origColorX := []float64{}
	origColorY := []float64{}
	origState := []bool{}
	origBrightness := []int{}

	for _, light := range n.lights {
		origColorX = append(origColorX, light.ColorX)
		origColorY = append(origColorY, light.ColorY)
		origState = append(origState, light.State)
		origBrightness = append(origBrightness, light.Brightness)
	}

	for {
		if len(n.alerts) == 0 {
			time.Sleep(50 * time.Millisecond)
		}

		alerts := []Alert{}
		lights := []*hue.Light{}

		n.alertsLock.Lock()
		for _, alert := range n.alerts {
			alerts = append(alerts, alert)
		}
		n.alertsLock.Unlock()

		n.lightsLock.Lock()
		for _, light := range n.lights {
			lights = append(lights, light)
		}
		n.lightsLock.Unlock()

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
				light.SetTransition(n.Transition)

				err = light.Commit()
				if err != nil {
					return
				}
			}
			time.Sleep(n.Transition + wait + alert.Duration)
		}

		for i, light := range lights {
			light.SetState(origState[i])
			light.SetBrightness(origBrightness[i])
			light.SetColorXY(origColorX[i], origColorY[i])
			light.SetTransition(n.Transition)

			err = light.Commit()
			if err != nil {
				return
			}
		}

		if n.interrupt {
			return
		}

		time.Sleep(n.Rate)

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
