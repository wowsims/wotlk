package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerKillCommandCD() {
	if hunter.pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 34026}
	hunter.pet.KillCommandAura = hunter.pet.RegisterAura(core.Aura{
		Label:     "Kill Command",
		ActionID:  actionID,
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})

	baseCost := 0.03 * hunter.BaseMana

	hunter.KillCommand = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoOnCastComplete,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute - time.Second*10*time.Duration(hunter.Talents.CatlikeReflexes),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.pet.KillCommandAura.Activate(sim)
			hunter.pet.KillCommandAura.SetStacks(sim, 3)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.KillCommand,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hunter.pet.IsEnabled() && hunter.CurrentMana() >= hunter.KillCommand.DefaultCast.Cost
		},
	})
}
