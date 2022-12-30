package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) registerShadowfiendSpell() {
	if !priest.UseShadowfiend {
		return
	}

	actionID := core.ActionID{SpellID: 34433}

	priest.ShadowfiendAura = priest.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Shadowfiend",
		Duration: time.Second * 15.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.ShadowfiendPet.Enable(sim, priest.ShadowfiendPet)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.ShadowfiendPet.Disable(sim)
		},
	})

	priest.Shadowfiend = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Minute * time.Duration(5-priest.Talents.VeiledShadows),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			priest.ShadowfiendAura.Activate(sim)
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.Shadowfiend,
		Priority: 1,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.5
		},
	})
}
