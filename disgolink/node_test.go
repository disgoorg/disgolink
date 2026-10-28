package disgolink

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func newTestWebSocketServer(t *testing.T, onRequest func(rq *http.Request)) *httptest.Server {
	t.Helper()

	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		if onRequest != nil {
			onRequest(rq)
		}

		conn, err := upgrader.Upgrade(rw, rq, nil)
		if err != nil {
			return
		}
		defer func() {
			_ = conn.Close()
		}()

		for {
			if _, _, err = conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	t.Cleanup(srv.Close)

	return srv
}

func newTestNode(t *testing.T, srv *httptest.Server, sessionID string) *Node {
	t.Helper()

	client := New(1)
	return newNode(client.logger, NodeConfig{
		Name:      "test",
		Address:   srv.Listener.Addr().String(),
		SessionID: sessionID,
	}, client, srv.Client())
}

func TestReconnectDelay(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: -1, want: time.Second},
		{attempt: 0, want: time.Second},
		{attempt: 1, want: 2 * time.Second},
		{attempt: 2, want: 4 * time.Second},
		{attempt: 3, want: 8 * time.Second},
		{attempt: 4, want: 16 * time.Second},
		{attempt: 5, want: 32 * time.Second},
		{attempt: 6, want: maximumConnectDelay},
		{attempt: 7, want: maximumConnectDelay},
		{attempt: 100, want: maximumConnectDelay},
	}

	for _, test := range tests {
		if got := reconnectDelay(test.attempt); got != test.want {
			t.Errorf("reconnectDelay(%d) = %v, want %v", test.attempt, got, test.want)
		}
	}
}

func TestNodeDoReconnectAttemptsImmediately(t *testing.T) {
	srv := newTestWebSocketServer(t, nil)
	node := newTestNode(t, srv, "")
	defer node.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 900*time.Millisecond)
	defer cancel()

	if err := node.Open(ctx); err != nil {
		t.Fatalf("Open() error = %v, want an immediate connection", err)
	}
}

func TestNodeOpenSendsSessionIDHeader(t *testing.T) {
	headers := make(chan string, 1)
	srv := newTestWebSocketServer(t, func(rq *http.Request) {
		headers <- rq.Header.Get("Session-Id")
	})
	node := newTestNode(t, srv, "resume-session")
	defer node.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := node.Open(ctx); err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	select {
	case header := <-headers:
		if header != "resume-session" {
			t.Errorf("Session-Id header = %q, want %q", header, "resume-session")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for the Session-Id header")
	}
}
