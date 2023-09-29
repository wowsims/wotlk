package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) getMiseryCoefficient() float64 {
	return 0.257 * (1 + 0.05*float64(priest.Talents.Misery))
}

func (priest *Priest) getMindFlayTickSpell(numTicks int32) *core.Spell {
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))
	miseryCoeff := priest.getMiseryCoefficient()

	return priest.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 58381}.WithTag(numTicks),
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskProc | core.ProcMaskNotInSpellbook,
		BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			float64(priest.Talents.MindMelt)*2*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetZabras, 4), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines),
		CritMultiplier:   priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 588.0/3 + miseryCoeff*spell.SpellPower()
			damage *= priest.MindFlayModifier
			result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
				if result.DidCrit() && hasGlyphOfShadow {
					priest.ShadowyInsightAura.Activate(sim)
				}
				if result.DidCrit() && priest.ImprovedSpiritTap != nil && sim.RandomFloat("Improved Spirit Tap") > 0.5 {
					priest.ImprovedSpiritTap.Activate(sim)
				}
			}
		},
	})
}

func (priest *Priest) getPainAndSufferingSpell() *core.Spell {
	return priest.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47948},
		ProcMask: core.ProcMaskSuppressedProc,
		Flags:    core.SpellFlagNoLogs,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			priest.ShadowWordPain.Dot(target).Rollover(sim)
		},
	})
}

func (priest *Priest) newMindFlaySpell(numTicksIdx int32) *core.Spell {
	numTicks := numTicksIdx
	flags := core.SpellFlagChanneled | core.SpellFlagNoMetrics
	if numTicksIdx == 0 {
		numTicks = 3
		flags |= core.SpellFlagAPL
	}

	var mfReducTime time.Duration
	if priest.HasSetBonus(ItemSetCrimsonAcolyte, 4) {
		mfReducTime = time.Millisecond * 170
	}
	tickLength := time.Second - mfReducTime
	channelTime := tickLength * time.Duration(numTicks)

	rolloverChance := float64(priest.Talents.PainAndSuffering) / 3.0
	shadowFocus := 0.02 * float64(priest.Talents.ShadowFocus)
	focusedMind := 0.05 * float64(priest.Talents.FocusedMind)
	miseryCoeff := priest.getMiseryCoefficient()

	painAndSufferingSpell := priest.getPainAndSufferingSpell()
	mindFlayTickSpell := priest.getMindFlayTickSpell(numTicksIdx)

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48156}.WithTag(numTicksIdx),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1 - (shadowFocus + focusedMind),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if spell.Unit.IsUsingAPL {
					return
				}
				// if our channel is longer than GCD it will have human latency to end it beause you can't queue the next spell.
				wait := priest.ApplyCastSpeed(channelTime)
				gcd := core.MaxDuration(core.GCDMin, priest.ApplyCastSpeed(core.GCDDefault))
				if wait > gcd && priest.Latency > 0 {
					base := priest.Latency * 0.67
					variation := base + sim.RandomFloat("spriest latency")*base // should vary from 0.66 - 1.33 of given latency
					cast.AfterCastDelay += time.Duration(variation) * time.Millisecond
					if sim.Log != nil {
						priest.Log(sim, "Latency: %0.02f, AfterCastDelay: %s", priest.Latency, cast.AfterCastDelay)
					}
				}
			},
		},
		BonusHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating: 0 +
			float64(priest.Talents.MindMelt)*2*core.CritRatingPerCritChance +
			core.TernaryFloat64(priest.HasSetBonus(ItemSetZabras, 4), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines),
		CritMultiplier: priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5),
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
				if priest.ShadowWordPain.Dot(target).IsActive() {
					if rolloverChance == 1 || sim.RandomFloat("Pain and Suffering") < rolloverChance {
						painAndSufferingSpell.Cast(sim, target)
					}
				}
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := 588.0/3 + miseryCoeff*spell.SpellPower()
			baseDamage *= float64(numTicks)

			if priest.Talents.Shadowform {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})
}

func (priest *Priest) MindFlayTickDuration() time.Duration {
	return priest.ApplyCastSpeed(time.Second - core.TernaryDuration(priest.T10FourSetBonus, time.Millisecond*170, 0))
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
