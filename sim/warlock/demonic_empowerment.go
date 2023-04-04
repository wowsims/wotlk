package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerDemonicEmpowermentSpell() {
	if !warlock.Talents.DemonicEmpowerment {
		return
	}

	var petAura core.Aura
	switch warlock.Options.Summon {
	case proto.Warlock_Options_Imp:
		petAura = core.Aura{
			Label:    "Demonic Empowerment Aura",
			ActionID: core.ActionID{SpellID: 47193},
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.Pet.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCrit: 20 * core.CritRatingPerCritChance,
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.Pet.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCrit: -20 * core.CritRatingPerCritChance,
				})
			},
		}
	case proto.Warlock_Options_Felguard:
		petAura = core.Aura{
			Label:    "Demonic Empowerment Aura",
			ActionID: core.ActionID{SpellID: 47193},
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.Pet.MultiplyAttackSpeed(sim, 1.2)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.Pet.MultiplyAttackSpeed(sim, 1/1.2)
			},
		}
	default:
		petAura = core.Aura{
			Label:    "Demonic Empowerment Aura",
			ActionID: core.ActionID{SpellID: 47193},
			Duration: time.Second * 15,
		}
	}

	if warlock.Options.Summon != proto.Warlock_Options_NoSummon {
		warlock.Pet.DemonicEmpowermentAura = warlock.Pet.RegisterAura(petAura)
	}

	warlock.DemonicEmpowerment = warlock.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 47193},

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Duration(180-6*warlock.Talents.Nemesis) * time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warlock.Pet.DemonicEmpowermentAura.Activate(sim)
		},
	})
}
