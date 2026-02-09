// File: test/helpers_test.go
package test

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

type PipeProps struct {
	Name               string
	DatabaseName       string
	SchemaName         string
	Definition         string
	NotificationChannel string
	Comment            string
}

func openSnowflake(t *testing.T) *sql.DB {
	t.Helper()

	orgName := mustEnv(t, "SNOWFLAKE_ORGANIZATION_NAME")
	accountName := mustEnv(t, "SNOWFLAKE_ACCOUNT_NAME")
	user := mustEnv(t, "SNOWFLAKE_USER")
	privateKeyPEM := mustEnv(t, "SNOWFLAKE_PRIVATE_KEY")
	role := os.Getenv("SNOWFLAKE_ROLE")

	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	require.NotNil(t, block, "Failed to decode PEM block from private key")

	var privateKey *rsa.PrivateKey
	var err error

	// Try PKCS8 first, then PKCS1
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		require.NoError(t, err, "Failed to parse private key")
	} else {
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		require.True(t, ok, "Private key is not RSA")
	}

	// Build account identifier: orgname-accountname
	account := fmt.Sprintf("%s-%s", orgName, accountName)

	config := gosnowflake.Config{
		Account:       account,
		User:          user,
		Authenticator: gosnowflake.AuthTypeJwt,
		PrivateKey:    privateKey,
	}

	if role != "" {
		config.Role = role
	}

	dsn, err := gosnowflake.DSN(&config)
	require.NoError(t, err, "Failed to build DSN")

	db, err := sql.Open("snowflake", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func pipeExists(t *testing.T, db *sql.DB, database, schema, pipeName string) bool {
	t.Helper()

	q := fmt.Sprintf("SHOW PIPES LIKE '%s' IN %s.%s;", escapeLike(pipeName), database, schema)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	return rows.Next()
}

func fetchPipeProps(t *testing.T, db *sql.DB, database, schema, pipeName string) PipeProps {
	t.Helper()

	q := fmt.Sprintf("SHOW PIPES LIKE '%s' IN %s.%s;", escapeLike(pipeName), database, schema)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	require.NoError(t, err)

	// Find column indices
	nameIdx, dbIdx, schemaIdx, defIdx, notifIdx, commentIdx := -1, -1, -1, -1, -1, -1
	for i, col := range cols {
		switch col {
		case "name":
			nameIdx = i
		case "database_name":
			dbIdx = i
		case "schema_name":
			schemaIdx = i
		case "definition":
			defIdx = i
		case "notification_channel":
			notifIdx = i
		case "comment":
			commentIdx = i
		}
	}
	require.NotEqual(t, -1, nameIdx, "name column not found in SHOW PIPES output")

	require.True(t, rows.Next(), "No pipe found matching %s", pipeName)

	// Create slice to hold all column values
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	err = rows.Scan(valuePtrs...)
	require.NoError(t, err)

	// Extract the values we need
	getString := func(idx int) string {
		if idx == -1 || values[idx] == nil {
			return ""
		}
		if s, ok := values[idx].(string); ok {
			return s
		}
		if b, ok := values[idx].([]byte); ok {
			return string(b)
		}
		return fmt.Sprintf("%v", values[idx])
	}

	return PipeProps{
		Name:               getString(nameIdx),
		DatabaseName:       getString(dbIdx),
		SchemaName:         getString(schemaIdx),
		Definition:         getString(defIdx),
		NotificationChannel: getString(notifIdx),
		Comment:            getString(commentIdx),
	}
}

func createTestDatabase(t *testing.T, db *sql.DB, dbName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	require.NoError(t, err, "Failed to create test database")
}

func createTestSchema(t *testing.T, db *sql.DB, dbName, schemaName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s.%s", dbName, schemaName))
	require.NoError(t, err, "Failed to create test schema")
}

func createTestTable(t *testing.T, db *sql.DB, dbName, schemaName, tableName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s.%s (id INT, data VARCHAR)", dbName, schemaName, tableName))
	require.NoError(t, err, "Failed to create test table")
}

func createTestStage(t *testing.T, db *sql.DB, dbName, schemaName, stageName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("CREATE STAGE IF NOT EXISTS %s.%s.%s", dbName, schemaName, stageName))
	require.NoError(t, err, "Failed to create test stage")
}

func dropTestDatabase(t *testing.T, db *sql.DB, dbName string) {
	t.Helper()
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	require.NoError(t, err, "Failed to drop test database")
}

func mustEnv(t *testing.T, key string) string {
	t.Helper()
	v := strings.TrimSpace(os.Getenv(key))
	require.NotEmpty(t, v, "Missing required environment variable %s", key)
	return v
}

func escapeLike(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
