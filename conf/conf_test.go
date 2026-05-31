package conf

import (
	"net/url"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func setDBEnv(t *testing.T, values map[string]string) {
	t.Helper()

	keys := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_CHARSET",
		"DB_PARSE_TIME",
		"DB_LOC",
	}

	for _, key := range keys {
		t.Setenv(key, values[key])
	}
}

func parseConfig(t *testing.T, dsn string) *mysql.Config {
	t.Helper()

	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("expected parseable DSN, got %v", err)
	}
	return config
}

func dsnQuery(t *testing.T, dsn string) url.Values {
	t.Helper()

	index := strings.Index(dsn, "?")
	if index < 0 {
		t.Fatalf("expected DSN query string, got %q", dsn)
	}

	query, err := url.ParseQuery(dsn[index+1:])
	if err != nil {
		t.Fatalf("expected parseable DSN query string, got %v", err)
	}
	return query
}

func TestDatabaseDSNRequiresUserAndName(t *testing.T) {
	setDBEnv(t, map[string]string{})

	dsn, err := DatabaseDSN()
	if err == nil {
		t.Fatal("expected missing config error")
	}
	if dsn != "" {
		t.Fatalf("expected empty DSN, got %q", dsn)
	}

	msg := err.Error()
	if !strings.Contains(msg, "DB_USER") || !strings.Contains(msg, "DB_NAME") {
		t.Fatalf("expected error to mention missing DB_USER and DB_NAME, got %q", msg)
	}
}

func TestDatabaseDSNUsesDefaults(t *testing.T) {
	setDBEnv(t, map[string]string{
		"DB_USER": "root",
		"DB_NAME": "singo",
	})

	dsn, err := DatabaseDSN()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config := parseConfig(t, dsn)
	if config.User != "root" {
		t.Fatalf("expected user root, got %q", config.User)
	}
	if config.Addr != "127.0.0.1:3306" {
		t.Fatalf("expected address 127.0.0.1:3306, got %q", config.Addr)
	}
	if config.DBName != "singo" {
		t.Fatalf("expected database singo, got %q", config.DBName)
	}
	if charset := dsnQuery(t, dsn).Get("charset"); charset != "utf8" {
		t.Fatalf("expected charset utf8, got %q", charset)
	}
	if !config.ParseTime {
		t.Fatal("expected parseTime to be true")
	}
	if config.Loc.String() != "Local" {
		t.Fatalf("expected location Local, got %q", config.Loc.String())
	}
}

func TestDatabaseDSNUsesCustomValues(t *testing.T) {
	setDBEnv(t, map[string]string{
		"DB_USER":       "app",
		"DB_PASSWORD":   "secret",
		"DB_HOST":       "mysql.local",
		"DB_PORT":       "3307",
		"DB_NAME":       "singo_test",
		"DB_CHARSET":    "utf8mb4",
		"DB_PARSE_TIME": "true",
		"DB_LOC":        "Asia/Shanghai",
	})

	dsn, err := DatabaseDSN()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config := parseConfig(t, dsn)
	if config.User != "app" {
		t.Fatalf("expected user app, got %q", config.User)
	}
	if config.Passwd != "secret" {
		t.Fatalf("expected password secret, got %q", config.Passwd)
	}
	if config.Addr != "mysql.local:3307" {
		t.Fatalf("expected address mysql.local:3307, got %q", config.Addr)
	}
	if config.DBName != "singo_test" {
		t.Fatalf("expected database singo_test, got %q", config.DBName)
	}
	if charset := dsnQuery(t, dsn).Get("charset"); charset != "utf8mb4" {
		t.Fatalf("expected charset utf8mb4, got %q", charset)
	}
	if !config.ParseTime {
		t.Fatal("expected parseTime to be true")
	}
	if config.Loc.String() != "Asia/Shanghai" {
		t.Fatalf("expected location Asia/Shanghai, got %q", config.Loc.String())
	}
}

func TestDatabaseDSNAllowsSpecialCharactersInPassword(t *testing.T) {
	setDBEnv(t, map[string]string{
		"DB_USER":     "app",
		"DB_PASSWORD": "pa:ss@word",
		"DB_HOST":     "127.0.0.1",
		"DB_PORT":     "3306",
		"DB_NAME":     "singo",
	})

	dsn, err := DatabaseDSN()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config := parseConfig(t, dsn)
	if config.Passwd != "pa:ss@word" {
		t.Fatalf("expected password to round-trip, got %q", config.Passwd)
	}
}

func TestDatabaseDSNAllowsAtSignInUser(t *testing.T) {
	setDBEnv(t, map[string]string{
		"DB_USER": "app@server",
		"DB_NAME": "singo",
	})

	dsn, err := DatabaseDSN()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config := parseConfig(t, dsn)
	if config.User != "app@server" {
		t.Fatalf("expected user to round-trip, got %q", config.User)
	}
}

func TestDatabaseDSNRejectsReservedCharacters(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]string
		want   string
	}{
		{
			name: "user colon",
			values: map[string]string{
				"DB_USER": "app:user",
				"DB_NAME": "singo",
			},
			want: "DB_USER",
		},
		{
			name: "database slash",
			values: map[string]string{
				"DB_USER": "app",
				"DB_NAME": "singo/test",
			},
			want: "DB_NAME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDBEnv(t, tt.values)

			dsn, err := DatabaseDSN()
			if err == nil {
				t.Fatal("expected reserved character error")
			}
			if dsn != "" {
				t.Fatalf("expected empty DSN, got %q", dsn)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error to mention %s, got %q", tt.want, err.Error())
			}
		})
	}
}
