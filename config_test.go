package main

import (
	"testing"
)

func validateTest(t *testing.T, expected string, prep func(c *Config)) {
	config := NewConfig()
	prep(config)
	err := config.validate()
	if (expected == "" && err != nil) || (expected != "" && err == nil) || (err != nil && (err.Error() != expected)) {
		t.Error("Expected: " + expected + "; Got: " + err.Error())
	}
}

func TestValidatesHeaderUsername(t *testing.T) {
	validateTest(t, "Configuration file requires header.username", func(c *Config) {})
}

func TestValidatesHeaderURI(t *testing.T) {
	validateTest(t, "Configuration file requires header.uri", func(c *Config) {
		c.HeaderUsername = "foo"
	})
}

func TestValidatesHTTPHostPort(t *testing.T) {
	validateTest(t, "Configuration file requires http.hostport", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
	})
}

func TestValidatesHTTPPath(t *testing.T) {
	validateTest(t, "Configuration file requires http.path", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
	})
}

func TestValidatesPasswdFilename(t *testing.T) {
	validateTest(t, "Configuration file requires htpasswd.filepath", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
	})
}

func TestValidatesTokenByteLength(t *testing.T) {
	validateTest(t, "Configuration file requires positive token.lengthBytes", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
	})
}

func TestValidatesTokenByteLengthPositive(t *testing.T) {
	validateTest(t, "Configuration file requires positive token.lengthBytes", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = -1
	})
}

func TestValidatesTokenByteLengthNonZero(t *testing.T) {
	validateTest(t, "Configuration file requires positive token.lengthBytes", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 0
	})
}

func TestValidatesTokenValidityHours(t *testing.T) {
	validateTest(t, "Configuration file requires positive token.durationHours", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
	})
}

func TestValidatesTokenValidityHoursPositive(t *testing.T) {
	validateTest(t, "Configuration file requires positive token.durationHours", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
		c.TokenValidityHours = -1
	})
}

func TestValidatesTokenValidityHoursNonZero(t *testing.T) {
	validateTest(t, "Configuration file requires positive token.durationHours", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
		c.TokenValidityHours = 0
	})
}

func TestValidatesMaintenanceIntervalSeconds(t *testing.T) {
	validateTest(t, "Configuration file requires positive maintenance.intervalSecs", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
		c.TokenValidityHours = 1
	})
}

func TestValidatesMaintenanceIntervalSecondsPositive(t *testing.T) {
	validateTest(t, "Configuration file requires positive maintenance.intervalSecs", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
		c.TokenValidityHours = 1
		c.MaintenanceIntervalSeconds = -1
	})
}

func TestValidatesMaintenanceIntervalSecondsNonZero(t *testing.T) {
	validateTest(t, "Configuration file requires positive maintenance.intervalSecs", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
		c.TokenValidityHours = 1
		c.MaintenanceIntervalSeconds = 0
	})
}

func TestValidatesSuccess(t *testing.T) {
	validateTest(t, "", func(c *Config) {
		c.HeaderUsername = "foo"
		c.HeaderURI = "something"
		c.HTTPHostPort = "80"
		c.HTTPPath = "bar"
		c.PasswdFilename = "file"
		c.TokenByteLength = 1
		c.TokenValidityHours = 1
		c.MaintenanceIntervalSeconds = 1
	})
}
