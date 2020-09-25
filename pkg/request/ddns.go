package request

type DDnsRequest interface {
	Update(domain string, recordType string, subDomain string, ipAddress string) error
}
