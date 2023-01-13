package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerUnbreakableArmorSpell() {
	if !dk.Talents.UnbreakableArmor {
		return
	}

	actionID := core.ActionID{SpellID: 51271}
	cdTimer := dk.NewTimer()
	cd := time.Minute*1 - dk.thassariansPlateCooldownReduction(dk.UnbreakableArmor)

	strDep := dk.NewDynamicMultiplyStat(stats.Strength, 1.2)
	armorDep := dk.NewDynamicMultiplyStat(stats.Armor, 1.25+core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfUnbreakableArmor), 0.3, 0.0))

	dk.UnbreakableArmorAura = dk.RegisterAura(core.Aura{
		Label:    "Unbreakable Armor",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, strDep)
			aura.Unit.EnableDynamicStatDep(sim, armorDep)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, strDep)
			aura.Unit.DisableDynamicStatDep(sim, armorDep)
		},
	})

	rs := &RuneSpell{}
	dk.UnbreakableArmor = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			// No GCD
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.UnbreakableArmorAura.Activate(sim)

			if !dk.Inputs.IsDps {
				rs.DoCost(sim)
			}
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.UnbreakableArmor.Spell,
			Type:     core.CooldownTypeSurvival,
			Priority: core.CooldownPriorityDefault,
			CanActivate: func(sim *core.Simulation, character *core.Character) bool {
				return dk.UnbreakableArmor.CanCast(sim)
			},
		})
	}
}
