package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadFromConfig(t *testing.T) {
	defer os.Clearenv()
	// create config.yaml file
	filePath := filepath.Join(t.TempDir(), "config.yaml")
	f, err := os.Create(filepath.Clean(filePath))
	require.NoError(t, err)
	defer func() {
		err := os.Remove(filePath)
		require.NoError(t, err)
	}()

	_, err = f.WriteString(`Host: localhost
Port: 8080
LogLevel: debug
ServiceName: test-service
IsDryRun: true
DB:
  User: testuser
  Password: testpassword
  Host: localhost
  Port: 5432
  DBName: testdb`)
	require.NoError(t, err)
	// Set up test environment variables
	_ = os.Setenv("CONFIG_FILE", filePath)
	// Clean up environment variables after test
	defer func() {
		_ = os.Unsetenv("CONFIG_FILE")
	}()

	// Test initConfig function
	config, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, filePath, viper.ConfigFileUsed())
	require.Equal(t, "localhost", config.Host)
	require.Equal(t, "8080", config.Port)
	require.Equal(t, "debug", config.LogLevel)
	require.Equal(t, "test-service", config.ServiceName)
	require.Equal(t, "testuser", config.DB.User)
	require.Equal(t, "testpassword", config.DB.Password)
}

func TestBindEnv(t *testing.T) {
	defer os.Clearenv()
	// Test bindEnv function
	_ = os.Setenv("HOST", "localhost")
	_ = os.Setenv("PORT", "8080")
	_ = os.Setenv("DB.USER", "testuser")
	c, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "localhost", c.Host)
	require.Equal(t, "8080", c.Port)
	require.Equal(t, "testuser", c.DB.User)

}
