package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// const JudgementDuration = time.Second * 20

// Shared conditions required to be able to cast any Judgement.
func (paladin *Paladin) canJudgement(sim *core.Simulation) bool {
	return paladin.CurrentSeal != nil && paladin.CurrentSeal.IsActive() && paladin.JudgementOfLight.IsReady(sim)
}

func (paladin *Paladin) registerJudgementOfWisdomSpell(cdTimer *core.Timer) {
	// paladin.JudgementOfLightAura = core.JudgementOfLightAura(paladin.CurrentTarget)

	baseCost := paladin.BaseMana * 0.05

	paladin.JudgementOfWisdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 53408},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        SpellFlagPrimaryJudgement,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: cdTimer,
				Duration: (time.Second * 10) -
					(time.Second * time.Duration(paladin.Talents.ImprovedJudgements)) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetRedemptionBattlegear, 4), 1*time.Second, 0) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetGladiatorsVindication, 4), 1*time.Second, 0),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Primary Judgements cannot crit or be dodged, parried, or blocked-- only miss. (Unless target is a hunter.)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)
		},
	})
}

func (paladin *Paladin) registerJudgementOfLightSpell(cdTimer *core.Timer) {
	// paladin.JudgementOfLightAura = core.JudgementOfLightAura(paladin.CurrentTarget)

	baseCost := paladin.BaseMana * 0.05

	paladin.JudgementOfLight = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 20271},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        SpellFlagPrimaryJudgement,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: cdTimer,
				Duration: (time.Second * 10) -
					(time.Second * time.Duration(paladin.Talents.ImprovedJudgements)) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetRedemptionBattlegear, 4), 1*time.Second, 0) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetGladiatorsVindication, 4), 1*time.Second, 0),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Primary Judgements cannot crit or be dodged, parried, or blocked-- only miss. (Unless target is a hunter.)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)
		},
	})
}

// Defines judgement refresh behavior from attacks
// Returns extra mana if a different pally applied Judgement of Wisdom
// func (paladin *Paladin) setupJudgementRefresh() {
// 	const mana = 74 / 2
// 	paladin.RegisterAura(core.Aura{
// 		Label:    "Refresh Judgement",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
// 			if spellEffect.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
// 				if paladin.CurrentJudgement != nil && paladin.CurrentJudgement.IsActive() {
// 					// Refresh the judgement
// 					paladin.CurrentJudgement.Refresh(sim)

// 					// Check if current judgement is not JoW and also that JoW is on the target
// 					if paladin.CurrentJudgement != paladin.JudgementOfWisdomAura && paladin.JudgementOfWisdomAura.IsActive() {
// 						// Just trigger a second JoW
// 						if paladin.JowManaMetrics == nil {
// 							paladin.JowManaMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 27164})
// 						}
// 						paladin.AddMana(sim, mana, paladin.JowManaMetrics, false)
// 					}
// 				}
// 			}
// 		},
// 	})
// }

func (paladin *Paladin) registerJudgements() {
	// Shared CD for all judgements.
	cdTimer := paladin.NewTimer()
	paladin.registerJudgementOfWisdomSpell(cdTimer)
	paladin.registerJudgementOfLightSpell(cdTimer)
}
