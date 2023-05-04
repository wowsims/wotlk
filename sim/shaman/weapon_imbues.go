package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TotemOfTheAstralWinds int32 = 27815
var TotemOfSplintering int32 = 40710

func (shaman *Shaman) RegisterOnItemSwapWithImbue(effectID int32, procMask *core.ProcMask, aura *core.Aura) {
	shaman.RegisterOnItemSwap(func(sim *core.Simulation) {
		mh := shaman.Equip[proto.ItemSlot_ItemSlotMainHand].TempEnchant == effectID
		oh := shaman.Equip[proto.ItemSlot_ItemSlotOffHand].TempEnchant == effectID
		*procMask = core.GetMeleeProcMaskForHands(mh, oh)

		if !mh && !oh {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	apBonus := 1250.0
	if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfTheAstralWinds {
		apBonus += 80
	} else if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfSplintering {
		apBonus += 212
	}

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58804}.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMelee,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: []float64{1, 1.13, 1.27, 1.4}[shaman.Talents.ElementalWeapons],
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := spell.BonusWeaponDamage()
			mAP := spell.MeleeAttackPower() + apBonus

			if isMH {
				baseDamage1 := constBaseDamage + spell.Unit.MHWeaponDamage(sim, mAP)
				baseDamage2 := constBaseDamage + spell.Unit.MHWeaponDamage(sim, mAP)
				result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
				result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
				spell.DealDamage(sim, result1)
				spell.DealDamage(sim, result2)
			} else {
				baseDamage1 := constBaseDamage + spell.Unit.OHWeaponDamage(sim, mAP)
				baseDamage2 := constBaseDamage + spell.Unit.OHWeaponDamage(sim, mAP)
				result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
				result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
				spell.DealDamage(sim, result1)
				spell.DealDamage(sim, result2)
			}
		},
	}

	return shaman.RegisterSpell(spellConfig)
}

func (shaman *Shaman) RegisterWindfuryImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	var proc = 0.2
	if mh && oh {
		proc = 0.36
	}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfWindfuryWeapon) {
		proc += 0.02 //TODO: confirm how this actually works
	}

	mhSpell := shaman.newWindfuryImbueSpell(true)
	ohSpell := shaman.newWindfuryImbueSpell(false)

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Second * 3,
	}

	if mh {
		shaman.Equip[proto.ItemSlot_ItemSlotMainHand].TempEnchant = 3787
	}

	if oh {
		shaman.Equip[proto.ItemSlot_ItemSlotOffHand].TempEnchant = 3787
	}

	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Windfury Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			isMHHit := spell.IsMH()
			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Windfury Imbue") > proc {
				return
			}
			icd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(3787, &procMask, aura)
}

func (shaman *Shaman) newFlametongueImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58790},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weapon := core.Ternary(isMH, shaman.GetMHWeapon(), shaman.GetOHWeapon())

			var damage float64 = 0
			if weapon != nil {
				baseDamage := weapon.SwingSpeed * 68.5
				spellCoeff := 0.1 * weapon.SwingSpeed / 2.6
				damage = baseDamage + spellCoeff*spell.SpellPower()
			}

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) ApplyFlametongueImbueToItem(item *core.Item, isDownranked bool) {
	if item == nil || item.TempEnchant == 3781 || item.TempEnchant == 3780 {
		return
	}

	spBonus := 211.0
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	id := 3781
	if isDownranked {
		spBonus = 186.0
		id = 3780
	}

	newStats := stats.Stats{stats.SpellPower: spBonus * spMod}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		newStats = newStats.Add(stats.Stats{stats.SpellCrit: 2 * core.CritRatingPerCritChance})
	}

	item.Stats = item.Stats.Add(newStats)
	item.TempEnchant = int32(id)
}

func (shaman *Shaman) RegisterFlametongueImbue(mh bool, oh bool) {
	if !mh && !oh && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if mh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetMHWeapon(), false)
	}
	if oh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetOHWeapon(), false)
	}

	ftIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueImbueSpell(true)
	ohSpell := shaman.newFlametongueImbueSpell(false)

	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			isMHHit := spell.IsMH()
			if !ftIcd.IsReady(sim) {
				return
			}
			ftIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(3781, &procMask, aura)
}

