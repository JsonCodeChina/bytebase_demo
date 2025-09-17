package database

import (
	"testing"
)

// TestTestConnection 测试MySQL数据库连接测试功能
func TestTestConnection(t *testing.T) {
	dm := NewDatabaseManager()

	tests := []struct {
		name        string
		config      *ConnectionConfig
		expectError bool
		description string
	}{
		{
			name: "Valid MySQL config - should succeed with real database",
			config: &ConnectionConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "team_knowledge_base",
				Username: "root",
				Password: "root",
				Engine:   "mysql",
			},
			expectError: false,
			description: "应该成功连接到本地MySQL数据库",
		},
		{
			name: "Invalid host - should fail",
			config: &ConnectionConfig{
				Host:     "invalid-host-12345.example.com",
				Port:     3306,
				Database: "test",
				Username: "root",
				Password: "password",
				Engine:   "mysql",
			},
			expectError: true,
			description: "无效主机地址应该连接失败",
		},
		{
			name: "Unsupported engine - should fail",
			config: &ConnectionConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "test",
				Username: "root",
				Password: "password",
				Engine:   "oracle",
			},
			expectError: true,
			description: "不支持的数据库引擎应该失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dm.TestConnection(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none. %s", tt.description)
				} else {
					t.Logf("Got expected error: %v (%s)", err, tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
