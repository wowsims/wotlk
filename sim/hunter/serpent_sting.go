package hunter

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	actionID := core.ActionID{SpellID: 49001}
	baseCost := 0.09 * hunter.BaseMana

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskEmpty,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(hunter.Talents.Efficiency)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		DamageMultiplierAdditive: 1 +
			0.1*float64(hunter.Talents.ImprovedStings) +
			core.TernaryFloat64(hunter.HasSetBonus(ItemSetScourgestalkerBattlegear, 2), .1, 0),
		CritMultiplier:   hunter.critMultiplier(false, false, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: hunter.OutcomeFuncRangedHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					hunter.SerpentStingDot.Apply(sim)
				}
			},
		}),
	})

	dotOutcome := hunter.OutcomeFuncTick()
	if hunter.HasSetBonus(ItemSetWindrunnersPursuit, 2) {
		dotOutcome = hunter.OutcomeFuncMeleeSpecialCritOnly()
	}

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
				hunter.AttackTables[aura.Unit.UnitIndex].DamageDealtMultiplier *= noxiousStingsMultiplier
				// Check for 1 because this aura will always be active inside OnGain.
				if aura.Unit.NumActiveAurasWithTag("SerpentSting") == 1 {
					for _, otherHunter := range huntersWithGlyphOfSteadyShot {
						otherHunter.SteadyShot.DamageMultiplierAdditive += .1
					}
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].DamageDealtMultiplier /= noxiousStingsMultiplier
				if !aura.Unit.HasActiveAuraWithTag("SerpentSting") {
					for _, otherHunter := range huntersWithGlyphOfSteadyShot {
						otherHunter.SteadyShot.DamageMultiplierAdditive -= .1
					}
				}
			},
		}),
		NumberOfTicks: 5 + int(core.TernaryInt32(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSerpentSting), 2, 0)),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic: true,

			BaseDamage: core.BuildBaseDamageConfig(func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
				return 242 + 0.04*spell.RangedAttackPower(spellEffect.Target)
			}, 0),
			OutcomeApplier: dotOutcome,
		}),
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
