package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerDevouringPlagueSpell() {
	actionID := core.ActionID{SpellID: 48300}
	baseCost := priest.BaseMana() * 0.25
	target := priest.CurrentTarget

	applier := priest.OutcomeFuncTick()
	if priest.Talents.Shadowform {
		applier = priest.OutcomeFuncMagicCrit(priest.SpellCritMultiplier(1, 1))
	}

	effect := core.SpellEffect{
		DamageMultiplier: 8 * 0.1 * float64(priest.Talents.ImprovedDevouringPlague) *
			(1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01 + 0.05*float64(priest.Talents.ImprovedDevouringPlague)) *
			core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1),
		ThreatMultiplier: 1 - 0.05*float64(priest.Talents.ShadowAffinity),
		BaseDamage:       core.BaseDamageConfigMagic(172.0, 172.0, 0.1849),
		OutcomeApplier:   priest.OutcomeFuncMagicHitAndCrit(priest.DefaultSpellCritMultiplier()),
		OnSpellHitDealt:  applyDotOnLanded(priest.DevouringPlagueDot),
		ProcMask:         core.ProcMaskSpellDamage,
	}

	priest.DevouringPlague = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(priest.Talents.MentalAgility)),
				GCD:  core.GCDDefault,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	priest.DevouringPlagueDot = core.NewDot(core.Dot{
		Spell: priest.DevouringPlague,
		Aura: target.RegisterAura(core.Aura{
			Label:    "DevouringPlague-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       8,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:             core.ProcMaskPeriodicDamage,
			BonusSpellCritRating: float64(priest.Talents.MindMelt) * 3 * core.CritRatingPerCritChance,
			DamageMultiplier: (1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01 + 0.05*float64(priest.Talents.ImprovedDevouringPlague)) *
				core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1),
			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
			IsPeriodic:       true,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(1376/8, 0.1849),
			OutcomeApplier:   applier,
		}),
	})
}

func applyDotOnLanded(dot *core.Dot) func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
	return func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		if spellEffect.Landed() {
			(*dot).Apply(sim)
		}
	}
}
