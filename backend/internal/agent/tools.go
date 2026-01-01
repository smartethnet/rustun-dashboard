package agent

import (
	"encoding/json"
	"fmt"

	"github.com/smartethnet/rustun-dashboard/internal/model"
	"github.com/smartethnet/rustun-dashboard/internal/service"
)

// ToolExecutor handles tool execution
type ToolExecutor struct {
	routeService *service.RouteService
}

// NewToolExecutor creates a new tool executor
func NewToolExecutor(routeService *service.RouteService) *ToolExecutor {
	return &ToolExecutor{
		routeService: routeService,
	}
}

// GetTools returns all available tools
func (te *ToolExecutor) GetTools() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "list_clusters",
				Description: "Get all clusters with their names and client counts",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "list_clients",
				Description: "Get client list, optionally filtered by cluster. Returns detailed client information including identity, name, IP address, etc.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Cluster name to filter clients. If not provided, returns all clients",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "get_client",
				Description: "Get detailed information of a single client by cluster name and client identity",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Cluster name where the client belongs",
						},
						"identity": map[string]interface{}{
							"type":        "string",
							"description": "Unique client identifier (UUID)",
						},
					},
					"required": []string{"cluster", "identity"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "create_client",
				Description: "Create a new VPN client. System will automatically generate client identity (UUID) and IP address configuration. If specified cluster doesn't exist, it will be created automatically",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Cluster name where the client belongs",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Friendly name for the client, e.g.: Headquarters, Branch, NAS, Laptop, Phone",
						},
						"ciders": map[string]interface{}{
							"type":        "array",
							"description": "CIDR route list for the client, e.g. [\"192.168.1.0/24\"]",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"cluster"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "update_client",
				Description: "Update existing client information. Can modify name and CIDR routes, but cannot modify cluster, identity and IP configuration",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Cluster name where the client belongs",
						},
						"identity": map[string]interface{}{
							"type":        "string",
							"description": "Unique client identifier (UUID)",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "New friendly name for the client",
						},
						"ciders": map[string]interface{}{
							"type":        "array",
							"description": "New CIDR route list",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"cluster", "identity"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "delete_client",
				Description: "Delete specified client. After deletion, the IP address occupied by the client will be released",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Cluster name where the client belongs",
						},
						"identity": map[string]interface{}{
							"type":        "string",
							"description": "Unique client identifier (UUID)",
						},
					},
					"required": []string{"cluster", "identity"},
				},
			},
		},
	}
}

// ExecuteTool executes a tool function
func (te *ToolExecutor) ExecuteTool(name string, arguments string) (string, error) {
	switch name {
	case "list_clusters":
		return te.listClusters()
	case "list_clients":
		return te.listClients(arguments)
	case "get_client":
		return te.getClient(arguments)
	case "create_client":
		return te.createClient(arguments)
	case "update_client":
		return te.updateClient(arguments)
	case "delete_client":
		return te.deleteClient(arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func (te *ToolExecutor) listClusters() (string, error) {
	clusters, err := te.routeService.GetAllClusters()
	if err != nil {
		return "", fmt.Errorf("获取集群列表失败: %w", err)
	}

	result, err := json.Marshal(clusters)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(result), nil
}

func (te *ToolExecutor) listClients(arguments string) (string, error) {
	var args struct {
		Cluster string `json:"cluster"`
	}

	if arguments != "" {
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return "", fmt.Errorf("解析参数失败: %w", err)
		}
	}

	var clients []model.Client
	var err error

	if args.Cluster != "" {
		clients, err = te.routeService.GetClientsByCluster(args.Cluster)
	} else {
		clients, err = te.routeService.GetAllClients()
	}

	if err != nil {
		return "", fmt.Errorf("获取客户端列表失败: %w", err)
	}

	result, err := json.Marshal(clients)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(result), nil
}

func (te *ToolExecutor) getClient(arguments string) (string, error) {
	var args struct {
		Cluster  string `json:"cluster"`
		Identity string `json:"identity"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	client, err := te.routeService.GetClient(args.Cluster, args.Identity)
	if err != nil {
		return "", fmt.Errorf("获取客户端失败: %w", err)
	}

	result, err := json.Marshal(client)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(result), nil
}

func (te *ToolExecutor) createClient(arguments string) (string, error) {
	var req model.ClientCreateRequest

	if err := json.Unmarshal([]byte(arguments), &req); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	// Convert ClientCreateRequest to Client for service layer
	// The service will handle identity and IP generation
	client := model.Client{
		Cluster: req.Cluster,
		Name:    req.Name,
		Ciders:  req.Ciders,
	}

	createdClient, err := te.routeService.CreateClient(client)
	if err != nil {
		return "", fmt.Errorf("创建客户端失败: %w", err)
	}

	result, err := json.Marshal(createdClient)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(result), nil
}

func (te *ToolExecutor) updateClient(arguments string) (string, error) {
	var args struct {
		Cluster  string   `json:"cluster"`
		Identity string   `json:"identity"`
		Name     string   `json:"name"`
		Ciders   []string `json:"ciders"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	// Get existing client first to preserve IP config
	existingClient, err := te.routeService.GetClient(args.Cluster, args.Identity)
	if err != nil {
		return "", fmt.Errorf("获取现有客户端失败: %w", err)
	}

	// Update only the allowed fields
	updatedClient := model.Client{
		Cluster:   args.Cluster,
		Identity:  args.Identity,
		Name:      args.Name,
		PrivateIP: existingClient.PrivateIP,
		Mask:      existingClient.Mask,
		Gateway:   existingClient.Gateway,
		Ciders:    args.Ciders,
	}

	if err := te.routeService.UpdateClient(args.Cluster, args.Identity, updatedClient); err != nil {
		return "", fmt.Errorf("更新客户端失败: %w", err)
	}

	result, err := json.Marshal(updatedClient)
	if err != nil {
		return "", fmt.Errorf("序列化结果失败: %w", err)
	}

	return string(result), nil
}

func (te *ToolExecutor) deleteClient(arguments string) (string, error) {
	var args struct {
		Cluster  string `json:"cluster"`
		Identity string `json:"identity"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	if err := te.routeService.DeleteClient(args.Cluster, args.Identity); err != nil {
		return "", fmt.Errorf("删除客户端失败: %w", err)
	}

	return `{"success": true, "message": "客户端已成功删除"}`, nil
}

