package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerPestilenceSpell() {

	hasGlyphOfDisease := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDisease)
	baseCost := float64(core.NewRuneCost(10, 1, 0, 0, 0))
	rs := &RuneSpell{}
	dk.Pestilence = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 50842},
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		},

		ApplyEffects: dk.withRuneRefund(rs, core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     0.0,
			ThreatMultiplier:     0.0,

			// Zero damage spell with a Hit mechanic, thanks blizz!
			BaseDamage:     core.BaseDamageConfigFlat(0),
			OutcomeApplier: dk.OutcomeFuncMagicHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Target == dk.CurrentTarget {
					dk.LastOutcome = spellEffect.Outcome
				}
				if spellEffect.Landed() {
					unitHit := spellEffect.Target
					// Main target
					if unitHit == dk.CurrentTarget {
						if hasGlyphOfDisease {
							// Update expire instead of Apply to keep old snapshotted value
							if dk.FrostFeverDisease[unitHit.Index].IsActive() {
								dk.FrostFeverDisease[unitHit.Index].Rollover(sim)
								if dk.Talents.IcyTalons > 0 {
									dk.IcyTalonsAura.Activate(sim)
								}
								dk.FrostFeverDebuffAura[unitHit.Index].Activate(sim)
							}

							if dk.BloodPlagueDisease[unitHit.Index].IsActive() {
								dk.BloodPlagueDisease[unitHit.Index].Rollover(sim)
							}
						}
					} else {
						applyCryptEbon := false
						// Apply diseases on every other target
						if dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() {
							dk.FrostFeverExtended[unitHit.Index] = 0
							dk.FrostFeverDisease[unitHit.Index].Apply(sim)
							applyCryptEbon = true
						}
						if dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() {
							dk.BloodPlagueExtended[unitHit.Index] = 0
							dk.BloodPlagueDisease[unitHit.Index].Apply(sim)
							applyCryptEbon = true
						}
						if applyCryptEbon {
							if dk.Talents.CryptFever > 0 {
								dk.CryptFeverAura[unitHit.Index].Activate(sim)
							}
							if dk.Talents.EbonPlaguebringer > 0 {
								dk.EbonPlagueAura[unitHit.Index].Activate(sim)
							}
						}
					}
				}
			},
		}, true),
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
