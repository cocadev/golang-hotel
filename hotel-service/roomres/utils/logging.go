package utils

import (
	"fmt"
	"sync"
	"time"
)

const (
	EventTypeNone    EventType = 0
	EventTypeError   EventType = 1
	EventTypeWarning EventType = 2
	EventTypeInfo1   EventType = 3
	EventTypeInfo2   EventType = 4
	EventTypeInfo3   EventType = 5
)

type EventType int

type ILogProvider interface {
	LogError(scopeRef string, source string, message string)
	LogWarning(scopeRef string, source string, message string)
	LogInformation(scopeRef string, source string, message string)
}

type ILog interface {
	AllowLogging(eventType EventType) bool
	LogEvent(eventType EventType, source string, message string)
	StartLogScope(scopeRef LogScopeRef) *LogScope
}

type LogScopeRef struct {
	Ref string
}

type LogSettings struct {
	AllowErrors    bool
	AllowWarnings  bool
	MaxAllowedInfo EventType
}

type Log struct {
	LogSettings LogSettings
	Providers   []ILogProvider
}

func NewLog(logSettings LogSettings, providers []ILogProvider) ILog {
	return ILog(&Log{Providers: providers, LogSettings: logSettings})
}

func (m *Log) AllowLogging(eventType EventType) bool {

	if eventType == EventTypeError {
		return m.LogSettings.AllowErrors
	} else if eventType == EventTypeWarning {
		return m.LogSettings.AllowWarnings
	} else if eventType <= m.LogSettings.MaxAllowedInfo {
		return true
	}

	return false
}

func (m *Log) LogEvent(eventType EventType, source string, message string) {

	m.LogEventByRef("", eventType, source, message)
}

func (m *Log) LogEventByRef(scopeRef string, eventType EventType, source string, message string) {

	if !m.AllowLogging(eventType) {
		return
	}

	for _, provider := range m.Providers {

		if eventType == EventTypeError {
			provider.LogError(scopeRef, source, message)
		} else if eventType == EventTypeWarning {
			provider.LogWarning(scopeRef, source, message)
		} else if eventType == EventTypeInfo1 || eventType == EventTypeInfo2 || eventType == EventTypeInfo3 {
			provider.LogInformation(scopeRef, source, message)
		}
	}
}

func (m *Log) StartLogScope(scopeRef LogScopeRef) *LogScope {

	logScope := &LogScope{Log: m, Ref: scopeRef.Ref}

	// if logScope.Ref == "" {

	// 	logScope.Ref = fmt.Sprintf("%s", uuid.NewV4())
	// }

	return logScope
}

type LogScope struct {
	Ref string
	Log *Log
}

func (m *LogScope) AllowLogging(eventType EventType) bool {
	return m.Log.AllowLogging(eventType)
}

func (m *LogScope) LogEvent(eventType EventType, source string, message string) {
	m.Log.LogEventByRef(m.Ref, eventType, source, message)
}

func (m *LogScope) StartLogScope(scopeRef LogScopeRef) *LogScope {

	return m.Log.StartLogScope(LogScopeRef{Ref: fmt.Sprintf("%[1]s:%[2]s", m.Ref, "" /*uuid.NewV4()*/)})
}

var consoleLogProviderMutex = &sync.Mutex{}

type ConsoleLogProvider struct {
}

func NewConsoleLogProvider() ILogProvider {
	return ILogProvider(&ConsoleLogProvider{})
}

func (m *ConsoleLogProvider) LogError(scopeRef string, source string, message string) {

	consoleLogProviderMutex.Lock()
	fmt.Printf("%s ERROR (%s): %s message: %s\n", time.Now().Format(time.RFC3339), scopeRef, source, message)
	consoleLogProviderMutex.Unlock()
}

func (m *ConsoleLogProvider) LogWarning(scopeRef string, source string, message string) {

	consoleLogProviderMutex.Lock()
	fmt.Printf("%s WARNING (%s): %s message: %s\n", time.Now().Format(time.RFC3339), scopeRef, source, message)
	consoleLogProviderMutex.Unlock()
}

func (m *ConsoleLogProvider) LogInformation(scopeRef string, source string, message string) {

	consoleLogProviderMutex.Lock()
	fmt.Printf("%s INFORMATION (%s): %s message: %s\n", time.Now().Format(time.RFC3339), scopeRef, source, message)
	consoleLogProviderMutex.Unlock()
}
