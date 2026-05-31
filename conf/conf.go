package conf

import (
	"fmt"
	"os"
	"singo/cache"
	"singo/model"
	"singo/util"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func envDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func validateDSNPart(name, value string, reserved ...string) error {
	for _, char := range reserved {
		if strings.Contains(value, char) {
			return fmt.Errorf("%s contains reserved DSN character %q", name, char)
		}
	}
	return nil
}

// DatabaseDSN builds a MySQL DSN from DB_* environment variables.
func DatabaseDSN() (string, error) {
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	missing := make([]string, 0, 2)
	if user == "" {
		missing = append(missing, "DB_USER")
	}
	if name == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return "", fmt.Errorf("missing required database config: %s", strings.Join(missing, ", "))
	}
	if err := validateDSNPart("DB_USER", user, ":"); err != nil {
		return "", err
	}
	if err := validateDSNPart("DB_NAME", name, "/", "?"); err != nil {
		return "", err
	}

	password := os.Getenv("DB_PASSWORD")
	host := envDefault("DB_HOST", "127.0.0.1")
	port := envDefault("DB_PORT", "3306")
	charset := envDefault("DB_CHARSET", "utf8")
	parseTime, err := strconv.ParseBool(envDefault("DB_PARSE_TIME", "True"))
	if err != nil {
		return "", fmt.Errorf("invalid DB_PARSE_TIME: %w", err)
	}

	loc, err := time.LoadLocation(envDefault("DB_LOC", "Local"))
	if err != nil {
		return "", fmt.Errorf("invalid DB_LOC: %w", err)
	}

	config := mysql.NewConfig()
	config.User = user
	config.Passwd = password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%s", host, port)
	config.DBName = name
	config.ParseTime = parseTime
	config.Loc = loc
	if err := config.Apply(mysql.Charset(charset, "")); err != nil {
		return "", fmt.Errorf("invalid DB_CHARSET: %w", err)
	}

	return config.FormatDSN(), nil
}

// Init initializes configuration.
func Init() error {
	// Load environment variables from local files.
	godotenv.Load()

	// Set the log level.
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// Load translation files.
	if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		return fmt.Errorf("load locales: %w", err)
	}

	// Connect to the database.
	dsn, err := DatabaseDSN()
	if err != nil {
		return err
	}

	model.Database(dsn)
	cache.Redis()
	return nil
}
