package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// The numbers in this file are VERY rough approximations based on logs.

func (mage *Mage) registerMirrorImageCD() {
	baseCost := mage.BaseMana * 0.1
	summonDuration := time.Second * 30

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

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 3,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				// Assume this is a pre-cast, so disable the GCD.
				// Probably should keep the mana cost but this matches the other mage sim.
				if sim.CurrentTime == 0 {
					cast.Cost = 0
					cast.GCD = 0
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.mirrorImage.EnableWithTimeout(sim, mage.mirrorImage, summonDuration)
			if t10Aura != nil {
				t10Aura.Activate(sim)
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    mage.MirrorImage,
		Priority: core.CooldownPriorityDrums + 1, // Always prefer to cast before drums or lust so the ele gets their benefits.
		Type:     core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentMana() >= mage.MirrorImage.DefaultCast.Cost
		},
	})
}

type MirrorImage struct {
	core.Pet

	// Water Ele almost never just stands still and spams like we want, it sometimes
	// does its own thing. This controls how much it does that.
	waitBetweenCasts time.Duration

	Frostbolt *core.Spell
	Fireblast *core.Spell
}

func (mage *Mage) NewMirrorImage() *MirrorImage {

	mirrorImage := &MirrorImage{
		Pet: core.NewPet(
			"Mirror Image",
			&mage.Character,
			mirrorImageBaseStats,
			mirrorImageInheritance,
			false,
			true,
		),
		waitBetweenCasts: time.Second * 0,
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

func (mi *MirrorImage) Reset(sim *core.Simulation) {
}

func (mi *MirrorImage) OnGCDReady(sim *core.Simulation) {
	spell := mi.Frostbolt
	if mi.Fireblast.CD.IsReady(sim) {
		spell = mi.Fireblast
	}

	if success := spell.Cast(sim, mi.CurrentTarget); !success {
		mi.Disable(sim)
	}
}

// These numbers are just rough guesses based on looking at some logs.
var mirrorImageBaseStats = stats.Stats{
	stats.Mana: 2000,
}

var mirrorImageInheritance = func(ownerStats stats.Stats) stats.Stats {
	// These numbers are just rough guesses based on looking at some logs.
	return ownerStats.DotProduct(stats.Stats{
		stats.Stamina:   1,
		stats.Intellect: 1,
		stats.Mana:      1,

		stats.SpellCrit: 1,
		stats.SpellHit:  1,

		stats.SpellPower: 0.33,
	})
}

func (mi *MirrorImage) registerFrostboltSpell() {
	baseCost := 90.0

	mi.Frostbolt = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 59638},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second*3 + mi.waitBetweenCasts, // extra wait time is pretty much cast time
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mi.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//3x damage for 3 mirror images
			baseDamage := 163*3 + 0.9*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (mi *MirrorImage) registerFireblastSpell() {
	baseCost := 120.0

	mi.Fireblast = mi.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 59637},
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDMin,
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
			baseDamage := 88*3 + 0.45*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
