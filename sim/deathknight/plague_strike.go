package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var PlagueStrikeActionID = core.ActionID{SpellID: 49921}

func (dk *Deathknight) newPlagueStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 378.0, 1.0, 0.5, true)
	if !isMH {
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 378.0, dk.nervesOfColdSteelBonus(), 0.5, true)
	}

	outbreakBonus := 1.0 + 0.1*float64(dk.Talents.Outbreak)
	glyphDmgBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfPlagueStrike), 1.2, 1.0)

	effect := core.SpellEffect{
		BonusCritRating:  (dk.annihilationCritBonus() + dk.scourgebornePlateCritBonus() + dk.viciousStrikesCritChanceBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: outbreakBonus,
		ThreatMultiplier: 1,

		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell) * dk.RoRTSBonus(hitEffect.Target) * glyphDmgBonus
			},
			TargetSpellCoefficient: 1,
		},

		OnSpellHitDealt: onhit,
	}

	dk.threatOfThassarianProcMasks(isMH, &effect, false, false, func(outcomeApplier core.OutcomeApplier) core.OutcomeApplier {
		return outcomeApplier
	})

	conf := core.SpellConfig{
		ActionID:     PlagueStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool:  core.SpellSchoolPhysical,
		Flags:        core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}
	rs := &RuneSpell{}
	if isMH { // only MH has cost & gcd
		rpGen := 10.0 + 2.5*float64(dk.Talents.Dirge)
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(uint8(rpGen), 0, 0, 1, 0))
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				Cost: conf.BaseCost,
				GCD:  core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			IgnoreHaste: true,
		}
		conf.ApplyEffects = dk.withRuneRefund(rs, effect, false)
	}

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.PlagueStrike.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerPlagueStrikeSpell() {
	dk.PlagueStrikeMhHit = dk.newPlagueStrikeSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if dk.Talents.ThreatOfThassarian > 0 && dk.GetOHWeapon() != nil && dk.threatOfThassarianWillProc(sim) {
			dk.PlagueStrikeOhHit.Cast(sim, spellEffect.Target)
		}
		dk.LastOutcome = spellEffect.Outcome
		if spellEffect.Landed() {
			dk.BloodPlagueExtended[spellEffect.Target.Index] = 0
			dk.BloodPlagueSpell.Cast(sim, spellEffect.Target)
			if dk.Talents.CryptFever > 0 {
				dk.CryptFeverAura[spellEffect.Target.Index].Activate(sim)
			}
			if dk.Talents.EbonPlaguebringer > 0 {
				dk.EbonPlagueAura[spellEffect.Target.Index].Activate(sim)
			}
		}
	})
	dk.PlagueStrikeOhHit = dk.newPlagueStrikeSpell(false, nil)
	dk.PlagueStrike = dk.PlagueStrikeMhHit
}
func (dk *Deathknight) registerDrwPlagueStrikeSpell() {
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 378.0, 1.0, 0.5, true)
	outbreakBonus := 1.0 + 0.1*float64(dk.Talents.Outbreak)

	dk.RuneWeapon.PlagueStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  (dk.annihilationCritBonus() + dk.scourgebornePlateCritBonus() + dk.viciousStrikesCritChanceBonus()) * core.CritRatingPerCritChance,
			DamageMultiplier: outbreakBonus,
			ThreatMultiplier: 1,
			OutcomeApplier:   dk.RuneWeapon.OutcomeFuncMeleeWeaponSpecialHitAndCrit(dk.RuneWeapon.MeleeCritMultiplier(1.0, 0.0)),
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return weaponBaseDamage(sim, hitEffect, spell)
				},
				TargetSpellCoefficient: 1,
			},
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					dk.RuneWeapon.BloodPlagueSpell.Cast(sim, spellEffect.Target)
				}
			},
		}),
	})
}
