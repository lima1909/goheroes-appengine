package gcloud

import (
	"context"
	"time"

	"google.golang.org/appengine/log"
)

// LogMessage the log message
type LogMessage struct {
	Time    time.Time
	Level   int
	Message string
}

// ShowLogs shows app-logs from the gcloud
func ShowLogs(c context.Context) []LogMessage {

	query := &log.Query{
		AppLogs: true,
		// Versions: []string{"1"},
	}
	log.Warningf(c, "---> LogMessageQuery: %v", query)

	msg := make([]LogMessage, 0)
	count := -1
	for results := query.Run(c); ; {
		count++
		record, err := results.Next()
		if err == log.Done {
			log.Warningf(c, "Done processing results: %v ", count)
			break
		}
		if err != nil {
			log.Errorf(c, "Failed to retrieve next log: %v", err)
			break
		}

		for _, al := range record.AppLogs {
			msg = append(msg, LogMessage{
				Level:   al.Level,
				Message: al.Message,
				Time:    al.Time,
			})
		}
		log.Infof(c, "Saw record %v", record)
	}

	return msg
}
