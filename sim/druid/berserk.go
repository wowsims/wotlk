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

	druid.Berserk = druid.RegisterSpell(core.SpellConfig{
		ActionID: actionId,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 180.0,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if druid.InForm(Cat) {
					cast.GCD = time.Second
				} else {
					cast.GCD = core.GCDDefault
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.BerserkAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Berserk,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return druid.InForm(Cat | Bear)
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			//druid.TigersFury.CD.TimeToReady(sim) > (druid.BerserkAura.Duration)
			// Manually handled in Feral Rotation
			return false
		},
	})

}
