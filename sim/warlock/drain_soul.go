package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// func (priest *Priest) MindFlayActionID(numTicks int) core.ActionID {
// 	return core.ActionID{SpellID: 48156, Tag: int32(numTicks)}
// }

func (warlock *Warlock) registerDrainSoulSpell(numTicks int) *core.Spell {
	baseCost := warlock.BaseMana * 0.14
	channelTime := 3 * time.Second * time.Duration(numTicks)
	epsilon := 1* time.Millisecond

	return warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47855},
		SpellSchool:  core.SpellSchoolShadow,
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
				// Everlasting Affliction Refresh
				if warlock.CorruptionDot.IsActive() {
					if sim.RandomFloat("EverlastingAffliction") < 0.2 * float64(warlock.Talents.EverlastingAffliction) {
						 warlock.CorruptionDot.Refresh(sim)
					}
				}
				warlock.DrainSoulDot[numTicks].Apply(sim)
				warlock.DrainSoulDot[numTicks].Aura.UpdateExpires(warlock.DrainSoulDot[numTicks].Aura.ExpiresAt() + epsilon)
			},
		}),
	})
}

func (warlock *Warlock) registerDrainSoulDot(numTicks int) *core.Dot {
	target := warlock.CurrentTarget
	afflictionSpellNumber:= 3.0

	effect := core.SpellEffect{
		DamageMultiplier:     1 + 0.03 * float64(warlock.Talents.SoulSiphon) * afflictionSpellNumber,
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		IsPeriodic:           true,
		OutcomeApplier:       warlock.OutcomeFuncTick(),
		ProcMask:             core.ProcMaskSpellDamage,
		BaseDamage:       	  core.BaseDamageConfigMagicNoRoll(710/5, 0.429),
	}

	return core.NewDot(core.Dot{
		Spell: warlock.DrainSoul[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "Drain Soul-" + strconv.Itoa(numTicks) + "-" + strconv.Itoa(int(warlock.Index)),
			ActionID: core.ActionID{SpellID: 47855},
		}),

		NumberOfTicks:       numTicks,
		TickLength:          3 * time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncSnapshot(target, effect),
	})
}

func (warlock *Warlock) setupDrainSoulExecutePhase() {
	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute20 bool) {
			if isExecute20 {
				for i := 1;  i<=5; i++ {
					warlock.DrainSoulDot[i].Spell.DamageMultiplier *= 2 //TODO : Fix (*=4) when DamageMultiplier is fixed
				}
			}
		})
	})
}
