# Terraform Provider Serverspace

## To use terraform provider, first of all:
1. Create an API key for the project that will work with Terraform.
2. Create and navigate to the directory which will be used to work with Terraform Provider.
3. Create and open the provider.tf configuration file.
4. Insert the provider information in the file, where <api key> is your API key, and save the changes.

## To use provider from terraform registry
1. Use template of configuration file:
```
terraform {
  required_providers {
    serverspace = {
      source = "itglobalcom/serverspace"
      version = "0.2.3"
    }
  }
}

variable "s2_token" {
  type = string
  default = "<api key>"
}

provider "serverspace" {
  key = var.s2_token
}
```

## To use local provider
1. Download source code of provider:
```
git clone https://github.com/itglobalcom/terraform-provider-serverspace.git
```

2. Open the provider's directory:
```
cd terraform-provider-serverspace
```

3. Run the following command to build the provider
```
go build -o terraform-provider-serverspace
```

4. Create directory to make it visible to terraform 
```
mkdir -p ~/.terraform.d/plugins/serverspace.local/local/serverspace/0.2.3/linux_amd64
```

5. Copy built provider in the directory
```
cp terraform-provider-serverspace ~/.terraform.d/plugins/serverspace.local/local/serverspace/0.2.3/linux_amd64
```

6. Use template of configuration file:
```
terraform {
  required_providers {
    serverspace = {
      source = "serverspace.local/local/serverspace"
      version = "0.2.3"
    }
  }
}

variable "s2_token" {
  type = string
  default = "<api key>"
}

provider "serverspace" {
  key = var.s2_token
}
```


## Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd example
terraform init && terraform apply
```

