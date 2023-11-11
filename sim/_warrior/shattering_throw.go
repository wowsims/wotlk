package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) RegisterShatteringThrowCD() {
	hasGlyph := warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfShatteringThrow)
	shattDebuffs := warrior.NewEnemyAuraArray(core.ShatteringThrowAura)

	ShatteringThrowSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 64382},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		RageCost: core.RageCostOptions{
			Cost: 25 - float64(warrior.Talents.FocusedRage),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: core.TernaryDuration(hasGlyph, time.Duration(0), time.Millisecond*1500),
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 5,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if warrior.AutoAttacks.MH().SwingSpeed == warrior.AutoAttacks.OH().SwingSpeed {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, true)
				} else {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, false)
				}
				if !hasGlyph && !warrior.StanceMatches(BattleStance) && warrior.BattleStance.IsReady(sim) {
					warrior.BattleStance.Cast(sim, nil)
				}
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance) || warrior.BattleStance.IsReady(sim) || hasGlyph
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 12 + 0.5*spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if result.Landed() {
				shattDebuffs.Get(target).Activate(sim)
			}
		},

		RelatedAuras: []core.AuraArray{shattDebuffs},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ShatteringThrowSpell,
		Type:  core.CooldownTypeDPS,
	})
}
