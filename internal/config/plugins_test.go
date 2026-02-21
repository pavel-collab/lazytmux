package config

import (
	"testing"
)

func TestGetPlugins(t *testing.T) {
	plugins := GetPlugins()

	if len(plugins) == 0 {
		t.Error("GetPlugins() should return at least one plugin")
	}

	// TPM should be the first plugin
	if plugins[0].Repo != "tmux-plugins/tpm" {
		t.Errorf("first plugin should be tpm, got %q", plugins[0].Repo)
	}
}

func TestGetPluginsHasRequiredFields(t *testing.T) {
	plugins := GetPlugins()

	for _, p := range plugins {
		t.Run(p.Name, func(t *testing.T) {
			if p.Name == "" {
				t.Error("plugin Name should not be empty")
			}
			if p.Repo == "" {
				t.Error("plugin Repo should not be empty")
			}
			if p.DescEN == "" {
				t.Error("plugin DescEN should not be empty")
			}
			if p.DescRU == "" {
				t.Error("plugin DescRU should not be empty")
			}
			if p.KeysEN == "" {
				t.Error("plugin KeysEN should not be empty")
			}
			if p.KeysRU == "" {
				t.Error("plugin KeysRU should not be empty")
			}
		})
	}
}

func TestGetPlugin(t *testing.T) {
	// Test existing plugin
	plugin, ok := GetPlugin("tmux-plugins/tpm")
	if !ok {
		t.Error("GetPlugin('tmux-plugins/tpm') should return true")
	}
	if plugin.Name != "TPM" {
		t.Errorf("plugin.Name = %q, expected 'TPM'", plugin.Name)
	}

	// Test non-existent plugin
	_, ok = GetPlugin("nonexistent/plugin")
	if ok {
		t.Error("GetPlugin('nonexistent/plugin') should return false")
	}
}

func TestTPMPlugin(t *testing.T) {
	plugin, ok := GetPlugin("tmux-plugins/tpm")
	if !ok {
		t.Fatal("TPM plugin not found")
	}

	// TPM should not require TPM itself
	if plugin.RequiresTPM {
		t.Error("TPM should not require TPM")
	}

	// TPM should have no requirements
	if len(plugin.Requires) > 0 {
		t.Error("TPM should have no requirements")
	}
}

func TestPluginsRequiringTPM(t *testing.T) {
	plugins := GetPlugins()

	for _, p := range plugins {
		if p.Repo == "tmux-plugins/tpm" {
			continue // Skip TPM itself
		}

		// All other plugins should require TPM
		if !p.RequiresTPM {
			t.Errorf("plugin %q should require TPM", p.Name)
		}
	}
}

func TestResurrectPlugin(t *testing.T) {
	plugin, ok := GetPlugin("tmux-plugins/tmux-resurrect")
	if !ok {
		t.Fatal("resurrect plugin not found")
	}

	// Should have settings
	if len(plugin.Settings) == 0 {
		t.Error("resurrect should have settings")
	}

	// Check specific settings
	settingKeys := make(map[string]bool)
	for _, s := range plugin.Settings {
		settingKeys[s.Key] = true
	}

	expectedSettings := []string{
		"@resurrect-capture-pane-contents",
		"@resurrect-strategy-vim",
	}

	for _, key := range expectedSettings {
		if !settingKeys[key] {
			t.Errorf("expected resurrect setting %q not found", key)
		}
	}
}

func TestContinuumPlugin(t *testing.T) {
	plugin, ok := GetPlugin("tmux-plugins/tmux-continuum")
	if !ok {
		t.Fatal("continuum plugin not found")
	}

	// Continuum requires resurrect
	found := false
	for _, req := range plugin.Requires {
		if req == "tmux-plugins/tmux-resurrect" {
			found = true
			break
		}
	}
	if !found {
		t.Error("continuum should require resurrect")
	}

	// Should have settings
	if len(plugin.Settings) == 0 {
		t.Error("continuum should have settings")
	}
}

func TestPluginSettingsHaveRequiredFields(t *testing.T) {
	plugins := GetPlugins()

	for _, p := range plugins {
		for _, s := range p.Settings {
			t.Run(p.Name+"_"+s.Key, func(t *testing.T) {
				if s.Key == "" {
					t.Error("setting Key should not be empty")
				}
				if s.DescEN == "" {
					t.Error("setting DescEN should not be empty")
				}
				if s.DescRU == "" {
					t.Error("setting DescRU should not be empty")
				}
				// Default can be empty for some settings
			})
		}
	}
}

