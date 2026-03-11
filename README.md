# Akmatori Terraform Provider

Terraform provider for managing resources on the [Akmatori](https://akmatori.com) AIOps platform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider plugin)

## Usage

```hcl
terraform {
  required_providers {
    akmatori = {
      source = "akmatori/akmatori"
    }
  }
}

provider "akmatori" {
  host  = "https://akmatori.example.com"
  token = var.akmatori_token
}
```

## Authentication

The provider supports two authentication methods:

### Token Authentication (recommended)

```hcl
provider "akmatori" {
  host  = "https://akmatori.example.com"
  token = var.akmatori_token
}
```

### Username/Password Authentication

```hcl
provider "akmatori" {
  host     = "https://akmatori.example.com"
  username = "admin"
  password = var.akmatori_password
}
```

### Environment Variables

All provider arguments can be set via environment variables:

| Variable | Description |
|----------|-------------|
| `AKMATORI_HOST` | The URL of the Akmatori instance |
| `AKMATORI_TOKEN` | JWT token for authentication |
| `AKMATORI_USERNAME` | Username for authentication |
| `AKMATORI_PASSWORD` | Password for authentication |

## Resources

| Resource | Description |
|----------|-------------|
| `akmatori_skill` | Manage AI skills |
| `akmatori_skill_tools` | Assign tool instances to skills |
| `akmatori_skill_script` | Manage skill scripts |
| `akmatori_tool_instance` | Manage tool instances |
| `akmatori_ssh_key` | Manage SSH keys for tool instances |
| `akmatori_alert_source` | Manage alert sources |
| `akmatori_context_file` | Manage context files |
| `akmatori_settings_slack` | Configure Slack integration |
| `akmatori_settings_llm` | Configure LLM provider |
| `akmatori_settings_proxy` | Configure proxy settings |
| `akmatori_settings_aggregation` | Configure alert aggregation |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `akmatori_skill` | Read skill data |
| `akmatori_tool_type` | Read tool type data |
| `akmatori_tool_types` | List all tool types |
| `akmatori_alert_source_types` | List all alert source types |

## Building from Source

```shell
git clone https://github.com/akmatori/terraform-provider-akmatori.git
cd terraform-provider-akmatori
make build
```

## Developing

```shell
# Run unit tests
make test

# Run acceptance tests (requires a running Akmatori instance)
make testacc

# Generate documentation
make docs

# Run linter
make lint
```

## License

Apache 2.0 - See [LICENSE](LICENSE) for more information.
