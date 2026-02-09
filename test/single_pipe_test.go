// File: test/single_pipe_test.go
package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// TestSinglePipe tests creating a single pipe via the module
func TestSinglePipe(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)
	schemaName := "TEST_SCHEMA"
	tableName := "TEST_TABLE"
	stageName := "TEST_STAGE"
	pipeName := fmt.Sprintf("TT_PIPE_%s", unique)

	// Setup: Create database, schema, table, and stage
	db := openSnowflake(t)
	createTestDatabase(t, db, dbName)
	createTestSchema(t, db, dbName, schemaName)
	createTestTable(t, db, dbName, schemaName, tableName)
	createTestStage(t, db, dbName, schemaName, stageName)
	_ = db.Close()

	tfDir := "../examples/single-pipe"

	copyStatement := fmt.Sprintf("COPY INTO %s.%s.%s FROM @%s.%s.%s", dbName, schemaName, tableName, dbName, schemaName, stageName)

	pipeConfigs := map[string]interface{}{
		"test_pipe": map[string]interface{}{
			"database":       dbName,
			"schema":         schemaName,
			"name":           pipeName,
			"copy_statement": copyStatement,
			"auto_ingest":    false,
			"comment":        "Terratest single pipe test",
		},
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"pipe_configs":                pipeConfigs,
			"snowflake_organization_name": os.Getenv("SNOWFLAKE_ORGANIZATION_NAME"),
			"snowflake_account_name":      os.Getenv("SNOWFLAKE_ACCOUNT_NAME"),
			"snowflake_user":              os.Getenv("SNOWFLAKE_USER"),
			"snowflake_role":              os.Getenv("SNOWFLAKE_ROLE"),
			"snowflake_private_key":       os.Getenv("SNOWFLAKE_PRIVATE_KEY"),
		},
	}

	defer func() {
		terraform.Destroy(t, tfOptions)
		// Cleanup: Drop test database
		db := openSnowflake(t)
		dropTestDatabase(t, db, dbName)
		_ = db.Close()
	}()

	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	db = openSnowflake(t)
	defer func() { _ = db.Close() }()

	exists := pipeExists(t, db, dbName, schemaName, pipeName)
	require.True(t, exists, "Expected pipe %q to exist in Snowflake", pipeName)

	props := fetchPipeProps(t, db, dbName, schemaName, pipeName)
	require.Equal(t, pipeName, props.Name)
	require.Equal(t, dbName, props.DatabaseName)
	require.Equal(t, schemaName, props.SchemaName)
	require.Contains(t, props.Comment, "Terratest single pipe test")
}
