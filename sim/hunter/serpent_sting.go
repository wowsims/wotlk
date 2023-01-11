package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	actionID := core.ActionID{SpellID: 49001}

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		Cost: core.NewManaCost(core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		}),
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		DamageMultiplierAdditive: 1 +
			0.1*float64(hunter.Talents.ImprovedStings) +
			core.TernaryFloat64(hunter.HasSetBonus(ItemSetScourgestalkerBattlegear, 2), .1, 0),
		CritMultiplier:   hunter.critMultiplier(false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)
			if result.Landed() {
				hunter.SerpentStingDot.Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	canCrit := hunter.HasSetBonus(ItemSetWindrunnersPursuit, 2)
	noxiousStingsMultiplier := 1 + 0.01*float64(hunter.Talents.NoxiousStings)
	huntersWithGlyphOfSteadyShot := hunter.GetAllHuntersWithGlyphOfSteadyShot()

	target := hunter.CurrentTarget
	hunter.SerpentStingDot = core.NewDot(core.Dot{
		Spell: hunter.SerpentSting,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SerpentSting-" + strconv.Itoa(int(hunter.Index)),
			Tag:      "SerpentSting",
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= noxiousStingsMultiplier
				// Check for 1 because this aura will always be active inside OnGain.
				if aura.Unit.NumActiveAurasWithTag("SerpentSting") == 1 {
					for _, otherHunter := range huntersWithGlyphOfSteadyShot {
						otherHunter.SteadyShot.DamageMultiplierAdditive += .1
					}
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= noxiousStingsMultiplier
				if !aura.Unit.HasActiveAuraWithTag("SerpentSting") {
					for _, otherHunter := range huntersWithGlyphOfSteadyShot {
						otherHunter.SteadyShot.DamageMultiplierAdditive -= .1
					}
				}
			},
		}),
		NumberOfTicks: 5 + core.TernaryInt32(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSerpentSting), 2, 0),
		TickLength:    time.Second * 3,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotBaseDamage = 242 + 0.04*dot.Spell.RangedAttackPower(target)
			if !isRollover {
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(target, attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			}
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			if canCrit {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			}
		},
	})
}

func (hunter *Hunter) GetAllHuntersWithGlyphOfSteadyShot() []*Hunter {
	allHunterAgents := hunter.Env.Raid.GetPlayersOfClass(proto.Class_ClassHunter)

	hunters := []*Hunter{}
	for _, agent := range allHunterAgents {
		h := agent.(HunterAgent).GetHunter()
		if h.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSteadyShot) {
			hunters = append(hunters, h)
		}
	}
	return hunters
}
