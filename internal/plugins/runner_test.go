package plugins

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunnerPreservesStructuredPluginError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cmdry-fixture")
	script := "#!/bin/sh\nprintf '%s' '{\"ok\":false,\"error\":{\"code\":\"COMMAND_MISSING\",\"message\":\"ss is unavailable\"}}'\nexit 1\n"
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	plugin := Registered{Manifest: validManifest(), Path: path, Status: StatusEnabled}
	response, err := (Runner{Timeout: time.Second}).Run(context.Background(), plugin, "list", nil)
	if err != nil {
		t.Fatal(err)
	}
	if response.OK || response.Error == nil || response.Error.Code != "COMMAND_MISSING" {
		t.Fatalf("unexpected response: %#v", response)
	}
}
