package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/graphtest"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"google.golang.org/grpc"
)

// GraphQLUpload describes a file attached to a gqlgen multipart request.
// VarPath is the JSON pointer used in the "map" part of the request, e.g.
// "variables.input.image" or "variables.input.images.0".
type GraphQLUpload struct {
	VarPath  string
	Filename string
	Content  []byte
}

// NewGraphQLHandler wires the apigateway GraphQL gateway (via its public
// graphtest helper) to the given in-process gRPC connections.
func NewGraphQLHandler(conns map[string]*grpc.ClientConn, cacheStore *cache.CacheStore, log logger.LoggerInterface) http.Handler {
	return graphtest.NewTestHandler(conns, cacheStore, log)
}

// NewAuthenticatedGraphQLHandler wires the gateway with the production auth
// middleware so protected operations (e.g. getMe) resolve the caller from the
// Authorization: Bearer header.
func NewAuthenticatedGraphQLHandler(conns map[string]*grpc.ClientConn, cacheStore *cache.CacheStore, log logger.LoggerInterface, tm auth.TokenManager) http.Handler {
	return graphtest.NewAuthenticatedTestHandler(conns, cacheStore, log, tm)
}

// DoGraphQLAuth executes a GraphQL document with an Authorization: Bearer
// header and returns the decoded JSON body.
func DoGraphQLAuth(h http.Handler, query string, bearerToken string) (map[string]any, error) {
	return doGraphQLBody(h, func() (*http.Request, error) {
		payload, err := json.Marshal(map[string]any{"query": query})
		if err != nil {
			return nil, err
		}
		req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		if bearerToken != "" {
			req.Header.Set("Authorization", "Bearer "+bearerToken)
		}
		return req, nil
	})
}

// DoGraphQL executes a GraphQL document (no variables) and returns the decoded
// JSON body. A non-200 response or undecodable body is returned as an error.
func DoGraphQL(h http.Handler, query string) (map[string]any, error) {
	return doGraphQLBody(h, func() (*http.Request, error) {
		payload, err := json.Marshal(map[string]any{"query": query})
		if err != nil {
			return nil, err
		}
		req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
}

// DoGraphQLWithUploads executes a GraphQL document over gqlgen's multipart
// transport so `Upload` scalar fields can be populated.
func DoGraphQLWithUploads(h http.Handler, query string, vars map[string]any, uploads []GraphQLUpload) (map[string]any, error) {
	return doGraphQLBody(h, func() (*http.Request, error) {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)

		operations, err := json.Marshal(map[string]any{
			"query":     query,
			"variables": vars,
		})
		if err != nil {
			return nil, err
		}
		if err := w.WriteField("operations", string(operations)); err != nil {
			return nil, err
		}

		fileMap := make(map[string][]string)
		for i, up := range uploads {
			fileMap[fmt.Sprintf("%d", i)] = []string{up.VarPath}
		}
		fileMapJSON, err := json.Marshal(fileMap)
		if err != nil {
			return nil, err
		}
		if err := w.WriteField("map", string(fileMapJSON)); err != nil {
			return nil, err
		}

		for i, up := range uploads {
			fw, err := w.CreateFormFile(fmt.Sprintf("%d", i), up.Filename)
			if err != nil {
				return nil, err
			}
			if _, err := fw.Write(up.Content); err != nil {
				return nil, err
			}
		}

		if err := w.Close(); err != nil {
			return nil, err
		}

		req := httptest.NewRequest(http.MethodPost, "/graphql", &buf)
		req.Header.Set("Content-Type", w.FormDataContentType())
		return req, nil
	})
}

func doGraphQLBody(h http.Handler, build func() (*http.Request, error)) (map[string]any, error) {
	req, err := build()
	if err != nil {
		return nil, err
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		return nil, fmt.Errorf("graphql request returned HTTP %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		return nil, fmt.Errorf("failed to decode graphql response: %w: %s", err, rec.Body.String())
	}
	return body, nil
}

// GraphQLErrorMessages extracts the messages of any top-level GraphQL errors.
func GraphQLErrorMessages(body map[string]any) []string {
	raw, ok := body["errors"].([]any)
	if !ok {
		return nil
	}
	messages := make([]string, 0, len(raw))
	for _, e := range raw {
		em, ok := e.(map[string]any)
		if !ok {
			continue
		}
		if msg, ok := em["message"].(string); ok {
			messages = append(messages, msg)
		}
	}
	return messages
}

// GraphQLOpData returns the `data.<op>` object of a GraphQL response.
func GraphQLOpData(body map[string]any, op string) map[string]any {
	data, ok := body["data"].(map[string]any)
	if !ok {
		return nil
	}
	res, _ := data[op].(map[string]any)
	return res
}

// GraphQLOpEnvelopeData returns the `data.<op>.data` payload, i.e. the inner
// domain payload of the standard {status, message, data} API envelope.
func GraphQLOpEnvelopeData(body map[string]any, op string) map[string]any {
	opRes := GraphQLOpData(body, op)
	if opRes == nil {
		return nil
	}
	payload, _ := opRes["data"].(map[string]any)
	return payload
}
