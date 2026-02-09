// File: test/multiple_pipes_test.go
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

// TestMultiplePipes tests creating multiple pipes via the module
func TestMultiplePipes(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)
	schemaName := "TEST_SCHEMA"

	pipe1Name := fmt.Sprintf("TT_ORDERS_%s", unique)
	pipe2Name := fmt.Sprintf("TT_CUSTOMERS_%s", unique)
	pipe3Name := fmt.Sprintf("TT_EVENTS_%s", unique)

	// Setup: Create database, schema, tables, and stages
	db := openSnowflake(t)
	createTestDatabase(t, db, dbName)
	createTestSchema(t, db, dbName, schemaName)
	createTestTable(t, db, dbName, schemaName, "ORDERS")
	createTestTable(t, db, dbName, schemaName, "CUSTOMERS")
	createTestTable(t, db, dbName, schemaName, "EVENTS")
	createTestStage(t, db, dbName, schemaName, "ORDERS_STAGE")
	createTestStage(t, db, dbName, schemaName, "CUSTOMERS_STAGE")
	createTestStage(t, db, dbName, schemaName, "EVENTS_STAGE")
	_ = db.Close()

	tfDir := "../examples/multiple-pipes"

	pipeConfigs := map[string]interface{}{
		"orders_pipe": map[string]interface{}{
			"database":       dbName,
			"schema":         schemaName,
			"name":           pipe1Name,
			"copy_statement": fmt.Sprintf("COPY INTO %s.%s.ORDERS FROM @%s.%s.ORDERS_STAGE", dbName, schemaName, dbName, schemaName),
			"auto_ingest":    false,
			"comment":        "Terratest orders pipe",
		},
		"customers_pipe": map[string]interface{}{
			"database":       dbName,
			"schema":         schemaName,
			"name":           pipe2Name,
			"copy_statement": fmt.Sprintf("COPY INTO %s.%s.CUSTOMERS FROM @%s.%s.CUSTOMERS_STAGE", dbName, schemaName, dbName, schemaName),
			"auto_ingest":    false,
			"comment":        "Terratest customers pipe",
		},
		"events_pipe": map[string]interface{}{
			"database":       dbName,
			"schema":         schemaName,
			"name":           pipe3Name,
			"copy_statement": fmt.Sprintf("COPY INTO %s.%s.EVENTS FROM @%s.%s.EVENTS_STAGE", dbName, schemaName, dbName, schemaName),
			"auto_ingest":    false,
			"comment":        "Terratest events pipe",
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

	// Verify all three pipes exist
	for _, pipeName := range []string{pipe1Name, pipe2Name, pipe3Name} {
		exists := pipeExists(t, db, dbName, schemaName, pipeName)
		require.True(t, exists, "Expected pipe %q to exist in Snowflake", pipeName)
	}

	// Verify properties of events pipe
	props := fetchPipeProps(t, db, dbName, schemaName, pipe3Name)
	require.Equal(t, pipe3Name, props.Name)
	require.Equal(t, dbName, props.DatabaseName)
	require.Equal(t, schemaName, props.SchemaName)
	require.Contains(t, props.Comment, "Terratest events pipe")
}
