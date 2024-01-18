package hunter

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (hunter *Hunter) ApplyTalents() {
	if hunter.pet != nil {
		hunter.applyFrenzy()
		hunter.registerBestialWrathCD()

		hunter.pet.AddStat(stats.MeleeCrit, core.CritRatingPerCritChance*3*float64(hunter.Talents.Ferocity))
		hunter.pet.AddStat(stats.SpellCrit, core.SpellCritRatingPerCritChance*3*float64(hunter.Talents.Ferocity))
		hunter.pet.PseudoStats.DamageDealtMultiplier *= 1 + 0.04*float64(hunter.Talents.UnleashedFury)

		hunter.pet.MultiplyStat(stats.Health, 1+(0.03*float64(hunter.Talents.EnduranceTraining)))
	}

}

func (hunter *Hunter) ApplyRunes() {

}

func (hunter *Hunter) critMultiplier(isRanged bool, target *core.Unit) float64 {
	primaryModifier := 1.0
	secondaryModifier := 0.0

	monsterMultiplier := 1.0 + 0.01*float64(hunter.Talents.MonsterSlaying)
	humanoidMultiplier := 1.0 + 0.01*float64(hunter.Talents.HumanoidSlaying)
	if target.MobType == proto.MobType_MobTypeBeast || target.MobType == proto.MobType_MobTypeGiant || target.MobType == proto.MobType_MobTypeDragonkin {
		primaryModifier *= monsterMultiplier
	} else if target.MobType == proto.MobType_MobTypeHumanoid {
		primaryModifier *= humanoidMultiplier
	}

	if isRanged {
		secondaryModifier += 0.06 * float64(hunter.Talents.MortalShots)
	}

	return hunter.MeleeCritMultiplier(primaryModifier, secondaryModifier)
}

func (hunter *Hunter) applyFrenzy() {
	if hunter.Talents.Frenzy == 0 {
		return
	}

	procChance := 0.2 * float64(hunter.Talents.Frenzy)

	procAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy Proc",
		ActionID: core.ActionID{SpellID: 19625},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.3)
		},
	})

	hunter.pet.RegisterAura(core.Aura{
		Label:    "Frenzy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellResult *core.SpellResult) {
			if !spellResult.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if procChance == 1 || sim.RandomFloat("Frenzy") < procChance {
				procAura.Activate(sim)
			}
		},
	})
}

func (hunter *Hunter) registerBestialWrathCD() {
	if !hunter.Talents.BestialWrath {
		return
	}

	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := hunter.pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 18,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.5
		},
	})

	bwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bestialWrathPetAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}
