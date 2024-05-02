package pkg

type Logger func(content string, params ...interface{})

func emptyLogger(content string, params ...interface{}) {}

var currentLogger Logger

func init() {
	currentLogger = emptyLogger
}

func SetLogger(logger Logger) {
	currentLogger = logger
}
