# coolify-tui

A read-only terminal dashboard for [Coolify](https://coolify.io), inspired by tools such as Lazygit.

Browse teams, projects, environments, resources, environment variables, deployments, and deployment logs without leaving your terminal.

> This is an independent project and is not officially affiliated with or maintained by Coolify.

## Features

- Lazygit-style multi-panel interface
- Browse Coolify teams and projects
- Navigate project environments and resources
- View application status and details
- View application environment variables
- Environment-variable values masked by default
- View deployment history
- View deployment details and logs
- Scroll through deployment logs
- Filter projects, environments, resources, variables, and deployments
- Keyboard help popup
- Read-only Coolify API access

## Requirements

- A running Coolify instance with API access enabled
- A Coolify API token
- A terminal with color support
- Go 1.24 or newer when installing from source

A terminal size of at least `80x24` is recommended.

## Installation

### Using Go

Install the latest version directly from GitHub:

```sh
go install github.com/micaelmcarvalho/coolify-tui/cmd/coolify-tui@latest
```

Make sure the Go binary directory is available in your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify the installation:

```sh
coolify-tui
```

If configuration has not been created yet, the application will explain how to configure it.

### Using a GitHub release

Download the archive for your operating system and architecture from the [GitHub Releases page](https://github.com/micaelmcarvalho/coolify-tui/releases).

Extract the archive and move the binary somewhere in your `PATH`.

On macOS or Linux:

```sh
chmod +x coolify-tui
sudo mv coolify-tui /usr/local/bin/
```

Verify that it is available:

```sh
coolify-tui
```

### From source

Clone and build the project:

```sh
git clone https://github.com/micaelmcarvalho/coolify-tui.git
cd coolify-tui

go build -o coolify-tui ./cmd/coolify-tui
./coolify-tui
```

## Configuration

The easiest way to configure the application is:

```sh
coolify-tui configure
```

You will be prompted for:

- Your Coolify URL
- Your Coolify API token

Example:

```text
Coolify URL: https://coolify.example.com
Coolify API token:
Configuration saved to ...
```

The API token remains hidden while you type or paste it.

After configuration, start the dashboard:

```sh
coolify-tui
```

### Configuration file location

The configuration is stored in your operating system's user configuration directory.

Typical locations are:

| Operating system | Location |
| --- | --- |
| macOS | `~/Library/Application Support/coolify-tui/config.json` |
| Linux | `~/.config/coolify-tui/config.json` |
| Windows | `%AppData%\coolify-tui\config.json` |

On macOS and Linux, the file is created with permissions restricted to the current user.

The file contains your Coolify URL and API token. Do not share or commit it to Git.

### Environment variables

You can use environment variables instead of the saved configuration:

```sh
export COOLIFY_URL="https://coolify.example.com"
export COOLIFY_TOKEN="your-api-token"

coolify-tui
```

Environment variables take precedence over the saved configuration.

### Local `.env` file

A local `.env` file is also supported, which is useful during development:

```dotenv
COOLIFY_URL=https://coolify.example.com
COOLIFY_TOKEN=your-api-token
```

Run the application from the directory containing the file:

```sh
coolify-tui
```

Never commit `.env` files or API tokens to Git.

Make sure your `.gitignore` contains:

```gitignore
.env
.env.*
!.env.example
```

### Configuration priority

When multiple configuration methods are present, the application uses the following priority:

1. Exported environment variables
2. Values from a local `.env` file
3. Values saved by `coolify-tui configure`

This allows release users to use the configuration command while developers and automation systems can continue using environment variables.

## Coolify API token

Create an API token from your Coolify dashboard.

For basic projects and resources, the token needs read access.

To view environment variables and deployment logs, it may also need:

```text
read:sensitive
```

Token access can depend on the Coolify team and permissions associated with the token.

See the [Coolify API authorization documentation](https://coolify.io/docs/api-reference/authorization) for more information.

The application currently performs only read-only API requests.

## Usage

Start the dashboard:

```sh
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

Filtering is available for:

- Projects
- Environments
- Resources
- Environment variables
- Deployments

### Environment variables

| Key | Action |
| --- | --- |
| `v` | Reveal or hide environment-variable values |

Environment-variable values are hidden by default.

Be careful when revealing values. Secrets may remain visible in terminal scrollback, screenshots, screen sharing, or screen recordings.

### Deployments

| Key | Action |
| --- | --- |
| `n` | Load the next deployments page |
| `p` | Load the previous deployments page |
| `Enter` | Open the selected deployment |
| `j` / `↓` | Scroll logs down |
| `k` / `↑` | Scroll logs up |
| `g` / `Home` | Jump to the beginning of the logs |
| `G` / `End` | Jump to the end of the logs |
| `r` | Refresh deployment details |
| `Esc` | Return to the dashboard |

### Help

| Key | Action |
| --- | --- |
| `?` | Open or close keyboard help |
| `Esc` | Close keyboard help |

## Development

Clone the repository:

```sh
git clone https://github.com/micaelmcarvalho/coolify-tui.git
cd coolify-tui
```

Download dependencies:

```sh
go mod download
```

Create a local configuration:

```sh
go run ./cmd/coolify-tui configure
```

Run the application:

```sh
go run ./cmd/coolify-tui
```

Run the project checks:

```sh
go fmt ./...
go vet ./...
go test ./...
```

Build a local binary:

```sh
go build -o coolify-tui ./cmd/coolify-tui
```

## Project structure

```text
cmd/coolify-tui/       Application entry point
internal/config/       Configuration loading and storage
internal/coolify/      Coolify API client and response types
internal/ui/           Bubble Tea terminal interface
```

## Security

- Use the least-privileged Coolify token possible.
- Treat the saved configuration file as sensitive.
- Never commit API tokens or `.env` files.
- Avoid passing tokens through command-line arguments because they can appear in shell history and process lists.
- Avoid revealing environment-variable values during screen sharing.
- Revoke and replace any token that may have been exposed.
- Keep the application and your Coolify installation updated.

## Roadmap

- Multiple Coolify profiles
- Multiple team-specific tokens
- Application runtime logs
- Copy UUIDs and other safe values to the clipboard
- Homebrew installation
- Optional operating-system keychain integration
- Confirmed deploy and restart actions

Any future action that changes Coolify resources should require explicit confirmation.

## Contributing

Issues and pull requests are welcome.

Before submitting a pull request, run:

```sh
go fmt ./...
go vet ./...
go test ./...
```

When reporting a bug, include:

- Your operating system
- Your terminal application
- Your Coolify version
- The steps needed to reproduce the problem

Do not include API tokens, environment-variable values, or other secrets in bug reports.

## License

Released under the [MIT License](LICENSE).
