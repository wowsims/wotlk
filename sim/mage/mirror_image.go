package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// The numbers in this file are VERY rough approximations based on logs.

func (mage *Mage) registerMirrorImageCD() {
	var t10Aura *core.Aura
	if mage.HasSetBonus(ItemSetBloodmagesRegalia, 4) {
		t10Aura = mage.RegisterAura(core.Aura{
			Label:    "Mirror Image Bonus Damage T10 4PC",
			ActionID: core.ActionID{SpellID: 70748},
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.PseudoStats.DamageDealtMultiplier *= 1.18
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				mage.PseudoStats.DamageDealtMultiplier /= 1.18
			},
		})
	}

	mage.MirrorImage = mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 55342},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.mirrorImage.EnableWithTimeout(sim, mage.mirrorImage, time.Second*30)
			if t10Aura != nil {
				t10Aura.Activate(sim)
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    mage.MirrorImage,
		Priority: core.CooldownPriorityDrums + 1000, // Always prefer to cast before drums or lust so the ele gets their benefits.
		Type:     core.CooldownTypeDPS,
	})
}

type MirrorImage struct {
	core.Pet

	mageOwner *Mage

	Frostbolt *core.Spell
	Fireblast *core.Spell
}

func (mage *Mage) NewMirrorImage() *MirrorImage {
	mirrorImage := &MirrorImage{
		Pet:       core.NewPet("Mirror Image", &mage.Character, mirrorImageBaseStats, createMirrorImageInheritance(mage), false, true),
		mageOwner: mage,
	}
	mirrorImage.EnableManaBar()

	mage.AddPet(mirrorImage)

	return mirrorImage
}

func (mi *MirrorImage) GetPet() *core.Pet {
	return &mi.Pet
}

func (mi *MirrorImage) Initialize() {
	mi.registerFireblastSpell()
	mi.registerFrostboltSpell()
}

func (mi *MirrorImage) Reset(_ *core.Simulation) {
}

func (mi *MirrorImage) ExecuteCustomRotation(sim *core.Simulation) {
	spell := mi.Frostbolt
	if mi.Fireblast.CD.IsReady(sim) && sim.RandomFloat("MirrorImage FB") < 0.5 {
		spell = mi.Fireblast
	}

	if success := spell.Cast(sim, mi.CurrentTarget); !success {
		mi.Disable(sim)
	}
}

var mirrorImageBaseStats = stats.Stats{
	stats.Mana: 3000, // Unknown
}

var createMirrorImageInheritance = func(mage *Mage) func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHit: ownerStats[stats.SpellHit] - float64(mage.Talents.Precision),
			// seems to be about 8% baseline
			stats.SpellCrit:  8 * core.CritRatingPerCritChance,
			stats.SpellPower: ownerStats[stats.SpellPower] * 0.33,
		}
	}
}

func (mi *MirrorImage) registerFrostboltSpell() {
	numImages := core.TernaryFloat64(mi.mageOwner.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMirrorImage), 4, 3)

	mi.Frostbolt = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 59638},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//3x damage for 3 mirror images
			baseDamage := (163 + 0.3*spell.SpellPower()) * numImages
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

func (mi *MirrorImage) registerFireblastSpell() {
	numImages := core.TernaryFloat64(mi.mageOwner.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMirrorImage), 4, 3)

	mi.Fireblast = mi.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 59637},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    mi.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//3x damage for 3 mirror images
			baseDamage := (88 + 0.15*spell.SpellPower()) * numImages
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
