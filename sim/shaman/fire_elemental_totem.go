package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

const fireTotemDuration time.Duration = time.Second * 120

func (shaman *Shaman) registerFireElementalTotem() {
	if !shaman.Totems.UseFireElemental {
		return
	}

	actionID := core.ActionID{SpellID: 2894}

	fireElementalAura := shaman.RegisterAura(core.Aura{
		Label:    "Fire Elemental Totem",
		ActionID: actionID,
		Duration: fireTotemDuration,
	})

	shaman.FireElementalTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.23,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * time.Duration(core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFireElementalTotem), 5, 10)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			// TODO: ToW needs a unique buff/debuff aura for each raidmember/target.
			//  Otherwise we will be possibly disabling another ele shaman's ToW debuff/buff.
			if shaman.Totems.Fire == proto.FireTotem_TotemOfWrath {
				shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + fireTotemDuration
			} else if shaman.Totems.Fire != proto.FireTotem_NoFireTotem && !shaman.Totems.UseFireMcd {
				shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + fireTotemDuration
			}
			shaman.MagmaTotem.AOEDot().Cancel(sim)
			shaman.SearingTotem.Dot(shaman.CurrentTarget).Cancel(sim)

			shaman.FireElemental.EnableWithTimeout(sim, shaman.FireElemental, fireTotemDuration)

			//TODO handle more then one swap if the fight is greater then 5 mins, for now will just do the one.
			if shaman.FireElementalTotem.SpellMetrics[target.Index].Casts == 1 {
				shaman.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}, true)
			}

			// Add a dummy aura to show in metrics
			fireElementalAura.Activate(sim)
		},
	})

	//Enh has 1.5seconds GCD also, so just going to wait the normal 1.5 instead of using the dynamic Spell GCD
	var castWindow = 1550 * time.Millisecond
	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.FireElementalTotem,
		Type:  core.CooldownTypeUnknown,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			success := shaman.fireElementalSnapShot.CanSnapShot(sim, castWindow)

			if (sim.Encounter.Duration <= 120*time.Second && sim.CurrentTime >= 10*time.Second) || (sim.Encounter.Duration > 120*time.Second && sim.CurrentTime >= 20*time.Second) {
				success = true
			}

			if success {
				shaman.castFireElemental = true
				shaman.fireElementalSnapShot.ActivateMajorCooldowns(sim)
				shaman.fireElementalSnapShot.ResetProcTrackers()
				shaman.castFireElemental = false
			}

			return success
		},
	})
}
