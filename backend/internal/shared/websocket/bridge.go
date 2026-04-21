package websocket

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// StartBridge subscribes to Redis and forwards messages to the WebSocket Hub
func StartBridge(ctx context.Context, hub *Hub, redisClient *redis.Client) {
	// Pattern: ws-events:{tenantID}
	pubsub := redisClient.PSubscribe(ctx, "ws-events:*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	log.Info().Msg("WebSocket bridge started, listening for Redis events...")

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				log.Warn().Msg("Redis pubsub channel closed, reconnecting in 5s...")
				time.Sleep(5 * time.Second)
				// Re-subscribe
				pubsub = redisClient.PSubscribe(ctx, "ws-events:*")
				ch = pubsub.Channel()
				continue
			}

			// Extract tenantID from channel name (ws-events:UUID)
			parts := strings.Split(msg.Channel, ":")
			if len(parts) < 2 {
				continue
			}
			tenantID := parts[1]

			var payload interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Error().Err(err).Str("payload", msg.Payload).Msg("Failed to unmarshal Redis event")
				continue
			}

			// Broadcast to local hub
			hub.broadcast <- BroadcastMessage{
				TenantID: tenantID,
				Type:     "realtime_update", // Generic type, could be more specific
				Payload:  payload,
			}
		}
	}
}
