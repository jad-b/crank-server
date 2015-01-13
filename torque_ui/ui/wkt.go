package ui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jad-b/torque"
)

// WktTimeLayout defines the common .wkt timestamp format
const WktTimeLayout = "2006 Jan 02 @ 1504"

// WorkoutToWkt converts a Workout into a .wkt textual representation
func WorkoutToWkt(w *Workout) (wkt string, err error) {
	var buf bytes.Buffer

	// Output timestamp header
	wktTimestamp := w.LastModified.Format(WktTimeLayout)
	if _, err := buf.WriteString(wktTimestamp + "\n"); err != nil {
		return "", err
	}

	// Workout tags
	if _, err := buf.WriteString(FormatTags(w.Tags)); err != nil {
		return "", err
	}

	// Exercises
	if _, err := buf.WriteString(FormatExercises(w.Exercises)); err != nil {
		return "", err
	}

	wkt = buf.String()
	return wkt, nil
}

// FormatTags outputs a string in tag form
func FormatTags(in string) string {
	var buf bytes.Buffer
	var err error
	for _, t := range strings.Split(in, ";") {
		kvp := strings.SplitN(t, "=", 2) // key-value pair(kvp)
		if len(kvp) > 1 {                // - key: value
			_, err = buf.WriteString(
				fmt.Sprintf("- %s: %s\n",
					strings.TrimSpace(kvp[0]), strings.TrimSpace(kvp[1])))
		} else { // - comment...
			_, err = buf.WriteString("- " + strings.TrimSpace(kvp[0]) + "\n")
		}
		if err != nil {
			return ""
		}
	}
	return buf.String()
}

// FormatExercises outputs a list of Exercises as a newline-delimited string
func FormatExercises(exs []Exercise) string {
	var buf bytes.Buffer
	var err error
	for _, ex := range exs {
		_, err = buf.WriteString(fmt.Sprintf("%s: %s", ex.Movement, ex.Sets) + "\n")
		_, err = buf.WriteString(FormatTags(ex.Tags))
		if err != nil {
			return ""
		}
	}
	return buf.String()
}

// StringToWktTime converts a normal timestamp into a .wkt timestamp
func StringToWktTime(ts string) (wktTime string, err error) {
	t, err := torque.ParseTimestamp(ts)
	if err != nil {
		return "", err
	}
	return t.Format(WktTimeLayout), nil
}
