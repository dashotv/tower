package app

import "time"

func (c *Connector) MessageList(page, limit int) ([]*Message, error) {
	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit
	list, err := c.Message.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
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
