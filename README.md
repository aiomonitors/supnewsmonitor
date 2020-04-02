# supnewsmonitor

A command line tool for [Supreme News](https://supremenewyork.com/news) to send the latest news or images to a webhook, or monitor to a webhook.

### Requirements:
    GoLang (Note: If you don't have Go, scroll down for instructions)

## Usage

### With Go installed

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

Run: 
```
go run supnewsmonitor.go
```
This should automatically download all the dependencies required for the monitor / tool.
Otherwise, run `go mod download`


### Without Go Installed
If you do not have Go installed on your system:

1. Create a filed named `config.json` in your current working directory
2. Edit the config.json so it looks like this:
```json
{
    "webhook" : "https://canary.discordapp.com/api/webhooks/WEBHOOK LINK",
    "hexcode" : "#F1B379",
    "icon" : "ICON URL HERE OR LEAVE EMPTY",
    "groupName" : "GROUP NAME HERE"
}
```

2. Download your operating system's release from [Releases](https://github.com/aiomonitors/supnewsmonitor/releases)

On Mac:
```
chmod +x supnewsmonitor-darwin
./supnewsmonitor-darwin
```

On Windows: `supnewsmonitor.exe` or click the executable
On Windows AMD64: `supnewsmonitor-windows-amd64.exe` or click the executable

#### Finished Product:
![supnews](https://www.shihab.dev/supnews.png)

Created by Shihab Chowdhury [Site](https://www.shihab.dev) | [Twitter](https://twitter.com/aiomonitors) | [Email](mailto:navr@discoders.us)


Questions:

**Q:** I ran the program, but it says "Error in program Error parsing news page, links not found", what do I do?
**A:** The reason for this could be due to you being banned on Supreme, or an error with your internet connection. If you are not experiencing either of those issues, 
please contact me.

**Q:** My executable is not running, wnat should I do?
**A:** I recommend navigating to the directory using Command Prompt (Windows) or Terminal (Mac) and runnign the command listed above. The program should tell you if there is an issue, which would most likely be caused by improper formatting of the config.json file.

*Any other questions, feel free to DM me on Twitter or my Discord: navr#0001*