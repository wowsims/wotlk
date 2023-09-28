package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) registerBlizzardSpell() {
	blizzardTickSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 42938},
		SpellSchool:      core.SpellSchoolFrost,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            SpellFlagMage,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3),
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 426 + (4.0/3.5/8)*spell.SpellPower()
			damage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	mage.Blizzard = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42940},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagChanneled | core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.74,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 8,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Blizzard",
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				blizzardTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
