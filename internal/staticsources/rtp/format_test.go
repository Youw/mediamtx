package rtp

import (
	"testing"

	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/stretchr/testify/require"
)

func TestFormatReorderQueueSize(t *testing.T) {
	uintPtr := func(v uint) *uint { return &v }

	for _, ca := range []struct {
		name        string
		queueSize   *uint
		initialized bool
		bufferSize  int
	}{
		{"default", nil, true, 64},
		{"disabled", uintPtr(0), false, 0},
		{"explicit", uintPtr(128), true, 128},
	} {
		t.Run(ca.name, func(t *testing.T) {
			f := &rtpFormat{
				desc: &format.H264{},
			}

			err := f.initialize(ca.queueSize)
			require.NoError(t, err)
			defer f.close()

			require.Equal(t, ca.initialized, f.reorderInitialized)
			require.Equal(t, ca.bufferSize, f.rtpReceiver.BufferSize)
		})
	}
}
