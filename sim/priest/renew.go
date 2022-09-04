package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerRenewSpell() {
	actionID := core.ActionID{SpellID: 48068}
	baseCost := 0.17 * priest.BaseMana

	priest.Renew = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			priest.RenewHot.Apply(sim)
		},
	})

}

func (priest *Priest) makeRenewHot(target *core.Unit) *core.Dot {
	return core.NewDot(core.Dot{
		Spell: priest.Renew,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Renew" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.Renew.ActionID,
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			IsPeriodic:       true,
			IsHealing:        true,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(280, 0.2),
			OutcomeApplier:   priest.OutcomeFuncTick(),
			//OnPeriodicDamageDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			//},
		}),
	})
}
