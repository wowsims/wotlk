package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) RegisterRecklessnessCD() {
	actionID := core.ActionID{SpellID: 1719}
	var affectedSpells []*core.Spell

	reckAura := warrior.RegisterAura(core.Aura{
		Label:     "Recklessness",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*core.Spell{
				warrior.HeroicStrikeOrCleave,
				warrior.Bloodthirst,
				warrior.Devastate,
				warrior.Execute,
				warrior.MortalStrike,
				warrior.Overpower,
				warrior.Revenge,
				warrior.ShieldSlam,
				warrior.Slam,
				warrior.ThunderClap,
				warrior.Whirlwind,
				warrior.WhirlwindOH,
				warrior.Shockwave,
				warrior.ConcussionBlow,
				warrior.Bladestorm,
				warrior.BladestormOH,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
			for _, spell := range affectedSpells {
				spell.BonusCritRating += 100 * core.CritRatingPerCritChance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
			for _, spell := range affectedSpells {
				spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) || result.Damage <= 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})

	reckSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: warrior.intensifyRageCooldown(time.Minute * 5),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BerserkerStance) || warrior.BerserkerStance.IsReady(sim)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if !warrior.StanceMatches(BerserkerStance) {
				warrior.BerserkerStance.Cast(sim, nil)
			}

			reckAura.Activate(sim)
			reckAura.SetStacks(sim, 3)
			warrior.WaitUntil(sim, sim.CurrentTime+core.GCDDefault)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: reckSpell,
		Type:  core.CooldownTypeDPS,
	})
}
