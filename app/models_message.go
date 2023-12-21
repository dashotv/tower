package app

import "time"

func (c *Connector) MessageList() ([]*Message, error) {
	list, err := c.Message.Query().Limit(10).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *Connector) MessageCreate(level string, message string, facility string, t time.Time) (*Message, error) {
	l := &Message{
		Level:    level,
		Message:  message,
		Facility: facility,
	}
	l.CreatedAt = time.Now()
	return l, c.Message.Save(l)
}
