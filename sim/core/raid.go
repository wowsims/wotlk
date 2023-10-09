package core

import (
	"slices"

	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	googleProto "google.golang.org/protobuf/proto"
)

type Party struct {
	Raid  *Raid
	Index int

	Players []Agent
	Pets    []PetAgent // Cached list of all the pets in the party.

	PlayersAndPets []Agent // Cached list of players + pets, concatenated.

	dpsMetrics DistributionMetrics
	hpsMetrics DistributionMetrics
}

func NewParty(raid *Raid, index int, partyConfig *proto.Party) *Party {
	party := &Party{
		Raid:       raid,
		Index:      index,
		dpsMetrics: NewDistributionMetrics(),
		hpsMetrics: NewDistributionMetrics(),
	}

	for playerIndex, playerConfig := range partyConfig.Players {
		if playerConfig != nil && playerConfig.Class != proto.Class_ClassUnknown {
			party.Players = append(party.Players, NewAgent(party, playerIndex, playerConfig))
		}
	}

	return party
}

func (party *Party) Size() int {
	return len(party.Players)
}

func (party *Party) IsFull() bool {
	return party.Size() >= 5
}

func (party *Party) GetPartyBuffs(basePartyBuffs *proto.PartyBuffs) *proto.PartyBuffs {
	// Compute the full party buffs for this party.
	partyBuffs := &proto.PartyBuffs{}
	if basePartyBuffs != nil {
		partyBuffs = googleProto.Clone(basePartyBuffs).(*proto.PartyBuffs)
	}
	for _, player := range party.Players {
		player.AddPartyBuffs(partyBuffs)
		player.GetCharacter().AddPartyBuffs(partyBuffs)
	}
	return partyBuffs
}

func (party *Party) AddStats(newStats stats.Stats) {
	for _, agent := range party.Players {
		agent.GetCharacter().AddStats(newStats)
	}
}

func (party *Party) AddStat(stat stats.Stat, amount float64) {
	for _, agent := range party.Players {
		agent.GetCharacter().AddStat(stat, amount)
	}
}

func (party *Party) reset(sim *Simulation) {
	for _, agent := range party.Players {
		agent.GetCharacter().reset(sim, agent)
	}

	party.dpsMetrics.reset()
	party.hpsMetrics.reset()
}

func (party *Party) doneIteration(sim *Simulation) {
	for _, agent := range party.Players {
		agent.GetCharacter().doneIteration(sim)
		party.dpsMetrics.Total += agent.GetCharacter().Metrics.dps.Total
		party.hpsMetrics.Total += agent.GetCharacter().Metrics.hps.Total
	}

	party.dpsMetrics.doneIteration(sim)
	party.hpsMetrics.doneIteration(sim)
}

func (party *Party) GetMetrics() *proto.PartyMetrics {
	metrics := &proto.PartyMetrics{
		Dps: party.dpsMetrics.ToProto(),
		Hps: party.hpsMetrics.ToProto(),
	}

	playerIdx := 0
	i := 0
	for playerIdx < len(party.Players) {
		player := party.Players[playerIdx]
		if player.GetCharacter().PartyIndex == i {
			metrics.Players = append(metrics.Players, player.GetCharacter().GetMetricsProto())
			playerIdx++
		} else {
			metrics.Players = append(metrics.Players, &proto.UnitMetrics{})
		}
		i++
	}

	return metrics
}

type Raid struct {
	Parties []*Party

	dpsMetrics DistributionMetrics
	hpsMetrics DistributionMetrics

	AllPlayerUnits []*Unit // Cached list of all Players in the raid.
	AllUnits       []*Unit // Cached list of all Units (players and pets) in the raid.

	nextPetIndex int32

	replenishmentUnits         []*Unit   // All units who can receive replenishment.
	curReplenishmentUnits      [][]*Unit // Units that currently have replenishment active, separated by source.
	leftoverReplenishmentUnits []*Unit   // Units without replenishment currently active.
}

func (raid *Raid) GetActiveUnits() []*Unit {
	activeUnits := []*Unit{}
	for _, unit := range raid.AllUnits {
		if unit.IsActive() {
			activeUnits = append(activeUnits, unit)
		}
	}
	return activeUnits
}

