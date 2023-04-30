# Golang-Book-Recommendation-Service
This repository contains a Golang service that provides book recommendations based on genre or author. It uses the Gorilla Mux library for routing and Dialogflow API to handle user input. The service can be deployed to a server and connected to a chatbot frontend to provide book recommendations to users.

## Prerequisites
- Golang 1.16 or later installed on your machine
- A Dialogflow account and API credentials
- A server to deploy the service
- ngrok a simplified API-first ingress-as-a-service that adds connectivity, security, and observability to your apps in one line 

## Setting up
Clone repository then move to directory:
```
git clone https://github.com/garywahyudi/Golang-Book-Recommendation-Service.git
cd Golang-Book-Recommendation-Service/
```

In my case, I uses ngrok for the webhook and as for the backend to hit dialogflow was just using local porting system, and for ngrok itself what you need is to set up ngrok locally then with that you can run following code:
```
ngrok http --inspect=true 8080
```

You will see 2 folders with the name website and webhook, the flow for this on a local setup would just be running both the backend under the folder `website` and the webook under the folder `webhook` with the following code:
```
go run main.go
```

In case of a dependency error try running this before executing the `go run ...` command:
```
go mod tidy
```



## Details related to port in case of online deployment:
- backend service under `website` uses port = 8000
- webhook service under `webhook` by default uses = 8080

## Usage
To use the chatbot, simply type your message in the input field and hit enter. The chatbot will respond with a message based on the natural language processing and recommendation logic provided by Dialogflow.

## Acknowledgements
This project uses the following third-party packages:
- using ngrok for local virtual web server for backend services (https://ngrok.com/)
- vscode's live server for quick local frontend services
