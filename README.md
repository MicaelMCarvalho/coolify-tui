# coolify-tui

A read-only terminal dashboard for [Coolify](https://coolify.io), inspired by tools such as Lazygit.

Browse teams, projects, environments, resources, environment variables, deployments, and deployment logs without leaving your terminal.

## Features

- Lazygit-style multi-panel interface
- Browse Coolify teams and projects
- Navigate project environments and resources
- View application status and details
- View application environment variables
- Environment-variable values masked by default
- View deployment history
- View deployment details and logs
- Filter projects, environments, resources, variables, and deployments
- Keyboard help popup
- Read-only API access

## Requirements

- A running Coolify instance with API access enabled
- A Coolify API token
- A terminal with color support
- Go 1.24 or newer when installing from source

The recommended terminal size is at least `80x24`.

## Installation

### Using Go

Install directly from GitHub:

```bash
go install github.com/micaelmcarvalho/coolify-tui/cmd/coolify-tui@latest
```

Make sure the Go binary directory is in your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then run:

```bash
coolify-tui
```

### From source

```bash
git clone https://github.com/micaelmcarvalho/coolify-tui.git
cd coolify-tui

go build -o coolify-tui ./cmd/coolify-tui
./coolify-tui
```

## Configuration

`coolify-tui` requires two environment variables:

```bash
export COOLIFY_URL="https://coolify.example.com"
export COOLIFY_TOKEN="your-api-token"
```

Then start the dashboard:

```bash
coolify-tui
```

You can also create a local `.env` file:

```dotenv
COOLIFY_URL=https://coolify.example.com
COOLIFY_TOKEN=your-api-token
```

Run the application from the directory containing that file:

```bash
coolify-tui
```

Do not commit `.env` files or API tokens to Git.

### Coolify token permissions

For basic projects and resources, the token needs:

```text
read
```

To view environment variables and deployment logs, it also needs:

```text
read:sensitive
```

The application currently performs only read-only API requests.

Coolify tokens are scoped to the team that was active when the token was created. Accessing multiple teams will eventually require a separate token for each team.

See the [Coolify API authorization documentation](https://coolify.io/docs/api-reference/authorization) for more information.

## Usage

```bash
coolify-tui
```

The dashboard contains the following panels:

1. Teams
2. Projects
3. Environments
4. Resources
5. Resource Details
6. Environment Variables
7. Deployments

Selecting a project loads its environments. Selecting an environment loads its resources. Selecting an application loads its environment variables and deployments.

## Keyboard shortcuts

### Navigation

| Key | Action |
| --- | --- |
| `Tab` | Focus the next panel |
| `Shift+Tab` | Focus the previous panel |
| `1`–`7` | Focus a specific panel |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` / `Home` | Select the first item |
| `G` / `End` | Select the last item |
| `Enter` | Open an item or focus the next panel |
| `Esc` | Go back or clear an active filter |
| `r` | Refresh the active panel |
| `q` / `Ctrl+C` | Quit |

### Filtering

| Key | Action |
| --- | --- |
| `/` | Start or edit a panel filter |
| `Enter` | Accept the filter |
| `Esc` | Cancel or clear the filter |
| `Ctrl+U` | Clear the filter input |

Filtering is available for projects, environments, resources, environment variables, and deployments.

### Environment variables

| Key | Action |
| --- | --- |
| `v` | Reveal or hide environment-variable values |

Values are hidden by default.

Be careful when revealing values: secrets may remain visible in terminal scrollback, screenshots, or screen recordings.

### Deployments

| Key | Action |
| --- | --- |
| `n` | Next deployments page |
| `p` | Previous deployments page |
| `Enter` | Open deployment details and logs |
| `j` / `k` | Scroll deployment logs |
| `g` / `G` | Jump to the top or bottom of logs |
| `Esc` | Return to the dashboard |

### Help

| Key | Action |
| --- | --- |
| `?` | Open or close keyboard help |
| `Esc` | Close keyboard help |

## Development

Clone the repository:

```bash
git clone https://github.com/micaelmcarvalho/coolify-tui.git
cd coolify-tui
```

Install dependencies and run:

```bash
go mod download
go run ./cmd/coolify-tui
```

Run the checks:

```bash
go fmt ./...
go vet ./...
go test ./...
```

Build a local binary:

```bash
go build -o coolify-tui ./cmd/coolify-tui
```

## Project structure

```text
cmd/coolify-tui/       Application entry point
internal/config/       Environment configuration
internal/coolify/      Coolify API client and response types
internal/ui/           Bubble Tea interface
```

## Roadmap

- Multiple Coolify team profiles
- Application runtime logs
- Copy values and UUIDs to the clipboard
- Config file support
- GitHub binary releases
- Homebrew installation
- Confirmed deploy and restart actions

## Security

- Use the least-privileged Coolify token possible.
- Store tokens in environment variables or an ignored `.env` file.
- Never commit API tokens.
- Avoid revealing environment-variable values during screen sharing.
- Revoke and replace any token that may have been exposed.

## Contributing

Issues and pull requests are welcome.

Before submitting a pull request, run:

```bash
go fmt ./...
go vet ./...
go test ./...
```

## License

Released under the [MIT License](LICENSE).

## Disclaimer

This is an independent community project and is not officially affiliated with or maintained by Coolify.
