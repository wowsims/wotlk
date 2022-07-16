package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (paladin *Paladin) registerDivinePleaSpell() {
	// Currently implemented as a self-targeted DoT that restores mana.
	// In future maybe expose aura for buff asserting (For prot paladins.)

	actionID := core.ActionID{SpellID: 54428} // Divine plea

	plea := core.NewDot(core.Dot{
		Spell: paladin.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolHoly,
			Cast: core.CastConfig{
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
				CD: core.Cooldown{
					Timer:    paladin.NewTimer(),
					Duration: time.Minute * 1,
				},
			},
		}),
		Aura: paladin.RegisterAura(core.Aura{
			Label:    "Divine Plea-" + strconv.Itoa(int(paladin.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks: 5,
		TickLength:    time.Second * 3,

		TickEffects: func(sim *core.Simulation, _ *core.Spell) func() {
			return func() {
				if paladin.PleaManaMetrics == nil {
					paladin.PleaManaMetrics = paladin.NewManaMetrics(actionID)
				}
				paladin.AddMana(sim, paladin.MaxMana()*0.05, paladin.PleaManaMetrics, false)
			}
		},
	})

	plea.Spell.ApplyEffects = func(sim *core.Simulation, unit *core.Unit, _ *core.Spell) {
		plea.Activate(sim)
	}

	paladin.DivinePlea = plea.Spell
}
