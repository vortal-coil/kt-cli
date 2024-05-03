package pkg

// Logger is a function type used to log messages
// This library doesn't log anything by default. You can set your own logger using SetLogger function
type Logger func(content string, params ...interface{})

// emptyLogger is a default logger that does nothing
func emptyLogger(content string, params ...interface{}) {}

// currentLogger is a singleton logger used by the library
var currentLogger Logger

func init() {
	currentLogger = emptyLogger
}

func SetLogger(logger Logger) {
	currentLogger = logger
}
