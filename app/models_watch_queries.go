package app

import "go.mongodb.org/mongo-driver/bson/primitive"

func (c *Connector) MediumWatched(id primitive.ObjectID) bool {
	// TODO: add user name to config
	watches, _ := App().DB.Watch.Query().Where("medium_id", id).Where("username", "xenonsoul").Run()

	if len(watches) > 0 {
		return true
	}

	return false
}
