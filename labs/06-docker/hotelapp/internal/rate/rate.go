package rate

import (
	"sort"
	"time"
)

// Rate implements the rate service
type Rate struct {
	dbsession *DatabaseSession
}

type RateRequest struct {
	HotelIds []string
	InDate   string
	OutDate  string
}

type RateResult struct {
	RatePlans []*RatePlan
}

// NewRate returns a new server
func NewRate(db *DatabaseSession) *Rate {
	return &Rate{
		dbsession: db,
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return (check.Equal(start) || check.After(start)) && (check.Equal(end) || check.Before(end))
}

// GetRates gets rates for hotels for specific date range.
func (s *Rate) GetRates(req *RateRequest) (*RateResult, error) {
	res := new(RateResult)

	ratePlans, err := s.dbsession.GetRates(req.HotelIds)
	if err != nil {
		return nil, err
	}
	finalRatePlans := make(RatePlans, 0)

	start, _ := time.Parse("2006-01-02", req.InDate)
	end, _ := time.Parse("2006-01-02", req.OutDate)

	sort.Sort(ratePlans)
	for _, rateplan := range ratePlans {
		in, _ := time.Parse("2006-01-02", rateplan.InDate)
		out, _ := time.Parse("2006-01-02", rateplan.OutDate)
		if inTimeSpan(in, out, start) && inTimeSpan(in, out, end) {
			finalRatePlans = append(finalRatePlans, rateplan)
		}
	}

	res.RatePlans = finalRatePlans

	return res, nil
}

type RatePlans []*RatePlan

func (r RatePlans) Len() int {
	return len(r)
}

func (r RatePlans) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
	return r[i].RoomType.TotalRate > r[j].RoomType.TotalRate
}
