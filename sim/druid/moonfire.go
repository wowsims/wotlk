package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (druid *Druid) registerMoonfireSpell() {
	actionID := core.ActionID{SpellID: 26988}
	baseCost := 495.0

	druid.Moonfire = druid.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.03*float64(druid.Talents.Moonglow)),
				GCD:  core.GCDDefault,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellCritRating: float64(druid.Talents.ImprovedMoonfire) * 5 * core.SpellCritRatingPerCritChance,
			DamageMultiplier:     1 * (1 + 0.05*float64(druid.Talents.ImprovedMoonfire)) * (1 + 0.02*float64(druid.Talents.Moonfury)),
			ThreatMultiplier:     1,
			BaseDamage:           core.BaseDamageConfigMagic(305, 357, 0.15),
			OutcomeApplier:       druid.OutcomeFuncMagicHitAndCrit(druid.SpellCritMultiplier(1, 0.2*float64(druid.Talents.Vengeance))),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.MoonfireDot.Apply(sim)
				}
			},
		}),
	})

	target := druid.CurrentTarget
	druid.MoonfireDot = core.NewDot(core.Dot{
		Spell: druid.Moonfire,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Moonfire-" + strconv.Itoa(int(druid.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 4 + core.TernaryInt(ItemSetThunderheartRegalia.CharacterHasSetBonus(&druid.Character, 2), 1, 0),
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 * (1 + 0.05*float64(druid.Talents.ImprovedMoonfire)) * (1 + 0.02*float64(druid.Talents.Moonfury)),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(600/4, 0.13),
			OutcomeApplier:   druid.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})
}

func (druid *Druid) ShouldCastMoonfire(sim *core.Simulation, target *core.Unit, rotation proto.BalanceDruid_Rotation) bool {
	return rotation.Moonfire && !druid.MoonfireDot.IsActive()
}
