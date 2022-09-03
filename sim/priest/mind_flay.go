package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) MindFlayActionID(numTicks int) core.ActionID {
	return core.ActionID{SpellID: 48156, Tag: int32(numTicks)}
}

func (priest *Priest) newMindFlaySpell(numTicks int) *core.Spell {
	baseCost := priest.BaseMana * 0.09

	channelTime := time.Second * time.Duration(numTicks)
	if priest.HasSetBonus(ItemSetCrimsonAcolyte, 4) {
		channelTime = channelTime - time.Duration(numTicks)*(time.Millisecond*170)
	}

	effect := core.SpellEffect{
		ProcMask:            core.ProcMaskEmpty,
		BonusSpellHitRating: float64(priest.Talents.ShadowFocus)*1*core.SpellHitRatingPerHitChance + 3*core.SpellHitRatingPerHitChance, //not sure if misery is applying to this bonus spell hit so adding it here
		ThreatMultiplier:    1 - 0.08*float64(priest.Talents.ShadowAffinity),
		OutcomeApplier:      priest.OutcomeFuncMagicHitBinary(),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if priest.ShadowWordPainDot.IsActive() {
				if priest.Talents.PainAndSuffering == 3 {
					priest.ShadowWordPainDot.Rollover(sim)
				} else if sim.RandomFloat("Pain and Suffering") < (float64(priest.Talents.PainAndSuffering) * 0.33) {
					priest.ShadowWordPainDot.Rollover(sim)
				}
			}
			priest.MindFlayDot[numTicks].Apply(sim)
		},
	}

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
					base := priest.Latency * 0.25
					variation := base + sim.RandomFloat("spriest latency")*base // should vary from 0.66 - 1.33 of given latency
					variation = core.MaxFloat(variation, 10)                    // no player can go under XXXms response time
					cast.AfterCastDelay += time.Duration(variation) * time.Millisecond
				}
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (priest *Priest) newMindFlayDot(numTicks int) *core.Dot {
	target := priest.CurrentTarget

	effect := core.SpellEffect{
		DamageMultiplier:     1,
		ThreatMultiplier:     1 - 0.08*float64(priest.Talents.ShadowAffinity),
		BonusSpellHitRating:  float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		IsPeriodic:           true,
		BonusSpellCritRating: float64(priest.Talents.MindMelt)*2*core.CritRatingPerCritChance + core.TernaryFloat64(priest.HasSetBonus(ItemSetZabras, 4), 5, 0)*core.CritRatingPerCritChance,
		OutcomeApplier:       priest.OutcomeFuncMagicHitAndCrit(1 + float64(priest.Talents.ShadowPower)*0.2),
		ProcMask:             core.ProcMaskSpellDamage,
		OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			if spellEffect.DidCrit() && priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow)) {
				priest.ShadowyInsightAura.Activate(sim)
			}
			if spellEffect.DidCrit() && priest.ImprovedSpiritTap != nil && sim.RandomFloat("Improved Spirit Tap") > 0.5 {
				priest.ImprovedSpiritTap.Activate(sim)
			}
		},
	}

	normalCalc := core.BaseDamageFuncMagic(588/3, 588/3, 0.257)
	miseryCalc := core.BaseDamageFuncMagic(588/3, 588/3, (1+float64(priest.Talents.Misery)*0.05)*0.257)

	normMod := (1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01) * // initialize modifier
		core.TernaryFloat64(priest.HasSetBonus(ItemSetIncarnate, 4), 1.05, 1)

	effect.BaseDamage = core.BaseDamageConfig{
		Calculator: func(sim *core.Simulation, effect *core.SpellEffect, spell *core.Spell) float64 {
			var dmg float64
			shadowWeavingMod := 1 + float64(priest.ShadowWeavingAura.GetStacks())*0.02
			glyphMod := 0.0

			if priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfMindFlay)) {
				glyphMod = 0.1
			}

			if priest.MiseryAura.IsActive() {
				dmg = miseryCalc(sim, effect, spell)
			} else {
				dmg = normalCalc(sim, effect, spell)
			}
			if priest.ShadowWordPainDot.IsActive() {
				dmg *= normMod * (1 + glyphMod + float64(priest.Talents.TwistedFaith)*0.02) // multiply the damage
			} else {
				dmg *= normMod // multiply the damage
			}
			return dmg * shadowWeavingMod
		},
		TargetSpellCoefficient: 0.0,
	}

	var mfReducTime time.Duration
	if priest.HasSetBonus(ItemSetCrimsonAcolyte, 4) {
		mfReducTime = time.Millisecond * 170
	}

	return core.NewDot(core.Dot{
		Spell: priest.MindFlay[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "MindFlay-" + strconv.Itoa(numTicks) + "-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.MindFlayActionID(numTicks),
		}),

		NumberOfTicks:       numTicks,
		TickLength:          time.Second - mfReducTime,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncSnapshot(target, effect),
	})
}

func (priest *Priest) MindFlayTickDuration() time.Duration {
	return priest.ApplyCastSpeed(time.Second - core.TernaryDuration(priest.T10FourSetBonus, time.Millisecond*170, 0))
}
