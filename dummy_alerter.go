package watchingdog

type DummyAlerter struct {
	MsgBag []string
}

func NewDummyAlerter() Alerter {
	return &DummyAlerter{
		MsgBag: make([]string, 0),
	}
}

func (d *DummyAlerter) Alert() error {
	msg := "some thing went wrong yo"
	d.MsgBag = append(d.MsgBag, msg)
	return nil
}
