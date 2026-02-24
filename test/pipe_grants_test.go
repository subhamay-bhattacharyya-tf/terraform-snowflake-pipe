// File: test/pipe_grants_test.go
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

// TestPipeWithGrants tests creating a pipe with role grants
func TestPipeWithGrants(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())
	dbName := fmt.Sprintf("TT_DB_%s", unique)
	schemaName := "TEST_SCHEMA"
	tableName := "TEST_TABLE"
	stageName := "TEST_STAGE"
	pipeName := fmt.Sprintf("TT_PIPE_%s", unique)
	roleName := fmt.Sprintf("TT_ROLE_%s", unique)

	// Setup: Create database, schema, table, external stage, and role
	db := openSnowflake(t)
	createTestDatabase(t, db, dbName)
	createTestSchema(t, db, dbName, schemaName)
	createTestTable(t, db, dbName, schemaName, tableName)
	createExternalStage(t, db, dbName, schemaName, stageName)
	createTestRole(t, db, roleName)
	_ = db.Close()

	tfDir := "../examples/single-pipe"

	copyStatement := fmt.Sprintf("COPY INTO %s.%s.%s FROM @%s.%s.%s FILE_FORMAT = (TYPE = CSV)",
		dbName, schemaName, tableName, dbName, schemaName, stageName)

	pipeConfigs := map[string]interface{}{
		"test_pipe": map[string]interface{}{
			"database":       dbName,
			"schema":         schemaName,
			"name":           pipeName,
			"copy_statement": copyStatement,
			"auto_ingest":    false,
			"comment":        "Terratest pipe with grants test",
			"grants": []interface{}{
				map[string]interface{}{
					"role_name":  roleName,
					"privileges": []interface{}{"MONITOR", "OPERATE"},
				},
			},
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
		// Cleanup: Drop test database and role
		db := openSnowflake(t)
		dropTestDatabase(t, db, dbName)
		dropTestRole(t, db, roleName)
		_ = db.Close()
	}()

	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	db = openSnowflake(t)
	defer func() { _ = db.Close() }()

	// Verify pipe exists
	exists := pipeExists(t, db, dbName, schemaName, pipeName)
	require.True(t, exists, "Expected pipe %q to exist in Snowflake", pipeName)

	// Verify grants exist
	hasMonitor := roleHasPipePrivilege(t, db, roleName, dbName, schemaName, pipeName, "MONITOR")
	require.True(t, hasMonitor, "Expected role %q to have MONITOR privilege on pipe %q", roleName, pipeName)

	hasOperate := roleHasPipePrivilege(t, db, roleName, dbName, schemaName, pipeName, "OPERATE")
	require.True(t, hasOperate, "Expected role %q to have OPERATE privilege on pipe %q", roleName, pipeName)
}
