package interfaces

type LoggerConfig struct {
	Level        string      `yaml:"level" default:"verbose" desc:"log:level"`
	Format       string      `yaml:"format" default:"default" desc:"log:format"`
	HideConsole  bool        `yaml:"hideconsole" default:"false" desc:"log:hideconsole"`
	DisableColor bool        `yaml:"disablecolor" default:"false" desc:"log:disablecolor"`
	File         LoggingFile `yaml:"file"`
}

type LoggingFile struct {
	Enable   bool   `yaml:"enable" default:"false" desc:"log:file:enable"`
	Output   string `yaml:"output" default:"./logs/app.log" desc:"log:file:output"`
	MaxSize  int    `yaml:"maxsize" default:"100" desc:"log:file:maxsize"`
	MaxAge   int    `yaml:"maxage" default:"28" desc:"log:file:maxage"`
	Compress bool   `yaml:"compress" default:"false" desc:"log:file:compress"`
	Json     bool   `yaml:"json" default:"false" desc:"log:file:json"`
}

var LoggerConfigManual = map[string]string{
	"log:level": `logging level, valid value is
		- trace
		- verbose
		- info
		- warning
		- error
	`,
	"log:format": `logging output format, valid value is
		- default
		- json
	`,
	"log:hideconsole":   `Don't print log to console`,
	"log:disablecolor":  `Don't use color in log`,
	"log:file:enable":   `Enable writing log to file`,
	"log:file:output":   `File log output location`,
	"log:file:maxsize":  `Max log file size (in megabytes) before retain`,
	"log:file:maxage":   `Max log file age (in days) before retain`,
	"log:file:compress": `If true, old file will be compressed`,
	"log:file:json":     `If true, always write in json format`,
}

func SetManual(man map[string]string) map[string]string {
	for k, v := range LoggerConfigManual {
		man[k] = v
	}
	return man
}
