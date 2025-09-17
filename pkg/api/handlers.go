package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shenbo/sql-review-learning-demo/pkg/advisor"
	"github.com/shenbo/sql-review-learning-demo/pkg/database"
)

// Server HTTP服务器
type Server struct {
	dbManager *database.DatabaseManager
	advisor   advisor.Advisor
}

// NewServer 创建HTTP服务器
func NewServer(dbManager *database.DatabaseManager, advisor advisor.Advisor) *Server {
	return &Server{
		dbManager: dbManager,
		advisor:   advisor,
	}
}

// ConnectionTestRequest 连接测试请求
type ConnectionTestRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Database string `json:"database" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Engine   string `json:"engine" binding:"required"`
}

// ConnectionSaveRequest 连接保存请求
type ConnectionSaveRequest struct {
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Database string `json:"database" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Engine   string `json:"engine" binding:"required"`
}

// SQLRequest SQL请求
type SQLRequest struct {
	SQL          string   `json:"sql" binding:"required"`
	ConnectionID string   `json:"connection_id" binding:"required"`
	DryRun       bool     `json:"dry_run"`
	Rules        []string `json:"rules"`
}

// SQLResponse SQL响应
type SQLResponse struct {
	ReviewResults []*advisor.Advice      `json:"review_results"`
	ExecuteResult *ExecuteResult         `json:"execute_result,omitempty"`
	Schema        *database.SchemaInfo   `json:"schema,omitempty"`
}

// ExecuteResult SQL执行结果
type ExecuteResult struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	RowsAffected int64       `json:"rows_affected,omitempty"`
	Data         [][]string  `json:"data,omitempty"`
	Columns      []string    `json:"columns,omitempty"`
}

// TestConnection 测试数据库连接
func (s *Server) TestConnection(c *gin.Context) {
	var req ConnectionTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := &database.ConnectionConfig{
		Host:     req.Host,
		Port:     req.Port,
		Database: req.Database,
		Username: req.Username,
		Password: req.Password,
		Engine:   req.Engine,
	}

	if err := s.dbManager.TestConnection(config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "连接测试成功",
	})
}

// SaveConnection 保存数据库连接
func (s *Server) SaveConnection(c *gin.Context) {
	var req ConnectionSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成连接ID
	id, err := generateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate connection ID"})
		return
	}

	config := &database.ConnectionConfig{
		ID:       id,
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Database: req.Database,
		Username: req.Username,
		Password: req.Password,
		Engine:   req.Engine,
	}

	if err := s.dbManager.AddConnection(config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "连接保存成功",
		"id":      id,
	})
}

// ListConnections 列出所有连接
func (s *Server) ListConnections(c *gin.Context) {
	connections := s.dbManager.ListConnections()

	// 隐藏密码信息
	for _, conn := range connections {
		conn.Password = "****"
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"connections": connections,
	})
}

// GetSchema 获取数据库schema
func (s *Server) GetSchema(c *gin.Context) {
	connectionID := c.Param("connection_id")

	db, err := s.dbManager.GetConnection(connectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	config, err := s.dbManager.GetConfig(connectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	schemaManager := database.NewSchemaManager(db, config.Engine)
	schema, err := schemaManager.GetSchemaInfo(config.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"schema":  schema,
	})
}

// ReviewSQL 审查SQL
func (s *Server) ReviewSQL(c *gin.Context) {
	var req SQLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := s.dbManager.GetConfig(req.ConnectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	db, err := s.dbManager.GetConnection(req.ConnectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 构建审查上下文
	checkCtx := &advisor.Context{
		SQL:          req.SQL,
		Engine:       config.Engine,
		DatabaseName: config.Database,
		Rules:        req.Rules,
		Connection:   db,
	}

	// 执行SQL审查
	advices, err := s.advisor.Check(c.Request.Context(), checkCtx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &SQLResponse{
		ReviewResults: advices,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  response,
	})
}

// ListRules 列出所有规则
func (s *Server) ListRules(c *gin.Context) {
	if defaultAdvisor, ok := s.advisor.(*advisor.DefaultAdvisor); ok {
		rules := defaultAdvisor.ListRules()
		rulesInfo := make([]map[string]interface{}, len(rules))

		for i, rule := range rules {
			rulesInfo[i] = map[string]interface{}{
				"id":          rule.ID(),
				"name":        rule.Name(),
				"description": rule.Description(),
				"level":       rule.Level(),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"rules":   rulesInfo,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to list rules"})
	}
}

// generateID 生成随机ID
func generateID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}