package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerAntiMagicShellSpell() {
	actionID := core.ActionID{SpellID: 48707}
	cdTimer := dk.NewTimer()
	cd := time.Second*45 - time.Second*time.Duration(core.TernaryInt32(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfAntiMagicShell), 2, 0))

	baseCost := float64(core.NewRuneCost(20, 0, 0, 0, 0))
	dk.AntiMagicShell = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.AntiMagicShellAura.Activate(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 20.0, 0, 0, 0) && dk.AntiMagicShell.IsReady(sim)
	}, nil)

	rpMetrics := dk.AntiMagicShell.RunicPowerMetrics()
	healthMetrics := dk.NewHealthMetrics(actionID)

	dk.AntiMagicShellAura = dk.RegisterAura(core.Aura{
		Label:    "Anti-Magic Shell",
		ActionID: actionID,
		Duration: time.Millisecond * 200,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//if dk.Inputs.IsDps {
			// Setup a PA that deals damage to unit
			//}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			//rpGain := 20.0*sim.RandomFloat("Anti Magic Shell RP") + 54.0
			//dk.AddRunicPower(sim, rpG  ain, rpMetrics)
		},

		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage > 0 {
				absorvedDamage := core.MinFloat(0.75*spellEffect.Damage, 0.5*dk.MaxHealth())
				dk.GainHealth(sim, absorvedDamage, healthMetrics)
				dk.AddRunicPower(sim, absorvedDamage/69.0, rpMetrics)
			}

			aura.Deactivate(sim)
		},
	})

	if dk.Rotation.UseAms {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.AntiMagicShell.Spell,
			Priority: core.CooldownPriorityLow, // Use low prio so other actives get used first.
			Type:     core.CooldownTypeDPS,
		})
	}
}
