package paladin

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) registerDivinePleaSpell() {
	// Currently implemented as a self-targeted DoT that restores mana.
	// In future maybe expose aura for buff asserting (For prot paladins.)

	actionID := core.ActionID{SpellID: 54428} // Divine plea
	hasGlyph := paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivinePlea)

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
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if hasGlyph {
					aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.97
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if hasGlyph {
					aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.97
				}
			},
		}),

		NumberOfTicks: 5,
		TickLength:    time.Second * 3,

		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if paladin.PleaManaMetrics == nil {
				paladin.PleaManaMetrics = paladin.NewManaMetrics(actionID)
			}
			paladin.AddMana(sim, paladin.MaxMana()*0.05, paladin.PleaManaMetrics, false)
		},
	})

	plea.Spell.ApplyEffects = func(sim *core.Simulation, unit *core.Unit, _ *core.Spell) {
		plea.Activate(sim)
	}

	paladin.DivinePleaAura = plea.Aura
	paladin.DivinePlea = plea.Spell
}
