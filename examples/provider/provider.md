# Provider Configuration

## Example Usage

```hcl
terraform {
  required_providers {
    statuspage = {
      source = "winebarrel/statuspage"
    }
  }
}

provider "statuspage" {
  api_key = "your-api-key"
}
```

## Authentication

The API key can be provided in one of the following ways:

1. Set the `api_key` attribute in the provider block.
2. Set the `STATUSPAGE_API_KEY` environment variable.

```shell
export STATUSPAGE_API_KEY="your-api-key"
```

## Argument Reference

- `api_key` (Optional, Sensitive) - The Statuspage.io API key. Can also be set via `STATUSPAGE_API_KEY` environment variable.
