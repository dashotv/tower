package app

//func (c *Connector) Upcoming() ([]Episode, error) {
//	//lookup := make(map[string]Series)
//	//q := c.Series.Query()
//	//series, err := q.
//	//	Where("_type", "Series").
//	//	Where("active", true).
//	//	Limit(1000).
//	//	Run()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//ids := make([]primitive.ObjectID, len(series))
//	//for _, s := range series {
//	//	ids = append(ids, s.ID)
//	//	//lookup[s.ID.String()] = s
//	//}
//
//	//q2 := c.Episode.Query()
//	//now := time.Now()
//	//since := time.Now().Add(-time.Hour * 24 * 21)
//	//fmt.Println("time between ", since, " and ", now)
//	//return q2.Where("_type", "Episode").
//	//	//In("series_id", ids).
//	//	Where("downloaded", false).
//	//	Where("completed", false).
//	//	Where("skipped", false).
//	//	//In("missing", []interface{}{false, nil}).
//	//	//GreaterThan("season_number", 0).
//	//	//GreaterThan("episode_number", 0).
//	//	LessThanEqual("release_date", now).
//	//	GreaterThanEqual("release_date", since).
//	//	Asc("release_date").Asc("season_number").Asc("episode_number").
//	//	Limit(5).
//	//	Run()
//
//}

//func (s *MediumStore) Upcoming() ([]Medium, error) {
//	lookup := make(map[string]Medium)
//	q := s.Query()
//	series, err := q.
//		Where("_type", "Series").
//		Where("active", true).
//		Limit(1000).
//		Run()
//	if err != nil {
//		return nil, err
//	}
//
//	ids := make([]primitive.ObjectID, len(series))
//	for _, s := range series {
//		ids = append(ids, s.ID)
//		lookup[s.ID.String()] = s
//	}
//
//	q2 := s.Query()
//	eps, err := q2.Where("_type", "Episode").
//		In("series_id", ids).
//		Where("downloaded", false).
//		Where("completed", false).
//		Where("skipped", false).
//		In("missing", []interface{}{false, nil}).
//		GreaterThan("season_number", 0).
//		GreaterThan("episode_number", 0).
//		LessThanEqual("release_date", time.Now()).
//		Asc("release_date").Asc("season_number").Asc("episode_number").
//		Limit(5).
//		Run()
//	if err != nil {
//		return nil, err
//	}
//
//	for _, ep := range eps {
//		ep.Series = lookup[ep.SeriesID]
//	}
//}

func (m *Medium) Background() string {
	for _, p := range m.Paths {
		if p.Type == "background" {
			return p.Local
		}
	}
	return ""
}

func (m *Medium) Cover() string {
	for _, p := range m.Paths {
		if p.Type == "cover" {
			return p.Local
		}
	}
	return ""
}
