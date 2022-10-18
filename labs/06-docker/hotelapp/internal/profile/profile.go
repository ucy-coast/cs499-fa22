package profile

// Profile implements the profile service
type Profile struct {
	dbsession *DatabaseSession
}

type Request struct {
	HotelIds []string
	Locale   string
}

type Result struct {
	Hotels []*Hotel
}

// NewProfile returns a new Profile service
func NewProfile(db *DatabaseSession) *Profile {
	return &Profile{
		dbsession: db,
	}
}

// GetProfiles returns hotel profiles for requested IDs
func (s *Profile) GetProfiles(req *Request) (*Result, error) {
	var err error
	res := new(Result)
	res.Hotels, err = s.dbsession.GetProfiles(req.HotelIds)
	return res, err
}
