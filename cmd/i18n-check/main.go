// i18n-check scans source files for i18n keys used via the locales package
// and reports which keys are missing from each locale YAML, including the
// default locale. Run with --fix to insert missing keys as TODO placeholders.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// matches l.T("key"), l.Tc("key"), l.N("key"), l.Nc("key")
var i18nCallRe = regexp.MustCompile(`l\.[TN]c?\(\s*"([^"]+)"`)

// matches a top-level section key under the locale root: "  sectionName:"
var sectionLineRe = regexp.MustCompile(`^  [a-zA-Z_][a-zA-Z0-9_]*:\s*$`)

// matches: SupportedLanguages = []string{"en", "fr"}
var supportedLangsRe = regexp.MustCompile(`SupportedLanguages\s*=\s*\[\]string\{([^}]+)\}`)
var langStringRe = regexp.MustCompile(`"([^"]+)"`)

func main() {
	fix := flag.Bool("fix", false, "insert missing keys into YAML files with a TODO placeholder")
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	usedKeys, err := extractKeys(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error extracting keys: %v\n", err)
		os.Exit(1)
	}

	supportedLangs, err := loadSupportedLanguages(filepath.Join(root, "locales", "locales.go"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading SupportedLanguages: %v\n", err)
		os.Exit(1)
	}

	localeKeys, err := loadLocales(filepath.Join(root, "locales", "lang"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading locales: %v\n", err)
		os.Exit(1)
	}

	hasMissing := false
	for _, lang := range supportedLangs {
		keys := localeKeys[lang] // nil if YAML file is missing — all keys reported

		var missing []string
		for key := range usedKeys {
			if val, exists := keys[key]; !exists || val == "TODO" {
				missing = append(missing, key)
			}
		}
		sort.Strings(missing)

		if len(missing) == 0 {
			continue
		}

		hasMissing = true
		fmt.Printf("Missing in %q (%d):\n", lang, len(missing))
		for _, k := range missing {
			fmt.Printf("  - %s\n", k)
		}

		if *fix {
			yamlPath := filepath.Join(root, "locales", "lang", lang+".yaml")
			if err := addMissingKeys(yamlPath, lang, missing); err != nil {
				fmt.Fprintf(os.Stderr, "  error fixing %s: %v\n", lang, err)
			} else {
				fmt.Printf("  → added %d key(s) to %s\n", len(missing), yamlPath)
			}
		}
		fmt.Println()
	}

	if !hasMissing {
		fmt.Println("All locales are complete.")
	} else if !*fix {
		os.Exit(1)
	}
}

// addMissingKeys inserts missing dot-notation keys into the locale's YAML file
// under the appropriate nested section, with a "TODO" placeholder value.
func addMissingKeys(path, lang string, keys []string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}

	localeRoot := findLocaleRoot(&doc, lang)
	if localeRoot == nil {
		return fmt.Errorf("locale %q not found in %s", lang, path)
	}

	for _, key := range keys {
		insertKey(localeRoot, strings.Split(key, "."))
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}

	// yaml.Encoder prepends "---\n"; strip it to match the original file style.
	result := bytes.TrimPrefix(buf.Bytes(), []byte("---\n"))

	// Re-insert blank lines between top-level sections (lost during encoding).
	result = normalizeSectionBlankLines(result)

	return os.WriteFile(path, result, 0644)
}

// findLocaleRoot returns the mapping node under the locale top-level key (e.g. "en:" or "fr:").
func findLocaleRoot(doc *yaml.Node, lang string) *yaml.Node {
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == lang {
			return root.Content[i+1]
		}
	}
	return nil
}

// insertKey navigates the dot-notation path inside a mapping node, creating
// intermediate maps as needed, and appends a "TODO" scalar at the leaf.
// Existing keys at any level are left unchanged.
func insertKey(node *yaml.Node, path []string) {
	if node.Kind != yaml.MappingNode || len(path) == 0 {
		return
	}

	key := path[0]
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			if len(path) > 1 {
				insertKey(node.Content[i+1], path[1:])
			}
			// leaf already exists — leave it untouched
			return
		}
	}

	// Key not found at this level — append it.
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
	if len(path) == 1 {
		valueNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "TODO"}
		node.Content = append(node.Content, keyNode, valueNode)
	} else {
		mapNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		node.Content = append(node.Content, keyNode, mapNode)
		insertKey(mapNode, path[1:])
	}
}

// normalizeSectionBlankLines ensures exactly one blank line precedes each
// top-level section key (e.g. "  nav:", "  products:") in the encoded YAML,
// restoring the readability that yaml.Encoder drops when re-serialising nodes.
func normalizeSectionBlankLines(data []byte) []byte {
	lines := strings.Split(string(data), "\n")

	// Strip any blank lines that sit immediately before a section header so we
	// can re-add exactly one regardless of what the encoder produced.
	cleaned := make([]string, 0, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			// Look ahead to the next non-blank line.
			next := ""
			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) != "" {
					next = lines[j]
					break
				}
			}
			if sectionLineRe.MatchString(next) {
				continue // drop this blank; we'll re-add one below
			}
		}
		cleaned = append(cleaned, line)
	}

	// Insert exactly one blank line before each section header except the first.
	result := make([]string, 0, len(cleaned)+10)
	firstSection := true
	for _, line := range cleaned {
		if sectionLineRe.MatchString(line) {
			if !firstSection {
				result = append(result, "")
			}
			firstSection = false
		}
		result = append(result, line)
	}

	return []byte(strings.Join(result, "\n"))
}

// loadSupportedLanguages parses the SupportedLanguages slice from locales.go.
func loadSupportedLanguages(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := supportedLangsRe.FindSubmatch(content)
	if m == nil {
		return nil, fmt.Errorf("SupportedLanguages not found in %s", path)
	}
	var langs []string
	for _, lm := range langStringRe.FindAllSubmatch(m[1], -1) {
		langs = append(langs, string(lm[1]))
	}
	return langs, nil
}

// extractKeys walks .go and .templ source files and collects all i18n keys.
func extractKeys(root string) (map[string]bool, error) {
	keys := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || path == filepath.Join(root, "cmd") {
				return filepath.SkipDir
			}
			if name != "." && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()
		ext := filepath.Ext(name)
		if ext != ".go" && ext != ".templ" {
			return nil
		}
		if strings.HasSuffix(name, "_templ.go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, m := range i18nCallRe.FindAllSubmatch(content, -1) {
			keys[string(m[1])] = true
		}
		return nil
	})
	return keys, err
}

// loadLocales reads all YAML files in the locales directory and returns
// a map of locale code → flattened dot-notation keys with their string values.
func loadLocales(dir string) (map[string]map[string]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]string{}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}

		var raw map[string]any
		if err := yaml.Unmarshal(content, &raw); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", e.Name(), err)
		}

		for lang, v := range raw {
			subtree, ok := v.(map[string]any)
			if !ok {
				continue
			}
			result[lang] = flattenKeys(subtree, "")
		}
	}
	return result, nil
}

// flattenKeys recursively flattens a YAML map into dot-notation keys mapped to
// their string values. Plural forms (maps with "one"/"other") are stored under
// the parent key with a synthetic value so they register as present.
func flattenKeys(m map[string]any, prefix string) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		if sub, ok := v.(map[string]any); ok {
			_, hasOne := sub["one"]
			_, hasOther := sub["other"]
			if hasOne || hasOther {
				result[key] = fmt.Sprintf("%v", sub)
			} else {
				maps.Copy(result, flattenKeys(sub, key))
			}
		} else {
			result[key] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
