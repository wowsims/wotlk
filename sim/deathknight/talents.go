package deathknight

import ( //	"time"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyTalents() {
	// Blood

	// Butchery

	// Subversion

	// Bladed Armor
	if deathKnight.Talents.BladedArmor > 0 {
		coeff := float64(deathKnight.Talents.BladedArmor)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Armor,
			ModifiedStat: stats.AttackPower,
			Modifier: func(armor float64, attackPower float64) float64 {
				return attackPower + coeff*armor/180.0
			},
		})
	}

	// Two Handed Specialization
	if deathKnight.Equip[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand {
		deathKnight.PseudoStats.PhysicalDamageDealtMultiplier *= 1 + 0.01*float64(deathKnight.Talents.TwoHandedWeaponSpecialization)
	}

	// Rune Tap

	// Dark Conviction
	deathKnight.PseudoStats.BonusMeleeCritRating += core.CritRatingPerCritChance * float64(deathKnight.Talents.DarkConviction)
	deathKnight.PseudoStats.BonusSpellCritRating += core.CritRatingPerCritChance * float64(deathKnight.Talents.DarkConviction)

	// Death Rune Mastery

	// Improved Rune Tap

	// Spell Deflection

	// Vendetta

	// Bloody Strikes

	// Veteran of the Third War
	if deathKnight.Talents.VeteranOfTheThirdWar > 0 {
		strengthCoeff := 0.02 * float64(deathKnight.Talents.VeteranOfTheThirdWar)
		staminaCoeff := 0.01 * float64(deathKnight.Talents.VeteranOfTheThirdWar)
		expertiseBonus := 2.0 * float64(deathKnight.Talents.VeteranOfTheThirdWar)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Strength,
			ModifiedStat: stats.Strength,
			Modifier: func(strength float64, _ float64) float64 {
				return strength * (1.0 + strengthCoeff)
			},
		})

		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(stamina float64, _ float64) float64 {
				return stamina * (1.0 + staminaCoeff)
			},
		})

		deathKnight.AddStat(stats.Expertise, expertiseBonus*core.ExpertisePerQuarterPercentReduction)
	}

	// Mark of Blood

	// Bloody Vengeance

	// Abomination's Might
	if deathKnight.Talents.AbominationsMight > 0 {
		strengthCoeff := 0.01 * float64(deathKnight.Talents.AbominationsMight)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Strength,
			ModifiedStat: stats.Strength,
			Modifier: func(strength float64, _ float64) float64 {
				return strength * (1.0 + strengthCoeff)
			},
		})
	}

	// Frost

	// Improved Icy Touch

	// Toughness
	if deathKnight.Talents.Toughness > 0 {
		armorCoeff := 0.02 * float64(deathKnight.Talents.Toughness)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Armor,
			ModifiedStat: stats.Armor,
			Modifier: func(armor float64, _ float64) float64 {
				return armor * (1.0 + armorCoeff)
			},
		})
	}

	// Icy Reach
	// Pointless to Implement

	// Black Ice
	deathKnight.PseudoStats.FrostDamageDealtMultiplier += 0.02 * float64(deathKnight.Talents.BlackIce)
	deathKnight.PseudoStats.ShadowDamageDealtMultiplier += 0.02 * float64(deathKnight.Talents.BlackIce)

	// Nerves Of Cold Steel
	deathKnight.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(deathKnight.Talents.NervesOfColdSteel))

	// Icy Talons
	// Pointless to Implement

	// Lichborne
	// Pointless to Implement

	// Annihilation

	// Killing Machine
	deathKnight.applyKillingMachine()

	// Chill of the Grave

	// Endless Winter
	if deathKnight.Talents.EndlessWinter > 0 {
		strengthCoeff := 0.02 * float64(deathKnight.Talents.EndlessWinter)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Strength,
			ModifiedStat: stats.Strength,
			Modifier: func(strength float64, _ float64) float64 {
				return strength * (1.0 + strengthCoeff)
			},
		})
	}

	//Unholy
}

func (deathKnight *DeathKnight) applyKillingMachine() {
	if deathKnight.Talents.KillingMachine == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 51130}
	weaponMH := deathKnight.GetMHWeapon()
	procChance := (weaponMH.SwingSpeed * 5.0 / 60.0) * float64(deathKnight.Talents.KillingMachine)

	deathKnight.KillingMachineAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Killing Machine Proc",
		ActionID: actionID,
		Duration: time.Second * 30.0,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			//TODO: add the other spells
			if spell == deathKnight.IcyTouch {
				aura.Deactivate(sim)
			}
		},
	})

	deathKnight.RegisterAura(core.Aura{
		Label:    "Killing Machine",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if spell != deathKnight.IcyTouch {
				return
			}

			if sim.RandomFloat("Killing Machine") < procChance {
				deathKnight.KillingMachineAura.Activate(sim)
			}
		},
	})
}
