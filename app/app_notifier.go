package app

import "time"

var notifier *Notifier

type Notifier struct {
	Log    *NotifierLog
	Notice *NotifierNotice
}

type NotifierLog struct{}
type NotifierNotice struct{}

func (n *Notifier) notice(level, title, message string) {
	e := &EventTowerNotice{
		Time:    time.Now().String(),
		Event:   "notice",
		Class:   title,
		Level:   level,
		Message: message,
	}
	if err := events.Send("tower.notices", e); err != nil {
		log.Errorf("sending notice: %s", err)
	}
}
func (n *Notifier) log(level, title, message string) {
	l := &Message{
		Level:    level,
		Message:  message,
		Facility: title,
	}
	l.CreatedAt = time.Now()
	if err := db.Message.Save(l); err != nil {
		events.Log.Errorf("error saving log: %s", err)
	}
	events.Send("tower.logs", &EventTowerLog{Event: "new", ID: l.ID.Hex(), Log: l})
}

func (n *Notifier) Notify(level, title, message string) {
	n.notice(level, title, message)
	n.log(level, title, message)
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

func (n *NotifierLog) Info(title, message string) {
	notifier.log("info", title, message)
}
func (n *NotifierLog) Warn(title, message string) {
	notifier.log("Warn", title, message)
}
func (n *NotifierLog) Error(title, message string) {
	notifier.log("error", title, message)
}
func (n *NotifierLog) Success(title, message string) {
	notifier.log("success", title, message)
}

func (n *NotifierNotice) Info(title, message string) {
	notifier.notice("info", title, message)
}
func (n *NotifierNotice) Warn(title, message string) {
	notifier.notice("Warn", title, message)
}
func (n *NotifierNotice) Error(title, message string) {
	notifier.notice("error", title, message)
}
func (n *NotifierNotice) Success(title, message string) {
	notifier.notice("success", title, message)
}

func setupNotifier() error {
	notifier = &Notifier{
		Log:    &NotifierLog{},
		Notice: &NotifierNotice{},
	}
	return nil
}
