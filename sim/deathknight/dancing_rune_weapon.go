package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerDancingRuneWeaponCD() {
	if !dk.Talents.DancingRuneWeapon {
		return
	}

	duration := time.Second * 12
	if dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDancingRuneWeapon) {
		duration += time.Second * 5
	}

	dancingRuneWeaponAura := dk.RegisterAura(core.Aura{
		Label:    "Dancing Rune Weapon",
		ActionID: core.ActionID{SpellID: 49028},
		Duration: duration,
		// Casts
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			switch spell {
			case dk.IcyTouch.Spell:
				dk.RuneWeapon.IcyTouch.Cast(sim, spell.Unit.CurrentTarget)
			case dk.PlagueStrike.Spell:
				dk.RuneWeapon.PlagueStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.DeathStrike.Spell:
				dk.RuneWeapon.DeathStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.HeartStrike.Spell:
				dk.RuneWeapon.HeartStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.DeathCoil.Spell:
				dk.RuneWeapon.DeathCoil.Cast(sim, spell.Unit.CurrentTarget)
			}
		},
	})

	baseCost := float64(core.NewRuneCost(60.0, 0, 0, 0, 0))
	dk.DancingRuneWeapon = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49028},

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.RuneWeapon.EnableWithTimeout(sim, dk.Gargoyle, duration)
			dk.RuneWeapon.CancelGCDTimer(sim)

			dk.RuneWeapon.PseudoStats.DamageDealtMultiplier = 0.5

			// What if?
			//dk.RuneWeapon.PseudoStats.DamageDealtMultiplier = 0.5 * dk.PseudoStats.DamageDealtMultiplier

			//dk.RuneWeapon.PseudoStats.BonusMeleeCritRating = dk.PseudoStats.BonusMeleeCritRating
			//dk.RuneWeapon.PseudoStats.BonusSpellCritRating = dk.PseudoStats.BonusSpellCritRating

			// dk.RuneWeapon.PseudoStats.PhysicalDamageDealtMultiplier = dk.PseudoStats.PhysicalDamageDealtMultiplier
			// dk.RuneWeapon.PseudoStats.DiseaseDamageDealtMultiplier = dk.PseudoStats.DiseaseDamageDealtMultiplier
			// dk.RuneWeapon.PseudoStats.ShadowDamageDealtMultiplier = dk.PseudoStats.ShadowDamageDealtMultiplier
			// dk.RuneWeapon.PseudoStats.FrostDamageDealtMultiplier = dk.PseudoStats.FrostDamageDealtMultiplier

			// dk.RuneWeapon.PseudoStats.BonusMHArmorPenRating = dk.PseudoStats.BonusMHArmorPenRating
			// dk.RuneWeapon.PseudoStats.BonusMHCritRating = dk.PseudoStats.BonusMHCritRating
			// dk.RuneWeapon.PseudoStats.BonusMHExpertiseRating = dk.PseudoStats.BonusMHExpertiseRating

			dancingRuneWeaponAura.Activate(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 60.0, 0, 0, 0) && dk.DancingRuneWeapon.IsReady(sim)
	}, nil)
}

func (runeWeapon *RuneWeaponPet) getImpurityBonus(hitEffect *core.SpellEffect, unit *core.Unit) float64 {
	return hitEffect.MeleeAttackPower(unit) + hitEffect.MeleeAttackPowerOnTarget()
}

type RuneWeaponPet struct {
	core.Pet

	dkOwner *Deathknight

	IcyTouch     *core.Spell
	PlagueStrike *core.Spell

	DeathStrike *core.Spell
	DeathCoil   *core.Spell

	HeartStrike       *core.Spell
	HeartStrikeOffHit *core.Spell

	// Diseases
	FrostFeverSpell    *core.Spell
	BloodPlagueSpell   *core.Spell
	FrostFeverDisease  []*core.Dot
	BloodPlagueDisease []*core.Dot
}

func (runeWeapon *RuneWeaponPet) Initialize() {
	runeWeapon.dkOwner.registerDrwDiseaseDots()

	runeWeapon.dkOwner.registerDrwIcyTouchSpell()
	runeWeapon.dkOwner.registerDrwPlagueStrikeSpell()
	runeWeapon.dkOwner.registerDrwDeathStrikeSpell()
	runeWeapon.dkOwner.registerDrwHeartStrikeSpell()
	runeWeapon.dkOwner.registerDrwDeathCoilSpell()
}

func (dk *Deathknight) NewRuneWeapon() *RuneWeaponPet {
	runeWeapon := &RuneWeaponPet{
		Pet: core.NewPet(
			"Rune Weapon",
			&dk.Character,
			runeWeaponBaseStats,
			runeWeaponStatInheritance,
			false,
			true,
		),
		dkOwner: dk,
	}

	runeWeapon.OnPetEnable = runeWeapon.enable
	runeWeapon.OnPetDisable = runeWeapon.disable

	runeWeapon.EnableAutoAttacks(runeWeapon, core.AutoAttackOptions{
		MainHand:       dk.WeaponFromMainHand(2),
		AutoSwingMelee: true,
	})

	runeWeapon.AutoAttacks.MH.SwingSpeed = 3.5
	runeWeapon.PseudoStats.DamageTakenMultiplier = 0

	dk.AddPet(runeWeapon)

	return runeWeapon
}

func (runeWeapon *RuneWeaponPet) GetPet() *core.Pet {
	return &runeWeapon.Pet
}

func (runeWeapon *RuneWeaponPet) Reset(sim *core.Simulation) {
}

func (runeWeapon *RuneWeaponPet) OnGCDReady(sim *core.Simulation) {
	// No GCD system on Rune Weapon
	runeWeapon.DoNothing()
}

func (runeWeapon *RuneWeaponPet) enable(sim *core.Simulation) {
	// Snapshot extra % speed modifiers from dk owner
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, runeWeapon.dkOwner.PseudoStats.MeleeSpeedMultiplier)
}

func (runeWeapon *RuneWeaponPet) disable(sim *core.Simulation) {
	// Clear snapshot speed
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, 1)
}

// These numbers are just rough guesses
var runeWeaponBaseStats = stats.Stats{
	stats.Stamina: 100,
}

var runeWeaponStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	return stats.Stats{
		stats.AttackPower:      ownerStats[stats.AttackPower],
		stats.MeleeHit:         ownerStats[stats.MeleeHit],
		stats.MeleeCrit:        ownerStats[stats.MeleeCrit],
		stats.SpellHit:         ownerStats[stats.SpellHit],
		stats.SpellCrit:        ownerStats[stats.SpellCrit],
		stats.Expertise:        ownerStats[stats.Expertise],
		stats.ArmorPenetration: ownerStats[stats.ArmorPenetration],
	}
}
