package command

import (
	"fmt"
	"lib/handler/bucket"
	"os"
	"strconv"
	"time"

	kakebo "github.com/bernardini687/kakebo-golang"
)

// TODO: refactor Stats() - ReadEntries() - GetEntries() - LookupKey()
func Stats(cmd Command) string {
	key := duesKey(cmd.Namespace)
	id := bucket.NewID(os.Getenv("BUCKET_NAME"), key)

	sess := bucket.NewSession()
	client := bucket.NewClient(sess)
	found, err := bucket.LookupKey(client, id)
	if err != nil {
		return err.Error()
	}

	if !found {
		return fmt.Sprintf("no data at `%s`", key)
	}

	dues, err := bucket.GetContents(client, id)
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
//     "15012019/dues.txt"
//
func duesKey(namespace string) string {
	return fmt.Sprintf("%s/%s.txt", namespace, "dues")
}
