package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// TODO Mind Flay (48156) now "periodically triggers" Mind Flay (58381), probably to allow haste to work.
// The first never deals damage, so the latter should probably be used as ActionID here.

func (priest *Priest) newMindFlaySpell(numTicks int32) *core.Spell {
	var mfReducTime time.Duration
	if priest.HasSetBonus(ItemSetCrimsonAcolyte, 4) {
		mfReducTime = time.Millisecond * 170
	}
	tickLength := time.Second - mfReducTime
	channelTime := tickLength * time.Duration(numTicks)

	rolloverChance := float64(priest.Talents.PainAndSuffering) / 3.0
	miseryCoeff := 0.257 * (1 + 0.05*float64(priest.Talents.Misery))
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48156}.WithTag(numTicks),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1 - 0.05*float64(priest.Talents.FocusedMind),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				// if our channel is longer than GCD it will have human latency to end it beause you can't queue the next spell.
				wait := priest.ApplyCastSpeed(channelTime)
				gcd := core.MaxDuration(core.GCDMin, priest.ApplyCastSpeed(core.GCDDefault))
				if wait > gcd && priest.Latency > 0 {
					base := priest.Latency * 0.25
					variation := base + sim.RandomFloat("spriest latency")*base // should vary from 0.66 - 1.33 of given latency
					variation = core.MaxFloat(variation, 10)                    // no player can go under XXXms response time
					cast.AfterCastDelay += time.Duration(variation) * time.Millisecond
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
		CritMultiplier:   priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindFlay-" + strconv.Itoa(int(numTicks)),
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 588.0/3 + miseryCoeff*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeMagicHitAndSnapshotCrit)
				dot.Spell.DealDamage(sim, result)

				if result.Landed() {
					priest.AddShadowWeavingStack(sim)
				}
				if result.DidCrit() && hasGlyphOfShadow {
					priest.ShadowyInsightAura.Activate(sim)
				}
				if result.DidCrit() && priest.ImprovedSpiritTap != nil && sim.RandomFloat("Improved Spirit Tap") > 0.5 {
					priest.ImprovedSpiritTap.Activate(sim)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				if priest.ShadowWordPain.Dot(target).IsActive() {
					if rolloverChance == 1 || sim.RandomFloat("Pain and Suffering") < rolloverChance {
						priest.ShadowWordPain.Dot(target).Rollover(sim)
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
	variation = core.MaxFloat(variation, 10)
	return time.Duration(variation) * time.Millisecond
}
