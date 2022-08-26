package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) RegisterPenanceSpell() {
	if !priest.Talents.Penance {
		return
	}

	actionID := core.ActionID{SpellID: 53007}
	baseCost := priest.BaseMana * 0.16

	var penanceDot *core.Dot

	priest.Penance = priest.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolHoly,
		Flags:        core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Duration(float64(time.Second*12-core.TernaryDuration(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPenance), time.Second*2, 0)) * (1 - .1*float64(priest.Talents.Aspiration))),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 0,
			OutcomeApplier:   priest.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					penanceDot.Apply(sim)
					// Do immediate tick
					penanceDot.TickOnce()
				}
			},
		}),
	})

	target := priest.CurrentTarget
	penanceDot = core.NewDot(core.Dot{
		Spell: priest.Penance,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Penance-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       2,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskPeriodicDamage,

			DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
			ThreatMultiplier: 0,

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(375, .4286),
			OutcomeApplier: priest.OutcomeFuncMagicHit(),
		})),
	})
}
