package torque

import "time"

// FlagParser is capable of parsing commmand-line arugments
type FlagParser interface {
	ParseFlags(action string, args []string) error
}

// TimestampFlag is a custom command-line flag for accepting timestamps
type TimestampFlag time.Time

func (ts *TimestampFlag) String() string {
	return time.Time(*ts).String()
}

// Set reads the raw string value into a TimestampFlag, or dies
// trying...actually it just returns nil.
func (ts *TimestampFlag) Set(value string) error {
	// TODO Change to a  list of valid timestamp formats
	MyTimeFormat := "2006Jan21504"
	t, err := time.Parse(MyTimeFormat, value)
	if err != nil {
		return err
	}
	*ts = TimestampFlag(t)
	return nil
}
