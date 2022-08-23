package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerForceOfNatureCD() {
	if !druid.Talents.ForceOfNature {
		return
	}

	forceOfNatureAura := druid.RegisterAura(core.Aura{
		Label:    "Force of Nature",
		ActionID: core.ActionID{SpellID: 65861},
		Duration: time.Second * 30,
	})
	baseCost := druid.BaseMana * 0.12
	druid.ForceOfNature = druid.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 65861},

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.Treant1.EnableWithTimeout(sim, druid.Treant1, time.Second*30)
			druid.Treant2.EnableWithTimeout(sim, druid.Treant2, time.Second*30)
			druid.Treant3.EnableWithTimeout(sim, druid.Treant3, time.Second*30)
			forceOfNatureAura.Activate(sim)

			// Animation delay, courtesy of our DK friends
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Second*1,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
				},
			}
			sim.AddPendingAction(&pa)
		},
	})
}

type TreantPet struct {
	core.Pet
	druidOwner *Druid
}

func (druid *Druid) NewTreant() *TreantPet {
	treant := &TreantPet{
		Pet: core.NewPet(
			"Treant",
			&druid.Character,
			treantBaseStats,
			func(ownerStats stats.Stats) stats.Stats {
				return stats.Stats{
					stats.AttackPower: ownerStats[stats.SpellPower] * 2,
				}
			},
			false,
			false,
		),
		druidOwner: druid,
	}

	treant.EnableAutoAttacks(treant, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  59,
			BaseDamageMax:  87,
			SwingSpeed:     2,
			SwingDuration:  time.Second * 2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	druid.AddPet(treant)
	return treant
}

func (treant *TreantPet) GetPet() *core.Pet {
	return &treant.Pet
}

func (treant *TreantPet) Initialize() {
}

func (treant *TreantPet) Reset(sim *core.Simulation) {
}

func (treant *TreantPet) OnGCDReady(sim *core.Simulation) {
	treant.DoNothing()
}

// Eyeballing those TODO: get more data
var treantBaseStats = stats.Stats{
	stats.Stamina:   9600,
	stats.MeleeHit:  4 * core.MeleeHitRatingPerHitChance,
	stats.Expertise: 14 * core.ExpertisePerQuarterPercentReduction,
}
