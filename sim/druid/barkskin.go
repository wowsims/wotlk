package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (druid *Druid) registerBarkskinCD() {
	if !druid.InForm(Bear) {
		return
	}

	actionId := core.ActionID{SpellID: 22812}

	setBonus := core.TernaryDuration(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 4), time.Second*3.0, 0.0)
	hasGlyph := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfBarkskin)
	cdSetBonus := core.TernaryDuration(druid.HasSetBonus(ItemSetMalfurionsBattlegear, 4), time.Second*12.0, 0.0)

	druid.BarkskinAura = druid.RegisterAura(core.Aura{
		Label:    "Barkskin",
		ActionID: actionId,
		Duration: (time.Second * 12) + setBonus,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageTakenMultiplier *= 0.8
			if hasGlyph {
				druid.PseudoStats.ReducedCritTakenChance += 0.25
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageTakenMultiplier /= 0.8
			if hasGlyph {
				druid.PseudoStats.ReducedCritTakenChance -= 0.25
			}
		},
	})

	druid.Barkskin = druid.RegisterSpell(core.SpellConfig{
		ActionID: actionId,
		Flags:    SpellFlagOmenTrigger,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: (time.Second * 60.0) - cdSetBonus,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.BarkskinAura.Activate(sim)
			druid.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Barkskin,
		Type:  core.CooldownTypeSurvival,
	})
}
