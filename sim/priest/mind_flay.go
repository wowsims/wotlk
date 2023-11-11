package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) getMindFlayTickSpell(numTicks int32) *core.Spell {

	return priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 18807}.WithTag(numTicks),
		SpellSchool:      core.SpellSchoolShadow,
		ProcMask:         core.ProcMaskProc | core.ProcMaskNotInSpellbook,
		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 426.0 / 3
			damage *= priest.MindFlayModifier
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
		},
	})
}

func (priest *Priest) newMindFlaySpell(numTicksIdx int32) core.SpellConfig {
	numTicks := numTicksIdx
	flags := core.SpellFlagChanneled | core.SpellFlagNoMetrics | core.SpellFlagAPL

	tickLength := time.Second
	channelTime := tickLength * time.Duration(numTicks)
	mindFlayTickSpell := priest.getMindFlayTickSpell(numTicksIdx)

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48156}.WithTag(numTicksIdx),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.09,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if spell.Unit.IsUsingAPL || priest.Latency == 0 {
					return
				}
				// if our channel is longer than GCD it will have human latency to end it because you can't queue the next spell.
				if float64(channelTime)*priest.CastSpeed > float64(core.GCDMin) {
					variation := priest.Latency * (0.66 + sim.RandomFloat("spriest latency")*(1.33-0.66)) // should vary from 0.66 - 1.33 of given latency
					cast.ChannelTime += time.Duration(variation / priest.CastSpeed * float64(time.Millisecond))
					if sim.Log != nil {
						priest.Log(sim, "Latency: %.3f, Applied Latency: %.3f", priest.Latency, variation)
					}
				}
			},
		},
		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindFlay-" + strconv.Itoa(int(numTicksIdx)),
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mindFlayTickSpell.Cast(sim, target)
				mindFlayTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			mindFlayTickSpell.SpellMetrics[target.UnitIndex].Casts += 1

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := 426.0/3 + spell.SpellPower()

			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
		},
	}
}

func (priest *Priest) MindFlayTickDuration() time.Duration {
	return priest.ApplyCastSpeed(time.Second)
}

func (priest *Priest) AverageMindFlayLatencyDelay(numTicks int, gcd time.Duration) time.Duration {
	wait := priest.ApplyCastSpeed(priest.MindFlay[numTicks].DefaultCast.ChannelTime)
	if wait <= gcd || priest.Latency == 0 {
		return 0
	}

	base := priest.Latency * 0.25
	variation := base + 0.5*base
	return time.Duration(variation) * time.Millisecond
}

func (priest *Priest) registerMindFlay() {
	priest.MindBlast = priest.GetOrRegisterSpell(priest.newMindFlaySpell(3))
}
