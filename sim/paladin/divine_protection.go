package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (paladin *Paladin) registerDivineProtectionSpell() {

	duration := time.Second*12 + core.TernaryDuration(paladin.HasSetBonus(ItemSetRedemptionPlate, 4), time.Second*3, 0)

	actionID := core.ActionID{SpellID: 498}
	paladin.DivineProtectionAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Protection",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier *= 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier /= 0.5
		},
	})

	cooldownDur := time.Minute*3 -
		30*time.Second*time.Duration(paladin.Talents.SacredDuty) -
		core.TernaryDuration(paladin.HasSetBonus(ItemSetTuralyonsPlate, 4), 30*time.Second, 0)

	paladin.DivineProtection = paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: cooldownDur,
			},
			SharedCD: core.Cooldown{
				Timer:    paladin.GetMutualLockoutDPAW(),
				Duration: 30*time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			paladin.DivineProtectionAura.Activate(sim)
			paladin.ForbearanceAura.Activate(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: paladin.DivineProtection,
		Type:  core.CooldownTypeSurvival,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			usable := !paladin.ForbearanceAura.IsActive()
			// Prevent Ret from screwing up their rotation by using this. TODO better logic
			if usable && paladin.Talents.TheArtOfWar > 0 {
				usable = false
			}
			return usable
		},
	})
}

func (paladin *Paladin) registerForbearanceDebuff() {
	
	actionID := core.ActionID{SpellID: 25771}
	duration := core.TernaryDuration(paladin.HasSetBonus(ItemSetTuralyonsPlate, 4), 90*time.Second, 120*time.Second)
	paladin.ForbearanceAura = paladin.RegisterAura(core.Aura{
		Label:    "Forbearance",
		ActionID: actionID,
		Duration: duration,
	})
	
}