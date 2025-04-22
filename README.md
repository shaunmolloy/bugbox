# BugBox

---

1. [Background](#background)
1. [Features](#features)
1. [Install](#install)
1. [Usage](#usage)
1. [Configuration](#configuration)
1. [Debugging](#debugging)
1. [License](#license)

---

## Background

BugBox: Git issue inbox from your terminal.

The aim for this project is to track and manage issues across multiple GitHub organizations and repositories.

---

## Features

- **List all issues**: Quickly view GitHub issues across multiple orgs and repos.
- **Polling**: Automatic polling for the latest open issues.
- **Search for Issues**: Search for issues using a built-in search bar.
- **Open Issues in Browser**: Directly open issues in your default browser from the terminal.

### Coming soon

- **Filter Issues by Org**: Supports filtering issues by GitHub organization names.

---

## Install

### Recommended

```bash
go install github.com/yourusername/bugbox@latest
```

---

## Usage

```bash
bugbox
```

---

## Configuration

Run setup to configure your GitHub token and select the organizations you want to track:

```bash
bugbox setup
```

Config files are saved to `~/.config/bugbox/`.

---

## Debugging

If you encounter any issues, you can find more detailed logs by running:

```bash
cat ~/.local/share/bugbox/bugbox.log
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
