package utils

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/ravvio/awst/ui/tlog"
)

func LogFromCloudwatchEvent(groupName *string, ev *types.FilteredLogEvent) tlog.Log {
	return tlog.Log{
		GroupName: groupName,
		Timestamp: ev.Timestamp,
		Message:   ev.Message,
	}
}
