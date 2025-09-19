package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shenbo/sql-review-learning-demo/pkg/advisor"
	"github.com/shenbo/sql-review-learning-demo/pkg/rules/mysql"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"

	// Global flags
	verbose bool
	format  string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "sql-review-demo",
	Short: "SQL Review Learning Demo - A simple SQL review tool",
	Long: `SQL Review Learning Demo is a command-line tool for reviewing SQL statements.
It checks SQL files against a set of configurable rules to help maintain
code quality and consistency.

This tool is based on Bytebase's SQL review system and serves as a
learning exercise for understanding enterprise-level rule engines.`,
	Version: version,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&format, "format", "text", "output format (text, json)")

	// Add subcommands
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(rulesCmd)
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [file]",
	Short: "Check SQL file against review rules",
	Long: `Check one or more SQL files against the configured review rules.
The command will parse the SQL statements and apply all enabled rules,
reporting any violations found.

Examples:
  sql-review-demo check examples/good_examples.sql
  sql-review-demo check examples/bad_examples.sql
  sql-review-demo check --format json examples/mixed_examples.sql`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	// Create advisor and register rules
	sqlAdvisor := advisor.NewDefaultAdvisor()
	sqlAdvisor.RegisterRule(mysql.NewTableRequirePKRule())

	for _, filePath := range args {
		if err := checkFile(sqlAdvisor, filePath); err != nil {
			return fmt.Errorf("failed to check file %s: %w", filePath, err)
		}
	}

	return nil
}

func checkFile(sqlAdvisor *advisor.DefaultAdvisor, filePath string) error {
	// Read SQL file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	sql := string(content)
	if strings.TrimSpace(sql) == "" {
		if verbose {
			fmt.Printf("Skipping empty file: %s\n", filePath)
		}
		return nil
	}

	// Create check context
	checkCtx := &advisor.Context{
		SQL:          sql,
		Engine:       "mysql",
		DatabaseName: "demo",
		Rules:        []string{}, // Empty means use all rules
	}

	// Execute review
	advices, err := sqlAdvisor.Check(context.Background(), checkCtx)
	if err != nil {
		return fmt.Errorf("failed to execute review: %w", err)
	}

	// Output results
	if err := outputResults(filePath, advices); err != nil {
		return fmt.Errorf("failed to output results: %w", err)
	}

	return nil
}

func outputResults(filePath string, advices []*advisor.Advice) error {
	fileName := filepath.Base(filePath)

	switch format {
	case "json":
		return outputJSON(fileName, advices)
	case "text":
		return outputText(fileName, advices)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func outputText(fileName string, advices []*advisor.Advice) error {
	fmt.Printf("\n=== SQL Review Results for %s ===\n", fileName)

	if len(advices) == 0 {
		fmt.Println("✅ No issues found! All rules passed.")
		return nil
	}

	fmt.Printf("Found %d issue(s):\n\n", len(advices))

	for i, advice := range advices {
		icon := getIcon(advice.Level)
		fmt.Printf("%d. %s [%s] %s\n", i+1, icon, advice.Level, advice.Title)
		fmt.Printf("   Rule: %s\n", advice.RuleID)
		fmt.Printf("   Message: %s\n", advice.Message)
		if advice.Line > 0 {
			fmt.Printf("   Location: Line %d", advice.Line)
			if advice.Column > 0 {
				fmt.Printf(", Column %d", advice.Column)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	return nil
}

func outputJSON(fileName string, advices []*advisor.Advice) error {
	// Simple JSON output (we could use encoding/json for more complex formatting)
	fmt.Printf(`{"file":"%s","issues":%d,"results":[`, fileName, len(advices))

	for i, advice := range advices {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf(`{"level":"%s","title":"%s","message":"%s","rule_id":"%s"`,
			advice.Level, advice.Title, advice.Message, advice.RuleID)
		if advice.Line > 0 {
			fmt.Printf(`,"line":%d`, advice.Line)
		}
		if advice.Column > 0 {
			fmt.Printf(`,"column":%d`, advice.Column)
		}
		fmt.Print("}")
	}

	fmt.Println("]}")
	return nil
}

func getIcon(level advisor.Level) string {
	switch level {
	case advisor.LevelError:
		return "❌"
	case advisor.LevelWarning:
		return "⚠️"
	case advisor.LevelInfo:
		return "ℹ️"
	default:
		return "•"
	}
}

// rulesCmd represents the rules command
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List available review rules",
	Long: `List all available SQL review rules with their descriptions and current status.
This helps you understand what rules are available and how they work.`,
	RunE: runRules,
}

func runRules(cmd *cobra.Command, args []string) error {
	// Create advisor and register rules
	sqlAdvisor := advisor.NewDefaultAdvisor()
	sqlAdvisor.RegisterRule(mysql.NewTableRequirePKRule())

	rules := sqlAdvisor.ListRules()

	if format == "json" {
		return outputRulesJSON(rules)
	}

	return outputRulesText(rules)
}

func outputRulesText(rules []advisor.Rule) error {
	fmt.Println("\n=== Available SQL Review Rules ===\n")

	if len(rules) == 0 {
		fmt.Println("No rules registered.")
		return nil
	}

	for i, rule := range rules {
		icon := getIcon(rule.Level())
		fmt.Printf("%d. %s [%s] %s\n", i+1, icon, rule.Level(), rule.Name())
		fmt.Printf("   ID: %s\n", rule.ID())
		fmt.Printf("   Description: %s\n", rule.Description())
		fmt.Println()
	}

	fmt.Printf("Total: %d rule(s) available\n", len(rules))
	return nil
}

func outputRulesJSON(rules []advisor.Rule) error {
	fmt.Print(`{"rules":[`)

	for i, rule := range rules {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf(`{"id":"%s","name":"%s","description":"%s","level":"%s"}`,
			rule.ID(), rule.Name(), rule.Description(), rule.Level())
	}

	fmt.Printf(`],"total":%d}`, len(rules))
	fmt.Println()
	return nil
}
