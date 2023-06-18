package core

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// This is just the static bonuses. Most professions are handled elsewhere.
func (character *Character) applyProfessionEffects() {
	if character.HasProfession(proto.Profession_Mining) {
		character.AddStat(stats.Stamina, 60)
	}

	if character.HasProfession(proto.Profession_Skinning) {
		character.AddStats(stats.Stats{stats.MeleeCrit: 40, stats.SpellCrit: 40})
	}

	if character.HasProfession(proto.Profession_Herbalism) {
		actionID := ActionID{SpellID: 55503}
		healthMetrics := character.NewHealthMetrics(actionID)

		spell := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolNature,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				amount := (3600 + character.MaxHealth()*0.016) / 5
				StartPeriodicAction(sim, PeriodicActionOptions{
					Period:   time.Second,
					NumTicks: 5,
					OnAction: func(sim *Simulation) {
						character.GainHealth(sim, amount*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
					},
				})
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Type:  CooldownTypeSurvival,
			Spell: spell,
		})
	}
}
