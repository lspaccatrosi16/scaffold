package error

import "fmt"

func Error() error {
	return fmt.Errorf("(internal command error)")
}
