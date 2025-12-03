package config

var TomlConfig = []byte(`
# This is a TOML document.
ssccprefix = "146024436369"
ssccstartnumber = 1
marktemplate = ""
packtemplate = ""
chunksize = 250
party="Б1"

[layouts]
timelayout = "2006-01-02T15:04:05-0700"
timelayoutclear = "2006.01.02 15:04:05"
timelayoutday = "2006.01.02"
timelayoututc = "2006-01-02T15:04:05"

`)

type Configuration struct {
	SsccPrefix      string              `mapstructure:"ssccprefix" toml:"ssccprefix" json:"ssccprefix"`
	SsccStartNumber int                 `mapstructure:"ssccstartnumber" toml:"ssccstartnumber" json:"ssccstartnumber"`
	MarkTemplate    string              `mapstructure:"marktemplate" toml:"marktemplate" json:"marktemplate"`
	PackTemplate    string              `mapstructure:"packtemplate" toml:"packtemplate" json:"packtemplate"`
	Party           string              `mapstructure:"party" toml:"party" json:"party"`
	ChunkSize       int                 `mapstructure:"chunksize" toml:"chunksize" json:"chunksize"`
	Layouts         LayoutConfiguration `mapstructure:"layouts" toml:"layouts" json:"layouts"`
}

type LayoutConfiguration struct {
	TimeLayout      string `mapstructure:"timelayout" toml:"timelayout" json:"timelayout"`
	TimeLayoutClear string `mapstructure:"timelayoutclear" toml:"timelayoutclear" json:"timelayoutclear"`
	TimeLayoutDay   string `mapstructure:"timelayoutday" toml:"timelayoutday" json:"timelayoutday"`
	TimeLayoutUTC   string `mapstructure:"timelayoututc" toml:"timelayoututc" json:"timelayoututc"`
}
