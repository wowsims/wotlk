package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const FrostFeverAuraLabel = "FrostFever-"
const BloodPlagueAuraLabel = "BloodPlague-"

func (dk *Deathknight) countActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.FrostFeverDisease[target.Index].IsActive() {
		count++
	}
	if dk.BloodPlagueDisease[target.Index].IsActive() {
		count++
	}
	if dk.Talents.CryptFever > 0 && dk.CryptFeverAura[target.Index].IsActive() {
		count++
	} else if dk.Talents.EbonPlaguebringer > 0 && dk.EbonPlagueAura[target.Index].IsActive() {
		count++
	}
	return float64(count)
}

func (dk *Deathknight) diseaseMultiplier(multiplier float64, hasDarkrune4P bool) float64 {
	if hasDarkrune4P {
		return multiplier * 1.2
	}
	return multiplier
}

func (dk *Deathknight) registerDiseaseDots() {
	dk.registerFrostFever()
	dk.registerBloodPlague()
}

func (dk *Deathknight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}

	flagTs := make([]bool, dk.Env.GetNumTargets())
	isRefreshing := make([]bool, dk.Env.GetNumTargets())

	dk.FrostFeverSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			if dk.FrostFeverDisease[unit.Index].IsActive() {
				isRefreshing[unit.Index] = true
			}
			dk.FrostFeverDisease[unit.Index].Apply(sim)
			isRefreshing[unit.Index] = false

			dk.FrostFeverDebuffAura[unit.Index].Activate(sim)

			if dk.IcyTalonsAura != nil {
				dk.IcyTalonsAura.Activate(sim)
			}
		},
	})

	dk.FrostFeverDisease = make([]*core.Dot, dk.Env.GetNumTargets())

	var wpWrapper func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)
	if dk.Talents.WanderingPlague > 0 {
		wpWrapper = dk.wpWrapper
	}
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.FrostFeverDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    FrostFeverAuraLabel + strconv.Itoa(int(dk.Index)),
				ActionID: actionID,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if !isRefreshing[aura.Unit.Index] {
						flagTs[aura.Unit.Index] = false
					}
				},
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:              core.ProcMaskPeriodicDamage,
				DamageMultiplier:      core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
				ThreatMultiplier:      1,
				IsPeriodic:            true,
				OnPeriodicDamageDealt: wpWrapper,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						firstTsApply := !flagTs[hitEffect.Target.Index]
						flagTs[hitEffect.Target.Index] = true
						return ((127.0 + 80.0*0.32) + dk.applyImpurity(hitEffect, spell.Unit)*0.055) *
							core.TernaryFloat64(firstTsApply, 1.0, dk.rageOfRivendareBonus(hitEffect.Target)*
								dk.tundraStalkerBonus(hitEffect.Target))
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: dk.OutcomeFuncAlwaysHit(),
			}),
		})

		dk.FrostFeverDisease[target.Index].Spell = dk.FrostFeverSpell
	}
}

func (dk *Deathknight) registerBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}

	flagRor := make([]bool, dk.Env.GetNumTargets())
	isRefreshing := make([]bool, dk.Env.GetNumTargets())

	dk.BloodPlagueSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			if dk.BloodPlagueDisease[unit.Index].IsActive() {
				isRefreshing[unit.Index] = true
			}
			dk.BloodPlagueDisease[unit.Index].Apply(sim)
			isRefreshing[unit.Index] = false
		},
	})

	dk.BloodPlagueDisease = make([]*core.Dot, dk.Env.GetNumTargets())

	// Tier9 4Piece
	outcomeApplier := dk.OutcomeFuncAlwaysHit()
	if dk.HasSetBonus(ItemSetThassariansBattlegear, 4) {
		outcomeApplier = dk.OutcomeFuncMagicCrit(dk.spellCritMultiplier())
	}

	var wpWrapper func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)
	if dk.Talents.WanderingPlague > 0 {
		wpWrapper = dk.wpWrapper
	}
	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.BloodPlagueDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    BloodPlagueAuraLabel + strconv.Itoa(int(dk.Index)),
				ActionID: actionID,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if !isRefreshing[aura.Unit.Index] {
						flagRor[aura.Unit.Index] = false
					}
				},
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:              core.ProcMaskPeriodicDamage,
				DamageMultiplier:      1,
				ThreatMultiplier:      1,
				IsPeriodic:            true,
				OnPeriodicDamageDealt: wpWrapper,
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						firstRorApply := !flagRor[hitEffect.Target.Index]
						flagRor[hitEffect.Target.Index] = true
						return ((127.0 + 80.0*0.32) + dk.applyImpurity(hitEffect, spell.Unit)*0.055) *
							core.TernaryFloat64(firstRorApply, 1.0, dk.rageOfRivendareBonus(hitEffect.Target)*
								dk.tundraStalkerBonus(hitEffect.Target))
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: outcomeApplier,
			}),
		})

		dk.BloodPlagueDisease[target.Index].Spell = dk.BloodPlagueSpell
	}
}

func (dk *Deathknight) wpWrapper(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	dk.doWanderingPlague(sim, spell, spellEffect)
}

func (dk *Deathknight) doWanderingPlague(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	if dk.Talents.WanderingPlague == 0 {
		return
	}

	critRating := spell.Unit.GetStats()[stats.MeleeCrit] + spellEffect.BonusCritRating + spellEffect.Target.PseudoStats.BonusCritRatingTaken
	critRating += spell.Unit.PseudoStats.BonusMeleeCritRating
	critChance := critRating / (core.CritRatingPerCritChance * 100)
	if sim.RandomFloat("Wandering Plague Roll") < critChance {
		dk.LastDiseaseDamage = spellEffect.Damage
		dk.WanderingPlague.Cast(sim, spellEffect.Target)
	}
}
