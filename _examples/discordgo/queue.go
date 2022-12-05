package main

import "github.com/disgoorg/disgolink/v2/lavalink"

type QueueType string

const (
	QueueTypeNormal      QueueType = "normal"
	QueueTypeRepeatTrack QueueType = "repeat_track"
	QueueTypeRepeatQueue QueueType = "repeat_queue"
)

func (q QueueType) String() string {
	switch q {
	case QueueTypeNormal:
		return "Normal"
	case QueueTypeRepeatTrack:
		return "Repeat Track"
	case QueueTypeRepeatQueue:
		return "Repeat Queue"
	default:
		return "unknown"
	}
}

type Queue struct {
	Tracks []lavalink.Track
	Type   QueueType
}

func (q *Queue) Add(track ...lavalink.Track) {
	q.Tracks = append(q.Tracks, track...)
}

func (q *Queue) Next() (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	track := q.Tracks[0]
	q.Tracks = q.Tracks[1:]
	return track, true
}

func (q *Queue) Clear() {
	q.Tracks = make([]lavalink.Track, 0)
}

type QueueManager struct {
	queues map[string]*Queue
}

func (q *QueueManager) Get(guildID string) *Queue {
	queue, ok := q.queues[guildID]
	if !ok {
		queue = &Queue{
			Tracks: make([]lavalink.Track, 0),
			Type:   QueueTypeNormal,
		}
		q.queues[guildID] = queue
	}
	return queue
}

func (q *QueueManager) Delete(guildID string) {
	delete(q.queues, guildID)
}
