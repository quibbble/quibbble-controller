package controller

type GameServerConfig struct {
	Host  string      `yaml:"host"`
	Image ImageConfig `yaml:"image"`
}

type ImageConfig struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}
