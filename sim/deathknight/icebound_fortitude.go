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
			dmgTakenMult = 1.0 - core.TernaryFloat64(hasGlyph, core.MaxFloat(0.4, 0.3+0.0015*(def-400)), 0.3+0.0015*(def-400))
			dk.IceboundFortitudeAura.Unit.PseudoStats.DamageTakenMultiplier *= dmgTakenMult
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.IceboundFortitudeAura.Unit.PseudoStats.DamageTakenMultiplier /= dmgTakenMult
		},
	})

	baseCost := float64(core.NewRuneCost(20.0, 0, 0, 0, 0))
	rs := &RuneSpell{}
	dk.IceboundFortitude = dk.RegisterSpell(rs, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.IceboundFortitudeAura.Activate(sim)
			rs.DoCost(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 20.0, 0, 0, 0) && dk.IceboundFortitude.IsReady(sim)
	}, nil)

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.IceboundFortitude.Spell,
			Type:     core.CooldownTypeSurvival,
			Priority: core.CooldownPriorityLow,
		})
	}
}
