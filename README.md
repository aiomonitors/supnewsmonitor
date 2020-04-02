# supnewsmonitor

A command line tool for [Supreme News](https://supremenewyork.com/news) to send the latest news or images to a webhook, or monitor to a webhook.

## Usage

1. Clone the repository using `git clone https://github.com/aiomonitors/supnewsmonitor` or download from here
2. Navigated to the repository using a file browser or `cd supnewsmonitor`
3. Open `config.json` and edit the file.

Your config.json should look like this after editing:
```json
{
    "webhook" : "https://canary.discordapp.com/api/webhooks/WEBHOOK LINK",
    "hexcode" : "#F1B379",
    "icon" : "ICON URL HERE OR LEAVE EMPTY",
    "groupName" : "GROUP NAME HERE"
}
```

After editing your config.json to look like this, you can proceed to one of the following ways:

### Without Go Installed
If you do not have Go installed on your system, there are executable files that you can run instead. 

On Mac:
```
chmod +x supnewsmonitor-darwin
./supnewsmonitor-darwin
```

On Windows: `supnewsmonitor.exe` or click the executable
On Windows AMD64: `supnewsmonitor-windows-amd64.exe` or click the executable

### With Go installed
Run: 
```
go run supnewsmonitor.go
```
This should automatically download all the dependencies required for the monitor / tool.
Otherwise, run `go mod download`

Created by Shihab Chowdhury [Site](https://www.shihab.dev) | [Twitter](https://twitter.com/aiomonitors) | [Email](mailto:navr@discoders.us)