package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}
var TotemOfTheDancingFlame int32 = 45169
var TotemOfDueling int32 = 40322

func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:     "Stormstrike-" + shaman.Label,
		ActionID:  StormstrikeActionID,
		Duration:  time.Second * 12,
		MaxStacks: 4,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AttackTables[aura.Unit.UnitIndex].NatureDamageTakenMultiplier *= core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfStormstrike), 1.28, 1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AttackTables[aura.Unit.UnitIndex].NatureDamageTakenMultiplier /= core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfStormstrike), 1.28, 1.2)

		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit != &shaman.Unit {
				return
			}
			if spell.SpellSchool != core.SpellSchoolNature {
				return
			}
			if !result.Landed() || result.Damage == 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) func(*core.Simulation, *core.Unit, *core.Spell) {
	var flatDamageBonus float64 = 0
	if shaman.Ranged().ID == TotemOfTheDancingFlame {
		flatDamageBonus += 155
	}

	var procMask core.ProcMask
	if isMH {
		procMask = core.ProcMaskMeleeMHSpecial
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		var baseDamage float64
		spell.ProcMask = procMask
		if isMH {
			baseDamage = flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
		} else {
			baseDamage = flatDamageBonus +
				spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()
		}

		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
	}
}

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	ssDebuffAuras := shaman.NewEnemyAuraArray(shaman.StormstrikeDebuffAura)

	var skyshatterAura *core.Aura
	if shaman.HasSetBonus(ItemSetSkyshatterHarness, 4) {
		skyshatterAura = shaman.NewTemporaryStatsAura("Skyshatter 4pc AP Bonus", core.ActionID{SpellID: 38432}, stats.Stats{stats.AttackPower: 70}, time.Second*12)
	}
	var totemOfDuelingAura *core.Aura
	if shaman.Ranged().ID == TotemOfDueling {
		totemOfDuelingAura = shaman.NewTemporaryStatsAura("Essense of the Storm", core.ActionID{SpellID: 60766},
			stats.Stats{stats.MeleeHaste: 60, stats.SpellHaste: 60}, time.Second*6)
	}

	manaMetrics := shaman.NewManaMetrics(core.ActionID{SpellID: 51522})

	cooldownTime := time.Duration(core.TernaryFloat64(shaman.HasSetBonus(ItemSetGladiatorsEarthshaker, 4), 6, 8))
	impSSChance := 0.5 * float64(shaman.Talents.ImprovedStormstrike)

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    StormstrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * cooldownTime,
			},
		},

		ThreatMultiplier: 1,
		DamageMultiplier: core.TernaryFloat64(shaman.HasSetBonus(ItemSetWorldbreakerBattlegear, 2), 1.2, 1),
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				if impSSChance > 0 && sim.RandomFloat("Improved Stormstrike") < impSSChance {
					shaman.AddMana(sim, 0.2*shaman.BaseMana, manaMetrics)
				}
				ssDebuffAura := ssDebuffAuras.Get(target)
				ssDebuffAura.Activate(sim)
				ssDebuffAura.SetStacks(sim, 4)

				if skyshatterAura != nil {
					skyshatterAura.Activate(sim)
				}
				if totemOfDuelingAura != nil {
					totemOfDuelingAura.Activate(sim)
				}

				if shaman.HasMHWeapon() {
					mhHit(sim, target, spell)
				}

				if shaman.AutoAttacks.IsDualWielding && shaman.HasOHWeapon() {
					ohHit(sim, target, spell)
				}

				shaman.Stormstrike.SpellMetrics[target.UnitIndex].Hits--
			}
			spell.DealOutcome(sim, result)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return shaman.HasMHWeapon() || shaman.HasOHWeapon()
		},
	})
}
