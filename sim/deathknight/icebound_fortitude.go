package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerIceboundFortitudeSpell() {
	actionID := core.ActionID{SpellID: 48792}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 2

	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIceboundFortitude)

	dmgTakenMult := 1.0
	dk.IceboundFortitudeAura = dk.RegisterAura(core.Aura{
		Label:    "Icebound Fortitude",
		ActionID: actionID,
		Duration: time.Second*12 + time.Second*2*time.Duration(float64(dk.Talents.GuileOfGorefiend)) + dk.scourgebornePlateIFDurationBonus(),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			def := dk.IceboundFortitudeAura.Unit.GetStat(stats.Defense)
			dmgTakenMult = 1.0 - core.TernaryFloat64(hasGlyph, max(0.4, 0.3+0.0015*(def-400)), 0.3+0.0015*(def-400))
			dk.IceboundFortitudeAura.Unit.PseudoStats.DamageTakenMultiplier *= dmgTakenMult
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.IceboundFortitudeAura.Unit.PseudoStats.DamageTakenMultiplier /= dmgTakenMult
		},
	})

	dk.IceboundFortitude = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 20,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.IceboundFortitudeAura.Activate(sim)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.IceboundFortitude,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
