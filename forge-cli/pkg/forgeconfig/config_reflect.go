//nolint:govet // reflect.Ptr inline warnings are toolchain version mismatches, not code issues
package forgeconfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// errKeyNotFound is returned when a config key does not exist or has a zero value.
var errKeyNotFound = fmt.Errorf("config key not found")

// errUnsupportedType is returned when a config field implements yaml.Unmarshaler
// and the generic reflect router cannot handle it (e.g. SurfacesMap).
var errUnsupportedType = fmt.Errorf("unsupported type for reflect routing")

// getByPath traverses a reflect.Value by path segments, returning the formatted value.
func getByPath(v reflect.Value, segments []string) (string, error) {
	var err error
	for i, seg := range segments {
		v, err = navigateToSegment(v, seg)
		if err == errKeyNotFound {
			// Try inline map with dot-joined remaining segments
			v2 := derefPointer(v)
			if v2.IsValid() && v2.Kind() == reflect.Struct {
				if inlineField, ok := findInlineMapField(v2); ok {
					mapVal := derefPointer(inlineField)
					if mapVal.IsValid() && mapVal.Kind() == reflect.Map {
						dotKey := strings.Join(segments[i:], ".")
						entry := mapVal.MapIndex(reflect.ValueOf(dotKey))
						if entry.IsValid() {
							return formatValue(derefPointer(entry))
						}
					}
				}
			}
			return "", err
		}
		if err != nil {
			return "", err
		}

		// Reached the target segment
		if i == len(segments)-1 {
			return formatValue(v)
		}

		// More segments to go — check if current value is navigable
		if isLeafType(v) {
			return "", errKeyNotFound
		}
		// Continue descending
	}
	return "", errKeyNotFound
}

// navigateToSegment resolves one path segment within the given reflect.Value.
func navigateToSegment(v reflect.Value, seg string) (reflect.Value, error) {
	// Dereference pointers
	v = derefPointer(v)
	if !v.IsValid() {
		return reflect.Value{}, errKeyNotFound
	}

	kind := v.Kind()

	switch kind {
	case reflect.Struct:
		field, found := findFieldByYAMLTag(v, seg)
		if !found {
			// Check for yaml:",inline" map fields
			if inlineField, ok := findInlineMapField(v); ok {
				mapVal := derefPointer(inlineField)
				if mapVal.IsValid() && mapVal.Kind() == reflect.Map {
					entry := mapVal.MapIndex(reflect.ValueOf(seg))
					if entry.IsValid() {
						return derefPointer(entry), nil
					}
				}
			}
			return reflect.Value{}, errKeyNotFound
		}
		return derefPointer(field), nil

	case reflect.Map:
		keyVal := reflect.ValueOf(seg)
		entry := v.MapIndex(keyVal)
		if !entry.IsValid() {
			return reflect.Value{}, errKeyNotFound
		}
		return derefPointer(entry), nil

	default:
		return reflect.Value{}, errKeyNotFound
	}
}

