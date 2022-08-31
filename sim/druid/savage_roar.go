package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) getSavageRoarMultiplier() float64 {
	glyphBonus := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfSavageRoar), 0.03, 0)
	return 1.3 + glyphBonus
}

func (druid *Druid) registerSavageRoarSpell() {
	actionID := core.ActionID{SpellID: 52610}
	baseCost := 25.0

	srm := druid.getSavageRoarMultiplier()
	durationBonus := core.TernaryDuration(druid.HasSetBonus(ItemSetNightsongBattlegear, 4), time.Second*8.0, 0.0)

	druid.SavageRoarAura = druid.RegisterAura(core.Aura{
		Label:    "Savage Roar Aura",
		ActionID: actionID,
		Duration: 9,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.PhysicalDamageDealtMultiplier *= srm
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.PhysicalDamageDealtMultiplier /= srm
		},
	})

	srSpell := druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Energy,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.SavageRoarAura.Duration = durationBonus + time.Duration(float64(time.Second)*(9.0+(5.0)*float64(druid.ComboPoints())))
			druid.SavageRoarAura.Activate(sim)
			druid.SpendComboPoints(sim, spell.ComboPointMetrics())
		},
	})

	druid.SavageRoar = srSpell
}

func (druid *Druid) CanSavageRoar() bool {
	return druid.InForm(Cat) && druid.ComboPoints() > 0 && (druid.CurrentEnergy() >= druid.CurrentSavageRoarCost())
}

func (druid *Druid) CurrentSavageRoarCost() float64 {
	return druid.SavageRoar.ApplyCostModifiers(druid.SavageRoar.BaseCost)
}
