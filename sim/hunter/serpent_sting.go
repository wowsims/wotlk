package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	actionID := core.ActionID{SpellID: 49001}
	baseCost := 0.09 * hunter.BaseMana

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskRangedSpecial,
			ThreatMultiplier: 1,
			OutcomeApplier:   hunter.OutcomeFuncRangedHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					hunter.SerpentStingDot.Apply(sim)
				}
			},
		}),
	})

	target := hunter.CurrentTarget
	hunter.SerpentStingDot = core.NewDot(core.Dot{
		Spell: hunter.SerpentSting,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SerpentSting-" + strconv.Itoa(int(hunter.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(core.TernaryInt32(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSerpentSting), 2, 0)),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + 0.1*float64(hunter.Talents.ImprovedStings),
			ThreatMultiplier: 1,
			IsPeriodic:       true,

			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
				attackPower := spellEffect.RangedAttackPower(spell.Unit) + spellEffect.RangedAttackPowerOnTarget()
				return 242 + attackPower*0.04
			}, 0),
			OutcomeApplier: hunter.OutcomeFuncTick(),
		}),
	})

	if hunter.Talents.NoxiousStings > 0 {
		multiplier := 1 + 0.01*float64(hunter.Talents.NoxiousStings)
		hunter.SerpentStingDot.Aura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
			hunter.AttackTables[aura.Unit.TableIndex].DamageDealtMultiplier *= multiplier
		}
		hunter.SerpentStingDot.Aura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
			hunter.AttackTables[aura.Unit.TableIndex].DamageDealtMultiplier /= multiplier
		}
	}
}
