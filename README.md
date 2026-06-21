# vim-go

A terminal-based, Vim-like text editor written in Go utilizing the `tcell/v3` library. The editor features a responsive, scrolling viewport, custom editing modes, dynamic text selection, and is backed by a highly efficient Gap Buffer data structure.

## Features

- **Multiple Input Modes**:
  - **Normal Mode**: Standard Vim navigation (`h`, `j`, `k`, `l` (customized mappings), `0`, `$`, `w`, `gg`, `G`), selection initialization (`v`), and deletion commands (`x`, `dd`, `dw`).
  - **Insert Mode**: Classic text editor insert/append and deletion (Backspace).
  - **Visual Mode**: Dynamic character/line selection with visible highlight styling (`tcell.ColorCadetBlue`), yanking (`y`), and deleting (`d`/`x`).
  - **Command Mode**: Run command-line commands (`:w` to save, `:q` to exit, `:wq` to save & exit).
- **Buffer Logic**:
  - Powered by an efficient **Gap Buffer** (`GapBuffer`) data structure for fast insertions and deletions in large files.
  - Automatic line boundary cache recomputation.
- **Rendering & TUI**:
  - Adaptive viewport scrolling keeping cursor in view.
  - Gutter showing line numbers that scales dynamically based on file length.
  - Enhanced status bar showcasing current mode, active file name, and cursor coordinate indices.

---

## Directory Structure

```
├── cmd/
│   └── main.go                  # Main entry point initializing Screen and Editor
├── internal/
│   ├── buffer/
│   │   ├── buffer.go            # Coordinate mapping wrapper and line metrics
│   │   ├── cursor.go            # Cursor positioning and movement boundaries
│   │   └── gap_buffer.go        # Low-level 1D Gap Buffer memory management
│   ├── editor/
│   │   ├── editor.go            # Core controller coordinating buffer, viewport, and rendering
│   │   └── events.go            # Event loop listening to resizing and keystrokes
│   ├── mode/
│   │   ├── mode.go              # Shared input Mode and Editor interfaces
│   │   ├── normal.go            # Key mapping parsing in Normal Mode
│   │   ├── insert.go            # Character insertions in Insert Mode
│   │   ├── command.go           # Command collection in Command Mode
│   │   └── visual.go            # Visual selection key combinations
│   └── ui/
│       ├── screen.go            # Tcell screen abstraction and draw primitives
│       └── viewport.go          # Viewport dimensions and scroll offsets
```

---

## Getting Started

### Prerequisites

- Go 1.18 or higher.
- `tcell/v3` library dependencies.

### Installation

Clone the repository:
```bash
git clone https://github.com/sanket9162/vim-go.git
cd vim-go
```

Build the editor binary:
```bash
go build -o vim-go ./cmd/main.go
```

### Usage

Launch the editor with an empty buffer:
```bash
./vim-go
```

Open or edit an existing file:
```bash
./vim-go path/to/file.txt
```

---

## Key Bindings

### Normal Mode (Default)
- `i` / `a` - Enter Insert Mode (insert / append).
- `v` - Enter Visual Mode.
- `:` - Enter Command Mode.
- `h` / `j` / `k` / `l` - Cursor navigation (mapped to custom layout).
- `0` / `$` - Jump to start / end of current line.
- `w` - Jump to start of next word.
- `g g` / `G` - Jump to start / end of file.
- `x` - Delete character under cursor.
- `d d` - Delete current line.
- `d w` - Delete word under cursor.

### Visual Mode
- `h` / `j` / `k` / `l` / `0` / `$` / `w` - Adjust selection boundaries.
- `y` - Yank (copy) selection.
- `d` / `x` - Delete (cut) selection and return to Normal Mode.
- `Esc` - Cancel selection and return to Normal Mode.

### Command Mode
- `w` + `Enter` - Save the file.
- `q` + `Enter` - Quit the editor.
- `wq` + `Enter` - Save and quit.
- `Esc` - Discard input and return to Normal Mode.
