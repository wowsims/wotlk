package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (mage *Mage) registerArcaneMissilesSpell() {
	spellCoeff := 1/3.5 + 0.03*float64(mage.Talents.ArcaneEmpowerment)
	t10ProcAura := mage.BloodmagesRegalia2pcAura()
	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)

	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42846},
		SpellSchool:  core.SpellSchoolArcane,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | core.SpellFlagChanneled,
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
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if mage.MissileBarrageAura.IsActive() {
					if t10ProcAura != nil {
						t10ProcAura.Activate(sim)
					}
				}
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
					dot.Spell.DealDamage(sim, result)
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
