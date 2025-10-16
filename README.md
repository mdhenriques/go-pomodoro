# go-pomodoro

A Pomodoro timer for the terminal, built to explore Go's concurrency model and cross-platform capabilities.
This project uses goroutines and channels for real-time keyboard input handling while the timer runs, and includes platform-specific implementations for Windows and Linux desktop notifications using Go's build tags.

## Installation

### Prerequisites
- Go 1.25.1 or higher

### Build from source

```bash
git clone https://github.com/mdhenriques/go-pomodoro.git
cd go-pomodoro
go build -o pomodoro
```

## Usage

### Interactive mode

Run without flags to enter values interactively:

```bash
./pomodoro
```

### Quick start with flags

```bash
./pomodoro --study 25 --rest 5 --cycles 4
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
```bash
./pomodoro --study 25 --rest 5 --cycles 4
```

Short sessions for testing:
```bash
./pomodoro --study 2 --rest 1 --cycles 3
```

Custom workflow:
```bash
./pomodoro --study 50 --rest 10 --cycles 3
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