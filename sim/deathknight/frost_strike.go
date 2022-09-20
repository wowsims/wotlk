package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var FrostStrikeActionID = core.ActionID{SpellID: 55268}

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	bonusBaseDamage := dk.sigilOfTheVengefulHeartFrostStrike()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 250+bonusBaseDamage, true)
	if !isMH {
		// SpellID 66962
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 125+bonusBaseDamage, true)
	}

	effect := core.SpellEffect{
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

	procMask := dk.threatOfThassarianProcMasks(isMH, dk.Talents.GuileOfGorefiend, &effect)
	effect.OutcomeApplier = dk.killingMachineOutcomeMod(effect.OutcomeApplier)

	conf := core.SpellConfig{
		ActionID:    FrostStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics,

		BonusCritRating: (dk.annihilationCritBonus() + dk.darkrunedBattlegearCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .55 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.bloodOfTheNorthCoeff(),
		ThreatMultiplier: 1,

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
			IgnoreHaste: true,
		}
		conf.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	}

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			runeCost := core.RuneCost(dk.FrostStrike.BaseCost)
			return dk.CastCostPossible(sim, float64(runeCost.RunicPower()), 0, 0, 0) && dk.FrostStrike.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerFrostStrikeSpell() {
	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		dk.LastOutcome = spellEffect.Outcome
		dk.threatOfThassarianProc(sim, spellEffect, dk.FrostStrikeOhHit)

		// KM Consume after OH
		if spellEffect.Landed() && dk.KillingMachineAura.IsActive() {
			dk.KillingMachineAura.Deactivate(sim)
		}
	})
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false, nil)
	dk.FrostStrike = dk.FrostStrikeMhHit
}
