package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// SlimeBlock is a block that slows down entities and bounces them up if they drop onto it.
type SlimeBlock struct {
	solid
	transparent
}

// EntityLand ...
func (SlimeBlock) EntityLand(_ cube.Pos, _ *world.World, e world.Entity) {
	if fallEntity, ok := e.(fallDistanceEntity); ok {
		fallEntity.ResetFallDistance()
	}
}

// BreakInfo ...
func (s SlimeBlock) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(s))
}

// Friction ...
func (SlimeBlock) Friction() float64 {
	return 0.8
}

// EncodeItem ...
func (SlimeBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:slime", 0
}

// EncodeBlock ...
func (SlimeBlock) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:slime", nil
}
