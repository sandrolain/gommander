# gommander

`gommander` is a terminal-based file manager written in Go. It provides a dual-pane interface for navigating and managing files and directories with ease.

## Features

- Dual-pane file navigation
- File and directory operations (copy, move, delete, create)
- Integration with VSCode for opening files
- Trash support for safe file deletion
- Keyboard shortcuts for efficient usage
- Real-time directory watching for updates

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/sandrolain/gommander.git
   cd gommander
   ```

2. Build the project:

   ```bash
   go build -o gommander
   ```

3. Run the application:

   ```bash
   ./gommander
   ```

## Usage

- Use the arrow keys to navigate files and directories.
- Press `Enter` to open a file or enter a directory.
- Use `Ctrl+C` to copy files, `Ctrl+X` to move files, and `Ctrl+D` to delete files.
- Press `Ctrl+K` to open the selected file in VSCode.
- Refer to the help menu (`Ctrl+H`) for a full list of keyboard shortcuts.

## Disclaimer

**Use at your own risk.** The authors of `gommander` are not responsible for any data loss, damage, or other issues that may arise from using this software. Always ensure you have backups of your important data before performing file operations.
