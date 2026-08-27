// Package developerpriority2 provides bounded, local-only developer utilities.
package developerpriority2

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/sottey/cmdry/plugins/cronexpr"
)

const maxDeveloperInput = 256 << 10

type PermissionMode struct {
	Octal, Symbolic    string
	User, Group, Other string
}

func ParsePermissionMode(input string) (PermissionMode, error) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "0") && len(input) == 4 {
		input = input[1:]
	}
	if len(input) == 3 && strings.IndexFunc(input, func(character rune) bool { return character < '0' || character > '7' }) == -1 {
		values := []int{int(input[0] - '0'), int(input[1] - '0'), int(input[2] - '0')}
		return permissionMode(values), nil
	}
	if len(input) == 9 {
		return parseSymbolicTriplets(input)
	}
	return parseSymbolicAssignments(input)
}
func permissionMode(values []int) PermissionMode {
	parts := make([]string, 3)
	labels := []string{"User", "Group", "Other"}
	for i, value := range values {
		parts[i] = permissionText(value, labels[i])
	}
	return PermissionMode{Octal: fmt.Sprintf("%d%d%d", values[0], values[1], values[2]), Symbolic: "-" + permissionTriplet(values[0]) + permissionTriplet(values[1]) + permissionTriplet(values[2]), User: parts[0], Group: parts[1], Other: parts[2]}
}
func parseSymbolicTriplets(input string) (PermissionMode, error) {
	values := make([]int, 3)
	for group := 0; group < 3; group++ {
		for bit, expected := range []byte{'r', 'w', 'x'} {
			character := input[group*3+bit]
			if character == '-' {
				continue
			}
			if character != expected {
				return PermissionMode{}, fmt.Errorf("symbolic mode must use rwx and - characters")
			}
			values[group] |= 4 >> bit
		}
	}
	return permissionMode(values), nil
}
func parseSymbolicAssignments(input string) (PermissionMode, error) {
	values := []int{0, 0, 0}
	for _, assignment := range strings.Split(input, ",") {
		pieces := strings.Split(assignment, "=")
		if len(pieces) != 2 || pieces[0] == "" {
			return PermissionMode{}, fmt.Errorf("enter octal (755), rwxr-xr-x, or assignments such as u=rwx,g=rx,o=rx")
		}
		permissions := 0
		for _, character := range pieces[1] {
			switch character {
			case 'r':
				permissions |= 4
			case 'w':
				permissions |= 2
			case 'x':
				permissions |= 1
			default:
				return PermissionMode{}, fmt.Errorf("symbolic permissions may contain only r, w, and x")
			}
		}
		for _, subject := range pieces[0] {
			index := strings.IndexRune("ugo", subject)
			if index < 0 {
				return PermissionMode{}, fmt.Errorf("symbolic assignments use u, g, and o")
			}
			values[index] = permissions
		}
	}
	return permissionMode(values), nil
}
func permissionTriplet(value int) string {
	var builder strings.Builder
	if value&4 != 0 {
		builder.WriteByte('r')
	} else {
		builder.WriteByte('-')
	}
	if value&2 != 0 {
		builder.WriteByte('w')
	} else {
		builder.WriteByte('-')
	}
	if value&1 != 0 {
		builder.WriteByte('x')
	} else {
		builder.WriteByte('-')
	}
	return builder.String()
}
func permissionText(value int, label string) string {
	names := []string{}
	if value&4 != 0 {
		names = append(names, "read")
	}
	if value&2 != 0 {
		names = append(names, "write")
	}
	if value&1 != 0 {
		names = append(names, "execute")
	}
	if len(names) == 0 {
		return label + ": no permissions"
	}
	return label + ": " + strings.Join(names, ", ")
}

func GenerateCron(minute, hour, dayOfMonth, month, dayOfWeek string) (string, []cronexpr.Field, string, error) {
	expression := strings.Join([]string{minute, hour, dayOfMonth, month, dayOfWeek}, " ")
	fields, note, err := cronexpr.Explain(expression)
	if err != nil {
		return "", nil, "", err
	}
	return expression, fields, note, nil
}

type ComposeDraft struct {
	YAML     string
	Image    string
	Warnings []string
}