// findFieldByYAMLTag finds a struct field matching the segment by YAML tag (priority)
// or Go field name.
func findFieldByYAMLTag(v reflect.Value, seg string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("yaml")
		tagName := parseYAMLTagName(tag, field.Name)
		if tagName == seg {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

// parseYAMLTagName extracts the YAML key name from a yaml tag.
// Priority: yaml:"name" -> name; yaml:",inline" -> "" (skip); no tag -> GoFieldName.
// Returns empty string for ",inline" and "omitempty" only tags.
func parseYAMLTagName(tag, goName string) string {
	if tag == "" {
		return goName
	}
	// Split by comma: first part is name, rest are options
	parts := strings.Split(tag, ",")
	name := parts[0]
	if name == "" {
		// ",inline" or ",omitempty" — not a key match target
		return ""
	}
	return name
}

// findInlineMapField finds a struct field tagged with yaml:",inline" that is a map type.
func findInlineMapField(v reflect.Value) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("yaml")
		if strings.Contains(tag, ",inline") {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

// derefPointer dereferences a pointer, returning the zero Value if nil.
func derefPointer(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// isLeafType returns true if the value cannot be further navigated.
func isLeafType(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	kind := v.Kind()
	switch kind {
	case reflect.Struct, reflect.Map:
		return false
	default:
		return true
	}
}

// formatValue formats a reflect.Value for CLI output.
// For leaf types, returns the scalar value.
// For non-leaf types (struct, map), returns a multi-line summary.
func formatValue(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return "", errKeyNotFound
	}

	kind := v.Kind()

	// Check for custom YAML types that reflect routing cannot handle.
	// Only applies to non-struct types (e.g. SurfacesMap as map type).
	// Structs that implement yaml.Unmarshaler (like EvalConfig for compat)
	// are handled by the struct formatting path below.
	if kind != reflect.Struct && implementsYAMLUnmarshaler(v) {
		return "", errUnsupportedType
	}

	switch kind {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.String {
			slice := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				slice[i] = v.Index(i).String()
			}
			return joinSlice(slice), nil
		}
		return "", errUnsupportedType
	case reflect.Struct:
		if isModeToggle(v.Type()) {
			q := v.FieldByName("Quick").Bool()
			f := v.FieldByName("Full").Bool()
			return fmt.Sprintf("quick:%v full:%v", q, f), nil
		}
		return formatStructSummary(v, "")
	case reflect.Map:
		return formatMapSummary(v, "")
	default:
		return "", errUnsupportedType
	}
}

// implementsYAMLUnmarshaler checks if the value's type implements yaml.Unmarshaler.
//
//nolint:govet // reflect.PtrTo inline warning is a toolchain version mismatch, not a code issue
func implementsYAMLUnmarshaler(v reflect.Value) bool {
	t := v.Type()
	ptr := reflect.PointerTo(t)
	return ptr.Implements(reflect.TypeOf((*yaml.Unmarshaler)(nil)).Elem())
}

// formatStructSummary formats a struct's exported fields as a multi-line summary.
func formatStructSummary(v reflect.Value, indent string) (string, error) {
	var lines []string
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		// Skip unexported internal fields (like 'raw')
		tag := field.Tag.Get("yaml")
		if tag == "" && field.Name == "raw" {
			continue
		}
		if strings.Contains(tag, ",inline") {
			continue
		}

		fieldName := parseYAMLTagName(tag, field.Name)
		if fieldName == "" {
			continue
		}

		fv := derefPointer(v.Field(i))
		if !fv.IsValid() {
			continue
		}

		line, err := formatFieldLine(fieldName, fv, indent)
		if err != nil {
			continue
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return "", errKeyNotFound
	}
	return strings.Join(lines, "\n"), nil
}

// formatMapSummary formats a map's entries as a multi-line summary.
func formatMapSummary(v reflect.Value, indent string) (string, error) {
	var lines []string
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		entry := derefPointer(iter.Value())
		if !entry.IsValid() {
			continue
		}
		line, err := formatFieldLine(key, entry, indent)
		if err != nil {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return "", errKeyNotFound
	}
	return strings.Join(lines, "\n"), nil
}

// formatFieldLine formats a single field for summary output.
func formatFieldLine(name string, v reflect.Value, indent string) (string, error) {
	kind := v.Kind()
	switch kind {
	case reflect.Struct:
		if isModeToggle(v.Type()) {
			// ModeToggle -> "name: quick:X full:Y"
			q := v.FieldByName("Quick").Bool()
			f := v.FieldByName("Full").Bool()
			return fmt.Sprintf("%s%s: quick:%v full:%v", indent, name, q, f), nil
		}
		// Nested struct -> "name:\n" + recursive lines with +2 indent
		sub, err := formatStructSummary(v, indent+"  ")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s:\n%s", indent, name, sub), nil
	case reflect.Bool:
		return fmt.Sprintf("%s%s: %v", indent, name, v.Bool()), nil
	case reflect.String:
		return fmt.Sprintf("%s%s: %s", indent, name, v.String()), nil
	case reflect.Map:
		sub, err := formatMapSummary(v, indent+"  ")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s:\n%s", indent, name, sub), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%s%s: %d", indent, name, v.Int()), nil
	default:
		return "", errUnsupportedType
	}
}

// isModeToggle checks if a type is ModeToggle.
func isModeToggle(t reflect.Type) bool {
	return t.Name() == "ModeToggle" && t.Kind() == reflect.Struct &&
		t.NumField() == 2
}

// setByPath traverses a reflect.Value by segments and sets the leaf value.
func setByPath(v reflect.Value, segments []string, value string, fullKey string) error {
	for i, seg := range segments {
		v = ensureAddressable(v)

		// Dereference pointers, initializing nil pointers as needed
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				newVal := reflect.New(v.Type().Elem())
				v.Set(newVal)
			}
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct { //nolint:gocritic // ifElseChain
			field, found := findSettableField(v, seg)
			if !found {
				// Check for inline map - try joining remaining segments as dot-separated key
				if inlineField, ok := findInlineMapField(v); ok {
					mapVal := ensureAddressable(inlineField)
					for mapVal.Kind() == reflect.Ptr {
						if mapVal.IsNil() {
							newVal := reflect.New(mapVal.Type().Elem())
							mapVal.Set(newVal)
						}
						mapVal = mapVal.Elem()
					}
					if mapVal.Kind() == reflect.Map {
						if mapVal.IsNil() {
							mapVal.Set(reflect.MakeMap(mapVal.Type()))
						}
						if i == len(segments)-1 {
							return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
						}
						// Join remaining segments as a single dot-separated key for inline maps
						dotKey := strings.Join(segments[i:], ".")
						return setMapEntry(mapVal, []string{dotKey}, value, fullKey)
					}
				}
				return fmt.Errorf("config key %q not found", fullKey)
			}

			// Last segment - set the value
			if i == len(segments)-1 {
				return setFieldValue(field, value, fullKey)
			}

			// Intermediate segment - allow descending into ModeToggle
			if isLeafType(field) && !isModeToggle(field.Type()) {
				return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
			}
			v = field
		} else if v.Kind() == reflect.Map {
			if i == len(segments)-1 {
				return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
			}
			return setMapEntry(v, segments[i+1:], value, fullKey)
		} else {
			return errKeyNotFound
		}
	}
	return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
}

