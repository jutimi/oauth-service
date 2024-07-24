package config

type JWT struct {
	Issuer              string `mapstructure:"issuer"`
	UserAccessTokenKey  string `mapstructure:"user_access_token_key"`
	UserRefreshTokenKey string `mapstructure:"user_refresh_token_key"`
	WSAccessTokenKey    string `mapstructure:"workspace_access_token_key"`
	WSRefreshTokenKey   string `mapstructure:"workspace_refresh_token_key"`
}

type Server struct {
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	SentryUrl string `mapstructure:"sentry_url"`
}

type GRPC struct {
	OAuthPort     int    `mapstructure:"oauth_port"`
	WorkspacePort string `mapstructure:"workspace_port"`
}
