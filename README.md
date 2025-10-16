# go-pomodoro

A Pomodoro timer for the terminal, built to explore Go's concurrency model and cross-platform capabilities.

This project uses **goroutines** and **channels** for real-time keyboard input handling while the timer runs, and includes platform-specific implementations for **Windows** and **Linux** desktop notifications using Go's build tags.

## Installation

### Prerequisites
- Go 1.25.1 or higher
- **Linux only**: `libnotify` for desktop notifications (usually pre-installed)

### Build from source

```bash
git clone https://github.com/mdhenriques/go-pomodoro.git
cd go-pomodoro
go build
```

This creates an executable:
- **Windows**: `go-pomodoro.exe`
- **Linux**: `go-pomodoro`

## Usage

### Interactive mode

Run without flags to enter values interactively:

**Windows:**
```bash
go-pomodoro.exe
```

**Linux:**
```bash
./go-pomodoro
```

### Quick start with flags

**Windows:**
```bash
go-pomodoro.exe --study 25 --rest 5 --cycles 4
```

**Linux:**
```bash
./go-pomodoro --study 25 --rest 5 --cycles 4
```

### Available flags

- `--study` : Study duration in minutes
- `--rest` : Rest duration in minutes  
- `--cycles` : Number of pomodoro cycles
- `--test-unit-seconds` : Use seconds instead of minutes (for testing)

### During a session

- Press `p` to pause the current timer
- Press `r` to resume
- Press `q` to quit the session early

## Examples

Standard pomodoro (25 min work, 5 min rest, 4 cycles):

**Windows:**
```bash
go-pomodoro.exe --study 25 --rest 5 --cycles 4
```

**Linux:**
```bash
./go-pomodoro --study 25 --rest 5 --cycles 4
```

Short sessions for testing:
```bash
# Windows
go-pomodoro.exe --study 2 --rest 1 --cycles 3

# Linux
./go-pomodoro --study 2 --rest 1 --cycles 3
```

Custom workflow:
```bash
# Windows
go-pomodoro.exe --study 50 --rest 10 --cycles 3

# Linux
./go-pomodoro --study 50 --rest 10 --cycles 3
```

## Session Summary

At the end of each session, you'll see a summary with:
- Start and end times
- Cycles completed
- Total study time
- Total rest time

## Platform Support

- **Windows**: Toast notifications via PowerShell
- **Linux**: Desktop notifications via `notify-send` (install `libnotify` if needed)

## License

MIT