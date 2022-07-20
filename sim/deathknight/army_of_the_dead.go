package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (deathKnight *DeathKnight) PrecastArmyOfTheDead(sim *core.Simulation) {
	// Mark the CD as used already
	deathKnight.ArmyOfTheDead.CD.Use(sim)
	deathKnight.ArmyOfTheDead.CD.Set(sim.CurrentTime + deathKnight.ArmyOfTheDead.CD.Duration - time.Second*10)

	for i := 0; i < 8; i++ {
		timeLeft := (40 - (10 - 0.5*float64(i)))
		if sim.Log != nil {
			sim.Log("Precasting ghoul " + strconv.Itoa(i) + " with duration " + strconv.FormatFloat(timeLeft, 'f', 2, 64))
		}
		deathKnight.ArmyGhoul[i].EnableWithTimeout(sim, deathKnight.ArmyGhoul[i], time.Duration(timeLeft*1000)*time.Millisecond)
	}
}

func (deathKnight *DeathKnight) registerArmyOfTheDeadCD() {
	aotdAura := deathKnight.RegisterAura(core.Aura{
		Label:    "Army of the Dead",
		ActionID: core.ActionID{SpellID: 42650},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.AutoAttacks.CancelAutoSwing(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.AutoAttacks.EnableAutoSwing(sim)
		},
	})

	var ghoulIndex = 0
	aotdDot := core.NewDot(core.Dot{
		Aura:                aotdAura,
		NumberOfTicks:       8,
		TickLength:          time.Millisecond * 500,
		AffectedByCastSpeed: true,
		TickEffects: core.TickFuncApplyEffects(func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			deathKnight.ArmyGhoul[ghoulIndex].EnableWithTimeout(sim, deathKnight.ArmyGhoul[ghoulIndex], time.Second*40)
			ghoulIndex++
		}),
	})

	deathKnight.ArmyOfTheDead = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 42650},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				ChannelTime: time.Second * 4,
				GCD:         core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = deathKnight.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    deathKnight.NewTimer(),
				Duration: time.Minute*10 - time.Minute*2*time.Duration(deathKnight.Talents.NightOfTheDead),
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 1, 1, 1)
			deathKnight.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 15.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

			ghoulIndex = 0
			aotdDot.Apply(sim)
		},
	})

	aotdDot.Spell = deathKnight.ArmyOfTheDead

	// Temp stuff for testing
	if deathKnight.Talents.SummonGargoyle && deathKnight.Rotation.ArmyOfTheDead == proto.DeathKnight_Rotation_AsMajorCd {
		deathKnight.AddMajorCooldown(core.MajorCooldown{
			Spell:    deathKnight.ArmyOfTheDead,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
			CanActivate: func(sim *core.Simulation, character *core.Character) bool {
				if deathKnight.Gargoyle != nil && !deathKnight.Gargoyle.IsEnabled() {
					return false
				}
				if !deathKnight.CanArmyOfTheDead(sim) {
					return false
				}
				return true
			},
		})
	}
}

func (deathKnight *DeathKnight) CanArmyOfTheDead(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 1, 1, 1) && deathKnight.ArmyOfTheDead.IsReady(sim)
}

func (deathKnight *DeathKnight) CastArmyOfTheDead(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanArmyOfTheDead(sim) {
		deathKnight.CastArmyOfTheDead(sim, target)
		return true
	}
	return false
}
