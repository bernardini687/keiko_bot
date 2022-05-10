package main

import (
	"lib/handler/command"
	"lib/handler/telegram"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("1) received body: %#v\n", req.Body)

	update, err := telegram.ParseUpdate(req.Body)
	if err != nil {
		log.Println(err)

		return httpOK()
	}
	log.Printf("2) parsed update: %#v\n", update)

	cmd := command.NewCommand(update.Message.From.ID, update.Message.Text) // TODO: should this return (*Command, error)?
	result := command.Process(cmd)

	result, err = telegram.SendMessage(update.Message.Chat.ID, result, "HTML")
	if err != nil {
		log.Println(err)

		return httpOK()
	}
	log.Printf("99) telegram acknowledgement: %#v\n", result)

	return httpOK()
}

func main() {
	lambda.Start(handler)
}

func httpOK() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}
