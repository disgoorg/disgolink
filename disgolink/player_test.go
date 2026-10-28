package disgolink

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

func newUpdatePlayerTestPlayer(t *testing.T, handler http.HandlerFunc) *Player {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client := New(1)
	node := newNode(client.logger, NodeConfig{
		Name:      "test",
		Address:   srv.Listener.Addr().String(),
		SessionID: "session",
	}, client, srv.Client())

	return newPlayer(client.logger, client, node, 1)
}

func TestPlayerUpdateReadsTrackFromResponse(t *testing.T) {
	player := newUpdatePlayerTestPlayer(t, func(rw http.ResponseWriter, rq *http.Request) {
		if rq.Method != http.MethodPatch {
			t.Errorf("unexpected method %s", rq.Method)
		}
		if rq.URL.Path != "/v4/sessions/session/players/1" {
			t.Errorf("unexpected path %s", rq.URL.Path)
		}
		_, _ = rw.Write([]byte(`{"guildId":"1","track":{"encoded":"new-encoded","info":{"identifier":"new-id","author":"author","length":1000,"title":"title","sourceName":"http"}},"volume":80,"paused":true}`))
	})

	player.Track = &lavalink.Track{Encoded: "old-encoded"}

	if err := player.Update(context.Background()); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if player.Track == nil {
		t.Fatal("Track is nil, expected track from response")
	}
	if player.Track.Encoded != "new-encoded" {
		t.Errorf("Track.Encoded = %q, want %q", player.Track.Encoded, "new-encoded")
	}
	if player.Volume != 80 {
		t.Errorf("Volume = %d, want 80", player.Volume)
	}
	if !player.Paused {
		t.Error("Paused = false, want true")
	}
}

func TestPlayerUpdateClearsTrackWhenResponseHasNone(t *testing.T) {
	player := newUpdatePlayerTestPlayer(t, func(rw http.ResponseWriter, rq *http.Request) {
		_, _ = rw.Write([]byte(`{"guildId":"1","track":null,"volume":50}`))
	})

	player.Track = &lavalink.Track{Encoded: "old-encoded"}

	if err := player.Update(context.Background()); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if player.Track != nil {
		t.Errorf("Track = %+v, want nil", player.Track)
	}
}

func TestPlayerUpdateHandlesNoContentResponse(t *testing.T) {
	player := newUpdatePlayerTestPlayer(t, func(rw http.ResponseWriter, rq *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	})

	player.Track = &lavalink.Track{Encoded: "old-encoded"}
	player.Volume = 42

	if err := player.Update(context.Background()); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if player.Track == nil || player.Track.Encoded != "old-encoded" {
		t.Errorf("Track = %+v, want unchanged", player.Track)
	}
	if player.Volume != 42 {
		t.Errorf("Volume = %d, want unchanged 42", player.Volume)
	}
}
