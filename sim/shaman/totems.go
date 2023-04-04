package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) newTotemSpellConfig(baseCost float64, spellID int32) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellID},
		Flags:    SpellFlagTotem,

		ManaCost: core.ManaCostOptions{
			BaseCost: baseCost,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.TotemicFocus) -
				0.02*float64(shaman.Talents.MentalQuickness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Second,
				GCD:      time.Second,
			},
			IgnoreHaste: true,
		},
	}
}

func (shaman *Shaman) registerWrathOfAirTotemSpell() {
	config := shaman.newTotemSpellConfig(0.11, 3738)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.WrathOfAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerWindfuryTotemSpell() {
	config := shaman.newTotemSpellConfig(0.11, 8512)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.WindfuryTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerManaSpringTotemSpell() {
	config := shaman.newTotemSpellConfig(0.04, 58774)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[WaterTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.ManaSpringTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerHealingStreamTotemSpell() {
	config := shaman.newTotemSpellConfig(0.03, 58757)
	hsHeal := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 52042},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,
		DamageMultiplier: 1 + (.02 * float64(shaman.Talents.Purification)) + 0.15*float64(shaman.Talents.RestorativeTotems),
		CritMultiplier:   1,
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: find healing stream coeff
			healing := 25 + spell.HealingPower(target)*0.08272
			spell.CalcAndDealHealing(sim, target, healing, spell.OutcomeHealing)
		},
	})
	config.Hot = core.DotConfig{
		Aura: core.Aura{
			Label: "HealingStreamHot",
		},
		NumberOfTicks: 150,
		TickLength:    time.Second * 2,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			hsHeal.Cast(sim, target)
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.NextTotemDrops[WaterTotem] = sim.CurrentTime + time.Second*300
		for _, agent := range shaman.Party.Players {
			spell.Hot(&agent.GetCharacter().Unit).Activate(sim)
		}
	}
	shaman.HealingStreamTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTotemOfWrathSpell() {
	config := shaman.newTotemSpellConfig(0.05, 57722)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*300
		shaman.applyToWDebuff(sim)
	}
	shaman.TotemOfWrath = shaman.RegisterSpell(config)
}

func (shaman *Shaman) applyToWDebuff(sim *core.Simulation) {
	for _, target := range sim.Encounter.TargetUnits {
		auraDef := core.TotemOfWrathDebuff(target)
		auraDef.Activate(sim)
	}
}

func (shaman *Shaman) registerFlametongueTotemSpell() {
	config := shaman.newTotemSpellConfig(0.11, 58656)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.FlametongueTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStrengthOfEarthTotemSpell() {
	config := shaman.newTotemSpellConfig(0.1, 58643)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.StrengthOfEarthTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTremorTotemSpell() {
	config := shaman.newTotemSpellConfig(0.02, 8143)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.TremorTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStoneskinTotemSpell() {
	config := shaman.newTotemSpellConfig(0.1, 58753)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.StoneskinTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) NextTotemAt(_ *core.Simulation) time.Duration {
	nextTotemAt := core.MinDuration(
		core.MinDuration(shaman.NextTotemDrops[0], shaman.NextTotemDrops[1]),
		core.MinDuration(shaman.NextTotemDrops[2], shaman.NextTotemDrops[3]))

	return nextTotemAt
}

// TryDropTotems will check to see if totems need to be re-cast.
//
//	Returns whether we tried to cast a totem, regardless of whether it succeeded.
func (shaman *Shaman) TryDropTotems(sim *core.Simulation) bool {
	var spell *core.Spell

	for totemTypeIdx, totemExpiration := range shaman.NextTotemDrops {
		if spell != nil {
			break
		}
		nextDrop := shaman.NextTotemDropType[totemTypeIdx]
		if sim.CurrentTime >= totemExpiration {
			switch totemTypeIdx {
			case AirTotem:
				switch proto.AirTotem(nextDrop) {
				case proto.AirTotem_WrathOfAirTotem:
					spell = shaman.WrathOfAirTotem
				case proto.AirTotem_WindfuryTotem:
					spell = shaman.WindfuryTotem
				}

			case EarthTotem:
				switch proto.EarthTotem(nextDrop) {
				case proto.EarthTotem_StrengthOfEarthTotem:
					spell = shaman.StrengthOfEarthTotem
				case proto.EarthTotem_TremorTotem:
					spell = shaman.TremorTotem
				case proto.EarthTotem_StoneskinTotem:
					spell = shaman.StoneskinTotem
				}

			case FireTotem:
				switch proto.FireTotem(nextDrop) {
				case proto.FireTotem_TotemOfWrath:
					spell = shaman.TotemOfWrath
				case proto.FireTotem_SearingTotem:
					spell = shaman.SearingTotem
				case proto.FireTotem_MagmaTotem:
					spell = shaman.MagmaTotem
				case proto.FireTotem_FlametongueTotem:
					spell = shaman.FlametongueTotem
				}

			case WaterTotem:
				switch proto.WaterTotem(nextDrop) {
				case proto.WaterTotem_ManaSpringTotem:
					spell = shaman.ManaSpringTotem
				case proto.WaterTotem_HealingStreamTotem:
					spell = shaman.HealingStreamTotem
				}
			}
		}
	}

	if spell != nil {
		if success := spell.Cast(sim, shaman.CurrentTarget); !success {
			shaman.WaitForMana(sim, spell.CurCast.Cost)
		}
		return true
	}
	return false
}
