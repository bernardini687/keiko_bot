package command

import (
	"fmt"
	"lib/handler/bucket"
	"log"
	"os"
	"strconv"
	"time"

	kakebo "github.com/bernardini687/kakebo-golang"
)

// AddEntries accepts an entries command and, if it is formatted properly, attempt to append it
// to the bucket object `CHAT_ID/YYYY/MM.txt`.
// If the object does not yet exist, it is created.
func AddEntries(cmd Command) string {
	contents, err := kakebo.FormatEntries(cmd.Text)
	if err != nil {
		return err.Error()
	}

	// return "<pre>" + newEntries + "</pre>"

	result, err := appendOrUpload(cmd.Namespace, contents)
	if err != nil {
		return err.Error()
	}

	return result
}

func ReadEntries(cmd Command) string {
	var num int
	var err error

	if cmd.Arg != "" {
		num, err = strconv.Atoi(cmd.Arg)
	}
	if err != nil {
		return err.Error()
	}

	date := designatedTime(num)
	log.Printf("5) date UTC: %v\n", date)

	entries, err := GetEntries(date, cmd)
	if err != nil {
		return err.Error()
	}

	tot, err := kakebo.CalcMonth(entries)
	if err != nil {
		return err.Error()
	}

	return kakebo.DisplayMonth(date, entries, tot)
}

func GetEntries(date time.Time, cmd Command) (string, error) {
	key := entriesKey(cmd.Namespace, date)

	entries, err := bucket.GetContentFromKey(key)
	if err != nil {
		return "", err
	}

	return entries, nil
}

func appendOrUpload(namespace, newContents string) (string, error) {
	now := time.Now() // .UTC()
	key := entriesKey(namespace, now)
	id := bucket.NewID(os.Getenv("BUCKET_NAME"), key)

	sess := bucket.NewSession()
	client := bucket.NewClient(sess)

	found, err := bucket.LookupKey(client, id)
	if err != nil {
		return "", err
	}

	contents := newContents

	if found {
		// we get the object's contents so we can append data
		entries, err := bucket.GetContent(client, id)
		if err != nil {
			return "", err
		}

		contents = entries + newContents
	}

	log.Printf("6) contents: %#v\n", contents)
	log.Printf("7) now: %v\n", now)

	uploader := bucket.NewUploader(sess)

	err = bucket.UploadContents(uploader, id, contents)
	if err != nil {
		return "", err
	}

	return contents, nil
}

func designatedTime(num int) time.Time {
	now := time.Now().UTC()
	log.Printf("4) now UTC: %v\n", now)

	return now.AddDate(0, num, 0)
}

// entriesKey build the S3 key to get to the month of given date's entries
//
// Example output:
//
//	"15012019/2006/01.txt"
func entriesKey(namespace string, date time.Time) string {
	return fmt.Sprintf("%s/%s.txt", namespace, date.Format("2006/01"))
}
