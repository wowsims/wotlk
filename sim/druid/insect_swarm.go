package druid

import (
	"github.com/wowsims/wotlk/sim/core/proto"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerInsectSwarmSpell() {
	actionID := core.ActionID{SpellID: 48468}
	baseCost := 0.08 * druid.BaseMana

	target := druid.CurrentTarget
	missAura := core.InsectSwarmAura(target)

	// T7-2P
	dreamwalkerGrab := core.TernaryFloat64(druid.SetBonuses.balance_t7_2, 1.1, 1.0)

	// Glyph of Insect Swarm
	glyphMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfInsectSwarm), 1.3, 1.0)

	druid.InsectSwarm = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskSpellDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			OutcomeApplier:   druid.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.InsectSwarmDot.Apply(sim)
					if !druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfInsectSwarm) {
						missAura.Activate(sim)
					}
				}
			},
		}),
	})

	druid.InsectSwarmDot = core.NewDot(core.Dot{
		Spell: druid.InsectSwarm,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Insect Swarm",
			ActionID: actionID,
		}),
		NumberOfTicks: 6 + druid.TalentsBonuses.naturesSplendorTick,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * druid.TalentsBonuses.genesisMultiplier * dreamwalkerGrab * glyphMultiplier,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(215, 0.2),
			OutcomeApplier:   druid.OutcomeFuncTick(),
			OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if sim.RandomFloat("Elune's Wrath proc") > (1-0.08) && druid.SetBonuses.balance_t8_4 {
					tierProc := druid.GetOrRegisterAura(core.Aura{
						Label:    "Elune's Wrath",
						ActionID: core.ActionID{SpellID: 64823},
						Duration: time.Second * 10,
					})
					tierProc.Activate(sim)
				}
			},
		}),
	})
}
