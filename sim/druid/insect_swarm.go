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

	druid.InsectSwarm = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		DamageMultiplier: 1 +
			druid.talentBonuses.genesis +
			core.TernaryFloat64(druid.setBonuses.balance_t7_2, 0.1, 0) +
			core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfInsectSwarm), 0.3, 0),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: druid.OutcomeFuncMagicHit(),
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
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				druid.Wrath.DamageMultiplier *= 1 + 0.01*float64(druid.Talents.ImprovedInsectSwarm)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				druid.Wrath.DamageMultiplier /= 1 + 0.01*float64(druid.Talents.ImprovedInsectSwarm)
			},
		}),
		NumberOfTicks: 6 + druid.talentBonuses.naturesSplendor,
		TickLength:    time.Second * 2,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = 215 + 0.2*dot.Spell.SpellPower()
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

			if sim.RandomFloat("Elune's Wrath proc") > (1-0.08) && druid.setBonuses.balance_t8_4 {
				tierProc := druid.GetOrRegisterAura(core.Aura{
					Label:    "Elune's Wrath",
					ActionID: core.ActionID{SpellID: 64823},
					Duration: time.Second * 10,
				})
				tierProc.Activate(sim)
			}
		},
	})
}
