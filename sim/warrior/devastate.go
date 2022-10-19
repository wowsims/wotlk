package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerDevastateSpell() {
	if !warrior.Talents.Devastate {
		return
	}

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

	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfDevastate)

	cost := 15.0 - float64(warrior.Talents.FocusedRage) - float64(warrior.Talents.Puncture)
	refundAmount := cost * 0.8
	flatThreatBonus := core.TernaryFloat64(hasGlyph, 630, 315)
	dynaThreatBonus := core.TernaryFloat64(hasGlyph, 0.1, 0.05)

	normalBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 0, false)
	weaponMulti := 1.2

	warrior.Devastate = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47498},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
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

		BonusCritRating:  5*core.CritRatingPerCritChance*float64(warrior.Talents.SwordAndBoard) + core.TernaryFloat64(warrior.HasSetBonus(ItemSetSiegebreakerPlate, 2), 10*core.CritRatingPerCritChance, 0),
		DamageMultiplier: weaponMulti * core.TernaryFloat64(warrior.HasSetBonus(ItemSetWrynnsPlate, 2), 1.05, 1.0),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,
		FlatThreatBonus:  flatThreatBonus,
		DynamicThreatBonus: func(spellEffect *core.SpellEffect, spell *core.Spell) float64 {
			return dynaThreatBonus * spell.MeleeAttackPower()
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
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
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if !warrior.ExposeArmorAura.IsActive() {
						warrior.SunderArmorDevastate.Cast(sim, spellEffect.Target)
					}
				} else {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanDevastate(sim *core.Simulation) bool {
	if warrior.Devastate != nil {
		return warrior.CurrentRage() >= warrior.Devastate.DefaultCast.Cost
	} else {
		return false
	}
}
