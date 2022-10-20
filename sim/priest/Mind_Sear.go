package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO see Mind Flay: Mind Sear (53023) now "periodically triggers" Mind Sear (53022).
//  Since Mind Flay no longer is a binary spell, Mind Sear likely isn't, either.
func (priest *Priest) MindSearActionID(numTicks int) core.ActionID {
	return core.ActionID{SpellID: 53023, Tag: int32(numTicks)}
}

func (priest *Priest) newMindSearSpell(numTicks int) *core.Spell {
	baseCost := priest.BaseMana * 0.28
	channelTime := time.Second * time.Duration(numTicks)

	effect := core.SpellEffect{
		OutcomeApplier: priest.OutcomeFuncMagicHit(),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			priest.MindSearDot[numTicks].Apply(sim)
		},
	}

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:     priest.MindSearActionID(numTicks),
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  float64(priest.Talents.MindMelt) * 2 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(priest.Env, effect),
	})
}

func (priest *Priest) newMindSearDot(numTicks int) *core.Dot {
	target := priest.CurrentTarget

	effect := core.SpellEffect{
		IsPeriodic:     true,
		OutcomeApplier: priest.OutcomeFuncMagicHit(),
		OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			if spellEffect.DidCrit() && priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow)) {
				priest.ShadowyInsightAura.Activate(sim)
			}
		},
	}

	normalCalc := core.BaseDamageFuncMagic(212, 228, 0.2861)
	miseryCalc := core.BaseDamageFuncMagic(212, 228, (1+float64(priest.Talents.Misery)*0.05)*0.2861)

	normMod := (1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01) // initialize modifier

	effect.BaseDamage = core.BaseDamageConfig{
		Calculator: func(sim *core.Simulation, effect *core.SpellEffect, spell *core.Spell) float64 {
			var dmg float64
			shadowWeavingMod := 1 + float64(priest.ShadowWeavingAura.GetStacks())*0.02

			if priest.MiseryAura.IsActive() {
				dmg = miseryCalc(sim, effect, spell)
			} else {
				dmg = normalCalc(sim, effect, spell)
			}
			dmg *= normMod // multiply the damage
			return dmg * shadowWeavingMod
		},
	}

	return core.NewDot(core.Dot{
		Spell: priest.MindSear[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "MindSear-" + strconv.Itoa(numTicks) + "-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.MindSearActionID(numTicks),
		}),

		NumberOfTicks:       numTicks,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,
		TickEffects:         core.TickFuncAOESnapshotCapped(priest.Env, effect),
	})
}
