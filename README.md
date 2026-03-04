# terraform-provider-statuspage

[![CI](https://github.com/winebarrel/terraform-provider-statuspage/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-statuspage/actions/workflows/ci.yml)
[![terraform registry](https://img.shields.io/badge/terraform-registry-blueviolet)](https://registry.terraform.io/providers/winebarrel/statuspage/latest)

Terraform provider for [Statuspage.io](https://www.atlassian.com/software/statuspage).

## Usage

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

resource "statuspage_component" "api" {
  page_id     = "abc123def456"
  name        = "API"
  description = "API service"
  status      = "operational"
  showcase    = true
}
```

## Authentication

The API key can be provided in one of the following ways:

1. Set the `api_key` attribute in the provider block.
2. Set the `STATUSPAGE_API_KEY` environment variable.

```shell
export STATUSPAGE_API_KEY="your-api-key"
```

## Resources

- [statuspage_page](docs/resources/page.md)
- [statuspage_component](docs/resources/component.md)
- [statuspage_component_group](docs/resources/component_group.md)
- [statuspage_incident](docs/resources/incident.md)
- [statuspage_incident_template](docs/resources/incident_template.md)
- [statuspage_incident_postmortem](docs/resources/incident_postmortem.md)
- [statuspage_metric](docs/resources/metric.md)
- [statuspage_metrics_provider](docs/resources/metrics_provider.md)
- [statuspage_subscriber](docs/resources/subscriber.md)
- [statuspage_user](docs/resources/user.md)
- [statuspage_user_permissions](docs/resources/user_permissions.md)
- [statuspage_page_access_user](docs/resources/page_access_user.md)
- [statuspage_page_access_group](docs/resources/page_access_group.md)
- [statuspage_status_embed_config](docs/resources/status_embed_config.md)

> **Note:** `statuspage_page` and `statuspage_status_embed_config` cannot be created or deleted via the API. Use `terraform import` to manage existing resources.

## Data Sources

- [statuspage_page](docs/data-sources/page.md)
- [statuspage_component](docs/data-sources/component.md)
- [statuspage_component_group](docs/data-sources/component_group.md)
- [statuspage_incident](docs/data-sources/incident.md)
- [statuspage_incident_template](docs/data-sources/incident_template.md)
- [statuspage_metric](docs/data-sources/metric.md)
- [statuspage_metrics_provider](docs/data-sources/metrics_provider.md)
- [statuspage_subscriber](docs/data-sources/subscriber.md)
- [statuspage_page_access_user](docs/data-sources/page_access_user.md)
- [statuspage_page_access_group](docs/data-sources/page_access_group.md)

## Development

```shell
# Build
make build

# Test
make test

# Lint
make lint

# Generate documentation
make docs
```
