package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (mage *Mage) registerArcaneMissilesSpell() {
	spellCoeff := 1/3.5 + 0.03*float64(mage.Talents.ArcaneEmpowerment)
	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)

	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42846},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | core.SpellFlagChanneled | core.SpellFlagAPL,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.31,
			Multiplier: 1 - .01*float64(mage.Talents.ArcaneFocus),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 5,
			},
		},

		BonusHitRating:   float64(mage.Talents.ArcaneFocus) * core.SpellHitRatingPerHitChance,
		BonusCritRating:  core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 + .04*float64(mage.Talents.TormentTheWeak),
		DamageMultiplierAdditive: 1 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage+core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneMissiles), .25, 0)),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissiles",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.MissileBarrageAura.IsActive() {
						if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
							mage.MissileBarrageAura.Deactivate(sim)
						}
					}

					// TODO: This check is necessary to ensure the final tick occurs before
					// Arcane Blast stacks are dropped. To fix this, ticks need to reliably
					// occur before aura expirations.
					dot := mage.ArcaneMissiles.Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickOnce(sim)
					}
					mage.ArcaneBlastAura.Deactivate(sim)
				},
			},

			NumberOfTicks:       5,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := 362 + spellCoeff*dot.Spell.SpellPower()
				result := dot.Spell.CalcDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				dot.Spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					// TODO: THIS IS A HACK TRY TO FIGURE OUT A BETTER WAY TO DO THIS.
					// Arcane Missiles is like mind flay in that its dmg ticks can proc things like a normal cast would.
					// However, ticks do not proc JoW. Since the dmg portion and the initial application are the same Spell
					//  we can't set one without impacting the other.
					// For now as a hack, set proc mask to prevent JoW, cast the tick dmg, and then unset it.
					// This also handles trinkets that can proc from proc (or not)
					oldMask := dot.Spell.ProcMask
					dot.Spell.ProcMask = core.ProcMaskProc
					dot.Spell.DealDamage(sim, result)
					dot.Spell.ProcMask = oldMask
				})
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
