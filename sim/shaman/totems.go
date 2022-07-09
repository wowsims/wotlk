package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) newTotemSpellConfig(baseCost float64, spellID int32) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellID},
		Flags:    SpellFlagTotem,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost -
					(baseCost * float64(shaman.Talents.TotemicFocus) * 0.05) -
					(baseCost * float64(shaman.Talents.MentalQuickness) * 0.02),
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
	}
}

func (shaman *Shaman) registerWrathOfAirTotemSpell() {
	config := shaman.newTotemSpellConfig(320.0, 3738)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.WrathOfAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerGraceOfAirTotemSpell() {
	config := shaman.newTotemSpellConfig(310.0, 25359)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.GraceOfAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTranquilAirTotemSpell() {
	baseCost := shaman.BaseMana() * 0.06
	config := shaman.newTotemSpellConfig(baseCost, 25908)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.TranquilAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerWindfuryTotemSpell() {
	config := shaman.newTotemSpellConfig(baseMana*0.11, 8512)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.WindfuryTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerManaSpringTotemSpell() {
	config := shaman.newTotemSpellConfig(baseMana*0.04, 58774)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[WaterTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.ManaSpringTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTotemOfWrathSpell() {
	config := shaman.newTotemSpellConfig(baseMana*0.05, 57722)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.TotemOfWrath = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStrengthOfEarthTotemSpell() {
	config := shaman.newTotemSpellConfig(300, 25528)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.StrengthOfEarthTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTremorTotemSpell() {
	config := shaman.newTotemSpellConfig(60, 8143)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.TremorTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) NextTotemAt(sim *core.Simulation) time.Duration {
	nextTotemAt := core.MinDuration(
		shaman.NextTotemDrops[0],
		core.MinDuration(
			shaman.NextTotemDrops[1],
			core.MinDuration(
				shaman.NextTotemDrops[2],
				shaman.NextTotemDrops[3])))

	return nextTotemAt
}

// TryDropTotems will check to see if totems need to be re-cast.
//  Returns whether we tried to cast a totem, regardless of whether it succeeded.
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
				case proto.AirTotem_TranquilAirTotem:
					spell = shaman.TranquilAirTotem
				}

			case EarthTotem:
				switch proto.EarthTotem(nextDrop) {
				case proto.EarthTotem_StrengthOfEarthTotem:
					spell = shaman.StrengthOfEarthTotem
				case proto.EarthTotem_TremorTotem:
					spell = shaman.TremorTotem
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
					// spell = shaman.FlametongueTotem
				}

			case WaterTotem:
				spell = shaman.ManaSpringTotem
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
