package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (warrior *Warrior) RegisterRecklessnessCD() {
	actionID := core.ActionID{SpellID: 1719}
	reckAura := warrior.RegisterAura(core.Aura{
		Label:     "Recklessness",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.BonusMeleeSpellCritRating += 100 * core.CritRatingPerCritChance
			warrior.PseudoStats.DamageTakenMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.BonusMeleeSpellCritRating -= 100 * core.CritRatingPerCritChance
			warrior.PseudoStats.DamageTakenMultiplier /= 1.2
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
				aura.RemoveStack(sim)
			}
		},
	})

	cooldownDur := time.Minute * 30
	if warrior.Talents.IntensifyRage == 1 {
		cooldownDur *= (100 - 11) / 100
	} else if warrior.Talents.IntensifyRage == 2 {
		cooldownDur *= (100 - 22) / 100
	} else if warrior.Talents.IntensifyRage == 3 {
		cooldownDur *= (100 - 33) / 100
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
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return warrior.StanceMatches(BerserkerStance)
		},
	})
}
