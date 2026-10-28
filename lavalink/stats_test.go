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

func TestStats_Better(t *testing.T) {
	tests := []struct {
		name  string
		self  Stats
		other Stats
		want  bool
	}{
		{
			name:  "lower load wins",
			self:  makeStats(0.1, 4),
			other: makeStats(0.8, 4),
			want:  true,
		},
		{
			name:  "higher load loses",
			self:  makeStats(0.8, 4),
			other: makeStats(0.1, 4),
			want:  false,
		},
		{
			name:  "equal load does not displace",
			self:  makeStats(0.5, 4),
			other: makeStats(0.5, 4),
			want:  false,
		},
		{
			name:  "same load across more cores is better",
			self:  makeStats(0.5, 8),
			other: makeStats(0.5, 2),
			want:  true,
		},
		{
			name:  "unmeasured node is never better",
			self:  makeStats(0, 0),
			other: makeStats(0.5, 4),
			want:  false,
		},
		{
			name:  "measured node is better than unmeasured node",
			self:  makeStats(0.5, 4),
			other: makeStats(0, 0),
			want:  true,
		},
		{
			name:  "fractional load difference is not truncated",
			self:  makeStats(0.501, 1),
			other: makeStats(0.509, 1),
			want:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.self.Better(test.other); got != test.want {
				t.Errorf("Better() = %v, want %v", got, test.want)
			}
		})
	}
}
