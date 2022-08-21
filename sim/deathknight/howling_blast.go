package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var HowlingBlastActionID = core.ActionID{SpellID: 51411}

func (dk *Deathknight) registerHowlingBlastSpell() {
	if !dk.Talents.HowlingBlast {
		return
	}

	rpBonus := 2.5 * float64(dk.Talents.ChillOfTheGrave)
	baseCost := float64(core.NewRuneCost(15, 0, 1, 1, 0))

	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfHowlingBlast)

	howlingBlast := &RuneSpell{}
	dk.HowlingBlast = dk.RegisterSpell(howlingBlast, core.SpellConfig{
		ActionID:     HowlingBlastActionID,
		SpellSchool:  core.SpellSchoolFrost,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
				if dk.RimeAura.IsActive() {
					cast.Cost = 0 // no runes, no regen
					dk.RimeAura.Deactivate(sim)
				}
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 8.0 * time.Second,
			},
		},

		ApplyEffects: dk.withRuneRefund(howlingBlast, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     1.0,
			ThreatMultiplier:     1.0,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					roll := (562.0-518.0)*sim.RandomFloat("Howling Blast") + 518.0
					return (roll + dk.getImpurityBonus(hitEffect, spell.Unit)*0.2) *
						dk.glacielRotBonus(hitEffect.Target) *
						dk.RoRTSBonus(hitEffect.Target) *
						dk.mercilessCombatBonus(sim)
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: dk.killingMachineOutcomeMod(dk.OutcomeFuncMagicHitAndCrit(dk.spellCritMultiplierGoGandMoM())),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == dk.CurrentTarget {
					dk.LastOutcome = spellEffect.Outcome
				}
				if dk.Talents.ChillOfTheGrave > 0 && spellEffect.Outcome.Matches(core.OutcomeLanded) {
					dk.AddRunicPower(sim, rpBonus, spell.RunicPowerMetrics())
				}

				// KM Consume after OH
				if spellEffect.Landed() && dk.KillingMachineAura.IsActive() {
					dk.KillingMachineAura.Deactivate(sim)
				}

				if hasGlyph {
					dk.FrostFeverSpell.Cast(sim, spellEffect.Target)
					if dk.Talents.CryptFever > 0 {
						dk.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
					}
					if dk.Talents.EbonPlaguebringer > 0 {
						dk.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
					}
				}
			},
		}, true),
	}, func(sim *core.Simulation) bool {
		if dk.RimeAura.IsActive() {
			return dk.HowlingBlast.IsReady(sim)
		}
		return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.HowlingBlast.IsReady(sim)
	}, nil)
}
