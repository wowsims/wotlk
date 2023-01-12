package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var SliceAndDiceActionID = core.ActionID{SpellID: 6774}

const SliceAndDiceEnergyCost = 25.0

func (rogue *Rogue) makeSliceAndDice(comboPoints int32) *core.Spell {
	actionID := SliceAndDiceActionID
	actionID.Tag = comboPoints
	duration := rogue.sliceAndDiceDurations[comboPoints]

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagFinisher,

		EnergyCost: core.EnergyCostOptions{
			Cost: SliceAndDiceEnergyCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			rogue.SliceAndDiceAura.Duration = duration
			rogue.SliceAndDiceAura.Activate(sim)
			rogue.ApplyFinisher(sim, spell)
		},
	})
}

func (rogue *Rogue) registerSliceAndDice() {
	durationMultiplier := 1.0 + 0.25*float64(rogue.Talents.ImprovedSliceAndDice)
	durationBonus := time.Duration(0)
	if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfSliceAndDice) {
		durationBonus += time.Second * 3
	}
	rogue.sliceAndDiceDurations = [6]time.Duration{
		0,
		time.Duration(float64(time.Second*9+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*12+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*15+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*18+durationBonus) * durationMultiplier),
		time.Duration(float64(time.Second*21+durationBonus) * durationMultiplier),
	}

	hasteBonus := 1.4
	if rogue.HasSetBonus(ItemSetSlayers, 2) {
		hasteBonus += 0.05
	}
	inverseHasteBonus := 1.0 / hasteBonus

	rogue.SliceAndDiceAura = rogue.RegisterAura(core.Aura{
		Label:    "Slice and Dice",
		ActionID: SliceAndDiceActionID,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, hasteBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MultiplyMeleeSpeed(sim, inverseHasteBonus)
		},
	})

	rogue.SliceAndDice = [6]*core.Spell{
		nil,
		rogue.makeSliceAndDice(1),
		rogue.makeSliceAndDice(2),
		rogue.makeSliceAndDice(3),
		rogue.makeSliceAndDice(4),
		rogue.makeSliceAndDice(5),
	}
}
