package ctxstore

const (
	ClientIP = Key("clientIp")
	TraceID  = Key("traceId")
	User     = Key("user")
)

type Key string

func (k Key) String() string {
	return string(k)
}