func (shaman *Shaman) newFlametongueDownrankImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58789},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			weapon := core.Ternary(isMH, shaman.GetMHWeapon(), shaman.GetOHWeapon())

			var damage float64 = 0
			if weapon != nil {
				baseDamage := weapon.SwingSpeed * 64
				spellCoeff := 0.1 * weapon.SwingSpeed / 2.6
				damage = baseDamage + spellCoeff*spell.SpellPower()
			}

			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFlametongueDownrankImbue(mh bool, oh bool) {
	if !mh && !oh && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if mh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetMHWeapon(), true)
	}
	if oh {
		shaman.ApplyFlametongueImbueToItem(shaman.GetOHWeapon(), true)
	}

	ftDownrankIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueDownrankImbueSpell(true)
	ohSpell := shaman.newFlametongueDownrankImbueSpell(false)
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue (downranked)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			isMHHit := spell.IsMH()
			if !ftDownrankIcd.IsReady(sim) {
				return
			}
			ftDownrankIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(3780, &procMask, aura)
}

func (shaman *Shaman) FrostbrandDebuffAura(target *core.Unit) *core.Aura {
	multiplier := 1 + 0.05*float64(shaman.Talents.FrozenPower)
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Frostbrand Attack-" + shaman.Label,
		ActionID: core.ActionID{SpellID: 58799},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.LightningBolt.DamageMultiplier *= multiplier
			shaman.ChainLightning.DamageMultiplier *= multiplier
			if shaman.LavaLash != nil {
				shaman.LavaLash.DamageMultiplier *= multiplier
			}
			shaman.EarthShock.DamageMultiplier *= multiplier
			shaman.FlameShock.DamageMultiplier *= multiplier
			shaman.FrostShock.DamageMultiplier *= multiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.LightningBolt.DamageMultiplier /= multiplier
			shaman.ChainLightning.DamageMultiplier /= multiplier
			if shaman.LavaLash != nil {
				shaman.LavaLash.DamageMultiplier /= multiplier
			}
			shaman.EarthShock.DamageMultiplier /= multiplier
			shaman.FlameShock.DamageMultiplier /= multiplier
			shaman.FrostShock.DamageMultiplier /= multiplier
		},
	})
}

func (shaman *Shaman) newFrostbrandImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58796},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 530 + 0.1*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFrostbrandImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	mhSpell := shaman.newFrostbrandImbueSpell(true)
	ohSpell := shaman.newFrostbrandImbueSpell(false)
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	ppmm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	if mh {
		shaman.Equip[proto.ItemSlot_ItemSlotMainHand].TempEnchant = 3784
	} else {
		shaman.Equip[proto.ItemSlot_ItemSlotOffHand].TempEnchant = 3784
	}

	fbDebuffAuras := shaman.NewEnemyAuraArray(shaman.FrostbrandDebuffAura)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if !ppmm.Proc(sim, spell.ProcMask, "Frostbrand Weapon") {
				return
			}

			if spell.IsMH() {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
			fbDebuffAuras.Get(result.Target).Activate(sim)
		},
	})

	shaman.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3784, 9.0, &ppmm, aura)
}

func (shaman *Shaman) newEarthlivingImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51994},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Earthliving",
				ActionID: core.ActionID{SpellID: 52000},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 280 + 0.171*dot.Spell.HealingPower(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
			spell.Hot(target).Apply(sim)
		},
	})
}

func (shaman *Shaman) ApplyEarthlivingImbueToItem(item *core.Item, isDownranked bool) {
	if item == nil || item.TempEnchant == 3350 || item.TempEnchant == 3349 {
		// downranking not implemented yet but put the temp enchant ID there.
		return
	}

	spBonus := 150.0
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	id := 3350

	newStats := stats.Stats{stats.SpellPower: spBonus * spMod}
	item.Stats = item.Stats.Add(newStats)
	item.TempEnchant = int32(id)
}

func (shaman *Shaman) RegisterEarthlivingImbue(mh bool, oh bool) {
	if !mh && !oh && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if mh {
		shaman.ApplyEarthlivingImbueToItem(shaman.GetMHWeapon(), false)
	}
	if oh {
		shaman.ApplyEarthlivingImbueToItem(shaman.GetOHWeapon(), false)
	}

	mhSpell := shaman.newEarthlivingImbueSpell(true)
	ohSpell := shaman.newEarthlivingImbueSpell(false)

	procChance := 0.2
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfEarthlivingWeapon) {
		procChance += 0.05
	}
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	aura := shaman.RegisterAura(core.Aura{
		Label:    "Earthliving Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != shaman.ChainHeal && spell != shaman.LesserHealingWave && spell != shaman.HealingWave && spell != shaman.Riptide {
				return
			}
			if mhSpell != nil && sim.RandomFloat("earthliving") < procChance {
				mhSpell.Cast(sim, result.Target)
			} else if ohSpell != nil && sim.RandomFloat("earthliving") < procChance {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(3350, &procMask, aura)
}
