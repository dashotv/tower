package app

import "time"

func (c *Connector) MessageList(page, limit int) ([]*Message, int64, error) {
	skip := (page - 1) * limit
	total, err := c.Message.Query().Count()
	if err != nil {
		return nil, 0, err
	}

	list, err := c.Message.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
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
