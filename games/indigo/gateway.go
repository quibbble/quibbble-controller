package indigo

type gateway struct {
	Locations [3][2]int
	Edges     string
	Teams     []string
}

func newGateway(locations [3][2]int, edges string, teams ...string) *gateway {
	return &gateway{
		Locations: locations,
		Edges:     edges,
		Teams:     teams,
	}
}
