package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (paladin *Paladin) registerHolyShieldSpell() {
	actionID := core.ActionID{SpellID: 27179}

	numCharges := 4 + 2*paladin.Talents.ImprovedHolyShield
	damage := 155 + core.TernaryFloat64(ItemSetJusticarArmor.CharacterHasSetBonus(&paladin.Character, 4), 15, 0)

	procSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID.WithTag(1),
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagBinary,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1 + 0.1*float64(paladin.Talents.ImprovedHolyShield),
			ThreatMultiplier: 1.35,

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(damage, 0.05),
			OutcomeApplier: paladin.OutcomeFuncMagicHitBinary(),
		}),
	})

	blockBonus := 30*core.BlockRatingPerBlockChance + core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 29388, 42, 0)

	paladin.HolyShieldAura = paladin.RegisterAura(core.Aura{
		Label:     "Holy Shield",
		ActionID:  actionID,
		Duration:  time.Second * 10,
		MaxStacks: numCharges,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, blockBonus)
			aura.SetStacks(sim, numCharges)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.Block, -blockBonus)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeBlock) {
				procSpell.Cast(sim, spell.Unit)
				aura.RemoveStack(sim)
			}
		},
	})

	baseCost := 280.0

	paladin.HolyShield = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.HolyShieldAura.Activate(sim)
		},
	})
}
