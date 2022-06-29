package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (priest *Priest) MindFlayActionID(numTicks int) core.ActionID {
	return core.ActionID{SpellID: 25387, Tag: int32(numTicks)}
}

func (priest *Priest) newMindFlaySpell(numTicks int) *core.Spell {
	baseCost := 230.0
	channelTime := time.Second * time.Duration(numTicks)

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:     priest.MindFlayActionID(numTicks),
		SpellSchool:  core.SpellSchoolShadow,
		Flags:        core.SpellFlagBinary | core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				// if our channel is longer than GCD it will have human latency to end it beause you can't queue the next spell.
				wait := priest.ApplyCastSpeed(channelTime)
				gcd := core.MaxDuration(core.GCDMin, priest.ApplyCastSpeed(core.GCDDefault))
				if wait > gcd && priest.Latency > 0 {
					base := priest.Latency * 0.66
					variation := base + sim.RandomFloat("spriest latency")*base // should vary from 0.66 - 1.33 of given latency
					variation = core.MaxFloat(variation, 10)                    // no player can go under XXXms response time
					cast.AfterCastDelay += time.Duration(variation) * time.Millisecond
				}
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: float64(priest.Talents.ShadowFocus) * 2 * core.SpellHitRatingPerHitChance,
			ThreatMultiplier:    1 - 0.08*float64(priest.Talents.ShadowAffinity),
			OutcomeApplier:      priest.OutcomeFuncMagicHitBinary(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					priest.MindFlayDot[numTicks].Apply(sim)
				}
			},
		}),
	})
}

func (priest *Priest) newMindFlayDot(numTicks int) *core.Dot {
	target := priest.CurrentTarget
	return core.NewDot(core.Dot{
		Spell: priest.MindFlay[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "MindFlay-" + strconv.Itoa(numTicks) + "-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.MindFlayActionID(numTicks),
		}),

		NumberOfTicks:       numTicks,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1 *
				(1 + float64(priest.Talents.Darkness)*0.02) *
				core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1) *
				core.TernaryFloat64(ItemSetIncarnate.CharacterHasSetBonus(&priest.Character, 4), 1.05, 1),
			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
			IsPeriodic:       true,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(528/3, 0.19),
			OutcomeApplier:   priest.OutcomeFuncTick(),
			ProcMask:         core.ProcMaskPeriodicDamage,
		}),
	})
}
