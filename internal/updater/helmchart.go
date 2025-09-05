package updater

import (
	"regexp"
	"strings"
)

// HelmChart creates an updater that modifies the version field in Chart.yaml files
func HelmChart() Updater {
	return helmchart{}
}

type helmchart struct{}

func (h helmchart) Files() []string {
	return []string{"Chart.yaml"}
}

func (h helmchart) CreateNewFiles() bool {
	return false
}

func (h helmchart) Update(info ReleaseInfo) func(content string) (string, error) {
	return func(content string) (string, error) {
		// We strip the "v" prefix to match Helm versioning convention
		version := strings.TrimPrefix(info.Version, "v")

		// Split content into lines, process each line individually, then rejoin
		lines := strings.Split(content, "\n")
		versionRegex := regexp.MustCompile(`^(\s*version\s*:\s*)([^\s#]*)(.*)$`)
		
		foundVersion := false
		for i, line := range lines {
			if versionRegex.MatchString(line) {
				foundVersion = true
				matches := versionRegex.FindStringSubmatch(line)
				prefix := matches[1]
				currentVersion := matches[2]
				suffix := matches[3]
				
				// If the current version is empty and the prefix doesn't end with a space, add one
				if currentVersion == "" && !strings.HasSuffix(prefix, " ") {
					prefix += " "
				}
				
				lines[i] = prefix + version + suffix
			}
		}
		
		// If no version field was found, return content unchanged
		if !foundVersion {
			return content, nil
		}
		
		return strings.Join(lines, "\n"), nil
	}
}