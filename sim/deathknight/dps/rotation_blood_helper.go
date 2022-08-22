package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) blDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *deathknight.RuneSpell, costRunes bool, casts int) bool {
	return dk.shDiseaseCheck(sim, target, spell, costRunes, casts, 0)
}

func (dk *DpsDeathknight) blSpreadDiseases(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.blDiseaseCheck(sim, target, dk.Pestilence, true, 1) {
		casted := dk.Pestilence.Cast(sim, target)
		landed := dk.LastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.sr.recastedFF = !(casted && landed)
		dk.sr.recastedBP = !(casted && landed)
		return -1
	} else {
		dk.blRecastDiseasesSequence(sim)
		return sim.CurrentTime
	}
}

// Save up Runic Power for DRW - Allow casts above 100 RP when DRW is ready or above 85 (for death strike glyph) when not
func (dk *DpsDeathknight) blDeathCoilCheck(sim *core.Simulation) bool {
	canCastDrw := dk.DancingRuneWeapon.IsReady(sim) || dk.DancingRuneWeapon.CD.TimeToReady(sim) < 5*time.Second
	currentRP := dk.CurrentRunicPower()
	return (!canCastDrw && currentRP >= 65) || (canCastDrw && dk.CurrentRunicPower() >= 100)
}
