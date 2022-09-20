package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

//Totem IDs
const (
	VentureCoFlameSlicer                      = 38367
	DeadlyGladiatorsTotemOfIndomitability     = 42607
	FuriousGladiatorsTotemOfIndomitability    = 42608
	RelentlessGladiatorsTotemOfIndomitability = 42609
	WrathfulGladiatorsTotemOfIndomitability   = 51507
)

func (shaman *Shaman) registerLavaLashSpell() {
	if !shaman.Talents.LavaLash {
		return
	}

	manaCost := 0.04 * shaman.BaseMana

	flatDamageBonus := core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == VentureCoFlameSlicer, 25, 0)
	offhandFlametongueImbued := false
	if shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon || shaman.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeaponDownrank {
		offhandFlametongueImbued = true
	}

	var indomitabilityAura *core.Aura
	switch shaman.Equip[items.ItemSlotRanged].ID {
	case DeadlyGladiatorsTotemOfIndomitability:
		indomitabilityAura = shaman.NewTemporaryStatsAura("Deadly Aggression", core.ActionID{SpellID: 60549}, stats.Stats{stats.AttackPower: 120}, time.Second*10)
	case FuriousGladiatorsTotemOfIndomitability:
		indomitabilityAura = shaman.NewTemporaryStatsAura("Furious Gladiator's Libram of Fortitute", core.ActionID{SpellID: 60551}, stats.Stats{stats.AttackPower: 144}, time.Second*10) //wowhead is wierd about this one, might be the same in-game idk
	case RelentlessGladiatorsTotemOfIndomitability:
		indomitabilityAura = shaman.NewTemporaryStatsAura("Relentless Aggression", core.ActionID{SpellID: 60553}, stats.Stats{stats.AttackPower: 172}, time.Second*10)
	case WrathfulGladiatorsTotemOfIndomitability:
		indomitabilityAura = shaman.NewTemporaryStatsAura("Fury of the Gladiator", core.ActionID{SpellID: 60555}, stats.Stats{stats.AttackPower: 204}, time.Second*10)
	}

	shaman.LavaLash = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 60103},
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1 + core.TernaryFloat64(offhandFlametongueImbued,
			core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLavaLash), 0.35, 0.25), 0)*
			core.TernaryFloat64(shaman.HasSetBonus(ItemSetWorldbreakerBattlegear, 2), 1.2, 1),
		ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeOHSpecial,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.OffHand, false, flatDamageBonus, true),
			OutcomeApplier: shaman.OutcomeFuncMeleeSpecialHitAndCrit(shaman.ElementalCritMultiplier(0)),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() { //TODO: verify that it actually needs to hit
					return
				}
				if indomitabilityAura != nil {
					indomitabilityAura.Activate(sim)
				}
			},
		}),
	})
}

func (shaman *Shaman) IsLavaLashCastable(sim *core.Simulation) bool {
	return shaman.LavaLash.IsReady(sim)
}
