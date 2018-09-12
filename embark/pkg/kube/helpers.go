package kube

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/thoas/go-funk"

	"github.com/mitchellh/mapstructure"

	"github.com/containerum/containerum/embark/pkg/utils/ym"
	"github.com/ghodss/yaml"
)

var (
	_ KubectlConfigProvider = StdKubectConfig
)

func StdKubectConfig() (KubectlConfig, error) {
	var config KubectlConfig
	return config, ym.LoadYAML(autoKubectlConfigPath(), &config)
}

func LoadKubectlConfigFromPath(configPath string) KubectlConfigProvider {
	return func() (KubectlConfig, error) {
		var config KubectlConfig
		return config, ym.LoadYAML(configPath, &config)
	}
}

func KubectlConfigFromReader(re io.Reader) KubectlConfigProvider {
	var config KubectlConfig
	var buf = &bytes.Buffer{}
	if _, err := buf.ReadFrom(re); err != nil {
		return config.AsProviderWithErr(err)
	}
	return config.AsProviderWithErr(yaml.Unmarshal(buf.Bytes(), &config))
}

func DecodeConfig(data []byte) (KubectlConfig, error) {
	var tree map[string]interface{}
	if err := yaml.Unmarshal(data, &tree); err != nil {
		return KubectlConfig{}, err
	}
	var config KubectlConfig
	var meta = &mapstructure.Metadata{}
	if err := mapstructure.WeakDecodeMetadata(tree, &config, meta); err != nil {
		return KubectlConfig{}, err
	}
	switch cert := funk.Get(tree, "clusters").(type) {
	case nil:
		// pass
	case string:
		var certData, err = base64.StdEncoding.DecodeString(cert)
		if err != nil {
			return KubectlConfig{}, err
		}
		config.Clusters[0].Cluster.CertificateAuthorityData = certData
	default:
		fmt.Printf("%T\n%v", cert, cert)
	}

	return config, nil
}

func autoKubectlConfigPath() string {
	var configPathFromEnv, configPathFromEnvExists = os.LookupEnv("KUBECONFIG")
	if configPathFromEnvExists {
		return configPathFromEnv
	}
	return os.ExpandEnv("$HOME/.kube/config")
}
