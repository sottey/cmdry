package developerpriority2

import (
	"encoding/base64"
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunChmodCalculator() {
	cmdry.Run(plugin("chmod-calculator", "Chmod Calculator", "Translate octal Unix file permission modes locally.", []string{"chmod", "permissions", "octal", "unix"}, chmodForm, chmodRun))
}
func RunCrontabGenerator() {
	cmdry.Run(plugin("crontab-generator", "Crontab Generator", "Build and explain a standard five-field cron schedule locally.", []string{"cron", "crontab", "schedule", "generator"}, cronForm, cronRun))
}
func RunDockerRunToCompose() {
	cmdry.Run(plugin("docker-run-to-compose", "Docker Run to Compose Converter", "Convert a bounded docker run command into a Compose draft locally.", []string{"docker", "compose", "container", "convert"}, dockerForm, dockerRun))
}
func RunEmailNormalizer() {
	cmdry.Run(plugin("email-normalizer", "Email Normalizer", "Normalize a pasted email address with explicit safe transformations.", []string{"email", "address", "normalize", "validate"}, emailForm, emailRun))
}
func RunRegexTester() {
	cmdry.Run(plugin("regex-tester", "Regex Tester", "Test a bounded Go RE2 regular expression locally.", []string{"regex", "regexp", "re2", "test"}, regexForm, regexRun))
}
func RunSQLPrettifier() {
	cmdry.Run(plugin("sql-prettifier", "SQL Prettifier", "Format pasted SQL locally without connecting to a database.", []string{"sql", "format", "prettify", "database"}, sqlForm, sqlRun))
}
func RunGitCheatsheet() {
	cmdry.Run(plugin("git-cheatsheet", "Git Cheatsheet", "Search a bundled reference for common Git commands.", []string{"git", "cheatsheet", "commands", "reference"}, gitForm, gitRun))
}
func RunRegexCheatsheet() {
	cmdry.Run(plugin("regex-cheatsheet", "Regex Cheatsheet", "Search a bundled regular-expression syntax reference.", []string{"regex", "regexp", "cheatsheet", "reference"}, regexCheatForm, regexCheatRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "developer", Icon: "developer", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}
func form(title, description, submit string, fields []cmdry.Field) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}, nil
}
func restart(name string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: name}}}
}
func table(title, message, id string, rows []map[string]any, again string) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: rows}, restart(again)}}, nil
}
func code(title, message, output, again string) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "code", Title: "Output", Text: output}, restart(again)}}, nil
}

