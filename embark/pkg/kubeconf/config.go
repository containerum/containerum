package kubeconf

import (
	"encoding/base64"
	"fmt"

	"github.com/ericchiang/k8s"
)

var (
	b64 = base64.StdEncoding
)

// diffs:
//		- AuthInfo.ClientCertificateData is a base64 string
//		- AuthInfo.ClientKeyData is a base64 string
//		- CLuster.CertificateAuthorityData is a base64 string
// Top level config objects and all values required for proper functioning are not "omitempty".  Any truly optional piece of config is allowed to be omitted.

// Config is adapter for kubect config, convertable to github.com/ericchiang/k8s.Config
// Config holds the information needed to build connect to remote kubernetes clusters as a given user
type Config struct {
	// Legacy field from pkg/api/types.go TypeMeta.
	// TODO(jlowdermilk): remove this after eliminating downstream dependencies.
	// +optional
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
	// DEPRECATED: APIVersion is the preferred api version for communicating with the kubernetes cluster (v1, v2, etc).
	// Because a cluster can run multiple API groups and potentially multiple versions of each, it no longer makes sense to specify
	// a single value for the cluster version.
	// This field isn't really needed anyway, so we are deprecating it without replacement.
	// It will be ignored if it is present.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	// Preferences holds general information to be use for cli interactions
	Preferences k8s.Preferences `json:"preferences" yaml:"preferences"`
	// Clusters is a map of referencable names to cluster configs
	Clusters []NamedCluster `json:"clusters" yaml:"clusters"`
	// AuthInfos is a map of referencable names to user configs
	AuthInfos []NamedAuthInfo `json:"users" yaml:"users"`
	// Contexts is a map of referencable names to context configs
	Contexts []k8s.NamedContext `json:"contexts" yaml:"contexts"`
	// CurrentContext is the name of the context that you would like to use by default
	CurrentContext string `json:"current-context" yaml:"current-context"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions []k8s.NamedExtension `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (config Config) ToK8S() (k8s.Config, error) {
	var void = k8s.Config{}
	var k8sConfig = k8s.Config{
		Kind:           config.Kind,
		APIVersion:     config.APIVersion,
		Preferences:    config.Preferences,
		Contexts:       append([]k8s.NamedContext{}, config.Contexts...),
		CurrentContext: config.CurrentContext,
		Extensions:     append([]k8s.NamedExtension{}, config.Extensions...),
	}
	for _, cluster := range config.Clusters {
		var k8scluster, k8sclusterErr = cluster.Cluster.ToK8S()
		if k8sclusterErr != nil {
			return void, k8sclusterErr
		}
		k8sConfig.Clusters = append(k8sConfig.Clusters, k8s.NamedCluster{
			Name:    cluster.Name,
			Cluster: k8scluster,
		})
	}
	for _, authInfo := range config.AuthInfos {
		var k8sauthInfo, k8sAuthInfoConversionErr = authInfo.AuthInfo.ToK8S()
		if k8sAuthInfoConversionErr != nil {
			return void, k8sAuthInfoConversionErr
		}
		k8sConfig.AuthInfos = append(k8sConfig.AuthInfos, k8s.NamedAuthInfo{
			Name:     authInfo.Name,
			AuthInfo: k8sauthInfo,
		})
	}
	return k8sConfig, nil
}

// Cluster contains information about how to communicate with a kubernetes cluster
type Cluster struct {
	// Server is the address of the kubernetes cluster (https://hostname:port).
	Server string `json:"server" yaml:"server"`
	// APIVersion is the preferred api version for communicating with the kubernetes cluster (v1, v2, etc).
	// +optional
	APIVersion string `json:"api-version,omitempty" yaml:"api-version,omitempty"`
	// InsecureSkipTLSVerify skips the validity check for the server's certificate. This will make your HTTPS connections insecure.
	// +optional
	InsecureSkipTLSVerify bool `json:"insecure-skip-tls-verify,omitempty" yaml:"insecure-skip-tls-verify,omitempty"`
	// CertificateAuthority is the path to a cert file for the certificate authority.
	// +optional
	CertificateAuthority string `json:"certificate-authority,omitempty" yaml:"certificate-authority,omitempty"`
	// CertificateAuthorityData contains PEM-encoded certificate authority certificates. Overrides CertificateAuthority
	// +optional
	CertificateAuthorityData string `json:"certificate-authority-data,omitempty" yaml:"certificate-authority-data,omitempty"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions []k8s.NamedExtension `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (cluster Cluster) ToK8S() (k8s.Cluster, error) {
	var void = k8s.Cluster{}
	var decodedCertAuthData, decodeCertAuthDataErr = b64.DecodeString(cluster.CertificateAuthorityData)
	if decodeCertAuthDataErr != nil {
		return void, fmt.Errorf("unable to decode certificate authority data: %v", decodeCertAuthDataErr)
	}
	return k8s.Cluster{
		Server:                   cluster.Server,
		APIVersion:               cluster.APIVersion,
		InsecureSkipTLSVerify:    cluster.InsecureSkipTLSVerify,
		CertificateAuthorityData: decodedCertAuthData,
		CertificateAuthority:     cluster.CertificateAuthority,
		Extensions:               append([]k8s.NamedExtension{}, cluster.Extensions...),
	}, nil
}

// AuthInfo contains information that describes identity information.  This is use to tell the kubernetes cluster who you are.
type AuthInfo struct {
	// ClientCertificate is the path to a client cert file for TLS.
	// +optional
	ClientCertificate string `json:"client-certificate,omitempty" yaml:"client-certificate,omitempty"`
	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS. Overrides ClientCertificate
	// +optional
	ClientCertificateData string `json:"client-certificate-data,omitempty" yaml:"client-certificate-data,omitempty"`
	// ClientKey is the path to a client key file for TLS.
	// +optional
	ClientKey string `json:"client-key,omitempty" yaml:"client-key,omitempty"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS. Overrides ClientKey
	// +optional
	ClientKeyData string `json:"client-key-data,omitempty" yaml:"client-key-data,omitempty"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	// +optional
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
	// TokenFile is a pointer to a file that contains a bearer token (as described above).  If both Token and TokenFile are present, Token takes precedence.
	// +optional
	TokenFile string `json:"tokenFile,omitempty" yaml:"tokenFile,omitempty"`
	// Impersonate is the username to imperonate.  The name matches the flag.
	// +optional
	Impersonate string `json:"as,omitempty" yaml:"as,omitempty"`
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	// Password is the password for basic authentication to the kubernetes cluster.
	// +optional
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// AuthProvider specifies a custom authentication plugin for the kubernetes cluster.
	// +optional
	AuthProvider *k8s.AuthProviderConfig `json:"auth-provider,omitempty" yaml:"auth-provider,omitempty"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions []k8s.NamedExtension `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (authInfo AuthInfo) ToK8S() (k8s.AuthInfo, error) {
	var void = k8s.AuthInfo{}
	var decodedCertData, decodeCertDataErr = b64.DecodeString(authInfo.ClientCertificateData)
	if decodeCertDataErr != nil {
		return void, fmt.Errorf("unable to decode client certificate data: %v", decodeCertDataErr)
	}
	var decodedKeyData, decodeKeyDataErr = b64.DecodeString(authInfo.ClientKeyData)
	if decodeCertDataErr != nil {
		return void, fmt.Errorf("unable to decode client key data: %v", decodeKeyDataErr)
	}
	return k8s.AuthInfo{
		ClientCertificate:     authInfo.ClientCertificate,
		ClientCertificateData: decodedCertData,
		ClientKey:             authInfo.ClientKey,
		ClientKeyData:         decodedKeyData,
		Token:                 authInfo.Token,
		TokenFile:             authInfo.TokenFile,
		Impersonate:           authInfo.Impersonate,
		Username:              authInfo.Username,
		Password:              authInfo.Password,
		AuthProvider:          authInfo.AuthProvider,
		Extensions:            append([]k8s.NamedExtension{}, authInfo.Extensions...),
	}, nil
}

// NamedCluster relates nicknames to cluster information
type NamedCluster struct {
	// Name is the nickname for this Cluster
	Name string `json:"name" yaml:"name"`
	// Cluster holds the cluster information
	Cluster Cluster `json:"cluster" yaml:"cluster"`
}

// NamedAuthInfo relates nicknames to auth information
type NamedAuthInfo struct {
	// Name is the nickname for this AuthInfo
	Name string `json:"name" yaml:"name"`
	// AuthInfo holds the auth information
	AuthInfo AuthInfo `json:"user" yaml:"user"`
}
