package legacyparams_test

import (
	"errors"
	"testing"
	"time"

	"github.com/wakatime/wakatime-cli/cmd/legacy/legacyparams"
	"github.com/wakatime/wakatime-cli/pkg/api"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_OfflineDisabled_ConfigTakesPrecedence(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("disable-offline", false)
	v.Set("disableoffline", false)
	v.Set("settings.offline", false)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.True(t, params.OfflineDisabled)
}

func TestLoad_OfflineDisabled_FlagDeprecatedTakesPrecedence(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("disable-offline", false)
	v.Set("disableoffline", true)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.True(t, params.OfflineDisabled)
}

func TestLoad_OfflineDisabled_FromFlag(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("disable-offline", true)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.True(t, params.OfflineDisabled)
}

func TestLoad_OfflineSyncMax(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("sync-offline-activity", 42)
	v.SetDefault("sync-offline-activity", 100)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, 42, params.OfflineSyncMax)
}

func TestLoad_OfflineSyncMax_NoEntity(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("sync-offline-activity", 42)
	v.SetDefault("sync-offline-activity", 100)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, 42, params.OfflineSyncMax)
}

func TestLoad_OfflineSyncMax_None(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("sync-offline-activity", "none")
	v.SetDefault("sync-offline-activity", 100)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, 0, params.OfflineSyncMax)
}

func TestLoad_OfflineSyncMax_Default(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.SetDefault("sync-offline-activity", 100)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, 100, params.OfflineSyncMax)
}

func TestLoad_OfflineSyncMax_NegativeNumber(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("sync-offline-activity", -1)
	v.SetDefault("sync-offline-activity", 100)

	_, err := legacyparams.Load(v)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "--sync-offline-activity")
}

func TestLoad_OfflineSyncMax_NonIntegerValue(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("sync-offline-activity", "invalid")
	v.SetDefault("sync-offline-activity", 100)

	_, err := legacyparams.Load(v)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "--sync-offline-activity")
}

func TestLoad_API_APIKey(t *testing.T) {
	tests := map[string]struct {
		ViperAPIKey          string
		ViperAPIKeyConfig    string
		ViperAPIKeyConfigOld string
		Expected             legacyparams.Params
	}{
		"api key flag takes preceedence": {
			ViperAPIKey:          "00000000-0000-4000-8000-000000000000",
			ViperAPIKeyConfig:    "10000000-0000-4000-8000-000000000000",
			ViperAPIKeyConfigOld: "20000000-0000-4000-8000-000000000000",
			Expected: legacyparams.Params{
				API: legacyparams.API{
					Key: "00000000-0000-4000-8000-000000000000",
					URL: "https://api.wakatime.com/api",
				},
			},
		},
		"api from config takes preceedence": {
			ViperAPIKeyConfig:    "00000000-0000-4000-8000-000000000000",
			ViperAPIKeyConfigOld: "10000000-0000-4000-8000-000000000000",
			Expected: legacyparams.Params{
				API: legacyparams.API{
					Key: "00000000-0000-4000-8000-000000000000",
					URL: "https://api.wakatime.com/api",
				},
			},
		},
		"api key from config deprecated": {
			ViperAPIKeyConfigOld: "00000000-0000-4000-8000-000000000000",
			Expected: legacyparams.Params{
				API: legacyparams.API{
					Key: "00000000-0000-4000-8000-000000000000",
					URL: "https://api.wakatime.com/api",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := viper.New()
			v.Set("key", test.ViperAPIKey)
			v.Set("settings.api_key", test.ViperAPIKeyConfig)
			v.Set("settings.apikey", test.ViperAPIKeyConfigOld)

			params, err := legacyparams.Load(v)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, params)
		})
	}
}

func TestLoad_API_APIKeyInvalid(t *testing.T) {
	tests := map[string]string{
		"unset":            "",
		"invalid format 1": "not-uuid",
		"invalid format 2": "00000000-0000-0000-8000-000000000000",
		"invalid format 3": "00000000-0000-4000-0000-000000000000",
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			v := viper.New()
			v.Set("key", value)

			_, err := legacyparams.Load(v)
			require.Error(t, err)

			var errauth api.ErrAuth
			require.True(t, errors.As(err, &errauth))
		})
	}
}

