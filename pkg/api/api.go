package api

// API represents the footloose REST API.
type API struct {
	BaseURI string
	db      db
}

// New creates a new object able to answer footloose REST API.
func New(baseURI string) *API {
	api := &API{
		BaseURI: baseURI,
	}
	api.db.init()
	return api
}
