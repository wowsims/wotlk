package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerPestilenceSpell() {
	hasGlyphOfDisease := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)
	baseCost := float64(core.NewRuneCost(10, 1, 0, 0, 0))

	rs := &RuneSpell{
		Refundable: true,
	}

	dk.Pestilence = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 50842},
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.GetModifiedGCD()
			},
		},

		DamageMultiplier: 0,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.Targets {
				aoeUnit := &aoeTarget.Unit

				// Zero damage spell with a Hit mechanic, thanks blizz!
				result := spell.CalcAndDealDamage(sim, aoeUnit, 0, spell.OutcomeMagicHit)

				if aoeUnit == dk.CurrentTarget {
					rs.OnResult(sim, result)
					dk.LastOutcome = result.Outcome
				}
				if result.Landed() {
					// Main target
					if aoeUnit == dk.CurrentTarget {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.FrostFeverDisease[aoeUnit.Index].IsActive() {
								dk.FrostFeverDisease[aoeUnit.Index].Rollover(sim)
								if dk.Talents.IcyTalons > 0 {
									dk.IcyTalonsAura.Activate(sim)
								}
								dk.FrostFeverDebuffAura[aoeUnit.Index].Activate(sim)
							}

							if dk.BloodPlagueDisease[aoeUnit.Index].IsActive() {
								dk.BloodPlagueDisease[aoeUnit.Index].Rollover(sim)
							}
						}
					} else {
						applyCryptEbon := false
						// Apply diseases on every other target
						if dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() {
							dk.FrostFeverExtended[aoeUnit.Index] = 0
							dk.FrostFeverDisease[aoeUnit.Index].Apply(sim)
							applyCryptEbon = true
						}
						if dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() {
							dk.BloodPlagueExtended[aoeUnit.Index] = 0
							dk.BloodPlagueDisease[aoeUnit.Index].Apply(sim)
							applyCryptEbon = true
						}
						if applyCryptEbon {
							if dk.Talents.CryptFever > 0 {
								dk.CryptFeverAura[aoeUnit.Index].Activate(sim)
							}
							if dk.Talents.EbonPlaguebringer > 0 {
								dk.EbonPlagueAura[aoeUnit.Index].Activate(sim)
							}
						}
					}
				}
			}
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 1, 0, 0) && dk.Pestilence.IsReady(sim)
	}, nil)
	if dk.Talents.BloodOfTheNorth+dk.Talents.Reaping >= 3 {
		rs.DeathConvertChance = 1.0
	} else {
		rs.DeathConvertChance = float64(dk.Talents.BloodOfTheNorth+dk.Talents.Reaping) * 0.33
	}
	rs.ConvertType = RuneTypeBlood
}
