package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyFrostTalents() {
	// Improved Icy Touch
	// Implemented outside

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
	if deathKnight.Talents.NervesOfColdSteel == 1 {
		deathKnight.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.08, true)
	} else if deathKnight.Talents.NervesOfColdSteel == 2 {
		deathKnight.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.16, true)
	} else {
		deathKnight.AutoAttacks.OHEffect.BaseDamage.Calculator = core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.25, true)
	}

	// Icy Talons
	// Pointless to Implement

	// Lichborne
	// Pointless to Implement

	// Annihilation

	// TODO: Implement

	// Killing Machine
	deathKnight.applyKillingMachine()

	// Chill of the Grave
	// Implemented outside

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

	// Frigid Dreadplate
	// TODO: Implement

	// Glacier rot
	// Implemented outside

	// Deathchill
	// TODO: Implement

	// Improved Icy Talons
	deathKnight.applyIcyTalons()
	if deathKnight.Talents.ImprovedIcyTalons {
		deathKnight.PseudoStats.MeleeSpeedMultiplier *= 1.05
	}

	// Merciless Combat
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

func (deathKnight *DeathKnight) applyIcyTalons() {
	if deathKnight.Talents.IcyTalons == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 50887}

	deathKnight.IcyTalonsAura = deathKnight.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Icy Talons",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier *= 1.0 + 0.04*float64(deathKnight.Talents.IcyTalons)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.MeleeSpeedMultiplier /= 1.0 + 0.04*float64(deathKnight.Talents.IcyTalons)
		},
	})
}

func (deathKnight *DeathKnight) applyThreatOfThassarian() {

}
