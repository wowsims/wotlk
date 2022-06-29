package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Maybe could switch "rank" parameter type to some proto thing. Would require updates to proto files.
// Prot guys do whatever you want here I guess
func (paladin *Paladin) RegisterConsecrationSpell(rank int32) {
	var manaCost float64
	var actionID core.ActionID
	var baseDamage float64

	switch rank {
	case 6:
		manaCost = 660
		actionID.SpellID = 27173
		baseDamage = 64
	case 4:
		manaCost = 390
		actionID.SpellID = 20923
		baseDamage = 35
	case 1:
		manaCost = 120
		actionID.SpellID = 26573
		baseDamage = 8
	default:
		manaCost = 0.0
	}

	switch paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID {
	case 27917:
		baseDamage += (47 * 0.952) // applies 47 "spell power" to the spell
	}

	// Check for bad input
	if manaCost == 0.0 {
		panic("Undefined Consecration rank specified.")
	}

	consecrationDot := core.NewDot(core.Dot{
		Aura: paladin.RegisterAura(core.Aura{
			Label:    "Consecration",
			ActionID: actionID,
		}),
		NumberOfTicks: 8,
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncAOESnapshot(paladin.Env, core.SpellEffect{
			ProcMask:        core.ProcMaskEmpty,
			BonusSpellPower: core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 27917, 47*0.8, 0),

			DamageMultiplier: 1 *
				core.TernaryFloat64(ItemSetLightbringerArmor.CharacterHasSetBonus(&paladin.Character, 4), 1.1, 1),
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(baseDamage, 0.119),
			OutcomeApplier:   paladin.OutcomeFuncMagicHit(),
			IsPeriodic:       true,
		}),
	})

	paladin.Consecration = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDot(consecrationDot),
	})

	consecrationDot.Spell = paladin.Consecration
}
