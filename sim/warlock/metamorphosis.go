package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
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
			warlock.ImmolationAura.AOEDot().Deactivate(sim)
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
				return character.HasActiveAuraWithTag(core.BloodlustAuraTag) || sim.IsExecutePhase35()
			}

			return true
		},
	})
}

func (warlock *Warlock) registerImmolationAuraSpell() {
	warlock.ImmolationAura = warlock.RegisterSpell(core.SpellConfig{
		// the spellID that deals damage in the combat log is 50590, but we don't use it here
		ActionID:    core.ActionID{SpellID: 50589},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.64,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(30),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Immolation Aura",
			},
			NumberOfTicks:       15,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDmg := (251 + 20*11.5 + 0.143*dot.Spell.SpellPower()) * sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
