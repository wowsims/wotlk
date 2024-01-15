package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) registerEvocation() {
	actionID := core.ActionID{SpellID: 12051}
	maxTicks := core.TernaryInt32(mage.HasSetBonus(ItemSetTempestRegalia, 2), 5, 4)
	manaMetrics := mage.NewManaMetrics(actionID)
	manaPerTick := 0.0

	evocation := mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful | core.SpellFlagChanneled | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * time.Duration(4-mage.Talents.ArcaneFlows),
			},
		},

		Dot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Evocation",
			},
			NumberOfTicks:       maxTicks,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.AddMana(sim, manaPerTick, manaMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			manaPerTick = mage.MaxMana() * 0.15
			spell.SelfHot().Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocation,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if character.HasActiveAuraWithTag(core.InnervateAuraTag) || character.HasActiveAuraWithTag(core.ManaTideTotemAuraTag) {
				return false
			}

			if sim.GetRemainingDuration() < 12*time.Second {
				return false
			}

			return character.CurrentManaPercent() < 0.1
		},
	})
}
