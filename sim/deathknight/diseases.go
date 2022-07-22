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

func (deathKnight *DeathKnight) countActiveDiseases(target *core.Unit) int {
	count := 0
	if deathKnight.TargetHasDisease(FrostFeverAuraLabel, target) {
		count++
	}
	if deathKnight.TargetHasDisease(BloodPlagueAuraLabel, target) {
		count++
	}
	if deathKnight.TargetHasDisease(core.EbonPlaguebringerAuraLabel, target) || deathKnight.TargetHasDisease(core.CryptFeverAuraLabel, target) {
		count++
	}
	return count
}

func (deathKnight *DeathKnight) TargetHasDisease(label string, unit *core.Unit) bool {
	return unit.HasActiveAura(label + strconv.Itoa(int(deathKnight.Index)))
}

func (deathKnight *DeathKnight) diseaseMultiplierBonus(target *core.Unit, multiplier float64) float64 {
	return 1.0 + float64(deathKnight.countActiveDiseases(target))*deathKnight.darkrunedBattlegearDiseaseBonus(multiplier)
}

func (deathKnight *DeathKnight) registerDiseaseDots() {
	deathKnight.registerFrostFever()
	deathKnight.registerBloodPlague()
}

func (deathKnight *DeathKnight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}

	deathKnight.FrostFeverSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			deathKnight.FrostFeverDisease[unit.Index].Apply(sim)
		},
	})

	deathKnight.FrostFeverDisease = make([]*core.Dot, deathKnight.Env.GetNumTargets())

	for _, encounterTarget := range deathKnight.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		deathKnight.FrostFeverDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    FrostFeverAuraLabel + strconv.Itoa(int(deathKnight.Index)),
				ActionID: actionID,
			}),
			NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: core.TernaryFloat64(deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfIcyTouch), 1.2, 1.0),
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					deathKnight.doWanderingPlague(sim, spell, spellEffect)
				},
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						return ((127.0 + 80.0*0.32) + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.055) *
							deathKnight.rageOfRivendareBonus(hitEffect.Target) *
							deathKnight.tundraStalkerBonus(hitEffect.Target)
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
			}),
		})

		deathKnight.FrostFeverDisease[target.Index].Spell = deathKnight.FrostFeverSpell
	}
}

func (deathKnight *DeathKnight) registerBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}

	deathKnight.BloodPlagueSpell = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagDisease,
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			deathKnight.BloodPlagueDisease[unit.Index].Apply(sim)
		},
	})

	deathKnight.BloodPlagueDisease = make([]*core.Dot, deathKnight.Env.GetNumTargets())

	for _, encounterTarget := range deathKnight.Env.Encounter.Targets {
		target := &encounterTarget.Unit

		// Tier9 4Piece
		outcomeApplier := deathKnight.OutcomeFuncAlwaysHit()
		if deathKnight.HasSetBonus(ItemSetThassariansBattlegear, 4) {
			outcomeApplier = deathKnight.OutcomeFuncMagicCrit(deathKnight.spellCritMultiplier())
		}
		deathKnight.BloodPlagueDisease[target.Index] = core.NewDot(core.Dot{
			Aura: target.RegisterAura(core.Aura{
				Label:    BloodPlagueAuraLabel + strconv.Itoa(int(deathKnight.Index)),
				ActionID: actionID,
			}),
			NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
			TickLength:    time.Second * 3,

			TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
				ProcMask:         core.ProcMaskPeriodicDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				IsPeriodic:       true,
				OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
					deathKnight.doWanderingPlague(sim, spell, spellEffect)
				},
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
						return ((127.0 + 80.0*0.32) + deathKnight.applyImpurity(hitEffect, spell.Unit)*0.055) *
							deathKnight.rageOfRivendareBonus(hitEffect.Target) *
							deathKnight.tundraStalkerBonus(hitEffect.Target)
					},
					TargetSpellCoefficient: 1,
				},
				OutcomeApplier: outcomeApplier,
			}),
		})

		deathKnight.BloodPlagueDisease[target.Index].Spell = deathKnight.BloodPlagueSpell
	}
}

func (deathKnight *DeathKnight) doWanderingPlague(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	if deathKnight.Talents.WanderingPlague == 0 {
		return
	}

	critRating := spell.Unit.GetStats()[stats.MeleeCrit] + spellEffect.BonusCritRating + spellEffect.Target.PseudoStats.BonusCritRatingTaken
	critRating += spell.Unit.PseudoStats.BonusMeleeCritRating
	critChance := critRating / (core.CritRatingPerCritChance * 100)
	if sim.RandomFloat("Wandering Plague Roll") < critChance {
		deathKnight.LastDiseaseDamage = spellEffect.Damage
		deathKnight.WanderingPlague.Cast(sim, spellEffect.Target)
	}
}
