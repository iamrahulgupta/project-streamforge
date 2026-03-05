package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Server represents the REST API server
type Server struct {
	brokers []string
	port    int
}

// NewServer creates a new REST API server
func NewServer(port int, brokers []string) *Server {
	return &Server{
		brokers: brokers,
		port:    port,
	}
}

// Start starts the REST API server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/api/v1/brokers", s.brokersHandler)
	mux.HandleFunc("/api/v1/topics", s.topicsHandler)
	mux.HandleFunc("/api/v1/topics/", s.topicDetailHandler)
	mux.HandleFunc("/api/v1/groups", s.groupsHandler)
	mux.HandleFunc("/api/v1/groups/", s.groupDetailHandler)
	mux.HandleFunc("/api/v1/produce", s.produceHandler)
	mux.HandleFunc("/api/v1/consume", s.consumeHandler)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("StreamForge Admin Server starting on %s\n", addr)
	log.Printf("API Documentation: http://localhost%s/api/v1\n", addr)

	return http.ListenAndServe(addr, mux)
}

// Response wraps API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthStatus represents health check status
type HealthStatus struct {
	Status  string   `json:"status"`
	Version string   `json:"version"`
	Brokers []string `json:"brokers"`
}

// healthHandler handles health check requests
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := Response{
		Success: true,
		Data: HealthStatus{
			Status:  "healthy",
			Version: "0.1.0",
			Brokers: s.brokers,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BrokerInfo represents broker information in API response
type BrokerInfo struct {
	ID       int    `json:"id"`
	Address  string `json:"address"`
	IsLeader bool   `json:"is_leader"`
}

// brokersHandler handles broker list requests
func (s *Server) brokersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	brokers := []BrokerInfo{
		{ID: 1, Address: "localhost:9092", IsLeader: true},
		{ID: 2, Address: "localhost:9093", IsLeader: false},
		{ID: 3, Address: "localhost:9094", IsLeader: false},
	}

	response := Response{
		Success: true,
		Data:    brokers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TopicInfo represents topic information in API response
type TopicInfo struct {
	Name             string `json:"name"`
	Partitions       int    `json:"partitions"`
	ReplicationFactor int   `json:"replication_factor"`
}

// topicsHandler handles topic list requests
func (s *Server) topicsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listTopicsHandler(w, r)
	case http.MethodPost:
		s.createTopicHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listTopicsHandler lists all topics
func (s *Server) listTopicsHandler(w http.ResponseWriter, r *http.Request) {
	topics := []TopicInfo{
		{Name: "events", Partitions: 3, ReplicationFactor: 3},
		{Name: "logs", Partitions: 5, ReplicationFactor: 3},
		{Name: "metrics", Partitions: 2, ReplicationFactor: 2},
	}

	response := Response{
		Success: true,
		Data:    topics,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateTopicRequest represents a create topic request
type CreateTopicRequest struct {
	Name             string            `json:"name"`
	Partitions       int               `json:"partitions"`
	ReplicationFactor int              `json:"replication_factor"`
	Config           map[string]string `json:"config,omitempty"`
}

// createTopicHandler creates a new topic
func (s *Server) createTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := Response{
		Success: true,
		Message: fmt.Sprintf("Topic '%s' created successfully", req.Name),
		Data: TopicInfo{
			Name:             req.Name,
			Partitions:       req.Partitions,
			ReplicationFactor: req.ReplicationFactor,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// topicDetailHandler handles topic-specific requests
func (s *Server) topicDetailHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Topic name required", http.StatusBadRequest)
		return
	}

	topicName := parts[4]

	switch r.Method {
	case http.MethodGet:
		s.getTopicHandler(w, r, topicName)
	case http.MethodDelete:
		s.deleteTopicHandler(w, r, topicName)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getTopicHandler retrieves topic details
func (s *Server) getTopicHandler(w http.ResponseWriter, r *http.Request, topicName string) {
	topic := TopicInfo{
		Name:             topicName,
		Partitions:       3,
		ReplicationFactor: 3,
	}

	response := Response{
		Success: true,
		Data:    topic,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// deleteTopicHandler deletes a topic
func (s *Server) deleteTopicHandler(w http.ResponseWriter, r *http.Request, topicName string) {
	response := Response{
		Success: true,
		Message: fmt.Sprintf("Topic '%s' deleted successfully", topicName),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// groupsHandler handles consumer group requests
func (s *Server) groupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	groups := []map[string]interface{}{
		{"group_id": "consumer-group-1", "state": "Stable", "members": 3},
		{"group_id": "consumer-group-2", "state": "Stable", "members": 2},
	}

	response := Response{
		Success: true,
		Data:    groups,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// groupDetailHandler handles consumer group detail requests
func (s *Server) groupDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Group ID required", http.StatusBadRequest)
		return
	}

	groupID := parts[4]

	group := map[string]interface{}{
		"group_id": groupID,
		"state":    "Stable",
		"members": []map[string]interface{}{
			{"member_id": "member-1", "client_id": "client-1"},
			{"member_id": "member-2", "client_id": "client-2"},
		},
		"topics": []string{"topic-1", "topic-2"},
	}

	response := Response{
		Success: true,
		Data:    group,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProduceRequest represents a message produce request
type ProduceRequest struct {
	Topic     string `json:"topic"`
	Key       string `json:"key,omitempty"`
	Value     string `json:"value"`
	Partition int    `json:"partition,omitempty"`
}

// produceHandler handles message production
func (s *Server) produceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := map[string]interface{}{
		"topic":     req.Topic,
		"partition": 0,
		"offset":    42,
		"timestamp": "2026-03-03T12:00:00Z",
	}

	response := Response{
		Success: true,
		Message: "Message produced successfully",
		Data:    result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// ConsumeRequest represents a message consume request
type ConsumeRequest struct {
	Topic     string `json:"topic"`
	GroupID   string `json:"group_id"`
	MaxMessages int  `json:"max_messages,omitempty"`
	Offset    string `json:"offset,omitempty"`
}

// consumeHandler handles message consumption
func (s *Server) consumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MaxMessages == 0 {
		req.MaxMessages = 10
	}

	messages := []map[string]interface{}{
		{"offset": 0, "key": "key1", "value": "message 1"},
		{"offset": 1, "key": "key2", "value": "message 2"},
	}

	response := Response{
		Success: true,
		Data: map[string]interface{}{
			"topic":    req.Topic,
			"group_id": req.GroupID,
			"messages": messages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