// Makes a new raid.
func NewRaid(raidConfig *proto.Raid) *Raid {
	numParties := int(raidConfig.NumActiveParties)
	if numParties == 0 {
		numParties = len(raidConfig.Parties)
	}

	raid := &Raid{
		dpsMetrics:   NewDistributionMetrics(),
		hpsMetrics:   NewDistributionMetrics(),
		nextPetIndex: int32(numParties) * 5,
	}

	for partyIndex, partyConfig := range raidConfig.Parties {
		if partyConfig != nil && partyIndex < numParties {
			raid.Parties = append(raid.Parties, NewParty(raid, partyIndex, partyConfig))
		}
	}

	numDummies := min(24, int(raidConfig.TargetDummies))
	for i := 0; i < numDummies; i++ {
		party, partyIndex := raid.GetFirstEmptyRaidIndex()
		dummy := NewTargetDummy(i, party, partyIndex)
		party.Players = append(party.Players, dummy)
	}

	return raid
}

func (raid *Raid) Size() int {
	totalPlayers := 0
	for _, party := range raid.Parties {
		totalPlayers += party.Size()
	}
	return totalPlayers
}

// Returns (party, index within party)
func (raid *Raid) GetFirstEmptyRaidIndex() (*Party, int) {
	for _, party := range raid.Parties {
		if party.IsFull() {
			continue
		}

		for partyIndex := 0; partyIndex < 5; partyIndex++ {
			slotTaken := false
			for _, player := range party.Players {
				if player.GetCharacter().PartyIndex == partyIndex {
					slotTaken = true
				}
			}
			if !slotTaken {
				return party, partyIndex
			}
		}
	}

	panic("Raid is full")
}

func (raid *Raid) GetFirstTargetDummy() *TargetDummy {
	for _, party := range raid.Parties {
		for _, player := range party.Players {
			dummy, ok := player.(*TargetDummy)
			if ok {
				return dummy
			}
		}
	}
	return nil
}

func (raid *Raid) getNextPetIndex() int32 {
	petIndex := raid.nextPetIndex
	raid.nextPetIndex++
	return petIndex
}

func (raid *Raid) GetRaidBuffs(baseRaidBuffs *proto.RaidBuffs) *proto.RaidBuffs {
	// Compute the full raid buffs from the raid.
	raidBuffs := &proto.RaidBuffs{}
	if baseRaidBuffs != nil {
		raidBuffs = baseRaidBuffs
	}
	for _, party := range raid.Parties {
		for _, player := range party.Players {
			player.AddRaidBuffs(raidBuffs)
			player.GetCharacter().AddRaidBuffs(raidBuffs)
		}
	}
	return raidBuffs
}

// Precompute the playersAndPets array for each party.
func (raid *Raid) updatePlayersAndPets() {
	var raidPlayers []*Unit
	var raidPets []*Unit

	for _, party := range raid.Parties {
		party.Pets = []PetAgent{}
		for _, player := range party.Players {
			for _, petAgent := range player.GetCharacter().PetAgents {
				party.Pets = append(party.Pets, petAgent)
				raidPets = append(raidPets, &petAgent.GetPet().Unit)
			}
		}
		party.PlayersAndPets = make([]Agent, len(party.Players)+len(party.Pets))
		for i, player := range party.Players {
			party.PlayersAndPets[i] = player
			raidPlayers = append(raidPlayers, &player.GetCharacter().Unit)
		}
		for i, pet := range party.Pets {
			party.PlayersAndPets[len(party.Players)+i] = pet
		}
	}

	raid.AllUnits = append(raidPlayers, raidPets...)
	raid.AllPlayerUnits = raidPlayers

	slices.SortFunc(raid.AllUnits, func(u1, u2 *Unit) int {
		return int(u1.Index - u2.Index)
	})
	slices.SortFunc(raid.AllPlayerUnits, func(u1, u2 *Unit) int {
		return int(u1.Index - u2.Index)
	})
}

