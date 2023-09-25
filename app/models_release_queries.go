package app

import (
	"context"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *Connector) ReleaseSetting(id, setting string, value bool) error {
	release := &Release{}
	err := c.Release.Find(id, release)
	if err != nil {
		return err
	}

	switch setting {
	case "verified":
		release.Verified = value
	}

	return c.Release.Update(release)
}

func (c *Connector) ReleasesPopular(date time.Time, count int) ([]*Popular, error) {
	return ReleasesPopular(c.Release.Collection, date, count)
}

type Popular struct {
	Name  string `json:"name" bson:"_id"`
	Type  string `json:"type" bson:"type"`
	Count int    `json:"count" bson:"count"`
}

func ReleasesPopular(coll *mgm.Collection, date time.Time, count int) ([]*Popular, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	p := []bson.M{
		bson.M{
			"$project": bson.M{"name": 1, "type": 1, "published": "$published_at"},
		},
		bson.M{
			"$match": bson.M{"type": "tv", "published": bson.M{"$gte": date}},
		},
		bson.M{
			"$group": bson.M{"_id": "$name", "type": bson.M{"$first": "$type"}, "count": bson.M{"$sum": 1}},
		},
		bson.M{
			"$sort": bson.M{"count": -1},
		},
		bson.M{"$limit": count},
	}

	cursor, err := coll.Aggregate(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "aggregating popular releases")
	}

	results := make([]*Popular, count)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.Wrap(err, "decoding popular releases")
	}

	return results, nil
}
