# Git Local Stats

A lightweight Go CLI that generates a GitHub-style contribution graph from your local Git repositories.

Instead of relying on GitHub activity, Git Local Stats scans repositories stored on your machine and visualizes your commit history directly in the terminal.

## Features

* Scan directories recursively for Git repositories
* Ignore common dependency folders (`vendor`, `node_modules`)
* Store discovered repositories automatically
* Aggregate commits across multiple repositories
* Filter commits by author email
* Display a contribution heatmap for the last 183 days
* Works completely offline

## Installation

```bash
git clone https://github.com/LuiGi69-orz/miniature-bassoon.git
cd miniature-bassoon

go build -o gitstats .
```

## Usage

### 1. Scan for repositories

```bash
./gitstats -add ~/Projects
```

This recursively searches for Git repositories and stores their paths in:

```text
~/.gogitloccalstats
```

### 2. Generate contribution statistics

```bash
./gitstats -email your@email.com
```

Example:

```bash
./gitstats -email john@example.com
```

The tool will:

1. Load all previously scanned repositories.
2. Read commit history from each repository.
3. Count commits authored by the specified email.
4. Display a terminal contribution graph.

## Command Line Options

| Flag     | Description                                    |
| -------- | ---------------------------------------------- |
| `-add`   | Scan a directory and register Git repositories |
| `-email` | Author email used for commit filtering         |

Example:

```bash
./gitstats -add ~/workspace
./gitstats -email me@example.com
```

## How It Works

### Repository Discovery

The scanner:

* Traverses directories recursively
* Detects repositories by locating `.git`
* Skips:

  * `vendor`
  * `node_modules`

### Commit Collection

The application uses:

```go
gopkg.in/src-d/go-git.v4
```

to read commit history directly from local repositories.

### Contribution Grid

The contribution graph:

* Covers the last 183 days
* Organizes commits by week and weekday
* Uses different color intensities based on commit count
* Highlights the current day

## Example Output

```text
     Jan    Feb    Mar    Apr

Sun  -  -  1  3  0  5
Mon  -  2  0  4  7  1
Tue  -  -  3  2  0  8
Wed  1  4  2  0  3  6
Thu  -  -  -  1  2  3
Fri  2  1  0  4  5  7
Sat  -  3  1  0  2  4
```

## Future Improvements

* Multiple email support
* Repository management commands
* Custom date ranges
* Export to SVG or PNG
* Repository-specific statistics
* Commit streak tracking
* Language breakdowns

## License

MIT License
