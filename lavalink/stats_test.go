package lavalink

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

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

    assert.True(t, low.Better(high), "low-load node should be better than high-load node")
    assert.False(t, high.Better(low), "high-load node should not be better than low-load node")
}

func TestStats_Better_EqualLoad(t *testing.T) {
    a := makeStats(0.5, 4)
    b := makeStats(0.5, 4)

    assert.False(t, a.Better(b), "equal load nodes should not displace each other")
}

func TestStats_Better_DifferentCoreCount(t *testing.T) {
    // 50% load on 8 cores = 6.25% per core
    // 50% load on 2 cores = 25% per core
    manycores := makeStats(0.5, 8)
    fewcores := makeStats(0.5, 2)

    assert.True(t, manycores.Better(fewcores), "same load spread across more cores should be better")
}

