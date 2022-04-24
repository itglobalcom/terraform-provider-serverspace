---
page_title: "serverspace_domain Resource - terraform-provider-serverspace"
  
---

# serverspace_domain (Resource)

Serverspace Domain resource allows you to create domains.  

## Example Usage

### Create a new domain

```hcl
resource "serverspace_domain" "example_domain" {
  name = "example.com"
}
```
## Schema

### Required

- **name** (String) Domain name.