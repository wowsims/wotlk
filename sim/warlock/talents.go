package warlock

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) ApplyTalents() {
	// demonic embrace
	if warlock.Talents.DemonicEmbrace > 0 {
		bonus := 1 + (0.03)*float64(warlock.Talents.DemonicEmbrace)
		negative := 1 - (0.01)*float64(warlock.Talents.DemonicEmbrace)
		warlock.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(in float64, _ float64) float64 {
				return in * bonus
			},
		})
		warlock.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Spirit,
			ModifiedStat: stats.Spirit,
			Modifier: func(in float64, _ float64) float64 {
				return in * negative
			},
		})
	}

	// Add 1% crit per level of backlash.
	warlock.PseudoStats.BonusCritRating += float64(warlock.Talents.Backlash) * 1 * core.SpellCritRatingPerCritChance

	// fel intellect
	if warlock.Talents.FelIntellect > 0 {
		bonus := (0.01) * float64(warlock.Talents.FelIntellect)
		// Adding a second 3% bonus int->mana dependency
		warlock.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Mana,
			Modifier: func(intellect float64, mana float64) float64 {
				return mana + intellect*(15*bonus)
			},
		})
	}

	warlock.PseudoStats.BonusCritRating += float64(warlock.Talents.DemonicTactics) * 1 * core.SpellCritRatingPerCritChance

	//  TODO: fel stamina increases max health (might be useful for warlock tanking sim)

	if !warlock.Options.SacrificeSummon && warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		if warlock.Talents.MasterDemonologist > 0 {
			switch warlock.Options.Summon {
			case proto.Warlock_Options_Imp:
				warlock.PseudoStats.ThreatMultiplier *= 0.96 * float64(warlock.Talents.MasterDemonologist)
			case proto.Warlock_Options_Succubus:
				warlock.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.02*float64(warlock.Talents.MasterDemonologist)
			case proto.Warlock_Options_Felgaurd:
				warlock.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
				// 		Felguard - Increases all damage caused by 1% and all resistances by .1 per level.
				// 		Voidwalker - Reduces physical damage taken by 2%.
				// 		Felhunter - Increases all resistances by .2 per level.
			}
		}

		if warlock.Talents.SoulLink {
			warlock.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Extract stats for demonic knowledge
		petChar := warlock.Pets[0].GetCharacter()
		bonus := (petChar.GetStat(stats.Stamina) + petChar.GetStat(stats.Intellect)) * (0.04 * float64(warlock.Talents.DemonicKnowledge))
		warlock.AddStat(stats.SpellPower, bonus)
	}

	// demonic tactics, applies even without pet out
	warlock.AddStats(stats.Stats{
		stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 1 * core.MeleeCritRatingPerCritChance,
		stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 1 * core.SpellCritRatingPerCritChance,
	})

	warlock.applyShadowEmbrace()
	warlock.setupNightfall()
	warlock.setupAmplifyCurse()
}

func (warlock *Warlock) applyShadowEmbrace() {
	if warlock.Talents.ShadowEmbrace == 0 {
		return
	}

	var debuffAuras []*core.Aura
	for _, target := range warlock.Env.Encounter.Targets {
		debuffAuras = append(debuffAuras, core.ShadowEmbraceAura(&target.Unit, warlock.Talents.ShadowEmbrace))
	}

	warlock.RegisterAura(core.Aura{
		Label:    "Shadow Embrace Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			if spell == warlock.Corruption || spell == warlock.SiphonLife || spell == warlock.CurseOfAgony || spell.SameAction(warlock.Seeds[0].ActionID) {
				debuffAuras[spellEffect.Target.Index].Activate(sim)
			}
		},
	})
}

func (warlock *Warlock) setupAmplifyCurse() {
	if !warlock.Talents.AmplifyCurse {
		return
	}
	warlock.AmplifyCurseAura = warlock.RegisterAura(core.Aura{
		Label:    "Amplify Curse",
		ActionID: core.ActionID{SpellID: 18288},
		Duration: time.Second * 30,
	})
	warlock.AmplifyCurse = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 18288},
		SpellSchool: core.SpellSchoolShadow,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.AmplifyCurseAura.Activate(sim)
		},
	})
}

func (warlock *Warlock) setupNightfall() {
	if warlock.Talents.Nightfall == 0 {
		return
	}

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
		Label: "Nightfall",
		// ActionID: core.ActionID{SpellID: 18095},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell != warlock.Corruption { // TODO: also works on drain life...
				return
			}
			if sim.RandomFloat("nightfall") > 0.04 {
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
