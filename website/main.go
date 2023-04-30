package main

import (
	"context"
	"fmt"
	"log"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func main() {
	projectID := "digitallibrary-379908"
	ctx := context.Background()
	client, err := createDialogflowClient(ctx)
	if err != nil {
		log.Fatalf("Error creating Dialogflow client: %v", err)
	}
	defer client.Close()

	r := gin.Default()

	r.Use(corsMiddleware())

	r.POST("/api/chatbot", func(c *gin.Context) {
		// Parse user input from the request body
		var userInput struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&userInput); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Bad request"})
			return
		}

		response, err := detectIntentText(ctx, client, projectID, "unique-session-id", userInput.Message, "en")
		if err != nil {
			log.Printf("Error detecting intent: %v", err)
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		fulfillmentText := response.GetQueryResult().GetFulfillmentText()
		log.Printf("User: %s | Bot: %s", userInput.Message, fulfillmentText)

		// Send bot's response back to the frontend
		c.JSON(200, gin.H{"message": fulfillmentText})
	})

	log.Println("Starting server...")
	log.Fatal(r.Run(":8000"))
}

func createDialogflowClient(ctx context.Context) (*dialogflow.SessionsClient, error) {
	serviceAccountKeyPath := "digitallibrary-379908-1b714bb8b0dc.json"

	client, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func detectIntentText(ctx context.Context, client *dialogflow.SessionsClient, projectID, sessionID, text, languageCode string) (*dialogflowpb.DetectIntentResponse, error) {
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	textInput := dialogflowpb.TextInput{
		Text:         text,
		LanguageCode: languageCode,
	}

	queryInput := &dialogflowpb.QueryInput{
		Input: &dialogflowpb.QueryInput_Text{
			Text: &textInput,
		},
	}

	request := &dialogflowpb.DetectIntentRequest{
		Session:    sessionPath,
		QueryInput: queryInput,
	}

	response, err := client.DetectIntent(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
