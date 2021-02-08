package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("EVENT: %s", request.Body)

	// Rendo il json in arrivo un oggetto utilizzabile
	/* var reqBody interface{}
	json.Unmarshal([]byte(request.Body), &reqBody)

	// Costruisco l'oggetto che manderò come notifica a Slack
	msg := buildMessage(&reqBody)

	// Rendo il messaggio un json
	requestBody, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln(err)
	}
	// Mando la notifica a slack, leggendo l'indirizzo della app dalla variabile di sistema che ho settato su Netlify
	resp, err := http.Post(os.Getenv("SLACK_ENDPOINT"), "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	// Ottengo la risposta da Slack dopo l'invio della notifica
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	// log
	log.Println((string(body)))

	// Restituisco la risposta di Slack a Netlify
	return &events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}, nil*/
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string("OK"),
	}, nil
}

func buildMessage(reqBody *interface{}) map[string]interface{} {
	var buildingInfo []map[string]string
	var messageText string
	// TODO: bisognerebbe capire se la build è andata a buon fine oppure no
	//       e cambiare di conseguenza il tipo di messaggio da inviare a Slack
	//       per ora sospetto che vada controllata la variabile "state", che quando
	//       la build è andata a buon fine è "ready"
	if (*reqBody).(map[string]interface{})["manual_deploy"] == true {
		buildingInfo = []map[string]string{
			{
				"type": "mrkdwn",
				"text": ":gear: Deploy manuale",
			},
		}
		messageText = ":pencil2: Build fatta partire manualmente: probabilmente è per test.\nNon dovrebbero esserci novità nel sito."
	} else {
		buildingInfo = []map[string]string{
			{
				"type": "mrkdwn",
				"text": fmt.Sprintf(":gear: %s", (*reqBody).(map[string]interface{})["branch"]),
			}, {
				"type": "mrkdwn",
				"text": fmt.Sprintf("🦸 %s", (*reqBody).(map[string]interface{})["committer"]),
			},
		}
		messageText = fmt.Sprintf(":pencil2: _%s_\n\n\n:basketball: <%s|Vai al sito>", (*reqBody).(map[string]interface{})["title"], (*reqBody).(map[string]interface{})["ssl_url"])
	}
	return map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color": "#2eb886",
				"blocks": []map[string]interface{}{
					{
						"type": "section",
						"text": map[string]string{
							"type": "mrkdwn",
							"text": ":tada: *Nuova build del sito Basket Gardolo*",
						},
					},
					{
						"type":     "context",
						"elements": buildingInfo,
					},
					{
						"type": "section",
						"text": map[string]string{
							"type": "mrkdwn",
							"text": messageText,
						},
						"accessory": map[string]string{
							"type":      "image",
							"image_url": "https://www.basketgardolo.it/wp-content/uploads/2014/09/logo_basket_gardolo.png",
							"alt_text":  "logo BC Gardolo",
						},
					},
				},
			},
		},
	}
}

func main() {
	lambda.Start(handler)
}
