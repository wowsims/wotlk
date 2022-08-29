package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerDevastateSpell() {
	if warrior.Talents.SwordAndBoard > 0 {
		warrior.SwordAndBoardAura = warrior.GetOrRegisterAura(core.Aura{
			Label:    "Sword And Board",
			ActionID: core.ActionID{SpellID: 46953},
			Duration: 5 * time.Second,
		})

		core.MakePermanent(warrior.GetOrRegisterAura(core.Aura{
			Label: "Sword And Board Trigger",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if !(spell == warrior.Revenge || spell == warrior.Devastate) {
					return
				}

				if sim.RandomFloat("Sword And Board") <= 0.1*float64(warrior.Talents.SwordAndBoard) {
					warrior.SwordAndBoardAura.Activate(sim)
					warrior.ShieldSlam.CD.Reset()
				}
			},
		}))
	}

	cost := 15.0 - float64(warrior.Talents.FocusedRage) - float64(warrior.Talents.Puncture)
	refundAmount := cost * 0.8

	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfDevastate)

	normalBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 0, 1.2, 1.0, true)

	warrior.Devastate = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30022},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			BonusCritRating:  5 * core.CritRatingPerCritChance * float64(warrior.Talents.SwordAndBoard),
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			FlatThreatBonus:  100,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// Bonus 242 damage / stack of sunder. Counts stacks AFTER cast but only if stacks > 0.
					sunderBonus := 0.0
					saStacks := warrior.SunderArmorAura.GetStacks()
					if saStacks != 0 {
						sunderBonus = 242 * float64(core.MinInt32(saStacks+1, 5))
					}

					return normalBaseDamage(sim, hitEffect, spell) + sunderBonus
				},
				TargetSpellCoefficient: 0,
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if !warrior.ExposeArmorAura.IsActive() {
						warrior.SunderArmorDevastate.Cast(sim, spellEffect.Target)
						if hasGlyph {
							warrior.SunderArmorDevastate.Cast(sim, spellEffect.Target)
						}
					}
				} else {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanDevastate(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.Devastate.DefaultCast.Cost
}
