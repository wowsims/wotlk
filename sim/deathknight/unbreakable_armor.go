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

	baseCost := float64(core.NewRuneCost(10, 0, 1, 0, 0))
	dk.UnbreakableArmor = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				// TODO: does not invoke the GCD?
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.UnbreakableArmorAura.Activate(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0, 0, 1, 0) && dk.UnbreakableArmor.IsReady(sim)
	}, nil)
}
