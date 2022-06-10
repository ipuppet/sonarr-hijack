package config

type Notifyer interface {
	Callback(*Config)
}

type EasyNotifyer struct {
	Notifyer
	callback func(c *Config)
}

func (a *EasyNotifyer) Callback(c *Config) {
	a.callback(c)
}

func NewNotifyer(callback func(c *Config)) Notifyer {
	n := &EasyNotifyer{
		callback: callback,
	}
	return n
}

func LoggerNotifyer() Notifyer {
	n := &EasyNotifyer{
		callback: func(c *Config) {
			logger.Printf("config %s reloaded", c.Filename)
		},
	}
	return n
}
