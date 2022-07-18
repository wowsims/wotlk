package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (deathKnight *DeathKnight) registerPestilenceSpell() {
	deathKnight.Pestilence = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50842},
		SpellSchool: core.SpellSchoolShadow,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: 0.0,
			DamageMultiplier:     0.0,
			ThreatMultiplier:     0.0,

			// Zero damage spell with a Hit mechanic, thanks blizz!
			BaseDamage:     core.BaseDamageConfigFlat(0),
			OutcomeApplier: deathKnight.OutcomeFuncMagicHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				// Spread diseases outside the hit check...
				for _, encounterTarget := range deathKnight.Env.Encounter.Targets {
					unit := &encounterTarget.Unit

					// Only refresh diseases duration on main target if glyphed
					if unit == spellEffect.Target {
						if deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDisease) {
							// Update expire instead of Apply to keep old snapshotted value
							deathKnight.FrostFeverDisease[unit.Index].UpdateExpires(sim.CurrentTime + deathKnight.FrostFeverDisease[unit.Index].Duration)
							deathKnight.BloodPlagueDisease[unit.Index].UpdateExpires(sim.CurrentTime + deathKnight.BloodPlagueDisease[unit.Index].Duration)
						}

						continue
					}

					// Apply diseases on every other target
					// TODO: Snapshot the current values of main target disease and roll over to off targets
					if deathKnight.TargetHasDisease(FrostFeverAuraLabel, spellEffect.Target) {
						deathKnight.FrostFeverDisease[unit.Index].Apply(sim)
					}
					if deathKnight.TargetHasDisease(FrostFeverAuraLabel, spellEffect.Target) {
						deathKnight.BloodPlagueDisease[unit.Index].Apply(sim)
					}
				}

				// Do the hit check for spending the runes and gaining RP
				if spellEffect.Landed() {
					dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 0, 0)
					deathKnight.Spend(sim, spell, dkSpellCost)

					amountOfRunicPower := 10.0 + 2.5*float64(deathKnight.Talents.ChillOfTheGrave)
					deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())
				}
			},
		}),
	})
}

func (deathKnight *DeathKnight) CanPestilence(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 0, 0) && deathKnight.Pestilence.IsReady(sim)
}
