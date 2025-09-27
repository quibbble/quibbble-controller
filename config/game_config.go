package config

import corev1 "k8s.io/api/core/v1"

type GameConfig struct {
	SDKServer  SDKServerConfig  `yaml:"sdkServer"`
	GameServer GameServerConfig `yaml:"gameServer"`
	Gateway    GatewayConfig    `yaml:"gateway"`

	NodeSelector map[string]string   `yaml:"nodeSelector"`
	Affinity     *corev1.Affinity    `yaml:"affinity"`
	Tolerations  []corev1.Toleration `yaml:"tolerations"`
}

type SDKServerConfig struct {
	// ControllerID used by the controller to authenticate with sdk servers.
	ControllerID string                       `yaml:"controllerId"`
	Port         string                       `yaml:"port"`
	Image        Image                        `yaml:"image"`
	Resources    *corev1.ResourceRequirements `yaml:"resources"`
}

type GameServerConfig struct {
	Port                string                       `yaml:"port"`
	RepositoryAllowList []string                     `yaml:"repositoryAllowList"`
	Resources           *corev1.ResourceRequirements `yaml:"resources"`
}

type GatewayConfig struct {
	Enabled bool `yaml:"enabled"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}
