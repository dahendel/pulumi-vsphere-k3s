package k3s

type VMS map[string]*VMExportData

type VMExportData struct {
	UUid string
	IPAddress string
	Addresses []string
}