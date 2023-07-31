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
			case dk.BloodStrike:
				dk.RuneWeapon.BloodStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.HeartStrike:
				dk.RuneWeapon.HeartStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.RuneStrike:
				dk.RuneWeapon.RuneStrike.Cast(sim, spell.Unit.CurrentTarget)
			case dk.DeathCoil:
				dk.RuneWeapon.DeathCoil.Cast(sim, spell.Unit.CurrentTarget)
			case dk.Pestilence:
				dk.RuneWeapon.Pestilence.Cast(sim, spell.Unit.CurrentTarget)
			case dk.BloodBoil:
				dk.RuneWeapon.BloodBoil.Cast(sim, spell.Unit.CurrentTarget)
			}
		},
	})

	dk.DancingRuneWeapon = dk.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49028},
		Flags:    core.SpellFlagAPL,

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
			dk.RuneWeapon.EnableWithTimeout(sim, dk.RuneWeapon, duration)
			dk.RuneWeapon.CancelGCDTimer(sim)
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

	BloodStrike       *core.Spell
	HeartStrike       *core.Spell
	HeartStrikeOffHit *core.Spell

	RuneStrike *core.Spell

	Pestilence *core.Spell
	BloodBoil  *core.Spell

	// Diseases
	FrostFeverSpell  *core.Spell
	BloodPlagueSpell *core.Spell
}

func (runeWeapon *RuneWeaponPet) Initialize() {
	runeWeapon.dkOwner.registerDrwDiseaseDots()
	runeWeapon.dkOwner.registerDrwPestilenceSpell()
	runeWeapon.dkOwner.registerDrwBloodBoilSpell()

	runeWeapon.dkOwner.registerDrwIcyTouchSpell()
	runeWeapon.dkOwner.registerDrwPlagueStrikeSpell()
	runeWeapon.dkOwner.registerDrwDeathStrikeSpell()
	runeWeapon.dkOwner.registerDrwBloodStrikeSpell()
	runeWeapon.dkOwner.registerDrwHeartStrikeSpell()
	runeWeapon.dkOwner.registerDrwDeathCoilSpell()
	runeWeapon.dkOwner.registerDrwRuneStrikeSpell()
}

func (dk *Deathknight) DrwWeaponDamage(sim *core.Simulation, spell *core.Spell) float64 {
	if dk.Inputs.NewDrw {
		return spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
			spell.BonusWeaponDamage()
	} else {
		return spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
			spell.BonusWeaponDamage()
	}
}

func (dk *Deathknight) NewRuneWeapon() *RuneWeaponPet {
	// Remove any hit that would be given by NocS as it does not translate to pets
	nocsHit := 0.0
	nocsSpellHit := 0.0
	if dk.nervesOfColdSteelActive() {
		nocsHit = float64(dk.Talents.NervesOfColdSteel)
		nocsSpellHit = (float64(dk.Talents.NervesOfColdSteel) / 8.0) * 17.0
	}
	if dk.HasDraeneiHitAura {
		nocsHit = nocsHit + 1.0
		nocsSpellHit = nocsSpellHit + 1.0
	}

	runeWeapon := &RuneWeaponPet{
		Pet: core.NewPet("Rune Weapon", &dk.Character,
			stats.Stats{
				stats.Stamina:   100,
				stats.MeleeHit:  -nocsHit * core.MeleeHitRatingPerHitChance,
				stats.SpellHit:  -nocsSpellHit * core.SpellHitRatingPerHitChance,
				stats.Expertise: -nocsHit * PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,
			},
			func(ownerStats stats.Stats) stats.Stats {
				ownerHitChance := ownerStats[stats.MeleeHit] / core.MeleeHitRatingPerHitChance
				return stats.Stats{
					stats.AttackPower: ownerStats[stats.AttackPower],
					stats.MeleeHaste:  (ownerStats[stats.MeleeHaste] / dk.PseudoStats.MeleeHasteRatingPerHastePercent) * core.HasteRatingPerHastePercent,

					stats.MeleeHit: ownerHitChance * core.MeleeHitRatingPerHitChance,
					stats.SpellHit: ((ownerHitChance / 8.0) * 17.0) * core.SpellHitRatingPerHitChance,

					stats.Expertise: ownerHitChance * PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,

					stats.MeleeCrit: ownerStats[stats.MeleeCrit],
					stats.SpellCrit: ownerStats[stats.SpellCrit],
				}
			},
			nil, false, true),
		dkOwner: dk,
	}

	runeWeapon.OnPetEnable = runeWeapon.enable
	runeWeapon.OnPetDisable = runeWeapon.disable

	mhWeapon := dk.WeaponFromMainHand(dk.DefaultMeleeCritMultiplier())
	runeWeapon.EnableAutoAttacks(runeWeapon, core.AutoAttackOptions{
		MainHand:       mhWeapon,
		AutoSwingMelee: true,
	})

	runeWeapon.AutoAttacks.MH.SwingSpeed = 3.5
	runeWeapon.AutoAttacks.MH.NormalizedSwingSpeed = 3.3
	runeWeapon.AutoAttacks.MH.SwingDuration = time.Duration(runeWeapon.AutoAttacks.MH.SwingSpeed * float64(time.Second))
	runeWeapon.PseudoStats.DamageTakenMultiplier = 0

	if dk.Inputs.NewDrw {
		baseDamage := (mhWeapon.BaseDamageMin + mhWeapon.BaseDamageMax) / 2
		baseDamage = (baseDamage / mhWeapon.SwingSpeed) * 3.5
		runeWeapon.AutoAttacks.MH.BaseDamageMin = baseDamage - 150
		runeWeapon.AutoAttacks.MH.BaseDamageMax = baseDamage + 150
	}

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

	if runeWeapon.dkOwner.Inputs.NewDrw {
		runeWeapon.dkOwner.drwDmgSnapshot = runeWeapon.dkOwner.PseudoStats.DamageDealtMultiplier * 0.5
		runeWeapon.dkOwner.RuneWeapon.PseudoStats.DamageDealtMultiplier *= runeWeapon.dkOwner.drwDmgSnapshot
	} else {
		runeWeapon.dkOwner.drwDmgSnapshot = runeWeapon.dkOwner.PseudoStats.DamageDealtMultiplier - 0.5
		runeWeapon.dkOwner.RuneWeapon.PseudoStats.DamageDealtMultiplier *= runeWeapon.dkOwner.drwDmgSnapshot
	}

	runeWeapon.dkOwner.drwPhysSnapshot = runeWeapon.dkOwner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical]
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= runeWeapon.dkOwner.drwPhysSnapshot

}

func (runeWeapon *RuneWeaponPet) disable(sim *core.Simulation) {
	// Clear snapshot speed
	runeWeapon.PseudoStats.MeleeSpeedMultiplier = 1
	runeWeapon.MultiplyMeleeSpeed(sim, 1)

	// Clear snapshot damage multipliers
	runeWeapon.dkOwner.RuneWeapon.PseudoStats.DamageDealtMultiplier /= runeWeapon.dkOwner.drwDmgSnapshot
	runeWeapon.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= runeWeapon.dkOwner.drwPhysSnapshot
	runeWeapon.dkOwner.drwPhysSnapshot = 1
	runeWeapon.dkOwner.drwDmgSnapshot = 1
}
