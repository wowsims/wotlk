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

func (dk *Deathknight) countActiveDiseases(target *core.Unit) int {
	count := 0
	if dk.TargetHasDisease(FrostFeverAuraLabel, target) {
		count++
	}
	if dk.TargetHasDisease(BloodPlagueAuraLabel, target) {
		count++
	}
	if dk.TargetHasDisease(core.EbonPlaguebringerAuraLabel, target) || dk.TargetHasDisease(core.CryptFeverAuraLabel, target) {
		count++
	}
	return count
}

func (dk *Deathknight) TargetHasDisease(label string, unit *core.Unit) bool {
	return unit.HasActiveAura(label + strconv.Itoa(int(dk.Index)))
}

func (dk *Deathknight) diseaseMultiplierBonus(target *core.Unit, multiplier float64) float64 {
	return 1.0 + float64(dk.countActiveDiseases(target))*dk.darkrunedBattlegearDiseaseBonus(multiplier)
}

func (dk *Deathknight) registerDiseaseDots() {
	dk.registerFrostFever()
	dk.registerBloodPlague()
}

func (dk *Deathknight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}

	applyTS := make([]bool, dk.Env.GetNumTargets())
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

	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		dk.FrostFeverDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    FrostFeverAuraLabel + strconv.Itoa(int(dk.Index)),
				ActionID: actionID,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if !isRefreshing[aura.Unit.Index] {
						applyTS[aura.Unit.Index] = false
					}
				},
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					dk.doWanderingPlague(sim, spell, spellEffect)
				},
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						applyTsBonus := !applyTS[hitEffect.Target.Index]
						applyTS[hitEffect.Target.Index] = true
						return ((127.0 + 80.0*0.32) + dk.applyImpurity(hitEffect, spell.Unit)*0.055) *
							core.TernaryFloat64(applyTsBonus, 1.0, dk.rageOfRivendareBonus(hitEffect.Target)*
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

	applyRor := make([]bool, dk.Env.GetNumTargets())
	isRefreshing := make([]bool, dk.Env.GetNumTargets())

	dk.BloodPlagueSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.BloodPlagueDisease[unit.Index].Apply(sim)
		},
	})

	dk.BloodPlagueDisease = make([]*core.Dot, dk.Env.GetNumTargets())

	for _, encounterTarget := range dk.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		// Tier9 4Piece
		outcomeApplier := dk.OutcomeFuncAlwaysHit()
		if dk.HasSetBonus(ItemSetThassariansBattlegear, 4) {
			outcomeApplier = dk.OutcomeFuncMagicCrit(dk.spellCritMultiplier())
		}
		dk.BloodPlagueDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    BloodPlagueAuraLabel + strconv.Itoa(int(dk.Index)),
				ActionID: actionID,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if !isRefreshing[aura.Unit.Index] {
						applyRor[aura.Unit.Index] = false
					}
				},
			}),
			NumberOfTicks: 5 + int(dk.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					dk.doWanderingPlague(sim, spell, spellEffect)
				},
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						applyRorBonus := !applyRor[hitEffect.Target.Index]
						applyRor[hitEffect.Target.Index] = true
						return ((127.0 + 80.0*0.32) + dk.applyImpurity(hitEffect, spell.Unit)*0.055) *
							core.TernaryFloat64(applyRorBonus, 1.0, dk.rageOfRivendareBonus(hitEffect.Target)*
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
