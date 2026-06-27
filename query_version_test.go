package hypr

import (
	"context"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	v, err := Version(context.Background(), fixtureOpts(t, "version.json")...)
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v.Version != "0.56.2" || v.Tag != "v0.56.2" || v.Commit != "efb50993780079460b0cbed1363e2166a2de1d9f" || v.Dirty {
		t.Errorf("Version = %+v", v)
	}
	if v.CommitMessage != "[gha] Nix: update inputs" || v.BuildHyprlang != "0.6.8" || v.SystemAquamarine != "0.15.0" {
		t.Errorf("Version tags: %+v", v)
	}
	if v.Flags == nil || len(v.Flags) != 0 {
		t.Errorf("Flags = %v, want an empty list", v.Flags)
	}
}
