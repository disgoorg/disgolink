package lavalink

import "testing"

func makeStats(systemLoad float64, cores int) Stats {
	return Stats{
		CPU: CPU{
			SystemLoad: systemLoad,
			Cores:      cores,
		},
	}
}

func TestStats_Better_LowerLoadWins(t *testing.T) {
	low := makeStats(0.1, 4)  // 10% load
	high := makeStats(0.8, 4) // 80% load

	if !low.Better(high) {
		t.Error("low-load node should be better than high-load node")
	}
	if high.Better(low) {
		t.Error("high-load node should not be better than low-load node")
	}
}

func TestStats_Better_EqualLoad(t *testing.T) {
	a := makeStats(0.5, 4)
	b := makeStats(0.5, 4)

	if a.Better(b) {
		t.Error("equal load nodes should not displace each other")
	}
}

func TestStats_Better_DifferentCoreCount(t *testing.T) {
	// 50% load on 8 cores = 6.25% per core
	// 50% load on 2 cores = 25% per core
	manycores := makeStats(0.5, 8)
	fewcores := makeStats(0.5, 2)

	if !manycores.Better(fewcores) {
		t.Error("same load spread across more cores should be better")
	}
}
