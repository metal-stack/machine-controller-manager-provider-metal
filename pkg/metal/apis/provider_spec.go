package api

// MetalProviderSpec is the spec to be used while parsing the calls.
type MetalProviderSpec struct {
	Partition  string      `json:"partition,omitempty"` // required
	Size       string      `json:"size,omitempty"`      // required
	Image      string      `json:"image,omitempty"`     // required
	Project    string      `json:"project,omitempty"`   // required
	Network    string      `json:"network,omitempty"`   // required
	Tags       []string    `json:"tags,omitempty"`
	SSHKeys    []string    `json:"sshKeys,omitempty"`
	DNSServers []DNSServer `json:"dnsServers,omitempty"`
	NTPServers []NTPServer `json:"ntpServers,omitempty"`
}

type DNSServer struct {
	IP string `json:"ip"`
}

type NTPServer struct {
	Address string `json:"address"`
}
