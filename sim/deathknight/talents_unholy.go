package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyUnholyTalents() {
	// Vicious Strikes
	// TODO:

	// Virulence
	// TODO:

	// Epidemic
	// Implemented outside

	// Morbidity
	// TODO:

	// Ravenous Dead
	// TODO:

	// Outbreak
	// TODO: Add damage to SS when implemented. PS done

	// Necrosis
	// TODO:

	// Blood-Caked Blade
	// TODO:

	// Night of the Dead
	// TODO:

	// Unholy Blight
	// TODO:

	// Impurity
	// TODO:

	// Dirge
	// Implemented outside

	// Reaping
	// TODO:

	// Master of Ghouls
	// TODO:

	// Desolation
	deathKnight.applyDesolation()

	// Ghoul Frenzy
	// TODO:

	// Crypt Fever
	// TODO:

	// Bone Shield
	// TODO:

	// Wandering Plague
	// TODO:

	// Ebon Plaguebringer
	// TODO:

	// Scourge Strike
	// TODO:

	// Rage of Rivendare
	// TODO:

	// Summon Gargoyle
	// TODO:
}

func (deathKnight *DeathKnight) applyDesolation() {
	if deathKnight.Talents.Desolation == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 66803}

	deathKnight.DesolationAura = deathKnight.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Desolation",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(deathKnight.Talents.Desolation)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.0 + 0.01*float64(deathKnight.Talents.Desolation)
		},
	})
}
