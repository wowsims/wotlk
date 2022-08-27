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

			bonusID := core.ActionID{ItemID: 11133}
			bonusStats := stats.Stats{stats.Strength: druid.GetStat(stats.SpellPower) * 0.5}

			druid.Treant1.NewTemporaryStatsAura("SP Snapshot", bonusID, bonusStats, time.Second*30).Activate(sim)
			druid.Treant2.NewTemporaryStatsAura("SP Snapshot", bonusID, bonusStats, time.Second*30).Activate(sim)
			druid.Treant3.NewTemporaryStatsAura("SP Snapshot", bonusID, bonusStats, time.Second*30).Activate(sim)

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
				return stats.Stats{}
			},
			false,
			false,
		),
		druidOwner: druid,
	}
	treant.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	treant.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/83.3)
	treant.EnableAutoAttacks(treant, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  252,
			BaseDamageMax:  357,
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

// TODO : fix miss/dodge
var treantBaseStats = stats.Stats{
	stats.Strength:  331,
	stats.Agility:   113,
	stats.Stamina:   598,
	stats.Intellect: 281,
	stats.Spirit:    109,
	stats.MeleeCrit: 5 * core.CritRatingPerCritChance,
	stats.MeleeHit:  5 * core.MeleeHitRatingPerHitChance,
	stats.Expertise: 120,
}
