package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FrostStrikeActionID = core.ActionID{SpellID: 55268}

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	baseDamage := 250.0 + dk.sigilOfTheVengefulHeartFrostStrike()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, baseDamage, 0.55, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, baseDamage, 0.55*dk.nervesOfColdSteelBonus(), true)
	}

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.darkrunedBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: dk.bloodOfTheNorthCoeff(),
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) *
					dk.glacielRotBonus(hitEffect.Target) *
					dk.RoRTSBonus(hitEffect.Target) *
					dk.mercilessCombatBonus(sim)
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: onhit,
	}

	dk.threatOfThassarianProcMasks(isMH, &effect, true, false, dk.killingMachineOutcomeMod)

	conf := core.SpellConfig{
		ActionID:     FrostStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolFrost,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}

	rs := &RuneSpell{}
	if isMH {
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(
			core.Ternary(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfFrostStrike), uint8(32), 40), 0, 0, 0, 0,
		))
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: conf.BaseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
		}
		conf.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	}

	return dk.RegisterSpell(rs, conf)
}

func (dk *Deathknight) registerFrostStrikeSpell() {
	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		dk.LastOutcome = spellEffect.Outcome
		dk.threatOfThassarianProc(sim, spellEffect, dk.FrostStrikeMhHit, dk.FrostStrikeOhHit)

		// KM Consume after OH
		if spellEffect.Landed() {
			if dk.KillingMachineAura.IsActive() {
				dk.KillingMachineAura.Deactivate(sim)
			}
		}
	})
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false, nil)
	dk.FrostStrike = dk.FrostStrikeMhHit
}

func (dk *Deathknight) CanFrostStrike(sim *core.Simulation) bool {
	runeCost := core.RuneCost(dk.FrostStrike.BaseCost)
	return dk.CastCostPossible(sim, float64(runeCost.RunicPower()), 0, 0, 0) && dk.FrostStrike.IsReady(sim)
}

func (dk *Deathknight) CastFrostStrike(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanFrostStrike(sim) {
		return dk.FrostStrike.Cast(sim, target)
	}
	return false
}
