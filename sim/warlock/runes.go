package warlock

import (
	"math"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) ApplyRunes() {
	warlock.applyDemonicTactics()
	warlock.applyDemonicPact()
}

func (warlock *Warlock) EverlastingAfflictionRefresh(sim *core.Simulation, target *core.Unit) {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsEverlastingAffliction) {
		return
	}

	if warlock.Corruption.Dot(target).IsActive() {
		warlock.Corruption.Dot(target).Rollover(sim)
	}
}

func (warlock *Warlock) applyDemonicTactics() {
	if !warlock.HasRune(proto.WarlockRune_RuneChestDemonicTactics) {
		return
	}

	warlock.AddStat(stats.MeleeCrit, 10*core.CritRatingPerCritChance)
	warlock.AddStat(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)

	if warlock.Pet != nil {
		pet := warlock.Pet.GetPet()
		pet.AddStat(stats.MeleeCrit, 10*core.CritRatingPerCritChance)
		pet.AddStat(stats.SpellCrit, 10*core.SpellCritRatingPerCritChance)
	}
}

func (warlock *Warlock) applyDemonicPact() {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsDemonicPact) {
		return
	}

	if warlock.Options.Summon == proto.WarlockOptions_NoSummon {
		return
	}

	icd := core.Cooldown{
		Timer:    warlock.NewTimer(),
		Duration: 1 * time.Second,
	}

	spellPower := max(warlock.GetStat(stats.SpellPower)*0.1, float64(warlock.Level)/2.0)
	demonicPactAuras := warlock.NewAllyAuraArray(func(u *core.Unit) *core.Aura {
		return core.DemonicPactAura(u, spellPower)
	})

	warlock.Pet.RegisterAura(core.Aura{
		Label:    "Demonic Pact Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PreviousTime = 0
			aura.Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.DidCrit() || !icd.IsReady(sim) {
				return
			}

			icd.Use(sim)

			spBonus := max(math.Round(warlock.GetStat(stats.SpellPower)*0.1), math.Round(float64(warlock.Level)/2))
			for _, dpAura := range demonicPactAuras {
				if dpAura != nil {
					dpAura.ExclusiveEffects[0].SetPriority(sim, spBonus)

					// Force expire/gain because of new sp bonus
					dpAura.Deactivate(sim)
					dpAura.Activate(sim)
				}
			}
		},
	})
}
