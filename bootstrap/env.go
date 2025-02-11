package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type OIDC struct {
	IssuerURI        string `mapstructure:"issuerURI" json:"issuerURI"`
	SkipTLSVerify    bool   `mapstructure:"skipTLSVerify" json:"skipTLSVerify"`
	JWKURI           string `mapstructure:"jwkURI" json:"jwkURI"`
	PublicKey        string `mapstructure:"publicKey" json:"publicKey"`
	ClientID         string `mapstructure:"clientID" json:"clientID"`
	Audience         string `mapstructure:"audience" json:"audience"`
	UsernameClaim    string `mapstructure:"usernameClaim" json:"usernameClaim"`
	GroupsClaim      string `mapstructure:"groupsClaim" json:"groupsClaim"`
	RolesClaim       string `mapstructure:"rolesClaim" json:"rolesClaim"`
	ExtraQueryParams string `mapstructure:"extraQueryParams" json:"extraQueryParams"`
}

type Security struct {
	CORSAllowedOrigins []string `mapstructure:"corsAllowedOrigins" json:"corsAllowedOrigins"`
}

type K8SPublicEndpoint struct {
	OidcConfiguration struct {
		IssuerURI string `mapstructure:"issuerURI" json:"issuerURI"`
		ClientID  string `mapstructure:"clientID" json:"clientID"`
	} `mapstructure:"oidcConfiguration" json:"oidcConfiguration"`
	URL string `mapstructure:"URL" json:"URL"`
}

type Quota struct {
	RequestsMemory           string `mapstructure:"requests.memory" json:"requests.memory"`
	RequestsCPU              string `mapstructure:"requests.cpu" json:"requests.cpu"`
	LimitsMemory             string `mapstructure:"limits.memory" json:"limits.memory"`
	LimitsCPU                string `mapstructure:"limits.cpu" json:"limits.cpu"`
	RequestsStorage          string `mapstructure:"requests.storage" json:"requests.storage"`
	CountPods                string `mapstructure:"count/pods" json:"count/pods"`
	RequestsEphemeralStorage string `mapstructure:"requests.ephemeral-storage" json:"requests.ephemeral-storage"`
	LimitsEphemeralStorage   string `mapstructure:"limits.ephemeral-storage" json:"limits.ephemeral-storage"`
	RequestsGPU              string `mapstructure:"requests.nvidia.com/gpu" json:"requests.nvidia.com/gpu"`
	LimitsGPU                string `mapstructure:"limits.nvidia.com/gpu" json:"limits.nvidia.com/gpu"`
}

type Quotas struct {
	Enabled      bool  `mapstructure:"enabled" json:"enabled"`
	Default      Quota `mapstructure:"default" json:"default"`
	UserEnabled  bool  `mapstructure:"userEnabled" json:"userEnabled"`
	User         Quota `mapstructure:"user" json:"user"`
	GroupEnabled bool  `mapstructure:"groupEnabled" json:"groupEnabled"`
	Group        Quota `mapstructure:"group" json:"group"`
}

type Service struct {
	NamespacePrefix      string `mapstructure:"namespacePrefix" json:"namespacePrefix"`
	GroupNamespacePrefix string `mapstructure:"groupNamespacePrefix" json:"groupNamespacePrefix"`
	Quotas               Quotas `mapstructure:"quotas" json:"quotas"`
}

type Env struct {
	AppEnv             string            `mapstructure:"appEnv" json:"appEnv"`
	AuthenticationMode string            `mapstructure:"authenticationMode" json:"authenticationMode"`
	OIDC               OIDC              `mapstructure:"oidc" json:"oidc"`
	Security           Security          `mapstructure:"security" json:"security"`
	K8SPublicEndpoint  K8SPublicEndpoint `mapstructure:"k8sPublicEndpoint" json:"k8sPublicEndpoint"`
	Service            Service           `mapstructure:"service" json:"service"`
}

func NewEnv() *Env {
	env := Env{}

	viper.SetConfigName("env")  // Name of the file (without extension)
	viper.SetConfigType("yaml") // File type
	viper.AddConfigPath(".")    // Look in the current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ No env.yaml file found, relying on system environment variables")
	} else {
		log.Println("✅ Successfully loaded env.yaml")
	}

	viper.AutomaticEnv()

	// Map the environment variables to the Env struct
	if err := viper.Unmarshal(&env); err != nil {
		panic(err)
	}

	return &env
}
