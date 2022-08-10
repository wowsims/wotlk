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
	})

	rpMetrics := dk.AntiMagicShell.RunicPowerMetrics()

	dk.AntiMagicShellAura = dk.RegisterAura(core.Aura{
		Label:    "Anti-Magic Shell",
		ActionID: actionID,
		Duration: time.Millisecond * 200,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.UpdateExpires(sim.CurrentTime + time.Second*time.Duration(2.0*sim.RandomFloat("Anti Magic Shell Duration")))
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rpGain := 20.0*sim.RandomFloat("Anti Magic Shell RP") + 54.0
			dk.AddRunicPower(sim, rpGain, rpMetrics)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: dk.AntiMagicShell.Spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (dk *Deathknight) CanAntiMagicShell(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 20.0, 0, 0, 0) && dk.AntiMagicShell.IsReady(sim)
}

func (dk *Deathknight) CastAntiMagicShell(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanAntiMagicShell(sim) {
		return dk.AntiMagicShell.Cast(sim, target)
	}
	return false
}
