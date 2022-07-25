package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Maybe could switch "rank" parameter type to some proto thing. Would require updates to proto files.
// Prot guys do whatever you want here I guess
func (paladin *Paladin) registerConsecrationSpell() {
	// TODO: Properly implement max rank consecration.

	baseCost := 0.22 * paladin.BaseMana
	actionID := core.ActionID{SpellID: 48819}

	consecrationDot := core.NewDot(core.Dot{
		Aura: paladin.RegisterAura(core.Aura{
			Label:    "Consecration",
			ActionID: actionID,
		}),
		NumberOfTicks: 8 + core.TernaryInt(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration), 2, 0),
		TickLength:    time.Second * 1,
		TickEffects: core.TickFuncAOESnapshot(paladin.Env, core.SpellEffect{
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					// i = 113 + 0.04*HolP + 0.04*AP
					scaling := hybridScaling{
						AP: 0.04,
						SP: 0.04,
					}

					sp := hitEffect.SpellPower(spell.Unit, spell) +
						core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 27917, 47*0.8, 0) +
						core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40337, 141, 0) // Libram of Resurgence

					damage := 113 + (scaling.AP * hitEffect.MeleeAttackPower(spell.Unit)) + (scaling.SP * sp)

					return damage
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMagicHit(),
			IsPeriodic:     true,
		}),
	})

	paladin.Consecration = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: (time.Second * 8) + core.TernaryDuration(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration), time.Second*2, 0),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDot(consecrationDot),
	})

	consecrationDot.Spell = paladin.Consecration
}
