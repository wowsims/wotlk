package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) channelCheck(sim *core.Simulation, dot *core.Dot, maxTicks int) *core.Spell {

	if dot.IsActive() && dot.TickCount+1 < maxTicks {
		return warlock.DrainSoulChannelling
	} else {
		return dot.Spell
	}
}

func (warlock *Warlock) dynamicDrainSoulMultiplier(sim *core.Simulation) float64 {

	// Execute Multiplier - is now basekit for performance optimization
	// Additive with Death's Embrace so we need to remove its effect to add it again with the spell's own execution multiplier.
	// if sim.IsExecutePhase20() {
	// 	dynamicMultiplier *= (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace))/(1 + 0.04*float64(warlock.Talents.DeathsEmbrace))
	// }

	// Soul Siphon Multiplier
	afflictionSpellNumber := 1. + // you don't need to calculate that: core.TernaryFloat64(warlock.DrainSoulDot.IsActive(), 1, 0) +
		core.TernaryFloat64(warlock.CorruptionDot.IsActive(), 1, 0) + core.TernaryFloat64(warlock.UnstableAffDot.IsActive(), 1, 0) +
		core.TernaryFloat64(warlock.CurseOfAgonyDot.IsActive(), 1, 0) //you should doom, otherwise add "warlock.CurseOfDoomDot.IsActive() ||"
	dynamicMultiplier := 1 + 0.03*float64(warlock.Talents.SoulSiphon) * core.MinFloat(3., afflictionSpellNumber)

	return dynamicMultiplier
}

func (warlock *Warlock) registerDrainSoulSpell() {
	actionID := core.ActionID{SpellID: 47855}
	spellSchool := core.SpellSchoolShadow
	baseAdditiveMultiplier:= warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true)
	// For performance optimization, the execute modifier is basekit since we never use it before execute
	normalMultiplier := baseAdditiveMultiplier * (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)) / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace))
	baseCost := warlock.BaseMana * 0.14
	channelTime := 3 * time.Second
	epsilon := 1 * time.Millisecond

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		Flags:        core.SpellFlagBinary | core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
			OutcomeApplier:   warlock.OutcomeFuncMagicHitBinary(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				warlock.DrainSoulDot.Apply(sim)
				warlock.DrainSoulDot.Aura.UpdateExpires(warlock.DrainSoulDot.Aura.ExpiresAt() + epsilon)
			},
		}),
	})

	target := warlock.CurrentTarget

	effect := core.SpellEffect{
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		IsPeriodic:       true,
		OutcomeApplier:   warlock.OutcomeFuncTick(),
		ProcMask:         core.ProcMaskPeriodicDamage,
		BaseDamage:       core.BaseDamageConfigMagicNoRoll(710/5, 3./7.),
		OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			spellEffect.DamageMultiplier = normalMultiplier * warlock.dynamicDrainSoulMultiplier(sim)
		},
	}

	warlock.DrainSoulDot = core.NewDot(core.Dot{
		Spell: warlock.DrainSoul,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Drain Soul-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       1,
		TickLength:          3 * time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncSnapshot(target, effect),
	})

	warlock.DrainSoulChannelling = warlock.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
				CastTime:    0,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,
			FlatThreatBonus:  1,
			OutcomeApplier:   warlock.OutcomeFuncAlwaysHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				warlock.DrainSoulDot.Refresh(sim)
				warlock.DrainSoulDot.Aura.UpdateExpires(warlock.DrainSoulDot.Aura.ExpiresAt() + epsilon)
			},
		}),
	})
}

