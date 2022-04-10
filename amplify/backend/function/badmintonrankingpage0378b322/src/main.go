package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type player struct {
	name      string
	player_id string
}

type player_ranking struct {
	Name  string
	Grade string
}

type sessionid_secret struct {
	SessionId string `json:"ebadders_sessionid"`
}

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("GOBADDERS_SESSIONID")),
	}
	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		log.Fatalf("unable to get secret, %v", err)
	}

	res := sessionid_secret{}
	err = json.Unmarshal([]byte(*result.SecretString), &res)
	if err != nil {
		log.Fatalf("unable to unmarshal secret, %v", err)
	}

	sessionid := res.SessionId

	player_ids := [3]player{
		{name: "Seb", player_id: "8216"},
		{name: "Joe", player_id: "8186"},
		{name: "Ben", player_id: "8187"},
	}

	ranking := make([]player_ranking, 0, 3)

	for _, player := range player_ids {
		ranking = append(ranking, *GetRanking(player, sessionid))
	}

	json, err := json.Marshal(ranking)

	log.Println(string(json))

	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		StatusCode: 200,
	}, nil
}

func GetRanking(p player, sessionid string) *player_ranking {
	req, err := http.NewRequest("GET", "https://ebadders.com/clubs/152/players/"+p.player_id, nil)
	cookie := http.Cookie{
		Name:  "sessionid",
		Value: sessionid,
	}
	log.Println(sessionid)
	req.AddCookie(&cookie)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := doc.Html()
	log.Println(out)

	grade := doc.Find("#pane").Children().Last().Text()
	grade = strings.TrimSpace(grade)

	grade = strings.Split(grade, " ")[1]

	player_ranking := player_ranking{
		Name:  p.name,
		Grade: grade,
	}

	return &player_ranking
}

func main() {
	lambda.Start(HandleRequest)
}
