package forgeconfig

// ReadInterfaces resolves the effective interface types for a project.
// Reads only from config.yaml interfaces field — no auto-detection.
// Returns nil (no error) when interfaces is not configured.
func ReadInterfaces(projectRoot string) ([]string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		return nil, nil
	}
	return cfg.Interfaces, nil
}
