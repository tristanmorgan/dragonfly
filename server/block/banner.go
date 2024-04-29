package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// Banner is a tall decorative block that can be customized.
type Banner struct {
	empty
	transparent

	// Colour is the colour of the banner.
	Colour item.Colour
	// Attach is the attachment of the Banner. It is either of the type WallAttachment or StandingAttachment.
	Attach Attachment
	// Patterns represents the patterns the Banner should show when rendering.
	Patterns []BannerPatternLayer
	// Illager returns true if the banner is an illager banner.
	Illager bool
}

// MaxCount ...
func (Banner) MaxCount() int {
	return 16
}

// PistonBreakable ...
func (Banner) PistonBreakable() bool {
	return true
}

// BreakInfo ...
func (b Banner) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, axeEffective, oneOf(b))
}

// FuelInfo ...
func (Banner) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// UseOnBlock ...
func (b Banner) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(w, pos, face, b)
	if !used || face == cube.FaceDown {
		return false
	}

	if face == cube.FaceUp {
		b.Attach = StandingAttachment(user.Rotation().Orientation().Opposite())
		place(w, pos, b, user, ctx)
		return
	}
	b.Attach = WallAttachment(face.Direction())
	place(w, pos, b, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (b Banner) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if b.Attach.hanging {
		if _, ok := w.Block(pos.Side(b.Attach.facing.Opposite().Face())).(Air); ok {
			w.SetBlock(pos, nil, nil)
			w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: b})
			dropItem(w, item.NewStack(b, 1), pos.Vec3Centre())
		}
		return
	}
	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Air); ok {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: b})
		dropItem(w, item.NewStack(b, 1), pos.Vec3Centre())
	}
}

// EncodeItem ...
func (b Banner) EncodeItem() (name string, meta int16) {
	return "minecraft:banner", invertColour(b.Colour)
}

// EncodeBlock ...
func (b Banner) EncodeBlock() (name string, properties map[string]any) {
	if b.Attach.hanging {
		return "minecraft:wall_banner", map[string]any{"facing_direction": int32(b.Attach.facing + 2)}
	}
	return "minecraft:standing_banner", map[string]any{"ground_sign_direction": int32(b.Attach.o)}
}

// EncodeNBT ...
func (b Banner) EncodeNBT() map[string]any {
	patterns := make([]any, 0, len(b.Patterns))
	for _, p := range b.Patterns {
		patterns = append(patterns, p.EncodeNBT())
	}
	return map[string]any{
		"id":       "Banner",
		"Patterns": patterns,
		"Type":     int32(boolByte(b.Illager)),
		"Base":     int32(invertColour(b.Colour)),
	}
}

// DecodeNBT ...
func (b Banner) DecodeNBT(m map[string]any) any {
	b.Colour = invertColourID(int16(nbtconv.Int32(m, "Base")))
	b.Illager = nbtconv.Int32(m, "Type") == 1
	if patterns := nbtconv.Slice(m, "Patterns"); patterns != nil {
		b.Patterns = make([]BannerPatternLayer, len(patterns))
		for i, p := range b.Patterns {
			b.Patterns[i] = p.DecodeNBT(patterns[i].(map[string]any)).(BannerPatternLayer)
		}
	}
	return b
}

// invertColour converts the item.Colour passed and returns the colour ID inverted.
func invertColour(c item.Colour) int16 {
	return ^int16(c.Uint8()) & 0xf
}

// invertColourID converts the int16 passed the returns the item.Colour inverted.
func invertColourID(id int16) item.Colour {
	return item.Colours()[uint8(^id&0xf)]
}

// allBanners returns all possible banners.
func allBanners() (banners []world.Block) {
	for _, d := range cube.Directions() {
		banners = append(banners, Banner{Attach: WallAttachment(d)})
	}
	for o := cube.Orientation(0); o <= 15; o++ {
		banners = append(banners, Banner{Attach: StandingAttachment(o)})
	}
	return
}
