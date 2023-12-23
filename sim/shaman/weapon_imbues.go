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
		mask := core.ProcMaskUnknown
		if shaman.MainHand().TempEnchant == effectID {
			mask |= core.ProcMaskMeleeMH
		}
		if shaman.OffHand().TempEnchant == effectID {
			mask |= core.ProcMaskMeleeOH
		}
		*procMask = mask

		if mask == core.ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	apBonus := 1250.0
	if shaman.Ranged().ID == TotemOfTheAstralWinds {
		apBonus += 80
	} else if shaman.Ranged().ID == TotemOfSplintering {
		apBonus += 212
	}

	tag := 1
	procMask := core.ProcMaskMeleeMHSpecial
	weaponDamageFunc := shaman.MHWeaponDamage
	if !isMH {
		tag = 2
		procMask = core.ProcMaskMeleeOHSpecial
		weaponDamageFunc = shaman.OHWeaponDamage
		apBonus *= 2 // applied after 50% offhand penalty
	}

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58804, Tag: int32(tag)},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: []float64{1, 1.13, 1.27, 1.4}[shaman.Talents.ElementalWeapons],
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := spell.BonusWeaponDamage()
			mAP := spell.MeleeAttackPower() + apBonus

			baseDamage1 := constBaseDamage + weaponDamageFunc(sim, mAP)
			baseDamage2 := constBaseDamage + weaponDamageFunc(sim, mAP)
			result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
			result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DealDamage(sim, result1)
			spell.DealDamage(sim, result2)
		},
	}

	return shaman.RegisterSpell(spellConfig)
}

func (shaman *Shaman) RegisterWindfuryImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = 3787
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = 3787
	}

	var proc = 0.2
	if procMask == core.ProcMaskMelee {
		proc = 0.36
	}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfWindfuryWeapon) {
		proc += 0.02 //TODO: confirm how this actually works
	}

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Second * 3,
	}

	mhSpell := shaman.newWindfuryImbueSpell(true)
	ohSpell := shaman.newWindfuryImbueSpell(false)

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

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Windfury Imbue") < proc {
				icd.Use(sim)

				if spell.IsMH() {
					mhSpell.Cast(sim, result.Target)
				} else {
					ohSpell.Cast(sim, result.Target)
				}
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(3787, &procMask, aura)
}

func (shaman *Shaman) newFlametongueImbueSpell(weapon *core.Item, isDownranked bool) *core.Spell {
	spellID := 58790
	baseDamage := 68.5
	if isDownranked {
		spellID = 58789
		baseDamage = 64
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(spellID)},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskWeaponProc,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if weapon.SwingSpeed != 0 {
				damage := weapon.SwingSpeed * (baseDamage + 0.1/2.6*spell.SpellPower())
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (shaman *Shaman) ApplyFlametongueImbueToItem(item *core.Item, isDownranked bool) {
	if item == nil || item.TempEnchant == 3781 || item.TempEnchant == 3780 {
		return
	}

	spBonus := 211.0
	enchantID := 3781
	if isDownranked {
		spBonus = 186.0
		enchantID = 3780
	}

	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)

	newStats := stats.Stats{stats.SpellPower: spBonus * spMod}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		newStats = newStats.Add(stats.Stats{stats.SpellCrit: 2 * core.CritRatingPerCritChance})
	}

	item.Stats = item.Stats.Add(newStats)
	item.TempEnchant = int32(enchantID)
}

func (shaman *Shaman) ApplyFlametongueImbue(procMask core.ProcMask, isDownranked bool) {
	if procMask.Matches(core.ProcMaskMeleeMH) && shaman.HasMHWeapon() {
		shaman.ApplyFlametongueImbueToItem(shaman.MainHand(), isDownranked)
	}

	if procMask.Matches(core.ProcMaskMeleeOH) && shaman.HasOHWeapon() {
		shaman.ApplyFlametongueImbueToItem(shaman.OffHand(), isDownranked)
	}
}

func (shaman *Shaman) RegisterFlametongueImbue(procMask core.ProcMask, isDownranked bool) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueImbueSpell(shaman.MainHand(), isDownranked)
	ohSpell := shaman.newFlametongueImbueSpell(shaman.OffHand(), isDownranked)

	label := "Flametongue Imbue"
	enchantID := 3781
	if isDownranked {
		label = "Flametongue Imbue (downranked)"
		enchantID = 3780
	}

	aura := shaman.RegisterAura(core.Aura{
		Label:    label,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)

			if spell.IsMH() {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(int32(enchantID), &procMask, aura)
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

func (shaman *Shaman) newFrostbrandImbueSpell() *core.Spell {
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

func (shaman *Shaman) RegisterFrostbrandImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = 3784
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = 3784
	}

	ppmm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	mhSpell := shaman.newFrostbrandImbueSpell()
	ohSpell := shaman.newFrostbrandImbueSpell()

	fbDebuffAuras := shaman.NewEnemyAuraArray(shaman.FrostbrandDebuffAura)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if ppmm.Proc(sim, spell.ProcMask, "Frostbrand Weapon") {
				if spell.IsMH() {
					mhSpell.Cast(sim, result.Target)
				} else {
					ohSpell.Cast(sim, result.Target)
				}
				fbDebuffAuras.Get(result.Target).Activate(sim)
			}
		},
	})

	shaman.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(3784, 9.0, &ppmm, aura)
}

func (shaman *Shaman) newEarthlivingImbueSpell() *core.Spell {
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

func (shaman *Shaman) ApplyEarthlivingImbueToItem(item *core.Item) {
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

func (shaman *Shaman) RegisterEarthlivingImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskEmpty && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.ApplyEarthlivingImbueToItem(shaman.MainHand())
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.ApplyEarthlivingImbueToItem(shaman.OffHand())
	}

	procChance := 0.2
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfEarthlivingWeapon) {
		procChance += 0.05
	}

	imbueSpell := shaman.newEarthlivingImbueSpell()

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

			if procMask.Matches(core.ProcMaskMeleeMH) && sim.RandomFloat("earthliving") < procChance {
				imbueSpell.Cast(sim, result.Target)
			}

			if procMask.Matches(core.ProcMaskMeleeOH) && sim.RandomFloat("earthliving") < procChance {
				imbueSpell.Cast(sim, result.Target)
			}
		},
	})

	shaman.RegisterOnItemSwapWithImbue(3350, &procMask, aura)
}
