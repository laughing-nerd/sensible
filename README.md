# Sensible

**Sensible** is a modern, declarative task runner for developers. It helps automate your project workflows with a clean and readable config format, without the overhead of traditional tools like Make, Ansible, or Bash scripts.

---

## Features

- 🔧 Declarative configuration using [HCL](https://github.com/hashicorp/hcl)
- 🖥️ Run tasks locally or remotely via SSH
- 🌍 Multi-environment support (base, staging, production, etc.)
- 🔐 Built-in encrypted secrets management
- 📊 System facts (CPU, memory, disk, OS info) available in actions
- 📦 Built-in components: `shell`, `installer`, `cron`
- 🔑 SSH authentication via password or private key

---

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/laughing-nerd/sensible/main/install.sh | bash
```

You can also specify a version:

```bash
VERSION=v1.0.0 curl -sSL https://raw.githubusercontent.com/laughing-nerd/sensible/main/install.sh | bash
```

---

## Quick Start

### 1. Initialize a new project

```bash
sensible init
```

This creates the following structure:

```
.sensible/
└── base/
    ├── actions/
    │   └── sample-action.hcl
    ├── resources/
    │   ├── hosts.hcl
    │   └── values.hcl
    └── secrets/
        └── .secrets.enc
```

### 2. Run an action

```bash
sensible run -f sample-action
```

---

## Commands

### `sensible init`

Initialize sensible in your project directory.

```bash
sensible init              # Creates .sensible/base/ structure
sensible init -e staging   # Creates .sensible/staging/ structure
```

### `sensible run`

Run an action file.

```bash
sensible run -f <action-file>        # Run action from base environment
sensible run -f <action-file> -e staging  # Run action from staging environment
```

### `sensible secret`

Manage encrypted secrets.

```bash
# Set a secret
sensible secret set -k API_KEY -v "my-secret-value"
sensible secret set -k DB_PASSWORD -v "password123" -e staging

# Get a secret
sensible secret get -k API_KEY
sensible secret get -a              # Get all secrets

# Remove a secret
sensible secret remove -k API_KEY
```

---

## Configuration Files

### `values.hcl`

Define reusable variables for your actions.

```hcl
values {
  server_address = "192.168.1.100:22"
  app_name       = "my-app"
  deploy_path    = "/var/www/my-app"
}
```

Use in actions with `${values.variable_name}`:

```hcl
shell "deploy" {
  command = "rsync -avz ./dist/ ${values.deploy_path}"
}
```

### `hosts.hcl`

Define remote hosts organized into groups for SSH execution.

```hcl
group "Web Servers" {
  host "web1" {
    address    = "192.168.1.10:22"
    username   = "deploy"
    password   = "secret"           # or use private_key
    timeout    = 30                 # optional, in seconds
  }

  host "web2" {
    address     = "${values.server_address}"
    username    = "deploy"
    private_key = "/path/to/key"    # alternative to password
  }
}

group "DB Servers" {
  host "db1" {
    address  = "192.168.1.20:22"
    username = "admin"
    password = "${secrets.db_password}"
  }
}
```

### Secrets

Secrets are stored encrypted in `.secrets.enc`. Manage them with the CLI:

```bash
sensible secret set -k DB_PASSWORD -v "supersecret"
```

Use in actions with `${secrets.secret_name}`:

```hcl
shell "connect-db" {
  command = "mysql -u admin -p${secrets.db_password} mydb"
}
```

---

## Actions

Actions are defined in `.hcl` files inside the `actions/` directory. Each action can contain multiple components that execute sequentially.

### Basic Structure

```hcl
action "My Action" {
  # Optional: specify host groups for remote execution
  # Omit this line to run locally
  groups = ["Web Servers", "DB Servers"]

  shell "step 1" {
    command = "echo 'Hello World'"
  }

  installer "step 2" {
    packages = ["git", "curl"]
  }
}
```

---

## Components

### `shell`

Execute shell commands locally or remotely.

```hcl
shell "build app" {
  command = "npm run build"
}

shell "deploy" {
  command = <<-EOT
    cd /var/www/app
    git pull origin main
    npm install
    pm2 restart all
  EOT
}
```

### `installer`

Install packages using the system's package manager. Automatically detects the available package manager (apt, dnf, pacman, brew, etc.).

```hcl
installer "install dependencies" {
  packages  = ["nginx", "git", "curl"]
  preferred = "apt"  # optional: prefer a specific package manager
}
```

**Supported package managers:** `apt`, `apt-get`, `dnf`, `pacman`, `yum`, `apk`, `brew`, `port`, `nix-env`, `pkg`, `zypper`

### `cron`

Add or remove cron jobs.

```hcl
# Add a cron job
cron "backup job" {
  type       = "add"
  expression = "0 2 * * *"           # Every day at 2 AM
  job        = "/scripts/backup.sh"
  user       = "root"                # optional
}

# Remove a cron job
cron "remove old job" {
  type = "remove"
  job  = "/scripts/old-backup.sh"
}
```

---

## Variables

Sensible supports three types of variables that can be used in your actions:

### Values (`${values.*}`)

User-defined variables from `values.hcl`:

```hcl
shell "example" {
  command = "echo ${values.app_name}"
}
```

### Secrets (`${secrets.*}`)

Encrypted secrets managed via CLI:

```hcl
shell "example" {
  command = "echo ${secrets.api_key}"
}
```

### Facts (`${facts.*}`)

System information automatically gathered at runtime:

| Fact | Description |
|------|-------------|
| `facts.os` | Operating system (linux, darwin, etc.) |
| `facts.os_version` | OS version |
| `facts.kernel_version` | Kernel version |
| `facts.architecture` | System architecture |
| `facts.hostname` | Machine hostname |
| `facts.uptime` | System uptime in seconds |
| `facts.boot_time` | Boot timestamp |
| `facts.cpu_model` | CPU model name |
| `facts.cpu_cores` | Number of CPU cores |
| `facts.cpu_logical_cores` | Number of logical CPU cores |
| `facts.memory_total_mb` | Total memory in MB |
| `facts.memory_available_mb` | Available memory in MB |
| `facts.disk_total` | Total disk space in MB |
| `facts.disk_used` | Used disk space in MB |
| `facts.disk_free` | Free disk space in MB |

**Example:**

```hcl
shell "system info" {
  command = "echo 'Running on ${facts.os} with ${facts.memory_total_mb}MB RAM'"
}
```

---

## Multi-Environment Support

Sensible supports multiple environments. Each environment has its own configuration:

```bash
# Initialize different environments
sensible init -e base       # Default
sensible init -e staging
sensible init -e production
```

This creates:

```
.sensible/
├── base/
│   ├── actions/
│   ├── resources/
│   └── secrets/
├── staging/
│   ├── actions/
│   ├── resources/
│   └── secrets/
└── production/
    ├── actions/
    ├── resources/
    └── secrets/
```

Run actions in specific environments:

```bash
sensible run -f deploy -e staging
sensible run -f deploy -e production
```

**Note:** Values from `base` environment are loaded first, then overridden by environment-specific values.

---

## Example: Full Deployment Action

```hcl
action "Deploy Application" {
  groups = ["Web Servers"]

  shell "pull latest code" {
    command = "cd ${values.deploy_path} && git pull origin main"
  }

  shell "install dependencies" {
    command = "cd ${values.deploy_path} && npm ci --production"
  }

  shell "run migrations" {
    command = "cd ${values.deploy_path} && npm run migrate"
  }

  shell "restart services" {
    command = "sudo systemctl restart ${values.app_name}"
  }

  shell "health check" {
    command = "curl -f http://localhost:${values.app_port}/health || exit 1"
  }
}
```