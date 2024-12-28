package isotime_test

import (
	"testing"
	"time"

	"github.com/hsfzxjy/sdxtra/internal/isotime"
)

func BenchmarkEncode(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		_ = isotime.Encode(t)
	}
}

func BenchmarkEncodeCustom(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		_ = t.Format("2006-01-02T15:04:05.000Z")
	}
}

func TestEncode(t *testing.T) {
	nsec := 1
	times := []time.Time{time.Now()}
	for range 9 {
		times = append(times, time.Date(2020, 1, 1, 0, 0, 0, nsec, time.FixedZone("Asia/Tokyo", 9*60*60)))
		nsec *= 10
	}
	for _, time := range times {
		s := isotime.Encode(time)
		target := isotime.String(time.UTC().Format("2006-01-02T15:04:05.000Z"))
		if s != target {
			t.Errorf("expected %v, got %v", target, s)
		}
	}
}
