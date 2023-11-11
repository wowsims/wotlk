package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) registerDivinePleaSpell() {
	actionID := core.ActionID{SpellID: 54428}
	hasGlyph := paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivinePlea)
	manaMetrics := paladin.NewManaMetrics(actionID)
	var manaPA *core.PendingAction

	paladin.DivinePleaAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Plea",
		ActionID: actionID,
		Duration: time.Second*15 + 1, // Add 1 to make sure the last tick takes effect
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			manaPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(sim *core.Simulation) {
					paladin.AddMana(sim, 0.05*paladin.MaxMana(), manaMetrics)
				},
			})
			if hasGlyph {
				aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.97
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			manaPA.Cancel(sim)
			if hasGlyph {
				aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.97
			}
		},
	})

	paladin.DivinePlea = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			paladin.DivinePleaAura.Activate(sim)
		},
	})
}
