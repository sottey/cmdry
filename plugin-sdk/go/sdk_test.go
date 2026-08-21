package cmdry

import "testing"

func TestBundledBuildVersionOverridesManifestVersion(t *testing.T) {
	original := BuildVersion
	t.Cleanup(func() { BuildVersion = original })
	BuildVersion = "0.21.0"
	manifest := builtManifest(Manifest{Version: "0.1.0"})
	if manifest.Version != "0.21.0" {
		t.Fatalf("manifest version = %q", manifest.Version)
	}
}

func TestIndependentPluginKeepsDeclaredVersion(t *testing.T) {
	original := BuildVersion
	t.Cleanup(func() { BuildVersion = original })
	BuildVersion = ""
	manifest := builtManifest(Manifest{Version: "2.3.4"})
	if manifest.Version != "2.3.4" {
		t.Fatalf("manifest version = %q", manifest.Version)
	}
}
