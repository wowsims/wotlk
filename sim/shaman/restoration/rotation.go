package restoration

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (resto *RestorationShaman) OnGCDReady(sim *core.Simulation) {
	resto.tryUseGCD(sim)
}

func (resto *RestorationShaman) tryUseGCD(sim *core.Simulation) {

	// TODO: This could actually just be made as a PA that runs and triggers the shield instead of part of the rotation here.
	es := resto.EarthShield.Hot(resto.CurrentTarget)
	if es.IsActive() && resto.earthShieldPPM > 0 {
		procTime := time.Duration(60.0 / float64(resto.earthShieldPPM) * float64(time.Second))
		lastProc := resto.lastEarthShieldProc
		if lastProc+procTime < sim.CurrentTime {
			es.OnSpellHitTaken(es.Aura, sim, nil, &core.SpellResult{Outcome: core.OutcomeHit, Target: resto.CurrentTarget})
			resto.lastEarthShieldProc = sim.CurrentTime
		}
	}

	if resto.TryDropTotems(sim) {
		return
	}

	var spell *core.Spell
	switch resto.rotation.PrimaryHeal {
	case proto.ShamanHealSpell_AutoHeal:
		if len(resto.Party.Players) > 3 {
			spell = resto.ChainHeal
		} else {
			// TODO: lots of things to consider here...
			spell = resto.LesserHealingWave
		}
	case proto.ShamanHealSpell_LesserHealingWave:
		spell = resto.LesserHealingWave
	case proto.ShamanHealSpell_HealingWave:
		panic("healing wave not implemented yet")
		spell = resto.HealingWave
	case proto.ShamanHealSpell_ChainHeal:
		spell = resto.ChainHeal
	}

	if resto.rotation.UseEarthShield && !es.IsActive() {
		spell = resto.EarthShield
	} else if resto.rotation.UseRiptide && !resto.Riptide.Hot(resto.CurrentTarget).IsActive() && resto.Riptide.IsReady(sim) {
		spell = resto.Riptide
	}

	if !spell.Cast(sim, resto.CurrentTarget) {
		resto.WaitForMana(sim, spell.CurCast.Cost)
	}
}
