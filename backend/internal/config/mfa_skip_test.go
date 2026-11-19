package config

import "testing"

// Login admin tanpa MFA hanya boleh di development.
func TestAdminMFADevSkip(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x@localhost/x")
	cases := []struct {
		env, skip string
		want      bool
	}{
		{"development", "", true},
		{"development", "false", false},
		{"staging", "true", false},
		{"production", "true", false},
	}
	for _, c := range cases {
		t.Setenv("APP_ENV", c.env)
		t.Setenv("ADMIN_MFA_SKIP", c.skip)
		cfg, _ := Load()
		if cfg.AdminMFADevSkip != c.want {
			t.Errorf("APP_ENV=%s ADMIN_MFA_SKIP=%q: got %v want %v", c.env, c.skip, cfg.AdminMFADevSkip, c.want)
		}
	}
}
