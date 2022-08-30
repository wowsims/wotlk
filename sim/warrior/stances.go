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

func (warrior *Warrior) StanceMatches(other Stance) bool {
	return (warrior.Stance & other) != 0
}

func (warrior *Warrior) makeStanceSpell(stance Stance, aura *core.Aura, stanceCD *core.Timer) *core.Spell {
	maxRetainedRage := 10.0 + 5*float64(warrior.Talents.TacticalMastery)
	actionID := aura.ActionID
	rageMetrics := warrior.NewRageMetrics(actionID)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if warrior.Stance == stance {
				panic("Already in stance " + string(stance))
			}

			if warrior.CurrentRage() > maxRetainedRage {
				warrior.SpendRage(sim, warrior.CurrentRage()-maxRetainedRage, rageMetrics)
			}

			// Add new stance aura.
			aura.Activate(sim)
			warrior.Stance = stance
		},
	})
}

func (warrior *Warrior) registerBattleStanceAura() {
	actionID := core.ActionID{SpellID: 2457}
	threatMult := 0.8
	armorPenBonus := core.ArmorPenPerPercentArmor * (10 + core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsBattlegear, 2), 6, 0))

	warrior.BattleStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Battle Stance",
		Tag:      "Stance",
		Priority: 1,
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
}

func (warrior *Warrior) registerDefensiveStanceAura() {
	actionID := core.ActionID{SpellID: 71}
	threatMult := 2.0735
	// TODO: Imp def stance
	impDefStanceMultiplier := 1 - 0.03*float64(warrior.Talents.ImprovedDefensiveStance)

	if warrior.Talents.ImprovedDefensiveStance > 0 {
		enrageAura := warrior.GetOrRegisterAura(core.Aura{
			Label:    "Enrage",
			ActionID: core.ActionID{SpellID: 57516},
			Duration: 12 * time.Second,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.PhysicalDamageDealtMultiplier *= 1.0 + 0.05*float64(warrior.Talents.ImprovedDefensiveStance)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.PhysicalDamageDealtMultiplier /= 1.0 + 0.05*float64(warrior.Talents.ImprovedDefensiveStance)
			},
		})

		core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
			Label:    "Enrage Trigger",
			Duration: core.NeverExpires,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
					if sim.RandomFloat("Enrage Trigger Chance") <= 0.5*float64(warrior.Talents.ImprovedDefensiveStance) {
						enrageAura.Activate(sim)
					}
				}
			},
		}))

	}

	warrior.DefensiveStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Defensive Stance",
		Tag:      "Stance",
		Priority: 1,
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 0.9
			aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier *= impDefStanceMultiplier
			aura.Unit.PseudoStats.FireDamageTakenMultiplier *= impDefStanceMultiplier
			aura.Unit.PseudoStats.FrostDamageTakenMultiplier *= impDefStanceMultiplier
			aura.Unit.PseudoStats.HolyDamageTakenMultiplier *= impDefStanceMultiplier
			aura.Unit.PseudoStats.NatureDamageTakenMultiplier *= impDefStanceMultiplier
			aura.Unit.PseudoStats.ShadowDamageTakenMultiplier *= impDefStanceMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 0.9
			aura.Unit.PseudoStats.ArcaneDamageTakenMultiplier /= impDefStanceMultiplier
			aura.Unit.PseudoStats.FireDamageTakenMultiplier /= impDefStanceMultiplier
			aura.Unit.PseudoStats.FrostDamageTakenMultiplier /= impDefStanceMultiplier
			aura.Unit.PseudoStats.HolyDamageTakenMultiplier /= impDefStanceMultiplier
			aura.Unit.PseudoStats.NatureDamageTakenMultiplier /= impDefStanceMultiplier
			aura.Unit.PseudoStats.ShadowDamageTakenMultiplier /= impDefStanceMultiplier
		},
	})
}

func (warrior *Warrior) registerBerserkerStanceAura() {
	threatMult := 0.8 - 0.02*float64(warrior.Talents.ImprovedBerserkerStance)
	critBonus := core.CritRatingPerCritChance * (3 + core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsBattlegear, 2), 2, 0))

	warrior.BerserkerStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Berserker Stance",
		Tag:      "Stance",
		Priority: 1,
		ActionID: core.ActionID{SpellID: 2458},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= threatMult
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, critBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier /= threatMult
			aura.Unit.AddStatDynamic(sim, stats.MeleeCrit, -critBonus)
		},
	})
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
