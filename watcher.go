package watchingdog

type Watcher interface {
	Watch() error
}
