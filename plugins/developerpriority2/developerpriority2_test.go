package developerpriority2

import (
	"strings"
	"testing"
)

func TestPermissionModes(t *testing.T) {
	for _, input := range []string{"755", "rwxr-xr-x", "u=rwx,g=rx,o=rx"} {
		mode, err := ParsePermissionMode(input)
		if err != nil || mode.Octal != "755" || mode.Symbolic != "-rwxr-xr-x" {
			t.Fatalf("%q: mode=%#v err=%v", input, mode, err)
		}
	}
	if _, err := ParsePermissionMode("89a"); err == nil {
		t.Fatal("accepted invalid mode")
	}
}

func TestDeveloperTransforms(t *testing.T) {
	expression, _, _, err := GenerateCron("0", "9", "*", "*", "MON-FRI")
	if err != nil || expression != "0 9 * * MON-FRI" {
		t.Fatalf("cron=%q err=%v", expression, err)
	}
	draft, err := DockerRunToCompose("docker run --name web -p 8080:80 -e APP_ENV=production nginx:alpine")
	if err != nil || !strings.Contains(draft.YAML, "image: \"nginx:alpine\"") || !strings.Contains(draft.YAML, "\"8080:80\"") {
		t.Fatalf("draft=%q err=%v", draft.YAML, err)
	}
	multilineDraft, err := DockerRunToCompose(`docker run --name pdfding \
    -p 8000:8000 \
    -v sqlite_data:/home/nonroot/pdfding/db -v media:/home/nonroot/pdfding/media \
    -e HOST_NAME=127.0.0.1 -e SECRET_KEY=some_secret \
    -e CSRF_COOKIE_SECURE=FALSE -e SESSION_COOKIE_SECURE=FALSE \
    -d \
    mrmn/pdfding:latest`)
	if err != nil || !strings.Contains(multilineDraft.YAML, "image: \"mrmn/pdfding:latest\"") || strings.Count(multilineDraft.YAML, "      - ") != 7 {
		t.Fatalf("multiline draft=%q err=%v", multilineDraft.YAML, err)
	}
	email, err := NormalizeEmail(" Ada.Example@EXAMPLE.COM ")
	if err != nil || email.Normalized != "Ada.Example@example.com" {
		t.Fatalf("email=%#v err=%v", email, err)
	}
	matches, err := TestRegex(`a(da)`, "Ada ada", "i")
	if err != nil || len(matches) != 2 || matches[0].Groups[0] != "da" {
		t.Fatalf("matches=%#v err=%v", matches, err)
	}
	sql, err := FormatSQL("select id,name from users where active = true order by name;")
	if err != nil || !strings.Contains(sql, "SELECT id,") || !strings.Contains(sql, "FROM users") {
		t.Fatalf("sql=%q err=%v", sql, err)
	}
}

func TestReferences(t *testing.T) {
	if got := FilterReferences(GitReferences(), "branch"); len(got) != 1 || got[0].Topic != "Branch" {
		t.Fatalf("references=%#v", got)
	}
}
