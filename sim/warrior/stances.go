package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type Stance uint8

const (
	BattleStance Stance = 1 << iota
	DefensiveStance
	BerserkerStance
)

const stanceEffectCategory = "Stance"

func (warrior *Warrior) StanceMatches(other Stance) bool {
	return (warrior.Stance & other) != 0
}

func (warrior *Warrior) makeStanceSpell(stance Stance, aura *core.Aura, stanceCD *core.Timer) *core.Spell {
	maxRetainedRage := 10.0 + 5*float64(warrior.Talents.TacticalMastery)
	actionID := aura.ActionID
	rageMetrics := warrior.NewRageMetrics(actionID)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.Stance != stance
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if warrior.CurrentRage() > maxRetainedRage {
				warrior.SpendRage(sim, warrior.CurrentRage()-maxRetainedRage, rageMetrics)
			}

			if warrior.WarriorInputs.StanceSnapshot {
				// Delayed, so same-GCD casts are affected by the current aura.
				//  Alternatively, those casts could just (artificially) happen before the stance change.
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt:     sim.CurrentTime + 10*time.Millisecond,
					OnAction: aura.Activate,
				})
			} else {
				aura.Activate(sim)
			}

			warrior.Stance = stance
		},
	})
}

func (warrior *Warrior) registerBattleStanceAura() {
	const threatMult = 0.8

	actionID := core.ActionID{SpellID: 2457}
	armorPenBonus := core.ArmorPenPerPercentArmor * (10 + core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsBattlegear, 2), 6, 0))

	warrior.BattleStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Battle Stance",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.AddStatDynamic(sim, stats.ArmorPenetration, armorPenBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.AddStatDynamic(sim, stats.ArmorPenetration, -armorPenBonus)
		},
	})
	warrior.BattleStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
}

func (warrior *Warrior) registerDefensiveStanceAura() {
	const threatMult = 2.0735

	actionID := core.ActionID{SpellID: 71}
	if warrior.Talents.ImprovedDefensiveStance > 0 {
		enrageAura := warrior.GetOrRegisterAura(core.Aura{
			Label:    "Enrage",
			ActionID: core.ActionID{SpellID: 57516},
			Duration: 12 * time.Second,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.0 + 0.05*float64(warrior.Talents.ImprovedDefensiveStance)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.0 + 0.05*float64(warrior.Talents.ImprovedDefensiveStance)
			},
		})

		core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
			Label:    "Enrage Trigger",
			Duration: core.NeverExpires,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
					if sim.RandomFloat("Enrage Trigger Chance") <= 0.5*float64(warrior.Talents.ImprovedDefensiveStance) {
						enrageAura.Activate(sim)
					}
				}
			},
		}))
	}

	impDefStanceMultiplier := 1 - 0.03*float64(warrior.Talents.ImprovedDefensiveStance)
	tacMasteryThreatMultiplier := 1 + 0.21*float64(warrior.Talents.TacticalMastery)

	warrior.DefensiveStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Defensive Stance",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 0.95
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.90
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= impDefStanceMultiplier
			if warrior.Bloodthirst != nil {
				warrior.Bloodthirst.ThreatMultiplier *= tacMasteryThreatMultiplier
			}
			if warrior.MortalStrike != nil {
				warrior.MortalStrike.ThreatMultiplier *= tacMasteryThreatMultiplier
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 0.95
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.9
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= impDefStanceMultiplier
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= impDefStanceMultiplier
			if warrior.Bloodthirst != nil {
				warrior.Bloodthirst.ThreatMultiplier /= tacMasteryThreatMultiplier
			}
			if warrior.MortalStrike != nil {
				warrior.MortalStrike.ThreatMultiplier /= tacMasteryThreatMultiplier
			}
		},
	})
	warrior.DefensiveStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
}

func (warrior *Warrior) registerBerserkerStanceAura() {
	threatMult := 0.8 - 0.02*float64(warrior.Talents.ImprovedBerserkerStance)
	critBonus := core.CritRatingPerCritChance * (3 + core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsBattlegear, 2), 2, 0))

	var dep *stats.StatDependency
	if warrior.Talents.ImprovedBerserkerStance > 0 {
		// alternatively, this could be default on
		dep = warrior.NewDynamicMultiplyStat(stats.Strength, 1.0+0.04*float64(warrior.Talents.ImprovedBerserkerStance))
	}

	warrior.BerserkerStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Berserker Stance",
		ActionID: core.ActionID{SpellID: 2458},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, critBonus)
			if dep != nil {
				warrior.EnableDynamicStatDep(sim, dep)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, -critBonus)
			if dep != nil {
				warrior.DisableDynamicStatDep(sim, dep)
			}
		},
	})
	warrior.BerserkerStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
}

func (warrior *Warrior) registerStances() {
	stanceCD := warrior.NewTimer()
	warrior.registerBattleStanceAura()
	warrior.registerDefensiveStanceAura()
	warrior.registerBerserkerStanceAura()
	warrior.BattleStance = warrior.makeStanceSpell(BattleStance, warrior.BattleStanceAura, stanceCD)
	warrior.DefensiveStance = warrior.makeStanceSpell(DefensiveStance, warrior.DefensiveStanceAura, stanceCD)
	warrior.BerserkerStance = warrior.makeStanceSpell(BerserkerStance, warrior.BerserkerStanceAura, stanceCD)
}
