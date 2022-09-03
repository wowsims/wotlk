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
			affectedSpells = []*core.Spell{
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
				warrior.Shockwave,
				warrior.ConcussionBlow,
				warrior.Bladestorm,
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
			for _, spell := range affectedSpells {
				if spell != nil {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
			for _, spell := range affectedSpells {
				if spell != nil {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial) || spellEffect.Damage <= 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})

	cooldownDur := time.Minute * 30
	if warrior.Talents.IntensifyRage == 1 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.89)
	} else if warrior.Talents.IntensifyRage == 2 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.78)
	} else if warrior.Talents.IntensifyRage == 3 {
		cooldownDur = time.Duration(float64(cooldownDur) * 0.67)
	}
	reckSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			reckAura.Activate(sim)
			reckAura.SetStacks(sim, 3)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: reckSpell,
		Type:  core.CooldownTypeDPS,
		ActivationFactory: func(sim *core.Simulation) core.CooldownActivation {
			return func(sim *core.Simulation, character *core.Character) {
				if !warrior.StanceMatches(BerserkerStance) {
					if !warrior.BerserkerStance.IsReady(sim) {
						return
					}
					warrior.BerserkerStance.Cast(sim, nil)
				}
				reckSpell.Cast(sim, character.CurrentTarget)
			}
		},
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
	})
}
