package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerHauntSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneHandsHaunt) {
		return
	}

	actionID := core.ActionID{SpellID: 403501}
	debuffMult := 1.2

	spellCoeff := 0.714
	level := float64(warlock.GetCharacter().Level)
	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseLowDamage := baseCalc * 2.51
	baseHighDamage := baseCalc * 2.95

	warlock.HauntDebuffAuras = warlock.NewEnemyAuraArray(func(target *core.Unit, level int32) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Haunt-" + warlock.Label,
			ActionID: actionID,
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier *= debuffMult
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier /= debuffMult
			},
		})
	})

	warlock.Haunt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.12,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.SpellCritMultiplier(1, 0),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					warlock.HauntDebuffAuras.Get(result.Target).Activate(sim)
					warlock.EverlastingAfflictionRefresh(sim, target)
				}
			})
		},
	})
}
