package command

import (
	"fmt"
	"lib/handler/bucket"

	kakebo "github.com/bernardini687/kakebo-golang"
)

func Dues(cmd Command) string {
	dues, err := bucket.GetContentFromKey(fmt.Sprintf("%s/%s.txt", cmd.Namespace, "dues"))
	if err != nil {
		return err.Error()
	}

	return kakebo.DisplayDues(dues)
}
