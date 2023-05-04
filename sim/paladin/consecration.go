package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Maybe could switch "rank" parameter type to some proto thing. Would require updates to proto files.
// Prot guys do whatever you want here I guess
func (paladin *Paladin) registerConsecrationSpell() {
	// TODO: Properly implement max rank consecration.
	bonusSpellPower := 0 +
		core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 27917, 47*0.8, 0) +
		core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40337, 141, 0) // Libram of Resurgence

	paladin.Consecration = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48819},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.22,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: (time.Second * 8) + core.TernaryDuration(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration), time.Second*2, 0),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Consecration",
			},
			NumberOfTicks: 8 + core.TernaryInt32(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration), 2, 0),
			TickLength:    time.Second * 1,

			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := paladin.CurrentTarget

				// i = 113 + 0.04*HolP + 0.04*AP
				dot.SnapshotBaseDamage = 113 +
					.04*dot.Spell.MeleeAttackPower() +
					.04*(dot.Spell.SpellPower()+bonusSpellPower)

				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
