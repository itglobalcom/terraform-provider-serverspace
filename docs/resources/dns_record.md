---
page_title: "serverspace_dns_record Resource - terraform-provider-serverspace"
  
---

# serverspace_dns_record (Resource)

Serverspace Domain Name Service record resource allows you to create records.

## Example Usage

### Create a new DNS record

```hcl
resource "serverspace_dns_record" "example_host" {
  name = "a.example.com"
  type = "A"
  ip = "93.184.216.38"
  ttl = "2h"
}
```
## Schema

### Required

- **name** (String) Domain name.
- **type** (String) Record type. "A" "AAAA" "MX" "CNAME" "NS" "TXT" "SRV" types are supported.
- **ttl** (String) Record TTL. "1s" "5s" "30s" "1m" "5m" "10m" "15m" "30m" "1h" "2h" "6h" "12h" "1d" values are supported.

### Optional

- **id** (String) The ID of the record.
- **ip** (String) IP address, used for A and AAAA records.
- **mail_host** (String) Mail server, used for MX records.
- **priority** (Number) Record priority, used for MX and SRV records. The priority must be a number between 0 and 65535.
- **canonical_name** (String) Canonical name, used for CNAME records.
- **name_server_host** (String) Domain name of a host, used for NS records.
- **text** (String) Text, used for TXT records.
- **protocol** (String) Protocol, used for SRV record. "TCP" "UDP" "TLS" protocols are supported.
- **service** (String) Service name, used for SRV records.
- **weight** (Number) Record weight, used for SRV records.
- **port** (Number) Port, used for SRV records.
- **terget** (String) The canonical name of the machine providing the service, used for the SRV records.