func TestLoad_API_APIUrl(t *testing.T) {
	tests := map[string]struct {
		ViperAPIUrl       string
		ViperAPIUrlConfig string
		ViperAPIUrlOld    string
		Expected          legacyparams.Params
	}{
		"api url flag takes preceedence": {
			ViperAPIUrl:       "http://localhost:8080",
			ViperAPIUrlConfig: "http://localhost:8081",
			ViperAPIUrlOld:    "http://localhost:8082",
			Expected: legacyparams.Params{
				API: legacyparams.API{
					Key: "00000000-0000-4000-8000-000000000000",
					URL: "http://localhost:8080",
				},
			},
		},
		"api url deprecated flag takes preceedence": {
			ViperAPIUrlConfig: "http://localhost:8081",
			ViperAPIUrlOld:    "http://localhost:8082",
			Expected: legacyparams.Params{
				API: legacyparams.API{
					Key: "00000000-0000-4000-8000-000000000000",
					URL: "http://localhost:8082",
				},
			},
		},
		"api url from config": {
			ViperAPIUrlConfig: "http://localhost:8081",
			Expected: legacyparams.Params{
				API: legacyparams.API{
					Key: "00000000-0000-4000-8000-000000000000",
					URL: "http://localhost:8081",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := viper.New()
			v.Set("key", "00000000-0000-4000-8000-000000000000")
			v.Set("api-url", test.ViperAPIUrl)
			v.Set("apiurl", test.ViperAPIUrlOld)
			v.Set("settings.api_url", test.ViperAPIUrlConfig)

			params, err := legacyparams.Load(v)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, params)
		})
	}
}

func TestLoad_APIUrl_Default(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, api.BaseURL, params.API.URL)
}

func TestLoad_API_Plugin(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("plugin", "plugin/10.0.0")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, legacyparams.API{
		Key:    "00000000-0000-4000-8000-000000000000",
		URL:    "https://api.wakatime.com/api",
		Plugin: "plugin/10.0.0",
	}, params.API)
}

func TestLoad_API_Timeout_FlagTakesPreceedence(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("timeout", 5)
	v.Set("settings.timeout", 10)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, 5*time.Second, params.API.Timeout)
}

func TestLoad_API_Timeout_FromConfig(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("settings.timeout", 10)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, 10*time.Second, params.API.Timeout)
}

func TestLoad_API_DisableSSLVerify_FlagTakesPrecedence(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("no-ssl-verify", true)
	v.Set("settings.no_ssl_verify", false)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.True(t, params.API.DisableSSLVerify)
}

func TestLoad_API_DisableSSLVerify_FromConfig(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("settings.no_ssl_verify", true)

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.True(t, params.API.DisableSSLVerify)
}

func TestLoad_API_DisableSSLVerify_Default(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.False(t, params.API.DisableSSLVerify)
}

func TestLoad_API_ProxyURL(t *testing.T) {
	tests := map[string]string{
		"https":  "https://john:secret@example.org:8888",
		"http":   "http://john:secret@example.org:8888",
		"socks5": "socks5://john:secret@example.org:8888",
		"ntlm":   `domain\\john:123456`,
	}

	for name, proxyURL := range tests {
		t.Run(name, func(t *testing.T) {
			v := viper.New()
			v.Set("key", "00000000-0000-4000-8000-000000000000")
			v.Set("entity", "/path/to/file")
			v.Set("proxy", proxyURL)

			params, err := legacyparams.Load(v)
			require.NoError(t, err)

			assert.Equal(t, proxyURL, params.API.ProxyURL)
		})
	}
}

func TestLoad_API_ProxyURL_FlagTakesPrecedence(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("proxy", "https://john:secret@example.org:8888")
	v.Set("settings.proxy", "ignored")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, "https://john:secret@example.org:8888", params.API.ProxyURL)
}

func TestLoad_API_ProxyURL_FromConfig(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("settings.proxy", "https://john:secret@example.org:8888")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, "https://john:secret@example.org:8888", params.API.ProxyURL)
}

func TestLoad_API_ProxyURL_InvalidFormat(t *testing.T) {
	proxyURL := "ftp://john:secret@example.org:8888"

	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("proxy", proxyURL)

	_, err := legacyparams.Load(v)
	require.Error(t, err)
}

func TestLoad_API_SSLCertFilepath_FlagTakesPrecedence(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("ssl-certs-file", "/path/to/cert.pem")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, "/path/to/cert.pem", params.API.SSLCertFilepath)
}

func TestLoad_API_SSLCertFilepath_FromConfig(t *testing.T) {
	v := viper.New()
	v.Set("key", "00000000-0000-4000-8000-000000000000")
	v.Set("entity", "/path/to/file")
	v.Set("settings.ssl_certs_file", "/path/to/cert.pem")

	params, err := legacyparams.Load(v)
	require.NoError(t, err)

	assert.Equal(t, "/path/to/cert.pem", params.API.SSLCertFilepath)
}
