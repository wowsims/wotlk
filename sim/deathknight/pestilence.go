package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) registerPestilenceSpell() {

	hasGlyphOfDisease := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)

	dk.Pestilence = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50842},
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncAOEDamage(dk.Env, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     0.0,
			ThreatMultiplier:     0.0,

			// Zero damage spell with a Hit mechanic, thanks blizz!
			BaseDamage:     core.BaseDamageConfigFlat(0),
			OutcomeApplier: dk.OutcomeFuncMagicHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == dk.CurrentTarget {
					dk.LastCastOutcome = spellEffect.Outcome
				}
				if spellEffect.Landed() {
					unitHit := spellEffect.Target
					// Main target
					if unitHit == dk.CurrentTarget {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.FrostFeverDisease[unitHit.Index].IsActive() {
								dk.FrostFeverDisease[unitHit.Index].Rollover(sim)
								dk.FrostFeverDebuffAura[unitHit.Index].Activate(sim)
								if dk.IcyTalonsAura != nil {
									dk.IcyTalonsAura.Activate(sim)
								}
							}

							if dk.BloodPlagueDisease[unitHit.Index].IsActive() {
								dk.BloodPlagueDisease[unitHit.Index].Rollover(sim)
							}
						}

						dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_B)
						if !dk.bloodOfTheNorthProc(sim, spell, dkSpellCost) {
							if !dk.reapingProc(sim, spell, dkSpellCost) {
								dk.Spend(sim, spell, dkSpellCost)
							}
						}

						amountOfRunicPower := 10.0 + 2.5*float64(dk.Talents.ChillOfTheGrave)
						dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
					} else {
						// Apply diseases on every other target
						if dk.TargetHasDisease(FrostFeverAuraLabel, dk.CurrentTarget) {
							dk.FrostFeverDisease[unitHit.Index].Apply(sim)
						}
						if dk.TargetHasDisease(FrostFeverAuraLabel, dk.CurrentTarget) {
							dk.BloodPlagueDisease[unitHit.Index].Apply(sim)
						}
					}
				}
			},
		}),
	})
}

func (dk *Deathknight) CanPestilence(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.Pestilence.IsReady(sim)
}

func (dk *Deathknight) CastPestilence(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanPestilence(sim) {
		dk.Pestilence.Cast(sim, target)
		return true
	}
	return false
}
