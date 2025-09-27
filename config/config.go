package config

// Config holds all controller configurations.
type Config struct {
	// The namespace to deploy the controller and games into
	Namespace string `yaml:"namespace"`
	// Database configs for game storage.
	Persistence *PersistConfig `yaml:"persistence,omitempty"`
	// Game configurations.
	Game *GameConfig `yaml:"game"`
}
