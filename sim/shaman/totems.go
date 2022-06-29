package shaman

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
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
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*120
		shaman.tryTwistWindfury(sim)
	}
	shaman.WrathOfAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerGraceOfAirTotemSpell() {
	config := shaman.newTotemSpellConfig(310.0, 25359)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*120
		shaman.tryTwistWindfury(sim)
	}
	shaman.GraceOfAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTranquilAirTotemSpell() {
	baseCost := shaman.BaseMana() * 0.06
	config := shaman.newTotemSpellConfig(baseCost, 25908)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*120
		shaman.tryTwistWindfury(sim)
	}
	shaman.TranquilAirTotem = shaman.RegisterSpell(config)
}

var windfuryTotemBaseManaCosts = []float64{
	95,
	140,
	200,
	275,
	325,
}

func (shaman *Shaman) registerWindfuryTotemSpell(rank int32) {
	if rank == 0 {
		// This will happen if we're not casting windfury totem. Just return a rank 1
		// template so we don't error.
		rank = 1
	}

	baseCost := windfuryTotemBaseManaCosts[rank-1]
	spellID := core.WindfuryTotemSpellRanks[rank-1]
	config := shaman.newTotemSpellConfig(baseCost, spellID)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*120
		shaman.tryTwistWindfury(sim)
	}
	shaman.WindfuryTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) tryTwistWindfury(sim *core.Simulation) {
	if !shaman.Totems.TwistWindfury {
		return
	}

	if shaman.Metrics.WentOOM && shaman.CurrentManaPercent() < 0.2 {
		shaman.NextTotemDropType[AirTotem] = int32(shaman.Totems.Air)
		return
	}

	// Swap to WF if we didn't just cast it, otherwise drop the other air totem immediately.
	if shaman.NextTotemDropType[AirTotem] != int32(proto.AirTotem_WindfuryTotem) {
		shaman.NextTotemDropType[AirTotem] = int32(proto.AirTotem_WindfuryTotem)
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second*10 // 10s until you need to drop WF
	} else {
		shaman.NextTotemDropType[AirTotem] = int32(shaman.Totems.Air)
		shaman.NextTotemDrops[AirTotem] = sim.CurrentTime + time.Second // drop immediately
	}
}

func (shaman *Shaman) tryTwistFireNova(sim *core.Simulation) {
	if !shaman.Totems.TwistFireNova {
		return
	}

	if shaman.Metrics.WentOOM && shaman.CurrentManaPercent() < 0.2 {
		shaman.NextTotemDropType[FireTotem] = int32(shaman.Totems.Fire)
		return
	}

	if shaman.NextTotemDropType[FireTotem] != int32(proto.FireTotem_FireNovaTotem) ||
		shaman.Totems.Fire == proto.FireTotem_NoFireTotem {
		shaman.NextTotemDropType[FireTotem] = int32(proto.FireTotem_FireNovaTotem)
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + shaman.FireNovaTotem.TimeToReady(sim)
	} else {
		shaman.NextTotemDropType[FireTotem] = int32(shaman.Totems.Fire)
	}
}

func (shaman *Shaman) registerManaSpringTotemSpell() {
	config := shaman.newTotemSpellConfig(120, 25570)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[WaterTotem] = sim.CurrentTime + time.Second*120
	}
	shaman.ManaSpringTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTotemOfWrathSpell() {
	baseCost := shaman.BaseMana() * 0.05
	config := shaman.newTotemSpellConfig(baseCost, 30706)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + time.Second*120
		shaman.tryTwistFireNova(sim)
	}
	shaman.TotemOfWrath = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStrengthOfEarthTotemSpell() {
	config := shaman.newTotemSpellConfig(300, 25528)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*120
	}
	shaman.StrengthOfEarthTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTremorTotemSpell() {
	config := shaman.newTotemSpellConfig(60, 8143)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.NextTotemDrops[EarthTotem] = sim.CurrentTime + time.Second*120
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
				case proto.AirTotem_GraceOfAirTotem:
					spell = shaman.GraceOfAirTotem
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
				case proto.FireTotem_FireNovaTotem:
					spell = shaman.FireNovaTotem
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
