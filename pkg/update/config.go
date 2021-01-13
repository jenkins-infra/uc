package update

type Config struct {
	ConnectionCheckUrl  string                     `json:"connectionCheckUrl"`
	Core                ConfigInfo                 `json:"core"`
	Deprecations        map[string]DeprecationInfo `json:"deprecations"`
	Id                  string                     `json:"id"`
	Plugins             map[string]PluginInfo      `json:"plugins"`
	UpdateCenterVersion string                     `json:"updateCenterVersion"`
}

type ConfigInfo struct {
	BuildDate string `json:"buildDate"`
	Name      string `json:"core"`
	Sha1      string `json:"sha1"`
	Sha256    string `json:"sha256"`
	Url       string `json:"url"`
	Version   string `json:"version"`
}

type DeprecationInfo struct {
	Url string `json:"url"`
}

type PluginInfo struct {
	BuildDate    string       `json:"buildDate"`
	Name         string       `json:"name"`
	Sha1         string       `json:"sha1"`
	Sha256       string       `json:"sha256"`
	Url          string       `json:"url"`
	Version      string       `json:"version"`
	RequiredCore string       `json:"requiredCore"`
	Dependencies []Dependency `json:"dependencies"`
}

type WarningInfo struct {
	Id       string        `json:"id"`
	Message  string        `json:"message"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Url      string        `json:"url"`
	Versions []VersionInfo `json:"versions"`
}

type VersionInfo struct {
	LastVersion string `json:"lastVersion"`
	Pattern     string `json:"pattern"`
}

type Dependency struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
	Version  string `json:"version"`
}
