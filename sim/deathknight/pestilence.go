package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

var PestilenceActionID = core.ActionID{SpellID: 50842}

func (dk *Deathknight) registerPestilenceSpell() {
	hasGlyphOfDisease := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)
	deathConvertChance := float64(dk.Talents.BloodOfTheNorth+dk.Talents.Reaping) / 3

	dk.Pestilence = dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50842},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 0,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				aoeUnit := &aoeTarget.Unit

				// Zero damage spell with a Hit mechanic, thanks blizz!
				result := spell.CalcAndDealDamage(sim, aoeUnit, 0, spell.OutcomeMagicHit)

				if aoeUnit == target {
					spell.SpendRefundableCostAndConvertBloodRune(sim, result, deathConvertChance)
					dk.LastOutcome = result.Outcome
				}
				if result.Landed() {
					// Main target
					if aoeUnit == target {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.FrostFeverSpell.Dot(aoeUnit).IsActive() {
								dk.FrostFeverSpell.Dot(aoeUnit).Rollover(sim)
								if dk.Talents.IcyTalons > 0 {
									dk.IcyTalonsAura.Activate(sim)
								}
								dk.FrostFeverDebuffAura[aoeUnit.Index].Activate(sim)
							}

							if dk.BloodPlagueSpell.Dot(aoeUnit).IsActive() {
								dk.BloodPlagueSpell.Dot(aoeUnit).Rollover(sim)
							}
						}
					} else {
						// Apply diseases on every other target
						if dk.FrostFeverSpell.Dot(target).IsActive() {
							dk.FrostFeverExtended[aoeUnit.Index] = 0
							dk.FrostFeverSpell.Dot(aoeUnit).Apply(sim)
						}
						if dk.BloodPlagueSpell.Dot(target).IsActive() {
							dk.BloodPlagueExtended[aoeUnit.Index] = 0
							dk.BloodPlagueSpell.Dot(aoeUnit).Apply(sim)
						}
					}
				}
			}
		},
	})
}
func (dk *Deathknight) registerDrwPestilenceSpell() {
	hasGlyphOfDisease := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)
	dk.RuneWeapon.Pestilence = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    PestilenceActionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 0,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// DRW and Pestilence have a weird interaction where the drws Dots can be applied
			// with the spread effect from pestilence if the target has the Dks dots up but it
			// only works if there is a valid target for spread mechanic to happen (2+ mobs)
			shouldApplyDrwDots := dk.Env.GetNumTargets() > 1 || dk.Inputs.DrwPestiApply
			for _, aoeTarget := range sim.Encounter.Targets {
				aoeUnit := &aoeTarget.Unit

				// Zero damage spell with a Hit mechanic, thanks blizz!
				result := spell.CalcAndDealDamage(sim, aoeUnit, 0, spell.OutcomeMagicHit)

				if result.Landed() {
					// Main target
					if aoeUnit == dk.CurrentTarget {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.FrostFeverSpell.Dot(aoeUnit).IsActive() {
								if dk.RuneWeapon.FrostFeverSpell.Dot(aoeUnit).IsActive() {
									dk.RuneWeapon.FrostFeverSpell.Dot(aoeUnit).Rollover(sim)
								} else if shouldApplyDrwDots {
									dk.RuneWeapon.FrostFeverSpell.Dot(aoeUnit).Apply(sim)
								}
							}

							if dk.BloodPlagueSpell.Dot(aoeUnit).IsActive() {
								if dk.RuneWeapon.BloodPlagueSpell.Dot(aoeUnit).IsActive() {
									dk.RuneWeapon.BloodPlagueSpell.Dot(aoeUnit).Rollover(sim)
								} else if shouldApplyDrwDots {
									dk.RuneWeapon.BloodPlagueSpell.Dot(aoeUnit).Apply(sim)
								}
							}
						} else if shouldApplyDrwDots {
							if dk.FrostFeverSpell.Dot(aoeUnit).IsActive() {
								dk.RuneWeapon.FrostFeverSpell.Dot(aoeUnit).Apply(sim)
							}

							if dk.BloodPlagueSpell.Dot(aoeUnit).IsActive() {
								dk.RuneWeapon.BloodPlagueSpell.Dot(aoeUnit).Apply(sim)
							}
						}
					} else {
						// Apply diseases on every other target
						if dk.FrostFeverSpell.Dot(dk.CurrentTarget).IsActive() {
							dk.RuneWeapon.FrostFeverSpell.Dot(aoeUnit).Apply(sim)
						}
						if dk.BloodPlagueSpell.Dot(dk.CurrentTarget).IsActive() {
							dk.RuneWeapon.BloodPlagueSpell.Dot(aoeUnit).Apply(sim)
						}
					}
				}
			}
		},
	})
}
