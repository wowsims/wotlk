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
			case dk.IcyTouch:
				dk.RuneWeapon.IcyTouch.Cast(sim, spell.Unit.CurrentTarget)
			case dk.PlagueStrike:
				dk.RuneWeapon.PlagueStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.DeathStrike:
				dk.RuneWeapon.DeathStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.HeartStrike:
				dk.RuneWeapon.HeartStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.DeathCoil:
				dk.RuneWeapon.DeathCoil.Cast(sim, spell.Unit.CurrentTarget)
			case dk.Pestilence.Spell:
				dk.RuneWeapon.Pestilence.Cast(sim, spell.Unit.CurrentTarget)
			}
		},
	})

	dk.DancingRuneWeapon = dk.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49028},

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.RuneWeapon.EnableWithTimeout(sim, dk.Gargoyle, duration)
			dk.RuneWeapon.CancelGCDTimer(sim)

			// Auto attacks snapshot damage dealt multipliers at half
			dk.RuneWeapon.PseudoStats.DamageDealtMultiplier = 0.5 * dk.PseudoStats.DamageDealtMultiplier
			dk.RuneWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] = dk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]

			dancingRuneWeaponAura.Activate(sim)
		},
	})
}

func (runeWeapon *RuneWeaponPet) getImpurityBonus(spell *core.Spell) float64 {
	return spell.MeleeAttackPower()
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

	Pestilence *core.Spell

	// Diseases
	FrostFeverSpell    *core.Spell
	BloodPlagueSpell   *core.Spell
	FrostFeverDisease  []*core.Dot
	BloodPlagueDisease []*core.Dot
}

func (runeWeapon *RuneWeaponPet) Initialize() {
	runeWeapon.dkOwner.registerDrwDiseaseDots()
	runeWeapon.dkOwner.registerDrwPestilenceSpell()

	runeWeapon.dkOwner.registerDrwIcyTouchSpell()
	runeWeapon.dkOwner.registerDrwPlagueStrikeSpell()
	runeWeapon.dkOwner.registerDrwDeathStrikeSpell()
	runeWeapon.dkOwner.registerDrwHeartStrikeSpell()
	runeWeapon.dkOwner.registerDrwDeathCoilSpell()
}

func (dk *Deathknight) NewRuneWeapon() *RuneWeaponPet {
	runeWeapon := &RuneWeaponPet{
		Pet:     core.NewPet("Rune Weapon", &dk.Character, runeWeaponBaseStats, runeWeaponStatInheritance, nil, false, true),
		dkOwner: dk,
	}

	runeWeapon.OnPetEnable = runeWeapon.enable
	runeWeapon.OnPetDisable = runeWeapon.disable

	runeWeapon.EnableAutoAttacks(runeWeapon, core.AutoAttackOptions{
		MainHand:       dk.WeaponFromMainHand(dk.DefaultMeleeCritMultiplier()),
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
		stats.AttackPower: ownerStats[stats.AttackPower],
		stats.MeleeHaste:  ownerStats[stats.MeleeHaste],
		stats.MeleeHit:    ownerStats[stats.MeleeHit],
		stats.MeleeCrit:   ownerStats[stats.MeleeCrit],
		stats.SpellHit:    ownerStats[stats.SpellHit],
		stats.SpellCrit:   ownerStats[stats.SpellCrit],
		stats.Expertise:   ownerStats[stats.Expertise],
	}
}
