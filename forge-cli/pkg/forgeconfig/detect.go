package forgeconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadSurfaces reads the surfaces field from .forge/config.yaml.
// Returns nil (no error) when surfaces is not configured or empty.
func ReadSurfaces(projectRoot string) (map[string]string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil || cfg == nil {
		return nil, nil
	}
	return cfg.Surfaces, nil
}

// SurfaceTypes extracts deduplicated surface type values from a surfaces map.
// Returns nil for nil/empty maps.
func SurfaceTypes(surfaces map[string]string) []string {
	if len(surfaces) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var types []string
	for _, typ := range surfaces {
		if !seen[typ] {
			seen[typ] = true
			types = append(types, typ)
		}
	}
	return types
}

// ErrMultiInterfaceMigration is returned when auto-migration encounters
// a multi-interface config that cannot be automatically migrated.
type ErrMultiInterfaceMigration struct {
	Interfaces []string
}

func (e *ErrMultiInterfaceMigration) Error() string {
	return fmt.Sprintf(
		"interfaces contains multiple types %v; automatic migration not possible. Run forge init to configure path-level surfaces.",
		e.Interfaces,
	)
}

// MigrateInterfacesToSurfaces performs first-run auto-migration from the legacy
// `interfaces` field to the new `surfaces` field.
//
// Migration rules:
//   - No interfaces field found → no-op, nil
//   - interfaces has single type → auto-write surfaces as scalar, print migration notice, nil
//   - interfaces has multiple types → return ErrMultiInterfaceMigration
func MigrateInterfacesToSurfaces(projectRoot string) error {
	path := configPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read config for migration: %w", err)
	}

	// Parse raw YAML to detect legacy `interfaces` field
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil // malformed config, skip migration
	}

	interfacesNode := findMappingKey(&root, "interfaces")
	if interfacesNode == nil {
		return nil // no legacy field, nothing to migrate
	}

	// Parse the interfaces list
	var interfaces []string
	switch interfacesNode.Kind {
	case yaml.SequenceNode:
		for _, item := range interfacesNode.Content {
			interfaces = append(interfaces, item.Value)
		}
	case yaml.ScalarNode:
		// Single scalar value (unusual but handle it)
		interfaces = []string{interfacesNode.Value}
	default:
		return nil // unexpected format, skip
	}

	if len(interfaces) == 0 {
		return nil
	}

	// Read current config to check if surfaces already exists
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return fmt.Errorf("read config for migration: %w", err)
	}
	if cfg == nil {
		cfg = &Config{}
	}

	// Surfaces already configured — skip migration
	if len(cfg.Surfaces) > 0 {
		return nil
	}

	// Multi-interface: cannot auto-migrate
	if len(interfaces) > 1 {
		return &ErrMultiInterfaceMigration{Interfaces: interfaces}
	}

	// Single interface: auto-migrate to scalar form
	cfg.Surfaces = SurfacesMap{".": interfaces[0]}
	if err := writeConfig(projectRoot, cfg); err != nil {
		return fmt.Errorf("write migrated config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "migrated interfaces [%s] -> surfaces: %s\n", interfaces[0], interfaces[0])
	return nil
}
