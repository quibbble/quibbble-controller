package controller

import (
	"github.com/quibbble/quibbble-controller/pkg/k8s"
)

type ServerConfig struct {
	FullName  string            `yaml:"fullname"`
	Namespace string            `yaml:"namespace"`
	Pod       k8s.PodConfig     `yaml:"pod"`
	Ingress   k8s.IngressConfig `yaml:"ingress"`
}
