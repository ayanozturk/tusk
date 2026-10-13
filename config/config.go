package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/ayanozturk/tusk/overrides"

	"gopkg.in/yaml.v2"
)

var DefaultConfigFilenames = []string{
	"tusk.yaml",
}

const DefaultConfigContent = `path: .
# Composer vendor is always indexed for symbol resolution and never
# type-checked. Extra includes are additional index-only trees.
includes: []
extensions:
  - php
ignore:
  - vendor
  - node_modules
  - cache
  - .git
rules:
  - all
analysis_level: 0
`

// WriteDefaultConfig creates filename with default config content.
// Returns an error if the file already exists.
func WriteDefaultConfig(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("%s already exists", filename)
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(filename, []byte(DefaultConfigContent), 0644)
}

type Config struct {
	Path          string                  `yaml:"path"`
	Includes      []string                `yaml:"includes"`
	Extensions    []string                `yaml:"extensions"`
	Ignore        []string                `yaml:"ignore"`
	Rules         []string                `yaml:"rules"`
	AnalysisLevel *int                    `yaml:"analysis_level"`
	Overrides     overrides.RuleOverrides `yaml:"overrides"`
}

func DiscoverConfig(dir string) (string, error) {
	if dir == "" {
		dir = "."
	}

	for _, name := range DefaultConfigFilenames {
		candidate := filepath.Join(dir, name)
		info, err := os.Stat(candidate)
		if err == nil {
			if info.IsDir() {
				continue
			}
			return candidate, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", fmt.Errorf("no config file found; checked %s", DefaultConfigFilenames)
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	if config.AnalysisLevel != nil && *config.AnalysisLevel < 0 {
		return nil, fmt.Errorf("analysis_level must be zero or greater")
	}

	return &config, nil
}

func PrintEffectiveConfig(w io.Writer, cfg *Config, source string) {
	fmt.Fprintf(w, "config_file: %s\n", quoteYAMLString(source))
	fmt.Fprintf(w, "path: %s\n", quoteYAMLString(cfg.Path))
	writeStringList(w, "includes", cfg.Includes)
	writeStringList(w, "extensions", cfg.Extensions)
	writeStringList(w, "ignore", cfg.Ignore)
	writeStringList(w, "rules", cfg.Rules)
	if cfg.AnalysisLevel == nil {
		fmt.Fprintln(w, "analysis_level: null")
	} else {
		fmt.Fprintf(w, "analysis_level: %d\n", *cfg.AnalysisLevel)
	}
	writeOverrides(w, cfg.Overrides)
}

func writeStringList(w io.Writer, name string, values []string) {
	if len(values) == 0 {
		fmt.Fprintf(w, "%s: []\n", name)
		return
	}
	fmt.Fprintf(w, "%s:\n", name)
	for _, value := range values {
		fmt.Fprintf(w, "  - %s\n", quoteYAMLString(value))
	}
}

func writeOverrides(w io.Writer, ruleOverrides overrides.RuleOverrides) {
	if len(ruleOverrides) == 0 {
		fmt.Fprintln(w, "overrides: {}")
		return
	}

	codes := make([]string, 0, len(ruleOverrides))
	for code := range ruleOverrides {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	fmt.Fprintln(w, "overrides:")
	for _, code := range codes {
		fmt.Fprintf(w, "  %s:\n", quoteYAMLString(code))
		writeNestedStringList(w, "classes", ruleOverrides[code].Classes)
	}
}

func writeNestedStringList(w io.Writer, name string, values []string) {
	if len(values) == 0 {
		fmt.Fprintf(w, "    %s: []\n", name)
		return
	}
	fmt.Fprintf(w, "    %s:\n", name)
	for _, value := range values {
		fmt.Fprintf(w, "      - %s\n", quoteYAMLString(value))
	}
}

func quoteYAMLString(value string) string {
	return strconv.Quote(value)
}

func ignoreDirSet(ignore []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ignore))
	for _, dir := range ignore {
		set[dir] = struct{}{}
	}
	return set
}

func extSet(extensions []string) map[string]struct{} {
	set := make(map[string]struct{}, len(extensions))
	for _, ext := range extensions {
		set["."+ext] = struct{}{}
	}
	return set
}

func walkForFiles(root string, ignoreDirs, allowedExts map[string]struct{}) ([]string, error) {
	root = filepath.Clean(root)
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if filepath.Clean(path) != root {
				if _, ignored := ignoreDirs[d.Name()]; ignored {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if _, allowed := allowedExts[filepath.Ext(path)]; allowed {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func GetFilesToScan(config *Config) ([]string, error) {
	return walkForFiles(config.Path, ignoreDirSet(config.Ignore), extSet(config.Extensions))
}

// GetIncludeFiles walks Composer vendor under config.Path plus any extra
// directories listed in config.Includes. Those files are parsed and indexed
// for cross-file symbol resolution but are never scanned for diagnostics.
// The project-root vendor tree is always included when it exists, even if
// it is also listed in ignore (host scans still skip it).
func GetIncludeFiles(config *Config) ([]string, error) {
	ignoreDirs := ignoreDirSet(config.Ignore)
	allowedExts := extSet(config.Extensions)

	var files []string
	seen := make(map[string]struct{})
	for _, dir := range includeIndexDirs(config) {
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !info.IsDir() {
			if _, allowed := allowedExts[filepath.Ext(dir)]; allowed {
				if _, exists := seen[dir]; !exists {
					seen[dir] = struct{}{}
					files = append(files, dir)
				}
			}
			continue
		}
		found, err := walkForFiles(dir, ignoreDirs, allowedExts)
		if err != nil {
			return nil, err
		}
		for _, file := range found {
			if _, exists := seen[file]; exists {
				continue
			}
			seen[file] = struct{}{}
			files = append(files, file)
		}
	}
	return files, nil
}

func includeIndexDirs(config *Config) []string {
	if config == nil {
		return nil
	}
	seen := make(map[string]struct{})
	var dirs []string
	add := func(dir string) {
		if dir == "" {
			return
		}
		cleaned := filepath.Clean(dir)
		if _, exists := seen[cleaned]; exists {
			return
		}
		seen[cleaned] = struct{}{}
		dirs = append(dirs, cleaned)
	}
	root := config.Path
	if root == "" {
		root = "."
	}
	add(filepath.Join(root, "vendor"))
	for _, dir := range config.Includes {
		add(dir)
	}
	for _, dir := range composerAutoloadDirs(root) {
		add(dir)
	}
	return dirs
}

type composerAutoload struct {
	PSR4     map[string]json.RawMessage `json:"psr-4"`
	PSR0     map[string]json.RawMessage `json:"psr-0"`
	Classmap []string                   `json:"classmap"`
}

type composerJSON struct {
	Autoload    composerAutoload `json:"autoload"`
	AutoloadDev composerAutoload `json:"autoload-dev"`
}

// composerAutoloadDirs walks upward from root looking for the nearest
// composer.json and returns every directory referenced by its autoload /
// autoload-dev psr-4, psr-0, and classmap entries, resolved relative to the
// composer.json location. This lets legacy monorepo layouts (e.g. a sibling
// tree holding PSR-0/classmap-style underscore class names) get indexed for
// symbol resolution even when they sit outside config.Path.
func composerAutoloadDirs(root string) []string {
	composerDir, raw := findComposerJSON(root)
	if raw == nil {
		return nil
	}

	var parsed composerJSON
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil
	}

	var dirs []string
	collect := func(al composerAutoload) {
		for _, targets := range al.PSR4 {
			dirs = append(dirs, resolveAutoloadTargets(composerDir, targets)...)
		}
		for _, targets := range al.PSR0 {
			dirs = append(dirs, resolveAutoloadTargets(composerDir, targets)...)
		}
		for _, dir := range al.Classmap {
			dirs = append(dirs, filepath.Join(composerDir, dir))
		}
	}
	collect(parsed.Autoload)
	collect(parsed.AutoloadDev)
	return dirs
}

// resolveAutoloadTargets decodes a psr-4/psr-0 autoload map value, which may
// be a single path string or an array of path strings.
func resolveAutoloadTargets(composerDir string, raw json.RawMessage) []string {
	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		return []string{filepath.Join(composerDir, single)}
	}
	var multiple []string
	if err := json.Unmarshal(raw, &multiple); err == nil {
		dirs := make([]string, 0, len(multiple))
		for _, dir := range multiple {
			dirs = append(dirs, filepath.Join(composerDir, dir))
		}
		return dirs
	}
	return nil
}

// findComposerJSON walks upward from dir looking for composer.json, stopping
// at the filesystem root. It returns the directory containing composer.json
// and its raw bytes, or ("", nil) if none is found.
func findComposerJSON(dir string) (string, []byte) {
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	} else {
		dir = filepath.Clean(dir)
	}
	for {
		candidate := filepath.Join(dir, "composer.json")
		if data, err := os.ReadFile(candidate); err == nil {
			return dir, data
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

// StreamFilesToScan walks the configured path in a background goroutine and
// streams discovered file paths into the returned channel, which is closed
// when the walk completes. This allows callers to overlap I/O and parsing
// with the directory walk rather than waiting for the full file list first.
// Ignore-name matching matches GetFilesToScan / walkForFiles: the walk root
// itself is never skipped even when its basename appears in config.Ignore.
func StreamFilesToScan(config *Config) <-chan string {
	ignoreDirs := ignoreDirSet(config.Ignore)
	allowedExts := extSet(config.Extensions)
	root := filepath.Clean(config.Path)

	ch := make(chan string, 256)
	go func() {
		defer close(ch)
		filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error { //nolint:errcheck
			if err != nil {
				return nil // skip unreadable entries
			}
			if d.IsDir() {
				if filepath.Clean(path) != root {
					if _, ignored := ignoreDirs[d.Name()]; ignored {
						return filepath.SkipDir
					}
				}
				return nil
			}
			if _, allowed := allowedExts[filepath.Ext(path)]; allowed {
				ch <- path
			}
			return nil
		})
	}()
	return ch
}
