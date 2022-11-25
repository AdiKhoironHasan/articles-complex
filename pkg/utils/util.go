package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetBTBPrivKeySignature() string {
	// api key for integration auth
	key := viper.GetString("api.integration_key")

	return key
}

// get data integration from config
const urlIntegration = "integrations.externals.http.%s.endpoints.%s"
const host = "integrations.externals.http.%s.host"

// get integration url
func GetIntegURL(integration string, name string) string {
	getHost := fmt.Sprintf(host, integration)
	getString := fmt.Sprintf(urlIntegration, integration, name)
	return viper.GetString(getHost) + viper.GetString(getString)
}
