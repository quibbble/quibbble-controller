package indigo

type gateway struct {
	Locations [3][2]int `json:"locations"`
	Edges     string    `json:"edges"`
	Teams     []string  `json:"teams"`
}

func newGateway(locations [3][2]int, edges string, teams ...string) *gateway {
	return &gateway{
		Locations: locations,
		Edges:     edges,
		Teams:     teams,
	}
}
