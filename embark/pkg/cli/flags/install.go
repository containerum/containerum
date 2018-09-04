package flags

import "github.com/go-playground/validator"

type Install struct {
	KubeConfig string `json:"kube_config"`
	Namespace  string `json:"namespace" validate:"isdefault|alphanum"`
	Dir        string `json:"dir" validate:"required"`
	Values     string `json:"values"`
	Debug      bool
}

func (install Install) Validate() error {
	return validator.New().Struct(install)
}
