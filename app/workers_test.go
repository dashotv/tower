package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCleanLogs(t *testing.T) {
	funcs := []func() error{
		setupConfig,
		setupLogger,
		setupDb,
	}
	for _, f := range funcs {
		err := f()
		require.NoError(t, err)
	}

	before, err := db.Message.Count(bson.M{})
	assert.NoError(t, err)

	list := []struct {
		Message   string
		CreatedAt time.Time
	}{
		{Message: "old", CreatedAt: time.Now().UTC().AddDate(0, 0, -6)},
		{Message: "older", CreatedAt: time.Now().UTC().AddDate(0, 0, -7)},
	}
	for _, v := range list {
		m := &Message{Message: v.Message}
		err := db.Message.Save(m)
		require.NoError(t, err)
		m.CreatedAt = v.CreatedAt
		err = db.Message.Save(m)
		require.NoError(t, err)
	}

	err = CleanLogs(nil)
	assert.NoError(t, err)

	count, err := db.Message.Count(bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, before, count)

	// 	// check that only the old messages were deleted
	// 	count, err := collection.CountDocuments(context.Background(), bson.M{})
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, int64(len(newMessages)), count)
}
