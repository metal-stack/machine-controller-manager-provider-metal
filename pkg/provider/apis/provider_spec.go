package api

// MetalProviderSpec is the spec to be used while parsing the calls.
// nolint:musttag
type MetalProviderSpec struct {
	Partition string // required
	Size      string // required
	Image     string // required
	Project   string // required
	Network   string // required
	Tags      []string
	SSHKeys   []string
}
