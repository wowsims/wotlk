package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerMetamorphosisSpell() {

	warlock.MetamorphosisAura = warlock.RegisterAura(core.Aura{
		Label:    "Metamorphosis Aura",
		ActionID: core.ActionID{SpellID: 47241},
		Duration: time.Second * (30 + 6*core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfMetamorphosis), 1, 0)),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
			warlock.ImmolationAuraDot.Deactivate(sim)
		},
	})

	warlock.Metamorphosis = warlock.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47241},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(3*60.*(1.-0.1*float64(warlock.Talents.Nemesis))),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.MetamorphosisAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: warlock.Metamorphosis,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if !warlock.GetAura("Demonic Pact").IsActive() {
				return false
			}
			MetamorphosisNumber := (float64(sim.Duration) + float64(warlock.MetamorphosisAura.Duration)) / float64(warlock.Metamorphosis.CD.Duration)
			if MetamorphosisNumber < 1 {
				if character.HasActiveAuraWithTag(core.BloodlustAuraTag) || sim.IsExecutePhase35() {
					return true
				}
			} else if warlock.Metamorphosis.CD.IsReady(sim) {
				return true
			}
			return false
		},
	})
}

func (warlock *Warlock) registerImmolationAuraSpell() {
	// the spellID that deals damage in the combat log is 50590, but we don't use it here
	actionID := core.ActionID{SpellID: 50589}
	baseCost := 0.64 * warlock.BaseMana
	spellSchool := core.SpellSchoolFire

	warlock.ImmolationAuraDot = core.NewDot(core.Dot{
		Aura: warlock.RegisterAura(core.Aura{
			Label:    "Immolation Aura",
			ActionID: actionID,
		}),
		NumberOfTicks:       15,
		TickLength:          time.Second * 1,
		AffectedByCastSpeed: true,
		// TODO: obey the AoE cap
		TickEffects: func(sim *core.Simulation, dot *core.Dot) func() {
			effectsFunc := core.ApplyEffectFuncAOEDamage(warlock.Env, core.SpellEffect{
				// TODO: spell is flagged as "Treat As Periodic" but doesn't proc timbal's, so not
				// adding core.ProcMaskPeriodicDamage should be correct?
				ProcMask:         core.ProcMaskSpellDamage,
				DamageMultiplier: 1,
				ThreatMultiplier: 1,
				BaseDamage:       core.BaseDamageConfigMagicNoRoll(251+20*11.5, 0.143),
				OutcomeApplier:   warlock.OutcomeFuncMagicHit(),
				IsPeriodic:       false,
			})

			return func() {
				effectsFunc(sim, dot.Spell.Unit.CurrentTarget, dot.Spell)
			}
		},
	})

	warlock.ImmolationAura = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(30),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDot(warlock.ImmolationAuraDot),
	})
	warlock.ImmolationAuraDot.Spell = warlock.ImmolationAura
}
