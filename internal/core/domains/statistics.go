package domains

type Statistics struct {
	TotalPrice int
}

func (s *Statistics) CalcStatistics(subs []Sub) {
	for _, sub := range subs {
		s.TotalPrice += sub.Price
	}
}
