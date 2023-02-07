package watchingdog

type LeaderTaker interface {
	TakeLeader() (isLeader bool, err error)
	Wait()
	ExtendDuration() (err error)
	Revoke() (err error)
}
