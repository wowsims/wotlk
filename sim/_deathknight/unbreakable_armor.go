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
	armorDep := dk.NewDynamicMultiplyStat(stats.Armor, core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfUnbreakableArmor), 1.3, 1.25))

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

	dk.UnbreakableArmor = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

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
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.UnbreakableArmor,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
