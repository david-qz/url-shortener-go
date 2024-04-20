package urls

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var client *dynamodb.Client
var table string

func init() {
	table = os.Getenv("TABLE")
	if table == "" {
		log.Fatal("Missing TABLE environment variable")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

type url struct {
	Key string `dynamodbav:"key"`
	Url string `dynamodbav:"url"`
}

func (u url) getKey() map[string]types.AttributeValue {
	key, _ := attributevalue.Marshal(u.Key)
	return map[string]types.AttributeValue{"key": key}
}

func GetUrl(key string, ctx context.Context) (*string, error) {
	u := url{Key: key}
	response, err := client.GetItem(ctx, &dynamodb.GetItemInput{Key: u.getKey(), TableName: &table})
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &u)
	if err != nil {
		return nil, err
	}

	return &u.Url, nil
}
