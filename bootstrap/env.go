package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type OIDC struct {
	IssuerURI     string `mapstructure:"issuerURI"     json:"issuerURI"`
	SkipTLSVerify bool   `mapstructure:"skipTLSVerify" json:"skipTLSVerify"`
	ClientID      string `mapstructure:"clientID"      json:"clientID"`
	Audience      string `mapstructure:"audience"      json:"audience"`
	UsernameClaim string `mapstructure:"usernameClaim" json:"usernameClaim"`
	GroupsClaim   string `mapstructure:"groupsClaim"   json:"groupsClaim"`
	RolesClaim    string `mapstructure:"rolesClaim"    json:"rolesClaim"`
}

type Security struct {
	CORSAllowedOrigins []string `mapstructure:"corsAllowedOrigins" json:"corsAllowedOrigins"`
}

type K8SPublicEndpoint struct {
	OidcConfiguration struct {
		IssuerURI string `mapstructure:"issuerURI" json:"issuerURI"`
		ClientID  string `mapstructure:"clientID" json:"clientID"`
	} `mapstructure:"oidcConfiguration" json:"oidcConfiguration"`
	URL string `mapstructure:"URL"               json:"URL"`
}

type Quota struct {
	RequestsMemory           string `mapstructure:"requests.memory"            json:"requests.memory"`
	RequestsCPU              string `mapstructure:"requests.cpu"               json:"requests.cpu"`
	LimitsMemory             string `mapstructure:"limits.memory"              json:"limits.memory"`
	LimitsCPU                string `mapstructure:"limits.cpu"                 json:"limits.cpu"`
	RequestsStorage          string `mapstructure:"requests.storage"           json:"requests.storage"`
	CountPods                string `mapstructure:"count/pods"                 json:"count/pods"`
	RequestsEphemeralStorage string `mapstructure:"requests.ephemeral-storage" json:"requests.ephemeral-storage"`
	LimitsEphemeralStorage   string `mapstructure:"limits.ephemeral-storage"   json:"limits.ephemeral-storage"`
	RequestsGPU              string `mapstructure:"requests.nvidia.com/gpu"    json:"requests.nvidia.com/gpu"`
	LimitsGPU                string `mapstructure:"limits.nvidia.com/gpu"      json:"limits.nvidia.com/gpu"`
}

type Quotas struct {
	Enabled      bool  `mapstructure:"enabled"      json:"enabled"`
	Default      Quota `mapstructure:"default"      json:"default"`
	UserEnabled  bool  `mapstructure:"userEnabled"  json:"userEnabled"`
	User         Quota `mapstructure:"user"         json:"user"`
	GroupEnabled bool  `mapstructure:"groupEnabled" json:"groupEnabled"`
	Group        Quota `mapstructure:"group"        json:"group"`
}

type Service struct {
	NamespacePrefix      string `mapstructure:"namespacePrefix"      json:"namespacePrefix"`
	GroupNamespacePrefix string `mapstructure:"groupNamespacePrefix" json:"groupNamespacePrefix"`
	Quotas               Quotas `mapstructure:"quotas"               json:"quotas"`
}

type Env struct {
	AppEnv             string            `mapstructure:"appEnv"             json:"appEnv"`
	AuthenticationMode string            `mapstructure:"authenticationMode" json:"authenticationMode"`
	OIDC               OIDC              `mapstructure:"oidc"               json:"oidc"`
	Security           Security          `mapstructure:"security"           json:"security"`
	K8SPublicEndpoint  K8SPublicEndpoint `mapstructure:"k8sPublicEndpoint"  json:"k8sPublicEndpoint"`
	Service            Service           `mapstructure:"service"            json:"service"`
}

func NewEnv() (*Env, error) {
	env := Env{}

	// Define the list of config files in priority order (low -> high priority)
	configFiles := []string{"env.default", "env"}
	viper.SetConfigType("yaml") // File type

	for _, file := range configFiles {
		viper.SetConfigName(file) // Name without extension
		viper.AddConfigPath(".")  // Look in the current directory

		// Merge the configurations instead of replacing
		if err := viper.MergeInConfig(); err == nil {
			slog.Info("Loaded config file", slog.String("file", file))
		}
	}

	viper.AutomaticEnv()

	// Map the environment variables to the Env struct
	if err := viper.Unmarshal(&env); err != nil {
		slog.Error("Failed to parse environment configuration", slog.Any("error", err))
		return nil, fmt.Errorf("failed to parse environment configuration: %w", err)
	}

	return &env, nil
}
