package config

type JWT struct {
	Issuer                   string `mapstructure:"issuer"`
	UserAccessTokenKey       string `mapstructure:"user_access_token_key"`
	UserRefreshTokenKey      string `mapstructure:"user_refresh_token_key"`
	WorkspaceAccessTokenKey  string `mapstructure:"workspace_access_token_key"`
	WorkspaceRefreshTokenKey string `mapstructure:"workspace_refresh_token_key"`
}

type Server struct {
	Port        int    `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`
	SentryUrl   string `mapstructure:"sentry_url"`
	ServiceName string `mapstructure:"service_name"`
	UptraceDNS  string `mapstructure:"uptrace_dns"`
}

type GRPC struct {
	OAuthUrl     string `mapstructure:"oauth_url"`
	WorkspaceUrl string `mapstructure:"workspace_url"`
}
