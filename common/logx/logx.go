package logx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	DEBUG logLevelTyp = iota
	LOG
	ERROR
	CRITICAL
	CUSTOM
	NOLOG
	eol // end of level
)

type logLevelTyp int

func (l logLevelTyp) Prefix(p ...string) string {
	if l == CUSTOM {
		return customPrefix
	}
	if p == nil || p[0] == "" {
		return prefixes[l]
	}
	if l == eol {
		return p[0]
	}
	return fmt.Sprintf("%s-%s", prefixes[l], p[0])
}

func Write(prefix, format string, a ...any) { write(prefix, eol, format, a...) }

func Debugf(format string, a ...any)    { write("", DEBUG, format, a...) }
func Logf(format string, a ...any)      { write("", LOG, format, a...) }
func Errorf(format string, a ...any)    { write("", ERROR, format, a...) }
func Criticalf(format string, a ...any) { write("", CRITICAL, format, a...) }

func Debug(prefix, format string, a ...any)    { write(prefix, DEBUG, format, a...) }
func Log(prefix, format string, a ...any)      { write(prefix, LOG, format, a...) }
func Error(prefix, format string, a ...any)    { write(prefix, ERROR, format, a...) }
func Critical(prefix, format string, a ...any) { write(prefix, CRITICAL, format, a...) }

func Custom(format string, a ...any) { write("", CUSTOM, format, a...) }

var (
	customPrefix             = ""
	prefixes                 = [eol]string{"DEBUG", "LOG", "ERROR", "CRITICAL", "", "-"}
	logWriter    io.Writer   = os.Stdout
	logLevel     logLevelTyp = DEBUG
	timeFormat               = time.RFC3339
	timeLocation             = time.UTC
	mu           sync.Mutex
)

func timestamp() string {
	if timeFormat == "" {
		return ""
	}
	return time.Now().In(timeLocation).Format(timeFormat)
}

func write(prefix string, level logLevelTyp, msgOrFormat string, a ...any) {
	if level < logLevel || logLevel == NOLOG {
		return
	}
	var msg string
	if a == nil {
		msg = msgOrFormat
	} else {
		msg = fmt.Sprintf(msgOrFormat, a...)
	}
	fmt.Fprintf(logWriter, "%s [%s] %s\n", timestamp(), level.Prefix(prefix), msg)
}

func IsNoTimestamp() bool     { return timeFormat == "" }
func GetCustomPrefix() string { return customPrefix }
func GetLogWriter() io.Writer { return logWriter }
func GetLevel() logLevelTyp   { return logLevel }
func GetTimeFormat() string   { return timeFormat }
func GetTimeZone() string     { return timeLocation.String() }

func set(w io.Writer, l *logLevelTyp, cp *string, tf *string, loc *time.Location) error {
	mu.Lock()
	defer mu.Unlock()

	if w != nil {
		logWriter = w
	}
	if l != nil {
		logLevel = *l
	}
	if cp != nil {
		customPrefix = *cp
	}
	if tf != nil {
		timeFormat = *tf
	}
	if loc != nil {
		timeLocation = loc
	}
	return nil
}

func SetLogger(logger io.Writer) error {
	if logger == nil {
		return errors.New("logger must not be nil")
	}
	return set(logger, nil, nil, nil, nil)
}

func SetLevel(level logLevelTyp) error {
	return set(nil, &level, nil, nil, nil)
}

func SetCustomPrefix(prefix string) error {
	return set(nil, nil, &prefix, nil, nil)
}

func SetTimeFormat(format string) error {
	if format == "" {
		return set(nil, nil, nil, &format, nil)
	}
	if _, err := time.Parse(format, "2006-01-02 15:04:05"); err != nil {
		return err
	}
	return set(nil, nil, nil, &format, nil)
}

func SetTimeZone(zone string) error {
	location, err := time.LoadLocation(zone)
	if err != nil {
		return err
	}
	return set(nil, nil, nil, nil, location)
}