func chmodForm(cmdry.Request) (cmdry.View, error) {
	return form("Chmod Calculator", "Enter octal (755), triplets (rwxr-xr-x), or assignments (u=rwx,g=rx,o=rx).", "Calculate", []cmdry.Field{{Name: "mode", Label: "Permission mode", Type: "text", Value: "755", Required: true, Placeholder: "755 or rwxr-xr-x"}})
}
func chmodRun(r cmdry.Request) (cmdry.View, error) {
	mode, err := ParsePermissionMode(fmt.Sprint(r.Params["mode"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return table("Permission mode", "The permission mode was interpreted locally.", "permissions", []map[string]any{{"field": "Octal", "value": mode.Octal}, {"field": "Symbolic", "value": mode.Symbolic}, {"field": "User", "value": mode.User}, {"field": "Group", "value": mode.Group}, {"field": "Other", "value": mode.Other}}, "Calculate another mode")
}

func cronForm(cmdry.Request) (cmdry.View, error) {
	fields := []cmdry.Field{{Name: "minute", Label: "Minute", Type: "text", Value: "0", Required: true}, {Name: "hour", Label: "Hour", Type: "text", Value: "9", Required: true}, {Name: "day_month", Label: "Day of month", Type: "text", Value: "*", Required: true}, {Name: "month", Label: "Month", Type: "text", Value: "*", Required: true}, {Name: "day_week", Label: "Day of week", Type: "text", Value: "MON-FRI", Required: true}}
	return form("Crontab Generator", "Builds standard five-field crontab syntax; it does not support Quartz extensions.", "Generate schedule", fields)
}
func cronRun(r cmdry.Request) (cmdry.View, error) {
	expression, fields, note, err := GenerateCron(fmt.Sprint(r.Params["minute"]), fmt.Sprint(r.Params["hour"]), fmt.Sprint(r.Params["day_month"]), fmt.Sprint(r.Params["month"]), fmt.Sprint(r.Params["day_week"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"field": "Cron expression", "value": expression}}
	for _, field := range fields {
		rows = append(rows, map[string]any{"field": field.Name + " (" + field.Expression + ")", "value": field.Explanation})
	}
	rows = append(rows, map[string]any{"field": "Important", "value": note})
	return table("Cron schedule", "The schedule was validated locally.", "cron", rows, "Generate another schedule")
}

func dockerForm(cmdry.Request) (cmdry.View, error) {
	return form("Docker Run to Compose Converter", "Supports image, name, environment, ports, volumes, restart policy, and command. It never runs Docker.", "Convert command", []cmdry.Field{{Name: "input", Label: "docker run command", Type: "textarea", Value: "docker run --name web -p 8080:80 -e APP_ENV=production nginx:alpine", Required: true}})
}
func dockerRun(r cmdry.Request) (cmdry.View, error) {
	draft, err := DockerRunToCompose(fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Compose draft", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: "A Compose draft was generated locally; it has not been executed."}, {Type: "download", Filename: "compose.yaml", MIMEType: "application/yaml;charset=utf-8", Content: base64.StdEncoding.EncodeToString([]byte(draft.YAML))}, {Type: "code", Title: "compose.yaml", Text: draft.YAML}, restart("Convert another command")}}, nil
}

func emailForm(cmdry.Request) (cmdry.View, error) {
	return form("Email Normalizer", "Removes surrounding whitespace and lowercases only the domain.", "Normalize email", []cmdry.Field{{Name: "input", Label: "Email address", Type: "text", Value: "Ada.Example@EXAMPLE.COM ", Required: true}})
}
func emailRun(r cmdry.Request) (cmdry.View, error) {
	email, err := NormalizeEmail(fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return table("Normalized email", "Safe transformations were applied locally.", "email", []map[string]any{{"field": "Normalized", "value": email.Normalized}, {"field": "Local part", "value": email.LocalPart}, {"field": "Domain", "value": email.Domain}, {"field": "Note", "value": email.Note}}, "Normalize another email")
}

func regexForm(cmdry.Request) (cmdry.View, error) {
	return form("Regex Tester", "Uses Go's RE2 engine: no lookarounds or backreferences, with a 100-match result limit.", "Test regex", []cmdry.Field{{Name: "pattern", Label: "Pattern", Type: "text", Value: "\\b\\w+@\\w+\\.\\w+\\b", Required: true}, {Name: "flags", Label: "Flags (i, m, s)", Type: "text", Value: ""}, {Name: "input", Label: "Test text", Type: "textarea", Value: "Contact ada@example.com or lin@example.net.", Required: true}})
}
func regexRun(r cmdry.Request) (cmdry.View, error) {
	matches, err := TestRegex(fmt.Sprint(r.Params["pattern"]), fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["flags"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{}
	for _, match := range matches {
		groups := strings.Join(match.Groups, " | ")
		if groups == "" {
			groups = "—"
		}
		rows = append(rows, map[string]any{"field": fmt.Sprintf("#%d [%d:%d]", match.Number, match.Start, match.End), "value": match.Text + "\nGroups: " + groups})
	}
	if len(rows) == 0 {
		rows = append(rows, map[string]any{"field": "Matches", "value": "No matches found"})
	}
	return table("Regex results", fmt.Sprintf("%d match(es) found locally.", len(matches)), "matches", rows, "Test another regex")
}

func sqlForm(cmdry.Request) (cmdry.View, error) {
	return form("SQL Prettifier", "A lightweight local formatter for readable SQL. It does not validate dialect-specific syntax or connect to a database.", "Format SQL", []cmdry.Field{{Name: "input", Label: "SQL", Type: "textarea", Value: "select id,name from users where active = true order by name;", Required: true}})
}
func sqlRun(r cmdry.Request) (cmdry.View, error) {
	output, err := FormatSQL(fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return code("Formatted SQL", "The SQL was formatted locally.", output, "Format another query")
}

func gitForm(cmdry.Request) (cmdry.View, error) {
	return form("Git Cheatsheet", "Search the bundled command reference. This tool does not run Git.", "Search reference", []cmdry.Field{{Name: "query", Label: "Search", Type: "text", Value: "", Placeholder: "status, branch, commit..."}})
}
func gitRun(r cmdry.Request) (cmdry.View, error) {
	return referenceView("Git reference", FilterReferences(GitReferences(), fmt.Sprint(r.Params["query"])), "Search Git commands")
}
func regexCheatForm(cmdry.Request) (cmdry.View, error) {
	return form("Regex Cheatsheet", "Search the bundled RE2-oriented syntax reference.", "Search reference", []cmdry.Field{{Name: "query", Label: "Search", Type: "text", Value: "", Placeholder: "anchors, groups, flags..."}})
}
func regexCheatRun(r cmdry.Request) (cmdry.View, error) {
	return referenceView("Regex reference", FilterReferences(RegexReferences(), fmt.Sprint(r.Params["query"])), "Search regex syntax")
}
func referenceView(title string, items []Reference, again string) (cmdry.View, error) {
	rows := []map[string]any{}
	for _, item := range items {
		rows = append(rows, map[string]any{"field": item.Topic + ": " + item.Syntax, "value": item.Description})
	}
	if len(rows) == 0 {
		rows = append(rows, map[string]any{"field": "No results", "value": "Try a broader search."})
	}
	return table(title, "The reference is bundled with Cmdry.", "reference", rows, again)
}
