package duckduckgo

import (
	"errors"

	"github.com/Clinet/clinet_convos"
	"github.com/Clinet/clinet_features"
	"github.com/Clinet/clinet_storage"
	duckduckgo "github.com/JoshuaDoes/duckduckgolang"
)

var Feature = features.Feature{
	Name: "duckduckgo",
	Desc: "DuckDuckGo is available as a conversation service. You can @Clinet with a question, and DuckDuckGo may answer it!",
	ServiceConvo: &ClientDuckDuckGo{},
}

type ClientDuckDuckGo struct {
	Client *duckduckgo.Client
}

func (ddg *ClientDuckDuckGo) Login() error {
	cfg := &storage.Storage{}
	if err := cfg.LoadFrom("duckduckgo"); err != nil {
		return err
	}

	appName := "Clinet"
	rawAppName, err := cfg.ConfigGet("cfg", "appName")
	if err != nil {
		cfg.ConfigSet("cfg", "appName", appName)
	} else {
		appName = rawAppName.(string)
	}

	ddg.Client = &duckduckgo.Client{
		AppName: appName,
	}
	return nil
}

func (ddg *ClientDuckDuckGo) Query(query *convos.ConversationQuery, lastState *convos.ConversationState) (*convos.ConversationResponse, error) {
	resp := &convos.ConversationResponse{}

	queryResult, err := ddg.Client.GetQueryResult(query.Text)
	if err != nil {
		return nil, err
	}

	result := ""
	if queryResult.Definition != "" {
		result = queryResult.Definition
	} else if queryResult.Answer != "" {
		result = queryResult.Answer
	} else if queryResult.AbstractText != "" {
		result = queryResult.AbstractText
	}

	if result == "" {
		return nil, errors.New("duckduckgo: empty result")
	}
	resp.TextSimple = result

	if queryResult.Image != "" {
		resp.ImageURL = "https://duckduckgo.com" + queryResult.Image
	}

	return resp, nil
}