# A very simple tool to stream a log file to Azure Monitor

This tool has the sole purpose of tailing a JSON-formatted log file and push new lines to [Azure Monitor Logs](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/data-platform-logs) using the ingestion API.

## TODO

- [ ] (better) error handling
- [ ] test with long-time operation

## Setting up credentials

You need to [create a Microsoft Entra application](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#create-azure-ad-application) to authenticate against the API. With that done, export the following environment variables:

- `AZURE_TENANT_ID`
- `AZURE_CLIENT_ID`
- `AZURE_CLIENT_SECRET`

## Usage

### Setting up Azure Monitor Logs

Find the logfile, make sure it is formatted as JSON. Then set up the ingestion, again following the
related steps in the already mentioned tutorial to

1. [Create a data collection endpoint](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#create-data-collection-endpoint) and note down the **Logs ingestion** URI because you'll need it in a later step.
2. [Create a new table in the Log Analytics workspace](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#create-new-table-in-log-analytics-workspace)
3. [Parse and filter some sample data](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#parse-and-filter-sample-data) but instead of generating sample data in the first step, take one or two lines from your log file, wrap them into an array (`[{…},{…}]`) and use that as a sample.
4. [Collect information from the DCR](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#collect-information-from-the-dcr), you need the `immutableId` value later.
5. [Assign permissions to the DCR](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#assign-permissions-to-the-dcr)

### Pushing logs to Azure Monitor Logs

With the credential environment variables in place and the Azure Monitor setup done as described above,
using the following should push entries:

    logingestor --endpoint <endpointUri> --rule-id <ruleId> --stream-name <streamName> <path/to/logfile>

- The `<endpointUri>` is what you noted down in step 1 when setting up Azure Monitor Logs
- The `<ruleId>` is the `immutableID` you noted down in step 4
- The `<streamName>` is based on the table name from step 2 and follows the pattern `Custom-<table name>_CL`

Those flags can be omitted, if the following environment variables are used instead:

- `AZURE_MONITOR_ENDPOINT`
- `AZURE_MONITOR_RULE_ID`
- `AZURE_MONITOR_STREAM_NAME`

The tool also has a built-in help, at the time of writing it said:

    NAME:
       logingestor - Follow a log file and send lines to Azure Monitor Logs
    
    USAGE:
       logingestor [global options] command [command options] [arguments...]
    
    VERSION:
       9377315-dirty
    
    DESCRIPTION:
       A tool to follow a log file and send lines to Azure Monitor Logs
    
    COMMANDS:
       help, h  Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --endpoint ENDPOINT        The ENDPOINT to send data to [$AZURE_MONITOR_ENDPOINT]
       --rule-id RULE-ID          The RULE-ID to send data to [$AZURE_MONITOR_RULE_ID]
       --stream-name STREAM-NAME  The STREAM-NAME to send data to [$AZURE_MONITOR_STREAM_NAME]
       --help, -h                 show help
       --version, -v              print the version

## Some helpful bits & pieces

These are "stored" here for later reference.

### Links

- [Azure Monitor Ingestion client module for Go](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest)
- [Logs Ingestion API in Azure MonitorLogs Ingestion API in Azure Monitor](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/logs-ingestion-api-overview)
- [Sample code to send data to Azure Monitor using Logs ingestion API](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-code?tabs=go)

### Compile for Linux on i386

GOOS=linux GOARCH=386 go build logingestor.go
