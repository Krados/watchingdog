package watchingdog

type Alerter interface {
	Alert() error
}
