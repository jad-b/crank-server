package torque

import "time"

// FlagParser is capable of parsing commmand-line arugments
type FlagParser interface {
	ParseFlags(action string, args []string) error
}

// Custom timestamp flag
type timestamp time.Time

func (ts *timestamp) String() string {
	return time.Time(*ts).String()
}

func (ts *timestamp) Set(value string) error {
	MyTimeFormat := "2006Jan21504"
	t, err := time.Parse(MyTimeFormat, value)
	if err != nil {
		return err
	}
	*ts = timestamp(t)
	return nil
}
