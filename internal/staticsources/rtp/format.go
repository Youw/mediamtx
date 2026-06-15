package rtp

import (
	"time"

	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/bluenviron/gortsplib/v5/pkg/rtpreceiver"
	"github.com/pion/rtcp"
)

type rtpFormat struct {
	desc format.Format

	rtpReceiver        *rtpreceiver.Receiver
	reorderInitialized bool
}

// initialize sets up the RTP receiver.
// reorderQueueSize controls the reorder buffer:
//   - nil: enable reordering with the receiver's default buffer size (64 packets).
//   - explicit zero: disable reordering (packets are passed through in arrival
//     order, like before this option existed; lost packets are still counted).
//   - any other value: enable reordering with the given buffer size.
func (f *rtpFormat) initialize(reorderQueueSize *uint) error {
	f.rtpReceiver = &rtpreceiver.Receiver{
		ClockRate:            f.desc.ClockRate(),
		UnrealiableTransport: true,
		Period:               10 * time.Second,
		WritePacketRTCP: func(_ rtcp.Packet) {
		},
	}

	if reorderQueueSize != nil && *reorderQueueSize == 0 {
		return nil
	}

	if reorderQueueSize != nil {
		f.rtpReceiver.BufferSize = int(*reorderQueueSize)
	}

	err := f.rtpReceiver.Initialize()
	if err != nil {
		return err
	}

	f.reorderInitialized = true

	return nil
}

func (f *rtpFormat) close() {
	if f.reorderInitialized {
		f.rtpReceiver.Close()
	}
}
