package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

const EvocationId = 12051

func (mage *Mage) registerEvocation(numTicks int32) *core.Spell {
	actionID := core.ActionID{SpellID: EvocationId}
	manaMetrics := mage.NewManaMetrics(actionID)

	channelTime := time.Duration(numTicks) * time.Second * 2
	manaPerTick := 0.0
	mage.Env.RegisterPostFinalizeEffect(func() {
		manaPerTick = mage.MaxMana() * 0.15
	})

	return mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID: actionID.WithTag(numTicks),
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * time.Duration(4-mage.Talents.ArcaneFlows),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			period := spell.CurCast.ChannelTime / time.Duration(numTicks)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				NumTicks: int(numTicks),
				OnAction: func(sim *core.Simulation) {
					mage.AddMana(sim, manaPerTick, manaMetrics)
				},
			})
		},
	})
}

func (mage *Mage) registerEvocationCD() {
	maxTicks := core.TernaryInt32(mage.HasSetBonus(ItemSetTempestRegalia, 2), 5, 4)

	numTicks := core.MaxInt32(0, core.MinInt32(maxTicks, mage.Options.EvocationTicks))
	if numTicks == 0 {
		numTicks = maxTicks
	}

	manaThreshold := 0.0
	mage.Env.RegisterPostFinalizeEffect(func() {
		manaThreshold = mage.MaxMana() * 0.1
	})

	evocationSpell := mage.registerEvocation(numTicks)

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: evocationSpell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if character.HasActiveAuraWithTag(core.InnervateAuraTag) || character.HasActiveAuraWithTag(core.ManaTideTotemAuraTag) {
				return false
			}

			if sim.GetRemainingDuration() < 12*time.Second {
				return false
			}

			curMana := character.CurrentMana()

			return curMana < manaThreshold
		},
	})
}

func (mage *Mage) registerEvocationSpells() {
	maxTicks := core.TernaryInt32(mage.HasSetBonus(ItemSetTempestRegalia, 2), 5, 4)

	for i := int32(1); i <= maxTicks; i++ {
		mage.registerEvocation(i)
	}
}
