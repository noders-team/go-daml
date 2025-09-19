package codegen

type Manifest struct {
	Version    string
	CreatedBy  string
	Name       string
	SdkVersion string
	MainDalf   string
	Dalfs      []string
	Format     string
	Encryption string
}
