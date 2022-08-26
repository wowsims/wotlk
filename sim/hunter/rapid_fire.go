package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerRapidFireCD() {
	actionID := core.ActionID{SpellID: 3045}

	var manaMetrics *core.ResourceMetrics
	if hunter.Talents.RapidRecuperation > 0 {
		manaMetrics = hunter.NewManaMetrics(core.ActionID{SpellID: 53232})
	}

	hasteMultiplier := 1.4 + core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfRapidFire), 0.08, 0)

	hunter.RapidFireAura = hunter.RegisterAura(core.Aura{
		Label:    "Rapid Fire",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.RangedSpeedMultiplier *= hasteMultiplier

			if manaMetrics != nil {
				manaPerTick := 0.02 * float64(hunter.Talents.RapidRecuperation) * hunter.MaxMana()
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 3,
					NumTicks: 5,
					OnAction: func(sim *core.Simulation) {
						hunter.AddMana(sim, manaPerTick, manaMetrics, false)
					},
				})
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.RangedSpeedMultiplier /= hasteMultiplier
		},
	})

	baseCost := 0.03 * hunter.BaseMana
	hunter.RapidFire = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute*5 - time.Minute*time.Duration(hunter.Talents.RapidKilling),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.RapidFireAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.RapidFire,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Make sure we don't reuse after a Readiness cast.
			return !hunter.RapidFireAura.IsActive()
		},
	})
}
