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

	m := (50.0 - 35.0) / (680.0 - 540.0)
	dtMultiplier := 1.0
	dk.IceboundFortitudeAura = dk.RegisterAura(core.Aura{
		Label:    "Icebound Fortitude",
		ActionID: actionID,
		Duration: time.Second*12 + time.Second*2*time.Duration(float64(dk.Talents.GuileOfGorefiend)),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: Verify formula
			defRating := dk.IceboundFortitudeAura.Unit.GetStat(stats.Defense) * core.DefenseRatingPerDefense
			if defRating <= 306 {
				dtMultiplier = 1.1
			} else {
				dtMultiplier = 1.0 + (m*defRating - 22.8571428556)
			}

			if dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfIceboundFortitude) {
				dtMultiplier = core.MaxFloat(dtMultiplier, 1.4)
			}

			dk.IceboundFortitudeAura.Unit.PseudoStats.DamageTakenMultiplier *= dtMultiplier
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.IceboundFortitudeAura.Unit.PseudoStats.DamageTakenMultiplier /= dtMultiplier
		},
	})

	baseCost := float64(core.NewRuneCost(20.0, 0, 0, 0, 0))

	dk.IceboundFortitude = dk.RegisterSpell(nil, core.SpellConfig{
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
		},
	})
}

func (dk *Deathknight) CanIceboundFortitude(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 20.0, 0, 0, 0) && dk.IceboundFortitude.IsReady(sim)
}

func (dk *Deathknight) CastIceboundFortitude(sim *core.Simulation, target *core.Unit) bool {
	if dk.IceboundFortitude.IsReady(sim) {
		dk.LastCast = dk.IceboundFortitude
		return dk.IceboundFortitude.Cast(sim, target)
	}
	return false
}
