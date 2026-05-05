# cronwatch

Lightweight daemon that monitors cron job execution and sends alerts on failures or missed runs.

## Installation

```bash
go install github.com/cronwatch/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/cronwatch/cronwatch.git && cd cronwatch && make install
```

## Usage

Define your monitored jobs in `cronwatch.yaml`:

```yaml
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    timeout: 30m
    alert:
      email: ops@example.com

  - name: hourly-sync
    schedule: "0 * * * *"
    timeout: 5m
    alert:
      slack: "#alerts"
```

Start the daemon:

```bash
cronwatch --config cronwatch.yaml
```

Wrap an existing cron job to report its status:

```bash
# In your crontab
0 2 * * * cronwatch exec --job daily-backup -- /usr/local/bin/backup.sh
```

cronwatch will send an alert if the job exits with a non-zero status, exceeds its timeout, or fails to run within the expected schedule window.

## Configuration

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression defining the expected run interval |
| `timeout` | Maximum allowed execution duration |
| `alert` | Notification target (`email`, `slack`, or `webhook`) |

## License

MIT © cronwatch contributors