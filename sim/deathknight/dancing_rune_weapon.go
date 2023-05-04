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

			// RaiN:
			// Auto attacks snapshot damage dealt multipliers with weird formula
			// How it works in game is that 3 auras are applied one after the other:
			// https://wowclassicdb.com/wotlk/spell/51905 / no wowhead entry for some reason...
			// https://wowclassicdb.com/wotlk/spell/51906 / https://www.wowhead.com/wotlk/spell=51906/death-knight-rune-weapon-scaling-02
			// From the testing we could do (still more work to be done) it looks like the full
			// Damage dealt multiplier is transfered from the dk with the first aura (with AP and crit)
			// and then a 2nd -50% damage dealt multiplier is added from the second aura
			// Previous iteration we had made the dks full damage multiplier be applied at 50% to the RW
			// but comparing with logs the damage was way lower then what we saw in-game which lead us to
			// rethink all this and find the above mentioned auras
			dk.RuneWeapon.PseudoStats.DamageDealtMultiplier = dk.PseudoStats.DamageDealtMultiplier - 0.5
			// the second aura also transfers the DK owners physical school damage dealt to the Rune weapon
			// to make sure the dks own physical buffs also affect the rune weapon (tested in game and confirmed they do with UF)
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

	BloodStrike       *core.Spell
	HeartStrike       *core.Spell
	HeartStrikeOffHit *core.Spell

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
