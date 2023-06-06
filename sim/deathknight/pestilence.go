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
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				// Zero damage spell with a Hit mechanic, thanks blizz!
				result := spell.CalcAndDealDamage(sim, aoeTarget, 0, spell.OutcomeMagicHit)

				if aoeTarget == target {
					spell.SpendRefundableCostAndConvertBloodRune(sim, result, deathConvertChance)
					dk.LastOutcome = result.Outcome
				}
				if result.Landed() {
					// Main target
					if aoeTarget == target {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.FrostFeverSpell.Dot(aoeTarget).IsActive() {
								dk.FrostFeverSpell.Dot(aoeTarget).Rollover(sim)
								if dk.Talents.IcyTalons > 0 {
									dk.IcyTalonsAura.Activate(sim)
								}
								dk.FrostFeverDebuffAura[aoeTarget.Index].Activate(sim)
							}

							if dk.BloodPlagueSpell.Dot(aoeTarget).IsActive() {
								dk.BloodPlagueSpell.Dot(aoeTarget).Rollover(sim)
							}
						}
					} else {
						// Apply diseases on every other target
						if dk.FrostFeverSpell.Dot(target).IsActive() {
							dk.FrostFeverExtended[aoeTarget.Index] = 0
							dk.FrostFeverSpell.Dot(aoeTarget).Apply(sim)
						}
						if dk.BloodPlagueSpell.Dot(target).IsActive() {
							dk.BloodPlagueExtended[aoeTarget.Index] = 0
							dk.BloodPlagueSpell.Dot(aoeTarget).Apply(sim)
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
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				// Zero damage spell with a Hit mechanic, thanks blizz!
				result := spell.CalcAndDealDamage(sim, aoeTarget, 0, spell.OutcomeMagicHit)

				if result.Landed() {
					// Main target
					if aoeTarget == target {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.RuneWeapon.FrostFeverSpell.Dot(aoeTarget).IsActive() {
								dk.RuneWeapon.FrostFeverSpell.Dot(aoeTarget).Rollover(sim)
							} else if shouldApplyDrwDots && dk.FrostFeverSpell.Dot(aoeTarget).IsActive() {
								dk.RuneWeapon.FrostFeverSpell.Dot(aoeTarget).Apply(sim)
							}

							if dk.RuneWeapon.BloodPlagueSpell.Dot(aoeTarget).IsActive() {
								dk.RuneWeapon.BloodPlagueSpell.Dot(aoeTarget).Rollover(sim)
							} else if shouldApplyDrwDots && dk.BloodPlagueSpell.Dot(aoeTarget).IsActive() {
								dk.RuneWeapon.BloodPlagueSpell.Dot(aoeTarget).Apply(sim)
							}
						} else if shouldApplyDrwDots {
							if dk.FrostFeverSpell.Dot(aoeTarget).IsActive() {
								dk.RuneWeapon.FrostFeverSpell.Dot(aoeTarget).Apply(sim)
							}

							if dk.BloodPlagueSpell.Dot(aoeTarget).IsActive() {
								dk.RuneWeapon.BloodPlagueSpell.Dot(aoeTarget).Apply(sim)
							}
						}
					} else {
						// Apply diseases on every other target
						if dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive() {
							dk.RuneWeapon.FrostFeverSpell.Dot(aoeTarget).Apply(sim)
						}
						if dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive() {
							dk.RuneWeapon.BloodPlagueSpell.Dot(aoeTarget).Apply(sim)
						}
					}
				}
			}
		},
	})
}
