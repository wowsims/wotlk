package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (druid *Druid) registerBerserkCD() {
	if !druid.Talents.Berserk {
		return
	}

	actionId := core.ActionID{SpellID: 50334}

	glyphBonus := core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfBerserk), time.Second*5.0, 0.0)

	druid.BerserkAura = druid.RegisterAura(core.Aura{
		Label:    "Berserk",
		ActionID: actionId,
		Duration: (time.Second * 15) + glyphBonus,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if druid.InForm(Cat) {
				druid.PseudoStats.CostMultiplier /= 2.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if druid.InForm(Cat) {
				druid.PseudoStats.CostMultiplier *= 2.0
			}
		},
	})

	spell := druid.RegisterSpell(core.SpellConfig{
		ActionID: actionId,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 180.0,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.BerserkAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})

}