func (raid *Raid) applyCharacterEffects(raidConfig *proto.Raid) *proto.RaidStats {
	raidBuffs := raid.GetRaidBuffs(raidConfig.Buffs)
	raidStats := &proto.RaidStats{}

	for partyIdx, party := range raid.Parties {
		partyConfig := raidConfig.Parties[partyIdx]
		partyBuffs := party.GetPartyBuffs(partyConfig.Buffs)
		partyStats := &proto.PartyStats{
			Players: make([]*proto.PlayerStats, 5),
		}

		// Apply all buffs to the players in this party.
		for playerIdx, player := range party.Players {
			if playerIdx >= len(partyConfig.Players) {
				// This happens for target dummies.
				continue
			}
			playerConfig := partyConfig.Players[playerIdx]
			individualBuffs := &proto.IndividualBuffs{}
			if playerConfig.Buffs != nil {
				individualBuffs = playerConfig.Buffs
			}

			char := player.GetCharacter()
			char.EnableHealthBar()
			char.trackChanceOfDeath(playerConfig.HealingModel)
			partyStats.Players[char.PartyIndex] = char.applyAllEffects(player, raidBuffs, partyBuffs, individualBuffs)

			for _, pet := range char.Pets {
				pet.EnableHealthBar()
			}
		}

		raidStats.Parties = append(raidStats.Parties, partyStats)
	}

	return raidStats
}

func (raid *Raid) AddStats(s stats.Stats) {
	for _, party := range raid.Parties {
		party.AddStats(s)
	}
}

func (raid *Raid) GetPlayersOfClass(class proto.Class) []Agent {
	classPlayers := []Agent{}
	for _, party := range raid.Parties {
		for _, agent := range party.Players {
			if agent.GetCharacter().Class == class {
				classPlayers = append(classPlayers, agent)
			}
		}
	}
	return classPlayers
}

func (raid *Raid) GetPlayerFromUnit(unit *Unit) Agent {
	for _, party := range raid.Parties {
		for _, agent := range party.PlayersAndPets {
			if &agent.GetCharacter().Unit == unit {
				return agent
			}
		}
	}
	return nil
}

func (raid *Raid) GetFirstNPlayersOrPets(n int32) []*Unit {
	return raid.AllUnits[:min(n, int32(len(raid.AllUnits)))]
}

func (raid *Raid) GetPlayerFromUnitIndex(unitIndex int32) Agent {
	for _, party := range raid.Parties {
		for _, agent := range party.PlayersAndPets {
			if agent.GetCharacter().UnitIndex == unitIndex {
				return agent
			}
		}
	}
	return nil
}

func (raid *Raid) reset(sim *Simulation) {
	raid.resetReplenishment(sim)
	for _, party := range raid.Parties {
		party.reset(sim)
	}
	raid.dpsMetrics.reset()
	raid.hpsMetrics.reset()
}

func (raid *Raid) doneIteration(sim *Simulation) {
	for _, party := range raid.Parties {
		party.doneIteration(sim)
		raid.dpsMetrics.Total += party.dpsMetrics.Total
		raid.hpsMetrics.Total += party.hpsMetrics.Total
	}

	raid.dpsMetrics.doneIteration(sim)
	raid.hpsMetrics.doneIteration(sim)
}

func (raid *Raid) GetMetrics() *proto.RaidMetrics {
	metrics := &proto.RaidMetrics{
		Dps: raid.dpsMetrics.ToProto(),
		Hps: raid.hpsMetrics.ToProto(),
	}
	for _, party := range raid.Parties {
		metrics.Parties = append(metrics.Parties, party.GetMetrics())
	}
	return metrics
}

func SinglePlayerRaidProto(player *proto.Player, partyBuffs *proto.PartyBuffs, raidBuffs *proto.RaidBuffs, debuffs *proto.Debuffs) *proto.Raid {
	return &proto.Raid{
		Parties: []*proto.Party{
			{
				Players: []*proto.Player{
					player,
				},
				Buffs: partyBuffs,
			},
		},
		Buffs:   raidBuffs,
		Debuffs: debuffs,
	}
}

func RaidPlayersWithClass(raid *proto.Raid, class proto.Class) []*proto.Player {
	var players []*proto.Player
	for _, party := range raid.Parties {
		for _, player := range party.Players {
			if player != nil && player.Class == class {
				players = append(players, player)
			}
		}
	}
	return players
}
