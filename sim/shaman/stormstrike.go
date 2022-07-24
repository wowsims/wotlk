package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}

func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:     "Stormstrike-" + shaman.Label,
		ActionID:  StormstrikeActionID,
		Duration:  time.Second * 12,
		MaxStacks: 4,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.PseudoStats.NatureDamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.PseudoStats.NatureDamageDealtMultiplier /= 1.2
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell.Unit != &shaman.Unit {
				return
			}
			if spell.SpellSchool != core.SpellSchoolNature {
				return
			}
			if !spellEffect.Landed() || spellEffect.Damage == 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	effect := core.SpellEffect{
		DamageMultiplier: 1,
		ThreatMultiplier: core.TernaryFloat64(shaman.Talents.SpiritWeapons, 0.7, 1),
		OutcomeApplier:   shaman.OutcomeFuncMeleeSpecialCritOnly(shaman.DefaultMeleeCritMultiplier()),
	}

	flatDamageBonus := core.TernaryFloat64(shaman.HasSetBonus(ItemSetCycloneHarness, 4), 30, 0)
	if isMH {
		effect.ProcMask = core.ProcMaskMeleeMHSpecial
		effect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.MainHand, false, flatDamageBonus, 1, true)
	} else {
		effect.ProcMask = core.ProcMaskMeleeOHSpecial
		effect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.OffHand, false, flatDamageBonus, 1, true)
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    StormstrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	baseCost := baseMana * 0.08
	if shaman.Equip[items.ItemSlotRanged].ID == StormfuryTotem {
		baseCost -= 22
	}

	ssDebuffAura := shaman.StormstrikeDebuffAura(shaman.CurrentTarget)

	var skyshatterAura *core.Aura
	if shaman.HasSetBonus(ItemSetSkyshatterHarness, 4) {
		skyshatterAura = shaman.NewTemporaryStatsAura("Skyshatter 4pc AP Bonus", core.ActionID{SpellID: 38432}, stats.Stats{stats.AttackPower: 70}, time.Second*12)
	}

	manaMetrics := shaman.NewManaMetrics(core.ActionID{SpellID: 51522})

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    StormstrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			ThreatMultiplier: 1,

			OutcomeApplier: shaman.OutcomeFuncMeleeSpecialHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}

				if shaman.Talents.ImprovedStormstrike > 0 {
					shaman.AddMana(sim, baseMana*0.2, manaMetrics, true)
				}
				ssDebuffAura.Activate(sim)
				ssDebuffAura.SetStacks(sim, 4)

				if skyshatterAura != nil {
					skyshatterAura.Activate(sim)
				}

				mhHit.Cast(sim, spellEffect.Target)
				ohHit.Cast(sim, spellEffect.Target)
				shaman.Stormstrike.SpellMetrics[spellEffect.Target.TableIndex].Casts -= 2
				shaman.Stormstrike.SpellMetrics[spellEffect.Target.TableIndex].Hits--
			},
		}),
	})
}
