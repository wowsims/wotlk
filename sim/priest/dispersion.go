package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (priest *Priest) registerDispersionSpell() {
	if !priest.Talents.Dispersion {
		return
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 47585})
	glyphReduction := 0
	if priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfDispersion)) {
		glyphReduction = 45
	}

	priest.DispersionAura = priest.GetOrRegisterAura(core.Aura{
		Label:    "Dispersion",
		ActionID: core.ActionID{SpellID: 47585},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second,
				NumTicks: 6,
				OnAction: func(sim *core.Simulation) {
					manaGain := priest.MaxMana() * 0.06
					priest.AddMana(sim, manaGain, manaMetric)
				},
			})
		},
	})

	priest.Dispersion = priest.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47585},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second*120 - time.Second*time.Duration(glyphReduction),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			priest.DispersionAura.Activate(sim)
			priest.WaitUntil(sim, priest.DispersionAura.ExpiresAt())
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    priest.Dispersion,
		Priority: 1,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.01
		},
	})
}
