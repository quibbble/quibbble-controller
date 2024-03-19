package game

// Information provides additional info about a game
type Information struct {
	// Key is the unique key that differentiates this game from others - required
	Key string `json:"key"`

	// Min and Max represents the min and max teams allowed when creating a new game - required
	Min int `json:"min"`
	Max int `json:"max"`

	// Details allows for additional game specific info - optional
	Details interface{} `json:"details,omitempty"`
}

// Action represents an action that is performed on the game state
type Action struct {
	// Team is the team performing the action - required
	Team string `json:"team"`

	// Type is the key that determines what action to perform - required
	Type string `json:"type"`

	// Details allows for additional info to be passed - optional
	Details interface{} `json:"details,omitempty"`
}

// Snapshot represents the current state of the game that will be viewed by a player
type Snapshot struct {
	// Turn is the turn of the current team - required
	Turn string `json:"turn"`

	// Teams is a list of all teams playing the game - required
	Teams []string `json:"teams"`

	// Winners is a list of teams that have won the game - required
	Winners []string `json:"winners"`

	// Details allows for additional info to be returned that is unique to each game - optional
	Details interface{} `json:"details,omitempty"`

	// Actions are a list of all valid Actions that can be performed on the current game state - optional
	Actions []*Action `json:"actions,omitempty"`

	// History is a list of past game actions that have lead to the current game state - optional
	History []*Action `json:"history,omitempty"`

	// Message provides players with text info about what to do next - optional
	Message string `json:"message,omitempty"`
}
