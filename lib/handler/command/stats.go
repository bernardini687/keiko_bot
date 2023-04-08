package command

import (
	"fmt"
	"lib/handler/bucket"
	"os"
	"strconv"
	"time"

	kakebo "github.com/bernardini687/kakebo-golang"
)

func Stats(cmd Command) string {
	key := duesKey(cmd.Namespace)

	dues, err := bucket.GetContentFromKey(key)
	if err != nil {
		return err.Error()
	}

	bal, err := kakebo.CalcBalance(dues)
	if err != nil {
		return err.Error()
	}

	date := time.Now()

	entries, err := GetEntries(date, cmd)
	if err != nil {
		return err.Error()
	}

	monthTot, err := kakebo.CalcMonth(entries)
	if err != nil {
		return err.Error()
	}

	var savePercentage int
	if goal, found := os.LookupEnv("KEIKO_GOAL"); found {
		num, err := strconv.Atoi(goal)
		if err == nil {
			savePercentage = num
		}
	}

	return kakebo.DisplayStats(date, bal, monthTot, savePercentage)
}

// duesKey build the S3 key to get to the dues data
//
// Example output:
//
//	"15012019/dues.txt"
func duesKey(namespace string) string {
	return fmt.Sprintf("%s/%s.txt", namespace, "dues")
}
