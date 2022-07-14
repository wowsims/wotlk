package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyUnholyTalents() {
	// Vicious Strikes
	// Implemented outside

	// Virulence
	if deathKnight.Talents.Virulence > 0 {
		deathKnight.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*float64(deathKnight.Talents.Virulence))
	}

	// Epidemic
	// Implemented outside

	// Morbidity
	// TODO:

	// Ravenous Dead
	// TODO: Ghoul part
	if deathKnight.Talents.RavenousDead > 0 {
		strengthCoeff := 0.01 * float64(deathKnight.Talents.RavenousDead)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Strength,
			ModifiedStat: stats.Strength,
			Modifier: func(strength float64, _ float64) float64 {
				return strength * (1.0 + strengthCoeff)
			},
		})
	}

	// Outbreak
	// Implemented outside

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

	// Bone Shield
	// TODO:

	// Wandering Plague
	// TODO:

	// Crypt Fever
	// Ebon Plaguebringer
	// TODO: Diseases damage increase still missing
	if deathKnight.Talents.EbonPlaguebringer > 0 {
		deathKnight.PseudoStats.BonusMeleeCritRating += core.CritRatingPerCritChance * float64(deathKnight.Talents.EbonPlaguebringer)
		deathKnight.PseudoStats.BonusSpellCritRating += core.CritRatingPerCritChance * float64(deathKnight.Talents.EbonPlaguebringer)
	}

	// Scourge Strike
	// Implemented outside. Still missing shadow damage part

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
