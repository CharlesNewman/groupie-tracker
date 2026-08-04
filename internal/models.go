package internal

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationsResponse struct {
	Index []Relation `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type LocationResponse struct {
	Index []Location `json:"index"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DateResponse struct {
	Index []Date `json:"index"`
}

type PageData struct {
	Artists      []Artist
	CurrentPage  int
	PreviousPage int
	NextPage     int
	HasPrevious  bool
	HasNext      bool
	LocationList []string
}

type Find struct {
	Artist        Artist   `json:"artist"`
	Relation      Relation `json:"relation"`
	LocationCount int      `json:"locationCount"`
	ConcertCount  int      `json:"concertCount"`
}

type FilterData struct {
	Artists   []Artist
	Locations []Location
}
