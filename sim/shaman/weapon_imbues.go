package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TotemOfTheAstralWinds int32 = 27815
var TotemOfSplintering int32 = 40710

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	apBonus := 1250.0
	if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfTheAstralWinds {
		apBonus += 80
	} else if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfSplintering {
		apBonus += 212
	}

	actionID := core.ActionID{SpellID: 58804}

	baseEffect := core.SpellEffect{
		BonusAttackPower: apBonus,
		ProcMask:         core.ProcMaskMelee,
		DamageMultiplier: 1.0,
		ThreatMultiplier: 1,
		OutcomeApplier:   shaman.OutcomeFuncMeleeSpecialHitAndCrit(shaman.DefaultMeleeCritMultiplier()),
	}

	weaponDamageMultiplier := 1 + math.Round(float64(shaman.Talents.ElementalWeapons)*13.33)/100
	if isMH {
		actionID.Tag = 1
		baseEffect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, 1, weaponDamageMultiplier, true)
	} else {
		actionID.Tag = 2
		baseEffect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.OffHand, false, 0, 1, weaponDamageMultiplier, true)

		// For whatever reason, OH penalty does not apply to the bonus AP from WF OH
		// hits. Implement this by doubling the AP bonus we provide.
		baseEffect.BonusAttackPower += apBonus
	}

	effects := []core.SpellEffect{
		baseEffect,
		baseEffect,
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDamageMultipleTargeted(effects),
	})
}

func (shaman *Shaman) ApplyWindfuryImbue(mh bool, oh bool) {
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

	shaman.RegisterAura(core.Aura{
		Label:    "Windfury Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// ProcMask: 20
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			isMHHit := spellEffect.IsMH()
			if (!mh && isMHHit) || (!oh && !isMHHit) {
				return // cant proc if not enchanted
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Windfury Imbue") > proc {
				return
			}
			icd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

func (shaman *Shaman) newFlametongueImbueSpell(isMH bool) *core.Spell {
	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskEmpty,
		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		OutcomeApplier:   shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
	}

	if isMH {
		if weapon := shaman.GetMHWeapon(); weapon != nil {
			baseDamage := weapon.SwingSpeed * 68.5
			effect.BaseDamage = core.BaseDamageConfigMagic(baseDamage, baseDamage, (0.1 / 2.6 * weapon.SwingSpeed))
		}
	} else {
		if weapon := shaman.GetOHWeapon(); weapon != nil {
			baseDamage := weapon.SwingSpeed * 68.5
			effect.BaseDamage = core.BaseDamageConfigMagic(baseDamage, baseDamage, (0.1 / 2.6 * weapon.SwingSpeed))
		}
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 58790},
		SpellSchool:  core.SpellSchoolFire,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (shaman *Shaman) ApplyFlametongueImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	imbueCount := 1.0
	spBonus := 211.0
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	if mh && oh { // grant double SP+Crit bonuses for ft/ft (possible bug, but currently working on beta, its unclear)
		imbueCount += 1.0
	}
	shaman.AddStat(stats.SpellPower, spBonus*spMod*imbueCount)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		shaman.AddStat(stats.SpellCrit, 2*core.CritRatingPerCritChance*imbueCount)
	}

	ftIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueImbueSpell(true)
	ohSpell := shaman.newFlametongueImbueSpell(false)

	shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			isMHHit := spellEffect.IsMH()
			if (isMHHit && !mh) || (!isMHHit && !oh) {
				return // cant proc if not enchanted
			}
			if !ftIcd.IsReady(sim) {
				return
			}
			ftIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

func (shaman *Shaman) newFlametongueDownrankImbueSpell(isMH bool) *core.Spell {
	effect := core.SpellEffect{
		ProcMask:         core.ProcMaskEmpty,
		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		OutcomeApplier:   shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
	}

	if isMH {
		if weapon := shaman.GetMHWeapon(); weapon != nil {
			baseDamage := weapon.SwingSpeed * 64
			effect.BaseDamage = core.BaseDamageConfigMagic(baseDamage, baseDamage, (0.1 / 2.6 * weapon.SwingSpeed))
		}
	} else {
		if weapon := shaman.GetOHWeapon(); weapon != nil {
			baseDamage := weapon.SwingSpeed * 64
			effect.BaseDamage = core.BaseDamageConfigMagic(baseDamage, baseDamage, (0.1 / 2.6 * weapon.SwingSpeed))
		}
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 58789},
		SpellSchool:  core.SpellSchoolFire,
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (shaman *Shaman) ApplyFlametongueDownrankImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	imbueCount := 1.0
	spBonus := 186.0
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	if mh && oh { // grant double SP+Crit bonuses for ft/ft (possible bug, but currently working on beta, its unclear)
		imbueCount += 1.0
	}
	shaman.AddStat(stats.SpellPower, spBonus*spMod*imbueCount)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		shaman.AddStat(stats.SpellCrit, 2*core.CritRatingPerCritChance*imbueCount)
	}

	ftDownrankIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueDownrankImbueSpell(true)
	ohSpell := shaman.newFlametongueDownrankImbueSpell(false)

	shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue (downranked)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			isMHHit := spellEffect.IsMH()
			if (isMHHit && !mh) || (!isMHHit && !oh) {
				return // cant proc if not enchanted
			}
			if !ftDownrankIcd.IsReady(sim) {
				return
			}
			ftDownrankIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

func (shaman *Shaman) newFrostbrandImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58796},
		SpellSchool: core.SpellSchoolFrost,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskEmpty,
			BonusHitRating: float64(shaman.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMagic(530, 530, 0.1),
			OutcomeApplier: shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
		}),
	})
}

func (shaman *Shaman) ApplyFrostbrandImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	mhSpell := shaman.newFrostbrandImbueSpell(true)
	ohSpell := shaman.newFrostbrandImbueSpell(false)
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	ppmm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(procMask) {
				return
			}

			if !ppmm.Proc(sim, spellEffect.ProcMask, "Frostbrand Weapon") {
				return
			}

			if spellEffect.IsMH() {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

//earthliving? not important for dps sims though
