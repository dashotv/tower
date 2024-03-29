package app

import "github.com/dashotv/fae"

func (c *Connector) CollectionGet(id string) (*Collection, error) {
	collection, err := c.Collection.Get(id, &Collection{})
	if err != nil {
		return nil, err
	}

	// if err := c.processCollections([]*Collection{collection}); err != nil {
	// 	return nil, err
	// }

	return collection, nil
}

func (c *Connector) CollectionList(limit, skip int) ([]*Collection, error) {
	list, err := c.Collection.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, fae.Wrap(err, "query failed")
	}

	// if err := c.processCollections(list); err != nil {
	// 	return nil, fae.Wrap(err, "process collections failed")
	// }

	return list, nil
}

// func (c *Connector) processCollections(collections []*Collection) error {
// 	for _, collection := range collections {
// 		for _, cm := range collection.Media {
// 			m, err := c.Medium.Get(cm.MediumId.Hex(), &Medium{})
// 			if err != nil {
// 				c.Log.Warnf("failed to get medium %s: %s", cm.MediumId.Hex(), err)
// 				continue
// 			}
// 			cm.Medium = m
// 		}
// 	}
//
// 	return nil
// }
