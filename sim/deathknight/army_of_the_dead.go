package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) PrecastArmyOfTheDead(sim *core.Simulation) {
	dk.ArmyOfTheDead.CD.UsePrePull(sim, time.Second*10)
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
			dk.CancelGCDTimer(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.AutoAttacks.EnableAutoSwing(sim)
			dk.SetGCDTimer(sim, sim.CurrentTime)
		},
	})

	var ghoulIndex = 0
	aotdDot := core.NewDot(core.Dot{
		Aura:                aotdAura,
		NumberOfTicks:       8,
		TickLength:          time.Millisecond * 500,
		AffectedByCastSpeed: false,
		TickEffects: core.TickFuncApplyEffects(func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			dk.ArmyGhoul[ghoulIndex].EnableWithTimeout(sim, dk.ArmyGhoul[ghoulIndex], time.Second*40)
			ghoulIndex++
		}),
	})

	baseCost := float64(core.NewRuneCost(15, 1, 1, 1, 0))
	dk.ArmyOfTheDead = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 42650},
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				ChannelTime: time.Second * 4,
				GCD:         core.GCDDefault,
				Cost:        baseCost,
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
			ghoulIndex = 0
			aotdDot.Apply(sim)
		},
	})

	aotdDot.Spell = dk.ArmyOfTheDead.Spell
}

func (dk *Deathknight) CanArmyOfTheDead(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 1, 1, 1) && dk.ArmyOfTheDead.IsReady(sim)
}

func (dk *Deathknight) CastArmyOfTheDead(sim *core.Simulation, target *core.Unit) bool {
	if dk.ArmyOfTheDead.IsReady(sim) && dk.ArmyOfTheDead.Cast(sim, target) {
		dk.UpdateMajorCooldowns()
		return true
	}
	return false
}
