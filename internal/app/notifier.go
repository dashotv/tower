package app

import (
	"fmt"
	"time"
)

var notifier *Notifier

func init() {
	initializers = append(initializers, setupNotifier)
}

func setupNotifier(a *Application) error {
	notifier = &Notifier{
		Log:    &NotifierLog{},
		Notice: &NotifierNotice{},
	}
	return nil
}

type Notifier struct {
	Log    *NotifierLog
	Notice *NotifierNotice
}

type NotifierLog struct{}
type NotifierNotice struct{}

func (n *Notifier) notice(level, title, message string) {
	e := &EventNotices{
		Time:    time.Now().String(),
		Event:   "notice",
		Class:   title,
		Level:   level,
		Message: message,
	}
	if err := app.Events.Send("tower.notices", e); err != nil {
		app.Log.Errorf("sending notice: %s", err)
	}
}
func (n *Notifier) log(level, title, message string) {
	l := &Message{
		Level:    level,
		Message:  message,
		Facility: title,
	}
	l.CreatedAt = time.Now()
	if err := app.DB.Message.Save(l); err != nil {
		app.Events.Log.Errorf("error saving log: %s", err)
	}
	app.Log.Named("notifier").Debugf("[%s] %s: %s", level, title, message)
	if err := app.Events.Send("tower.logs", &EventLogs{Event: "new", ID: l.ID.Hex(), Log: l}); err != nil {
		app.Events.Log.Errorf("error sending log: %s", err)
	}
}

func (n *Notifier) Notify(level, title, message string) {
	n.notice(level, title, message)
	n.log(level, title, message)
}
func (n *Notifier) Debug(title, message string) {
	n.Notify("debug", title, message)
}
func (n *Notifier) Info(title, message string) {
	n.Notify("info", title, message)
}
func (n *Notifier) Warn(title, message string) {
	n.Notify("Warn", title, message)
}
func (n *Notifier) Error(title, message string) {
	n.Notify("error", title, message)
}
func (n *Notifier) Success(title, message string) {
	n.Notify("success", title, message)
}

func (n *NotifierLog) Debug(title, message string) {
	notifier.log("debug", title, message)
}
func (n *NotifierLog) Info(title, message string) {
	notifier.log("info", title, message)
}
func (n *NotifierLog) Warn(title, message string) {
	notifier.log("warning", title, message)
}
func (n *NotifierLog) Error(title, message string) {
	notifier.log("error", title, message)
}
func (n *NotifierLog) Success(title, message string) {
	notifier.log("success", title, message)
}
func (n *NotifierLog) Debugf(title, format string, a ...any) {
	notifier.log("debug", title, fmt.Sprintf(format, a...))
}
func (n *NotifierLog) Infof(title, format string, a ...any) {
	notifier.log("info", title, fmt.Sprintf(format, a...))
}
func (n *NotifierLog) Warnf(title, format string, a ...any) {
	notifier.log("warn", title, fmt.Sprintf(format, a...))
}
func (n *NotifierLog) Errorf(title, format string, a ...any) {
	notifier.log("error", title, fmt.Sprintf(format, a...))
}
func (n *NotifierLog) Successf(title, format string, a ...any) {
	notifier.log("success", title, fmt.Sprintf(format, a...))
}

func (n *NotifierNotice) Info(title, message string) {
	notifier.notice("info", title, message)
}
func (n *NotifierNotice) Warn(title, message string) {
	notifier.notice("warn", title, message)
}
func (n *NotifierNotice) Error(title, message string) {
	notifier.notice("error", title, message)
}
func (n *NotifierNotice) Success(title, message string) {
	notifier.notice("success", title, message)
}
