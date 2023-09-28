package client

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hasura/go-graphql-client"
)

type MonteCarloTransport struct {
	API_KEY_ID    string
	API_KEY_TOKEN string
	context       context.Context
}

type MonteCarloClient struct {
	client    *graphql.Client
	transport *MonteCarloTransport
}

func NewMonteCarloClient(context context.Context, api_key_id string, api_key_token string) (*MonteCarloClient, error) {
	transport := MonteCarloTransport{api_key_id, api_key_token, context}
	client := graphql.NewClient("https://api.getmontecarlo.com/graphql", &http.Client{Transport: transport})
	return &MonteCarloClient{client, &transport}, nil
}

func (mc *MonteCarloClient) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return mc.client.Mutate(ctx, m, variables, options...)
}

func (mc *MonteCarloClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return mc.client.Query(ctx, q, variables, options...)
}

func (mc *MonteCarloClient) ExecRaw(ctx context.Context, query string, variables map[string]interface{}, options ...graphql.Option) ([]byte, error) {
	return mc.client.ExecRaw(ctx, query, variables, options...)
}

func (transport MonteCarloTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(req, true)
	req.Header.Set("x-mcd-id", transport.API_KEY_ID)
	req.Header.Set("x-mcd-token", transport.API_KEY_TOKEN)
	resp, err := http.DefaultTransport.RoundTrip(req)
	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)
	tflog.Debug(transport.context, string(bytes))
	return resp, err
}

type Diagnostic struct {
	Cause           string
	FriendlyMessage string
	Resolution      string
}

type Warnings []Diagnostic
type Errors []Diagnostic

type TestBqCredentialsV2 struct {
	TestBqCredentialsV2 struct {
		Key              string
		ValidationResult struct {
			Success  bool
			Warnings Warnings
			Errors   Errors
		}
	} `graphql:"testBqCredentialsV2(validationName: $validationName, connectionDetails: $connectionDetails)"`
}

type AddConnection struct {
	AddConnection struct {
		Connection struct {
			Uuid      string
			Warehouse struct {
				Name string
				Uuid string
			}
		}
	} `graphql:"addConnection(dcId: $dcId, dwId: $dwId, key: $key, jobTypes: $jobTypes, name: $name, connectionType: $connectionType, createWarehouseType: $createWarehouseType)"`
}

type GetWarehouse struct {
	GetWarehouse *struct {
		Name        string `json:"name"`
		Connections []struct {
			Uuid string `json:"uuid"`
			Type string `json:"type"`
		} `json:"connections"`
	} `json:"getWarehouse"`
}

type RemoveConnection struct {
	RemoveConnection struct {
		Success bool
	} `graphql:"removeConnection(connectionId: $connectionId)"`
}

type SetWarehouseName struct {
	SetWarehouseName struct {
		Warehouse struct {
			Uuid string
			Name string
		}
	} `graphql:"setWarehouseName(dwId: $dwId, name: $name)"`
}

type UpdateCredentials struct {
	UpdateCredentials struct {
		Success bool
	} `graphql:"updateCredentials(changes: $changes, connectionId: $connectionId, shouldReplace: $shouldReplace, shouldValidate: $shouldValidate)"`
}