func TestPluginSettingsChoices(t *testing.T) {
	plugins := GetPlugins()

	for _, p := range plugins {
		for _, s := range p.Settings {
			if s.Type == TypeChoice {
				if len(s.Choices) == 0 {
					t.Errorf("setting %q is TypeChoice but has no choices", s.Key)
				}

				// Default should be one of the choices
				if s.Default != "" {
					found := false
					for _, choice := range s.Choices {
						if choice == s.Default {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("setting %q default %q is not in choices %v",
							s.Key, s.Default, s.Choices)
					}
				}
			}
		}
	}
}

func TestTPMInstallPath(t *testing.T) {
	path := TPMInstallPath()

	if path == "" {
		t.Error("TPMInstallPath() should not be empty")
	}

	if path != "~/.tmux/plugins/tpm" {
		t.Errorf("TPMInstallPath() = %q, expected '~/.tmux/plugins/tpm'", path)
	}
}

func TestPluginsDir(t *testing.T) {
	dir := PluginsDir()

	if dir == "" {
		t.Error("PluginsDir() should not be empty")
	}

	if dir != "~/.tmux/plugins" {
		t.Errorf("PluginsDir() = %q, expected '~/.tmux/plugins'", dir)
	}
}

func TestPluginStateStruct(t *testing.T) {
	state := PluginState{
		Repo:      "test/plugin",
		Enabled:   true,
		Installed: false,
		Settings:  map[string]string{"@key": "value"},
	}

	if state.Repo != "test/plugin" {
		t.Errorf("Repo = %q, expected 'test/plugin'", state.Repo)
	}
	if !state.Enabled {
		t.Error("Enabled should be true")
	}
	if state.Installed {
		t.Error("Installed should be false")
	}
	if state.Settings["@key"] != "value" {
		t.Errorf("Settings[@key] = %q, expected 'value'", state.Settings["@key"])
	}
}

func TestExpectedPlugins(t *testing.T) {
	expectedPlugins := []string{
		"tmux-plugins/tpm",
		"tmux-plugins/tmux-sensible",
		"tmux-plugins/tmux-resurrect",
		"tmux-plugins/tmux-continuum",
		"tmux-plugins/tmux-yank",
		"tmux-plugins/tmux-logging",
		"tmux-plugins/tmux-copycat",
		"tmux-plugins/tmux-open",
		"tmux-plugins/tmux-pain-control",
		"tmux-plugins/tmux-sessionist",
		"tmux-plugins/tmux-prefix-highlight",
		"tmux-plugins/tmux-cpu",
		"tmux-plugins/tmux-battery",
	}

	for _, repo := range expectedPlugins {
		t.Run(repo, func(t *testing.T) {
			_, ok := GetPlugin(repo)
			if !ok {
				t.Errorf("expected plugin %q to exist", repo)
			}
		})
	}
}

func TestPluginStruct(t *testing.T) {
	plugin := Plugin{
		Name:        "Test Plugin",
		Repo:        "test/repo",
		DescEN:      "Test description",
		DescRU:      "Тестовое описание",
		KeysEN:      "prefix + t: test",
		KeysRU:      "prefix + t: тест",
		RequiresTPM: true,
		Requires:    []string{"other/plugin"},
		Settings: []PluginSetting{
			{Key: "@test-key", DescEN: "Test", DescRU: "Тест", Type: TypeBool, Default: "off"},
		},
	}

	if plugin.Name != "Test Plugin" {
		t.Errorf("Name = %q, expected 'Test Plugin'", plugin.Name)
	}
	if plugin.Repo != "test/repo" {
		t.Errorf("Repo = %q, expected 'test/repo'", plugin.Repo)
	}
	if !plugin.RequiresTPM {
		t.Error("RequiresTPM should be true")
	}
	if len(plugin.Requires) != 1 {
		t.Errorf("Requires length = %d, expected 1", len(plugin.Requires))
	}
	if len(plugin.Settings) != 1 {
		t.Errorf("Settings length = %d, expected 1", len(plugin.Settings))
	}
}

func TestPluginSettingStruct(t *testing.T) {
	setting := PluginSetting{
		Key:     "@test-setting",
		DescEN:  "Test setting description",
		DescRU:  "Описание тестовой настройки",
		Type:    TypeChoice,
		Default: "option1",
		Choices: []string{"option1", "option2", "option3"},
	}

	if setting.Key != "@test-setting" {
		t.Errorf("Key = %q, expected '@test-setting'", setting.Key)
	}
	if setting.Type != TypeChoice {
		t.Errorf("Type = %v, expected TypeChoice", setting.Type)
	}
	if len(setting.Choices) != 3 {
		t.Errorf("Choices length = %d, expected 3", len(setting.Choices))
	}
}

func TestLoggingPluginSettings(t *testing.T) {
	plugin, ok := GetPlugin("tmux-plugins/tmux-logging")
	if !ok {
		t.Fatal("logging plugin not found")
	}

	// Should have path settings
	settingKeys := make(map[string]bool)
	for _, s := range plugin.Settings {
		settingKeys[s.Key] = true
	}

	if !settingKeys["@logging-path"] {
		t.Error("logging plugin should have @logging-path setting")
	}
	if !settingKeys["@screen-capture-path"] {
		t.Error("logging plugin should have @screen-capture-path setting")
	}
}