// DockerRunToCompose translates a limited docker run command without executing it.
func DockerRunToCompose(input string) (ComposeDraft, error) {
	if len(input) == 0 || len(input) > maxDeveloperInput {
		return ComposeDraft{}, fmt.Errorf("command must contain between 1 and %d bytes", maxDeveloperInput)
	}
	tokens, err := shellWords(input)
	if err != nil {
		return ComposeDraft{}, err
	}
	if len(tokens) < 3 || tokens[0] != "docker" || tokens[1] != "run" {
		return ComposeDraft{}, fmt.Errorf("enter a docker run command")
	}
	name, image, restart := "app", "", ""
	env, ports, volumes := []string{}, []string{}, []string{}
	var command []string
	for index := 2; index < len(tokens); index++ {
		token := tokens[index]
		if image != "" {
			command = append(command, token)
			continue
		}
		value := ""
		if strings.Contains(token, "=") {
			value = strings.SplitN(token, "=", 2)[1]
		}
		switch {
		case token == "-d" || token == "--detach" || token == "--rm":
		case token == "--name" || token == "-n":
			index++
			if index >= len(tokens) {
				return ComposeDraft{}, fmt.Errorf("%s requires a value", token)
			}
			name = tokens[index]
		case strings.HasPrefix(token, "--name="):
			name = value
		case token == "-e" || token == "--env":
			index++
			if index >= len(tokens) {
				return ComposeDraft{}, fmt.Errorf("%s requires a value", token)
			}
			env = append(env, tokens[index])
		case strings.HasPrefix(token, "--env="):
			env = append(env, value)
		case token == "-p" || token == "--publish":
			index++
			if index >= len(tokens) {
				return ComposeDraft{}, fmt.Errorf("%s requires a value", token)
			}
			ports = append(ports, tokens[index])
		case strings.HasPrefix(token, "--publish="):
			ports = append(ports, value)
		case token == "-v" || token == "--volume":
			index++
			if index >= len(tokens) {
				return ComposeDraft{}, fmt.Errorf("%s requires a value", token)
			}
			volumes = append(volumes, tokens[index])
		case strings.HasPrefix(token, "--volume="):
			volumes = append(volumes, value)
		case token == "--restart":
			index++
			if index >= len(tokens) {
				return ComposeDraft{}, fmt.Errorf("--restart requires a value")
			}
			restart = tokens[index]
		case strings.HasPrefix(token, "--restart="):
			restart = value
		case strings.HasPrefix(token, "-"):
			return ComposeDraft{}, fmt.Errorf("unsupported docker run option %q", token)
		default:
			image = token
		}
	}
	if image == "" {
		return ComposeDraft{}, fmt.Errorf("docker run command must include an image")
	}
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`).MatchString(name) {
		return ComposeDraft{}, fmt.Errorf("container name contains unsupported characters")
	}
	var out strings.Builder
	out.WriteString("services:\n  " + name + ":\n    image: " + yamlQuote(image) + "\n")
	if restart != "" {
		out.WriteString("    restart: " + yamlQuote(restart) + "\n")
	}
	writeYAMLList(&out, "environment", env)
	writeYAMLList(&out, "ports", ports)
	writeYAMLList(&out, "volumes", volumes)
	if len(command) > 0 {
		out.WriteString("    command:\n")
		for _, part := range command {
			out.WriteString("      - " + yamlQuote(part) + "\n")
		}
	}
	return ComposeDraft{YAML: out.String(), Image: image}, nil
}
func yamlQuote(value string) string { return strconv.Quote(value) }
func writeYAMLList(out *strings.Builder, label string, values []string) {
	if len(values) == 0 {
		return
	}
	out.WriteString("    " + label + ":\n")
	for _, value := range values {
		out.WriteString("      - " + yamlQuote(value) + "\n")
	}
}
func shellWords(input string) ([]string, error) {
	var words []string
	var word strings.Builder
	var quote rune
	escaped := false
	flush := func() {
		if word.Len() > 0 {
			words = append(words, word.String())
			word.Reset()
		}
	}
	for _, char := range input {
		if escaped {
			// A backslash immediately followed by a line break is a shell line
			// continuation, not an escaped newline within an argument.
			if char == '\n' || char == '\r' {
				escaped = false
				continue
			}
			word.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' && quote != '\'' {
			escaped = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				word.WriteRune(char)
			}
			continue
		}
		if char == '\'' || char == '"' {
			quote = char
			continue
		}
		if unicode.IsSpace(char) {
			flush()
			continue
		}
		word.WriteRune(char)
	}
	if escaped || quote != 0 {
		return nil, fmt.Errorf("command has an unfinished quote or escape")
	}
	flush()
	return words, nil
}

type Email struct{ Input, Normalized, LocalPart, Domain, Note string }

func NormalizeEmail(input string) (Email, error) {
	original := input
	input = strings.TrimSpace(input)
	if input == "" || len(input) > 320 {
		return Email{}, fmt.Errorf("enter an email address up to 320 characters")
	}
	if strings.Count(input, "@") != 1 {
		return Email{}, fmt.Errorf("email address must contain one @")
	}
	parts := strings.SplitN(input, "@", 2)
	if parts[0] == "" || parts[1] == "" || strings.ContainsAny(parts[0], " \t\r\n") || !strings.Contains(parts[1], ".") {
		return Email{}, fmt.Errorf("enter a conventional email address with a domain")
	}
	domain := strings.ToLower(parts[1])
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Contains(domain, "..") {
		return Email{}, fmt.Errorf("domain is malformed")
	}
	note := "Whitespace was removed and the domain was lowercased. The local part was preserved because its case can be mailbox-specific."
	return Email{Input: original, Normalized: parts[0] + "@" + domain, LocalPart: parts[0], Domain: domain, Note: note}, nil
}

type RegexMatch struct {
	Number, Start, End int
	Text               string
	Groups             []string
}

func TestRegex(pattern, input, flags string) ([]RegexMatch, error) {
	if len(pattern) == 0 || len(pattern) > 4096 {
		return nil, fmt.Errorf("pattern must contain between 1 and 4,096 bytes")
	}
	if len(input) > maxDeveloperInput {
		return nil, fmt.Errorf("test text must not exceed %d bytes", maxDeveloperInput)
	}
	if strings.Trim(flags, "ims") != "" {
		return nil, fmt.Errorf("flags may contain only i, m, and s")
	}
	expression := pattern
	if flags != "" {
		expression = "(?" + flags + ")" + pattern
	}
	compiled, err := regexp.Compile(expression)
	if err != nil {
		return nil, fmt.Errorf("compile regex: %w", err)
	}
	indexes := compiled.FindAllStringSubmatchIndex(input, 101)
	if len(indexes) > 100 {
		return nil, fmt.Errorf("regex produced more than 100 matches; narrow the pattern or input")
	}
	output := make([]RegexMatch, 0, len(indexes))
	for i, index := range indexes {
		match := RegexMatch{Number: i + 1, Start: index[0], End: index[1], Text: input[index[0]:index[1]]}
		for group := 2; group < len(index); group += 2 {
			if index[group] >= 0 {
				match.Groups = append(match.Groups, input[index[group]:index[group+1]])
			} else {
				match.Groups = append(match.Groups, "(unmatched)")
			}
		}
		output = append(output, match)
	}
	return output, nil
}

func FormatSQL(input string) (string, error) {
	if len(input) == 0 || len(input) > maxDeveloperInput {
		return "", fmt.Errorf("SQL must contain between 1 and %d bytes", maxDeveloperInput)
	}
	tokens, err := sqlTokens(input)
	if err != nil {
		return "", err
	}
	text := regexp.MustCompile(`\s+`).ReplaceAllString(strings.Join(tokens, " "), " ")
	text = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(text, ",\n  ")
	text = strings.ReplaceAll(text, " ;", ";")
	phrases := []string{"INSERT INTO", "DELETE FROM", "GROUP BY", "ORDER BY", "LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "SELECT", "FROM", "WHERE", "HAVING", "LIMIT", "VALUES", "UPDATE", "SET", "JOIN", "ON"}
	for _, phrase := range phrases {
		text = regexp.MustCompile(`(?i)\b`+strings.ReplaceAll(phrase, " ", `\s+`)+`\b`).ReplaceAllString(text, "\n"+phrase)
	}
	lines := strings.Split(strings.TrimSpace(text), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n") + "\n", nil
}
func sqlTokens(input string) ([]string, error) {
	tokens := []string{}
	var current strings.Builder
	var quote rune
	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, current.String())
			current.Reset()
		}
	}
	for _, char := range input {
		if quote != 0 {
			current.WriteRune(char)
			if char == quote {
				quote = 0
			}
			continue
		}
		if char == '\'' || char == '"' {
			current.WriteRune(char)
			quote = char
			continue
		}
		if unicode.IsSpace(char) {
			flush()
			continue
		}
		if char == ',' || char == ';' || char == '(' || char == ')' {
			flush()
			tokens = append(tokens, string(char))
			continue
		}
		current.WriteRune(char)
	}
	if quote != 0 {
		return nil, fmt.Errorf("SQL contains an unfinished quoted string")
	}
	flush()
	return tokens, nil
}

type Reference struct{ Topic, Syntax, Description string }

func FilterReferences(references []Reference, query string) []Reference {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return references
	}
	output := []Reference{}
	for _, reference := range references {
		if strings.Contains(strings.ToLower(reference.Topic+" "+reference.Syntax+" "+reference.Description), query) {
			output = append(output, reference)
		}
	}
	return output
}
func GitReferences() []Reference {
	return []Reference{{"Start", "git init", "Create a repository in the current directory."}, {"Status", "git status", "Show changed, staged, and untracked files."}, {"Stage", "git add <path>", "Stage selected changes for the next commit."}, {"Commit", "git commit -m \"message\"", "Create a local commit from staged changes."}, {"History", "git log --oneline --decorate", "Show compact commit history."}, {"Branch", "git switch -c <name>", "Create and switch to a branch."}, {"Sync", "git pull --ff-only", "Update safely when a fast-forward is possible."}, {"Inspect", "git diff --staged", "Review staged changes before committing."}}
}
func RegexReferences() []Reference {
	return []Reference{{"Anchors", "^  $", "Beginning and end of text or line in multiline mode."}, {"Character classes", "[abc]  [^abc]  \\d  \\w", "Match selected, excluded, digit, or word characters."}, {"Quantifiers", "*  +  ?  {m,n}", "Repeat the preceding expression."}, {"Groups", "(expr)  (?:expr)", "Capture or group expressions. Go RE2 does not support lookarounds."}, {"Alternation", "one|two", "Match either expression."}, {"Flags", "i m s", "Case-insensitive, multiline anchors, and dot-matches-newline."}}
}
func SortReferences(items []Reference) []Reference {
	out := append([]Reference(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].Topic < out[j].Topic })
	return out
}
