package dc

type Client interface {
	Connect(options ...interface{}) (chan struct{}, error)
	Services() *Remotes
	NotificationsChn() chan struct{}
}
