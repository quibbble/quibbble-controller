package controller

import "github.com/quibbble/quibbble-controller/pkg/k8s"

type GameServerConfig struct {
	Host    string            `yaml:"host"`
	Image   ImageConfig       `yaml:"image"`
	Ingress k8s.IngressConfig `yaml:"ingress"`
}

type ImageConfig struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}
