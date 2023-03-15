package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (rogue *Rogue) registerHemorrhageSpell() {
	actionID := core.ActionID{SpellID: 48660}

	var numPlayers int
	for _, u := range rogue.Env.Raid.AllUnits {
		if u.Type == core.PlayerUnit {
			numPlayers++
		}
	}

	var hemoAuras core.AuraArray

	if numPlayers >= 2 {
		bonusDamage := 75.0
		if rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfHemorrhage) {
			bonusDamage *= 1.4
		}

		hemoAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:     "Hemorrhage",
				ActionID:  actionID,
				Duration:  time.Second * 15,
				MaxStacks: 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.BonusPhysicalDamageTaken += bonusDamage
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.BonusPhysicalDamageTaken -= bonusDamage
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellSchool != core.SpellSchoolPhysical {
						return
					}
					if !result.Landed() || result.Damage == 0 {
						return
					}

					aura.RemoveStack(sim)
				},
			})
		})
	}

	daggerMH := rogue.Equip[proto.ItemSlot_ItemSlotMainHand].WeaponType == proto.WeaponType_WeaponTypeDagger
	rogue.Hemorrhage = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(35 - float64(rogue.Talents.SlaughterFromTheShadows)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance,

		DamageMultiplier: core.TernaryFloat64(daggerMH, 1.6, 1.1) * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0)) *
			(1 + 0.02*float64(rogue.Talents.SinisterCalling)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				if len(hemoAuras) > 0 {
					hemoAura := hemoAuras.Get(target)
					hemoAura.Activate(sim)
					hemoAura.SetStacks(sim, 10)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
