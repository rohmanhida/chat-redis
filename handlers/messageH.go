package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Message struct {
	Type    string `json:"type"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	// websocket init and upgrade
	websocketInit(w, r)
	addClient()

	// send message to let people know someone joined
	announcement := []string{"presence", "user", "Someone joined the chat!"}
	sendMessage(announcement, rdb)

	log.Println("New client connected")
	// Read messages from this client
	for {

		_, msg, err := Conn.ReadMessage()
		if err != nil {
			// send message to let people know someone left
			announcement := []string{"presence", "user", "Someone's just left"}
			sendMessage(announcement, rdb)

			log.Println("Client disconnected:", err)
			break
		}

		ip := getIP(r)
		limiter := getLimiter(ip)
		if limiter.Allow() {
			// validate json
			encoded, err := validateJSON(msg)
			if err != nil {
				log.Println("Invalid JSON message:", err)
				continue
			}

			// Publish message to Redis channel
			err = rdb.Publish(ctx, "chatroom", encoded).Err()
			if err != nil {
				log.Println("Redis publish error:", err)
			}
		}
	}

	removeClient(Conn)
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		log.Println("this is forwarded value", forwarded)
		return forwarded
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func validateJSON(input []byte) ([]byte, error) {

	// unmarshal json input
	var message Message
	if err := json.Unmarshal(input, &message); err != nil {
		return nil, err
	}

	// validation
	policy := bluemonday.UGCPolicy()
	safe := policy.Sanitize(message.Content)
	safe2 := html.EscapeString(safe)

	// marshal again after validation
	if safe2 == "" {
		return nil, errors.New("invalid message content")
	}

	validJSON, _ := json.Marshal(message)
	return validJSON, nil
}

func sendMessage(msg []string, rdb *redis.Client) {
	// build message and publish
	message := Message{Type: msg[0], Sender: msg[1], Content: msg[2]}
	encoded, _ := json.Marshal(message)
	rdb.Publish(ctx, "chatroom", encoded)
}

func SubscribeRedis(rdb *redis.Client) {
	sub := rdb.Subscribe(ctx, "chatroom")
	ch := sub.Channel()

	for msg := range ch {
		broadcast([]byte(msg.Payload))
	}
}

func broadcast(message []byte) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	for client := range Clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(Clients, client)
		}
	}
}
