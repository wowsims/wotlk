package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) PrecastArmyOfTheDead(sim *core.Simulation) {
	// Mark the CD as used already
	dk.ArmyOfTheDead.CD.Use(sim)
	dk.ArmyOfTheDead.CD.Set(sim.CurrentTime + dk.ArmyOfTheDead.CD.Duration - time.Second*10)

	dk.UpdateMajorCooldowns()

	for i := 0; i < 8; i++ {
		timeLeft := (40 - (10 - 0.5*float64(i)))
		if sim.Log != nil {
			sim.Log("Precasting ghoul " + strconv.Itoa(i) + " with duration " + strconv.FormatFloat(timeLeft, 'f', 2, 64))
		}
		dk.ArmyGhoul[i].EnableWithTimeout(sim, dk.ArmyGhoul[i], time.Duration(timeLeft*1000)*time.Millisecond)
	}
}

func (dk *Deathknight) registerArmyOfTheDeadCD() {
	if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_DoNotUse {
		return
	}

	aotdAura := dk.RegisterAura(core.Aura{
		Label:    "Army of the Dead",
		ActionID: core.ActionID{SpellID: 42650},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.AutoAttacks.CancelAutoSwing(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.AutoAttacks.EnableAutoSwing(sim)
		},
	})

	var ghoulIndex = 0
	aotdDot := core.NewDot(core.Dot{
		Aura:                aotdAura,
		NumberOfTicks:       8,
		TickLength:          time.Millisecond * 500,
		AffectedByCastSpeed: true,
		TickEffects: core.TickFuncApplyEffects(func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.ArmyGhoul[ghoulIndex].EnableWithTimeout(sim, dk.ArmyGhoul[ghoulIndex], time.Second*40)
			ghoulIndex++
		}),
	})

	dk.ArmyOfTheDead = dk.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 42650},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				ChannelTime: time.Second * 4,
				GCD:         core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute*10 - time.Minute*2*time.Duration(dk.Talents.NightOfTheDead),
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dkSpellCost := dk.DetermineCost(sim, core.DKCastEnum_BFU)
			dk.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 15.0
			dk.AddRunicPower(sim, amountOfRunicPower, spell.RunicPowerMetrics())

			ghoulIndex = 0
			aotdDot.Apply(sim)
		},
	})

	aotdDot.Spell = dk.ArmyOfTheDead

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell:    dk.ArmyOfTheDead,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			if dk.opener.IsOngoing() {
				return false
			}
			if dk.Gargoyle != nil && !dk.Gargoyle.IsEnabled() {
				return false
			}
			if !dk.CanArmyOfTheDead(sim) {
				return false
			}
			return true
		},
	})
}

func (dk *Deathknight) CanArmyOfTheDead(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 1, 1) && dk.ArmyOfTheDead.IsReady(sim)
}

func (dk *Deathknight) CastArmyOfTheDead(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanArmyOfTheDead(sim) {
		dk.ArmyOfTheDead.Cast(sim, target)
		dk.UpdateMajorCooldowns()
		return true
	}
	return false
}
