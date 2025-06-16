package session

import (
	"fmt"
	"time"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// EmoteHandler handles the Emote packet.
type EmoteHandler struct {
	LastEmote time.Time
}

// maxEmoteLength is the maximum allowed emote duration in ticks. It is equivalent to 7.5 seconds.
const maxEmoteLength = 7.5 * 20

// Handle ...
func (h *EmoteHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.Emote)

	if pk.EntityRuntimeID != selfEntityRuntimeID {
		return errSelfRuntimeID
	}
	if time.Since(h.LastEmote) < time.Second {
		return nil
	}
	h.LastEmote = time.Now()
	emote, err := uuid.Parse(pk.EmoteID)
	if err != nil {
		return err
	}

	// TODO: Do not rely on the emote length sent by the client.
	length := pk.EmoteLength
	if length <= 0 || length > maxEmoteLength {
		return fmt.Errorf("invalid emote length")
	}

	dur := time.Duration(pk.EmoteLength) * time.Millisecond * 50
	expiry := time.Now().Add(dur)
	s.emoteExpiration.Store(&expiry)

	for _, viewer := range tx.Viewers(c.Position()) {
		viewer.ViewEmote(c, emote)
	}
	return nil
}
