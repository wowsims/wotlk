package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	// demonic embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		bonus := 1.01 + float64(warlock.Talents.DemonicEmbrace)*0.03
		warlock.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(in float64, _ float64) float64 {
				return in * bonus
			},
		})
	}

	// Suppression
	warlock.AddStat(stats.SpellHit, float64(warlock.Talents.Suppression)*core.SpellHitRatingPerHitChance)

	// Shadow Mastery
	warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.03*float64(warlock.Talents.ShadowMastery)

	// Backlash (Add 1% crit per level)
	warlock.AddStat(stats.SpellCrit, float64(warlock.Talents.Backlash)*core.CritRatingPerCritChance)

	// Malediction (SP bonus)
	if warlock.Talents.Malediction > 0 {
		factor := 1 + 0.01*float64(warlock.Talents.Malediction)
		warlock.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.SpellPower,
			ModifiedStat: stats.SpellPower,
			Modifier: func(sp float64, _ float64) float64 {
				return sp * factor
			},
		})
	}

	// Fel Vitality
	if warlock.Talents.FelVitality > 0 {
		bonus := 0.01 * float64(warlock.Talents.FelVitality)
		// Adding a second 3% bonus int->mana dependency
		// TODO: increases max health
		warlock.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Mana,
			Modifier: func(intellect float64, mana float64) float64 {
				return mana + intellect*15*bonus
			},
		})
	}

	warlock.PseudoStats.BonusCritRating += float64(warlock.Talents.DemonicTactics) * 1 * core.CritRatingPerCritChance

	if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		if warlock.Talents.MasterDemonologist > 0 {
			switch warlock.Options.Summon {
			case proto.Warlock_Options_Imp:
				warlock.PseudoStats.FireDamageDealtMultiplier *= 1.0 + 0.01 * float64(warlock.Talents.MasterDemonologist)
				warlock.PseudoStats.BonusFireCritRating *= 1.0 + 0.01 * float64(warlock.Talents.MasterDemonologist)
			case proto.Warlock_Options_Succubus:
				warlock.PseudoStats.ShadowDamageDealtMultiplier *= 1.0 + 0.01 * float64(warlock.Talents.MasterDemonologist)
				warlock.PseudoStats.BonusShadowCritRating *= 1.0 + 0.01 * float64(warlock.Talents.MasterDemonologist)
			case proto.Warlock_Options_Felguard:
				warlock.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
			}
		}
		// Extract stats for demonic knowledge
		if warlock.Talents.DemonicKnowledge > 0 {
			petChar := warlock.Pets[0].GetCharacter()
			bonus := (petChar.GetStat(stats.Stamina) + petChar.GetStat(stats.Intellect)) * (0.04 * float64(warlock.Talents.DemonicKnowledge))
			warlock.AddStat(stats.SpellPower, bonus)
 		}
	}

	// demonic tactics, applies even without pet out
	if warlock.Talents.DemonicTactics > 0 {
		warlock.AddStats(stats.Stats{
			stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
			stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 2 * core.CritRatingPerCritChance,
		})
	}

	if warlock.Talents.DemonicPact > 0 {
		spBonus := 0.02 * float64(warlock.Talents.DemonicPact) * float64(stats.SpellPower)
		warlock.AddStat(stats.SpellPower, spBonus)
	}

 	if warlock.Talents.MoltenSkin > 0 {
 		warlock.PseudoStats.DamageTakenMultiplier /= 1 + 0.02 * float64(warlock.Talents.MoltenSkin)
 	}
 	
	if warlock.Talents.Nightfall > 0 {
		warlock.setupNightfall()
	}

	if warlock.Talents.ShadowEmbrace > 0 {
		warlock.setupShadowEmbrace()
	}

	if warlock.Talents.Eradication > 0 {
		warlock.setupEradication()
	}

	if warlock.Talents.DeathsEmbrace > 0 {
		warlock.applyDeathsEmbrace()
	}

}

func (warlock *Warlock) applyDeathsEmbrace() {
	multiplier := 1.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute20 bool) {
			if isExecute20 {
				warlock.PseudoStats.ShadowDamageDealtMultiplier *= multiplier
			}
		})
	})
}

func (warlock *Warlock) setupEradication() {
	hasteBonusPercent := float64(warlock.Talents.Eradication) * 6
	if warlock.Talents.Eradication == 3 {
		hasteBonusPercent += 2
	}
	warlock.EradicationAura = warlock.RegisterAura(core.Aura{
		Label:    "Eradication",
		ActionID: core.ActionID{SpellID: 64371},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 + hasteBonusPercent/100)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / (1 + hasteBonusPercent/100))
		},
	})

	warlock.RegisterAura(core.Aura{
		Label: "Eradication Talent Hidden Aura",
		// ActionID:  core.ActionID{SpellID: 47197},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Corruption {
				if sim.RandomFloat("Eradication") < 0.06 {
					warlock.EradicationAura.Activate(sim)
				}
			}
		},
	})
}

func (warlock *Warlock) setupShadowEmbrace() {
	warlock.ShadowEmbraceAura = warlock.RegisterAura(core.Aura{
		Label:     "Shadow Embrace",
		ActionID:  core.ActionID{SpellID: 32391},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier /= 1.0 + 0.01*float64(warlock.Talents.ShadowEmbrace)*float64(oldStacks)
			aura.Unit.PseudoStats.PeriodicShadowDamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.ShadowEmbrace)*float64(newStacks)
			// TO DO : Healing over time reduction part
		},
	})

	warlock.RegisterAura(core.Aura{
		Label: "Shadow Embrace Talent Hidden Aura",
		//		ActionID: core.ActionID{SpellID: 32394},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell == warlock.Shadowbolt || spell == warlock.Haunt {
				if !warlock.ShadowEmbraceAura.IsActive() {
					warlock.ShadowEmbraceAura.Activate(sim)
				}
				warlock.ShadowEmbraceAura.AddStack(sim)
			}
		},
	})
}

func (warlock *Warlock) setupNightfall() {
	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check for an instant cast shadowbolt to disable aura
			if spell != warlock.Shadowbolt || spell.CurCast.CastTime != 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	warlock.RegisterAura(core.Aura{
		Label: "Nightfall Hidden Aura",
		// ActionID: core.ActionID{SpellID: 18095},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell != warlock.Corruption { // TODO: also works on drain life...
				return
			}
			if sim.RandomFloat("nightfall") > 0.02*float64(warlock.Talents.Nightfall) {
				return
			}
			warlock.NightfallProcAura.Activate(sim)
		},
	})
}

func (warlock *Warlock) applyNightfall(cast *core.Cast) {
	if warlock.NightfallProcAura.IsActive() {
		cast.CastTime = 0
	}
}
