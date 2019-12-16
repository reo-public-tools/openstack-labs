package osutils

import (
    "log"
    "log/syslog"
)

// Set up package level sysLog var for logging use
var sysLog *syslog.Writer

func init() {

    // Set up a default local system logger
    var err error
    sysLog, err = syslog.New(syslog.LOG_EMERG | syslog.LOG_USER, "osutils(package)")
    if err != nil {
        log.Fatal(err)
    }

}

func OverrideLogWriter(logWriter *syslog.Writer) {

    sysLog = logWriter

}
