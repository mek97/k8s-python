package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kubernetes/mcp-server/config"
	"github.com/kubernetes/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func ReadstorageversionHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url := fmt.Sprintf("%s/apis/internal.apiserver.k8s.io/v1alpha1/storageversions/%s", cfg.BaseURL)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No specific authentication scheme defined - add fallback authentication
		if cfg.BearerToken != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.BearerToken)
		} else if cfg.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		} else if cfg.BasicAuth != "" {
			req.Header.Set("Authorization", "Basic "+cfg.BasicAuth)
		}
		// Note: If no auth tokens provided, requests will be made without authentication
		
		// Add custom headers if provided
		
		// Set client identification headers
		req.Header.Set("X-Request-Source", "Codeglide-MCP-generator")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateReadstorageversionTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_apis_internal_apiserver_k8s_io_v1alpha1_storageversions",
		mcp.WithDescription("read the specified StorageVersion"),
	)

	return models.Tool{
		Definition: tool,
		Handler:    ReadstorageversionHandler(cfg),
	}
}