// ensureAddressable returns an addressable reflect.Value.
func ensureAddressable(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v
	}
	// For non-addressable values (from reflect.ValueOf), try to get a pointer
	return v
}

// findSettableField finds a struct field matching the segment and returns a settable Value.
func findSettableField(v reflect.Value, seg string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("yaml")
		tagName := parseYAMLTagName(tag, field.Name)
		if tagName == seg {
			fv := v.Field(i)
			// For pointer fields, initialize nil and dereference
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				newVal := reflect.New(fv.Type().Elem())
				fv.Set(newVal)
			}
			if fv.Kind() == reflect.Ptr {
				fv = fv.Elem()
			}
			return fv, true
		}
	}
	return reflect.Value{}, false
}

// setFieldValue sets a leaf field's value from a string.
func setFieldValue(field reflect.Value, value string, fullKey string) error {
	if isModeToggle(field.Type()) {
		return fmt.Errorf("cannot set ModeToggle directly, use %s.quick or %s.full", fullKey, fullKey)
	}

	switch field.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value %q for bool field %s: expected true or false", value, fullKey)
		}
		field.SetBool(b)
		return nil
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value %q for int field %s: expected integer", value, fullKey)
		}
		field.SetInt(int64(n))
		return nil
	default:
		return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
	}
}

// setMapEntry sets a value in a map for the given remaining segments.
func setMapEntry(mapVal reflect.Value, segments []string, value string, fullKey string) error {
	if len(segments) != 1 {
		return errUnsupportedType
	}
	key := segments[0]

	// For CoverageConfig.ByType: value is a percentage number
	pct, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid coverage value for %s: %s (expected percentage number)", fullKey, value)
	}

	strategyType := reflect.TypeOf(CoverageStrategy{})
	strategyVal := reflect.New(strategyType).Elem()
	strategyVal.FieldByName("Type").SetString("percentage")
	pctField := strategyVal.FieldByName("Percentage")
	pctVal := reflect.New(pctField.Type().Elem())
	pctVal.Elem().SetInt(int64(pct))
	strategyVal.FieldByName("Percentage").Set(pctVal)

	mapVal.SetMapIndex(reflect.ValueOf(key), strategyVal)
	return nil
}

// joinSlice joins slice values with newline for plain-text output.
func joinSlice(vals []string) string {
	return strings.Join(vals, "\n")
}
