package mage

import (
	"time"

	"github.com/wowsims/classic/sod/sim/core"
	"github.com/wowsims/classic/sod/sim/core/proto"
)

// TODO: Classic verify Arcane Blast rune numbers
// https://www.wowhead.com/classic/news/patch-1-15-build-52124-ptr-datamining-season-of-discovery-runes-336044#news-post-336044
// https://www.wowhead.com/classic/spell=400574/arcane-blast
func (mage *Mage) registerArcaneBlastSpell() {
	if !mage.HasRune(proto.MageRune_RuneHandsArcaneBlast) {
		return
	}

	level := float64(mage.GetCharacter().Level)
	baseCalc := (13.828124 + 0.018012*level + 0.044141*level*level)
	baseLowDamage := baseCalc * 4.53
	baseHighDamage := baseCalc * 5.27

	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast Aura",
		ActionID:  core.ActionID{SpellID: 400573},
		Duration:  time.Second * 6,
		MaxStacks: 4,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Refresh(sim)
			mage.ArcaneBlast.CostMultiplier = 1.75 * float64(newStacks)
		},
	})

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 400574},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		CritMultiplier:   mage.DefaultHealingCritMultiplier(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.15*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + .714*spell.SpellPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				if !mage.ArcaneBlastAura.IsActive() {
					mage.ArcaneBlastAura.Activate(sim)
				}
				if mage.ArcaneBlastAura.GetStacks() == mage.ArcaneBlastAura.MaxStacks {
					mage.ArcaneBlastAura.Refresh(sim)
				}
				mage.ArcaneBlastAura.AddStack(sim)
			}
		},
	})
}
