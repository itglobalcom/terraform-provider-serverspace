# Terraform Provider Serverspace

## To use provider from terraform registry just use configuration below in your *.tf file
```
terraform {
    required_providers {
        serverspace = {
            source = "itglobalcom/serverspace"
            version = "~> 0.2.2"
        }
    }
}
```

## To use local provider

Run the following command to build the provider
```
go build -o terraform-provider-serverspace
```

Create directory to make it visible to terraform 
```
~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```

Copy built provider in the directory
```
cp terraform-provider-serverspace ~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```

Add provider configuration in your *.tf file
```
terraform {
	required_providers {
		serverspace = {
			source = "{host_name}/{namespace}/{type}"
			version = "{version}"
		}
	}
}
```


## Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd example
terraform init && terraform apply
```

