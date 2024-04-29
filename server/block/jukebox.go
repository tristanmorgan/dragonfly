package block

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"time"
)

// Jukebox is a block used to play music discs.
type Jukebox struct {
	solid
	bass

	// Item is the music disc played by the jukebox.
	Item item.Stack
}

// Source ...
func (Jukebox) Source() bool {
	return true
}

// WeakPower ...
func (j Jukebox) WeakPower(cube.Pos, cube.Face, *world.World, bool) int {
	if j.Item.Empty() {
		return 0
	}
	return 15
}

// StrongPower ...
func (Jukebox) StrongPower(cube.Pos, cube.Face, *world.World, bool) int {
	return 0
}

// FuelInfo ...
func (Jukebox) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// BreakInfo ...
func (j Jukebox) BreakInfo() BreakInfo {
	d := []item.Stack{item.NewStack(Jukebox{}, 1)}
	if !j.Item.Empty() {
		d = append(d, j.Item)
	}
	return newBreakInfo(0.8, alwaysHarvestable, axeEffective, simpleDrops(d...)).withBreakHandler(func(pos cube.Pos, w *world.World, u item.User) {
		if _, hasDisc := j.Disc(); hasDisc {
			w.PlaySound(pos.Vec3Centre(), sound.MusicDiscEnd{})
		}
		updateAroundRedstone(pos, w)
	})
}

// jukeboxUser represents an item.User that can use a jukebox.
type jukeboxUser interface {
	item.User
	// SendJukeboxPopup sends a jukebox popup to the item.User.
	SendJukeboxPopup(a ...any)
}

// Activate ...
func (j Jukebox) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User, ctx *item.UseContext) bool {
	if _, hasDisc := j.Disc(); hasDisc {
		dropItem(w, j.Item, pos.Side(cube.FaceUp).Vec3Middle())

		j.Item = item.Stack{}
		w.SetBlock(pos, j, nil)
		w.PlaySound(pos.Vec3Centre(), sound.MusicDiscEnd{})
	} else if held, _ := u.HeldItems(); !held.Empty() {
		if m, ok := held.Item().(item.MusicDisc); ok {
			j.Item = held

			w.SetBlock(pos, j, nil)
			w.PlaySound(pos.Vec3Centre(), sound.MusicDiscEnd{})
			ctx.SubtractFromCount(1)

			w.PlaySound(pos.Vec3Centre(), sound.MusicDiscPlay{DiscType: m.DiscType})
			if u, ok := u.(jukeboxUser); ok {
				u.SendJukeboxPopup(fmt.Sprintf("Now playing: %v - %v", m.DiscType.Author(), m.DiscType.DisplayName()))
			}
		}
	}
	return true
}

// Disc returns the currently playing music disc
func (j Jukebox) Disc() (sound.DiscType, bool) {
	if !j.Item.Empty() {
		if m, ok := j.Item.Item().(item.MusicDisc); ok {
			return m.DiscType, true
		}
	}
	return sound.DiscType{}, false
}

// EncodeNBT ...
func (j Jukebox) EncodeNBT() map[string]any {
	m := map[string]any{"id": "Jukebox"}
	if _, hasDisc := j.Disc(); hasDisc {
		m["RecordItem"] = nbtconv.WriteItem(j.Item, true)
	}
	return m
}

// DecodeNBT ...
func (j Jukebox) DecodeNBT(data map[string]any) any {
	s := nbtconv.MapItem(data, "RecordItem")
	if _, ok := s.Item().(item.MusicDisc); ok {
		j.Item = s
	}
	return j
}

// EncodeItem ...
func (Jukebox) EncodeItem() (name string, meta int16) {
	return "minecraft:jukebox", 0
}

// EncodeBlock ...
func (Jukebox) EncodeBlock() (string, map[string]any) {
	return "minecraft:jukebox", nil
}
