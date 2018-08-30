package flags

type Install struct {
	KubeConfig string `json:"kube_config"`
	Namespace  string `json:"namespace"`
	Host       string `json:"host"`
	Dir        string `json:"dir"`
	Values     string `json:"values"`

	Debug bool
}
