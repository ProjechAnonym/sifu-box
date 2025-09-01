package application

import (
	"go.uber.org/zap"
)

type SignalHook struct {
	PID  int
	Cron bool
}
type Signal struct {
	Operation int
	Cron      bool
}

func HookHandle(hook_chan *chan SignalHook, cron_chan *chan bool, web_chan *chan bool, logger *zap.Logger) {
	for hook_res := range *hook_chan {
		switch hook_res.Cron {
		case true:
			if hook_res.PID > 0 {
				*cron_chan <- true
			} else {
				*cron_chan <- false
			}
		case false:
			if hook_res.PID > 0 {
				*web_chan <- true
			} else {
				*web_chan <- false
			}
		}
	}
}
