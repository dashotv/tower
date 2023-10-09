package app

import "go.mongodb.org/mongo-driver/bson/primitive"

func (c *Connector) MediumWatched(id primitive.ObjectID) bool {
	// TODO: add user name to config
	watches, _ := db.Watch.Query().Where("medium_id", id).Where("username", "xenonsoul").Run()
	return len(watches) > 0
}
