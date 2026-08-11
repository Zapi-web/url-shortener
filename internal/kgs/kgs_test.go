package kgs

import (
	"sync"
	"testing"
)

func TestKGS_New(t *testing.T) {
	tests := []struct {
		name      string
		nodeID    int64
		wantError bool
	}{
		{
			name:   "Normal ID",
			nodeID: 1,
		},
		{
			name:      "Negative NodeID",
			nodeID:    -1,
			wantError: true,
		},
		{
			name:      "Exceeds Max NodeID",
			nodeID:    1024,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.nodeID)

			if tt.wantError {
				if err == nil {
					t.Fatalf("New() expected error; got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}
		})
	}
}

func TestKGS_Generate_Concurrency(t *testing.T) {
	const goroutines = 100
	const idsPerGoroutine = 1000

	gorGen, err := New(1)
	if err != nil {
		t.Fatalf("node init failed: %v", err)
	}

	var wg sync.WaitGroup
	idChan := make(chan uint64, goroutines*idsPerGoroutine)
	startBarrier := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startBarrier

			for j := 0; j < idsPerGoroutine; j++ {
				idChan <- gorGen.Generate()
			}
		}()
	}
	close(startBarrier)

	wg.Wait()
	close(idChan)

	seen := make(map[uint64]struct{}, goroutines*idsPerGoroutine)
	for id := range idChan {
		if _, ok := seen[id]; ok {
			t.Fatalf("Generate() duplicate ID found: %d", id)
		}
		seen[id] = struct{}{}
	}
}
