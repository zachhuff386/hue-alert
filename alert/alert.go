package alert

import (
	"github.com/Sirupsen/logrus"
	"github.com/zachhuff386/hue-alert/accounts"
	"github.com/zachhuff386/hue-alert/config"
	"github.com/zachhuff386/hue-alert/hue"
	"github.com/zachhuff386/hue-alert/notification"
	"time"
)

type Alert struct {
	interrupt  bool
	notf       *notification.Notification
	Hue        *hue.Hue
	Lights     []string
	Brightness int
	Rate       time.Duration
	Mode       string
}

func (a *Alert) runner() (err error) {
	lights, err := a.Hue.GetLightsById(a.Lights)
	if err != nil {
		return
	}

	a.notf = &notification.Notification{
		Mode:       a.Mode,
		Brightness: a.Brightness,
	}
	defer func() {
		if !a.interrupt && a.notf != nil {
			a.notf.Stop()
		}
	}()

	for _, light := range lights {
		a.notf.AddLight(light)
	}

	go a.notf.Run()

	accts, err := accounts.GetAccounts()
	if err != nil {
		return
	}

	alerts := map[string]notification.Alert{}

	for _, acct := range accts {
		client, e := acct.GetClient()
		if e != nil {
			err = e
			return
		}

		err = client.Update()
		if err != nil {
			return
		}

		err = config.Config.CommitAccount(acct)
		if err != nil {
			return
		}
	}

	for {
		for _, acct := range accts {
			client, e := acct.GetClient()
			if e != nil {
				err = e
				return
			}

			err = client.Sync()
			if err != nil {
				acct.Alert = false
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("alert: Error checking account")
			}

			err = config.Config.CommitAccount(acct)
			if err != nil {
				return
			}

			if acct.Alert {
				alrt := notification.Alert{
					Type:  acct.Type,
					Color: acct.GetColor(),
				}
				alerts[acct.Id] = alrt
				a.notf.AddAlert(alrt)
			} else {
				alrt, ok := alerts[acct.Id]
				if ok {
					a.notf.RemoveAlert(alrt)
					delete(alerts, acct.Id)
				}
			}
		}

		start := time.Now()
		for {
			time.Sleep(50 * time.Millisecond)
			if a.interrupt {
				return
			}
			if time.Since(start) >= a.Rate {
				break
			}
		}
	}
}

func (a *Alert) Run() {
	for {
		err := a.runner()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("alert: Error in runner")
		}

		if a.interrupt {
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func (a *Alert) Stop() {
	a.interrupt = true
	if a.notf != nil {
		a.notf.Stop()
	}
}
