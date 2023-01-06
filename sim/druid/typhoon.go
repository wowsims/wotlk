package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerTyphoonSpell() {
	if !druid.Talents.Typhoon {
		return
	}

	actionID := core.ActionID{SpellID: 61384}
	baseCost := 0.25 * druid.BaseMana
	spellCoeff := 0.193

	druid.Typhoon = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Flags:        SpellFlagOmenTrigger,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (20 - core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMonsoon), 3, 0)),
			},
		},

		DamageMultiplier: 1 +
			0.15*float64(druid.Talents.GaleWinds),
		ThreatMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				baseDamage := 1190 + (spellCoeff * spell.SpellPower())
				baseDamage *= sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.Targets {
					spell.CalcAndDealDamage(sim, &aoeTarget.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			})
		},
	})
}
