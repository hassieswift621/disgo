package gateway

import (
	"context"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"nhooyr.io/websocket"
)

// Gateway URL
const gatewayURL = "wss://gateway.discord.gg/?v=6&encoding=json"

// Session is a single connection to the Discord Gateway.
type Session struct {
	disconnect     chan bool
	done           chan struct{}
	heartbeatState *heartbeatState
	log            *zap.Logger
	sequence       uint64
	sessionID      string
	shardCount     uint
	shardID        uint
	token          string
	wg             sync.WaitGroup
	ws             *websocket.Conn
}

// NewSession creates a new session.
func NewSession(logger *zap.Logger, shardCount uint, shardID uint, token string) *Session {
	return &Session{
		log:        logger,
		shardCount: shardCount,
		shardID:    shardID,
		token:      token,
	}
}

// Open opens the session.
func (s *Session) Open(ctx context.Context) error {
	// Create channels.
	s.disconnect = make(chan bool)
	s.done = make(chan struct{})

	// Initialise websocket headers.
	headers := http.Header{}
	headers.Add("Accept-Encoding", "zlib")

	// Attempt to handshake.
	s.log.Debug("connecting to the gateway", zap.Uint("shard", s.shardID))
	ws, _, err := websocket.Dial(ctx, gatewayURL, websocket.DialOptions{HTTPHeader: headers})
	if err != nil {
		close(s.done)
		return xerrors.Errorf("failed to connect to the gateway: %w", err)
	}

	// Store websocket.
	s.ws = ws

	// Set message limit to 1GB.
	s.ws.SetReadLimit(1073741824)

	// Read hello payload.
	s.log.Debug("receiving hello payload", zap.Uint("shard", s.shardID))
	err = s.hello(ctx)
	if err != nil {
		close(s.done)
		return xerrors.Errorf("failed to read hello payload: %w", err)
	}

	// Start heartbeating.
	s.wg.Add(1)
	go s.heartbeat(ctx)

	// If we have a sequence number or a session ID, attempt to resume,
	// otherwise create a new session.
	if s.sequence != 0 || s.sessionID != "" {
		// Resume.
		s.log.Debug("resuming session")
		if err = s.resume(ctx); err != nil {
			s.log.Warn("failed to resume session", zap.Uint("shard", s.shardID))
			return err
		}
	} else {
		// Identify.
		if err = s.identify(ctx); err != nil {
			return xerrors.Errorf("failed to identify: %w", err)
		}
	}

	// Connected, start handling connection and events.
	go s.handleConnection(ctx)

	return nil
}

// HandleConnection handles the connection by listening to the signals.
func (s *Session) handleConnection(ctx context.Context) {
	// Var to store whether to reconnect after disconnecting.
	var reconnect bool

	select {
	case <-s.done:
		// Received disconnect signal from external source.
		// Stop handling connection.
		return
	case reconnect = <-s.disconnect:
	case <-ctx.Done():
	}

	// Signal goroutines to stop.
	close(s.done)

	// Wait for goroutines to stop.
	s.wg.Wait()

	// Close session.
	err := s.ws.Close(websocket.StatusInternalError, "")
	if err != nil {
		s.log.Info("failed to close websocket", zap.Uint("shard", s.shardID), zap.Error(err))
		return
	}

	if reconnect {
		// Reopen session.
		err := s.Open(ctx)
		if err != nil {
			s.log.Info("failed to reconnect websocket", zap.Uint("shard", s.shardID), zap.Error(err))
		}
	} else {
		// Reset session state.
		s.sequence = 0
		s.sessionID = ""
	}
}
