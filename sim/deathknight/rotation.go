package deathknight

import (
	//"math"
	//"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	//"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) OnAutoAttack(sim *core.Simulation, spell *core.Spell) {
	if deathKnight.GCD.IsReady(sim) {
		deathKnight.tryUseGCD(sim)
	}
}

func (deathKnight *DeathKnight) OnGCDReady(sim *core.Simulation) {
	deathKnight.tryUseGCD(sim)
}

const (
	DKRotation_Wait uint8 = iota
	DKRotation_IT
	DKRotation_PS
	DKRotation_Obli
	DKRotation_BS
	DKRotation_BT
	DKRotation_UA
	DKRotation_Pesti
	DKRotation_FS
)

func (deathKnight *DeathKnight) tryUseGCD(sim *core.Simulation) {
	//var spell *core.Spell
	var target = deathKnight.CurrentTarget

	if deathKnight.GCD.IsReady(sim) {
		if deathKnight.CanIcyTouch(sim) {
			deathKnight.IcyTouch.Cast(sim, target)
		} else if deathKnight.CanPlagueStrike(sim) {
			deathKnight.PlagueStrike.Cast(sim, target)
		} else {
			nextCD := deathKnight.IcyTouch.ReadyAt()

			if nextCD > sim.CurrentTime {
				deathKnight.WaitUntil(sim, nextCD)
			}
		}
		/*
			if !deathKnight.FrostFeverDisease.IsActive() {
				if deathKnight.CanIcyTouch(sim) {
					deathKnight.IcyTouch.Cast(sim, target)
				}
			} else if !deathKnight.BloodPlagueDisease.IsActive() {
				if deathKnight.CanPlagueStrike(sim) {
					deathKnight.PlagueStrike.Cast(sim, target)
				}
			} else if deathKnight.FrostFeverDisease.IsActive() &&
				deathKnight.BloodPlagueDisease.IsActive() &&
				deathKnight.FrostFeverDisease.ExpiresAt() > deathKnight.FrostRuneReadyAt(sim) &&
				deathKnight.BloodPlagueDisease.ExpiresAt() > deathKnight.UnholyRuneReadyAt(sim) &&
				deathKnight.CanObliterate(sim) {
				deathKnight.Obliterate.Cast(sim, target)
			} else {
				frostFeverExpireTime := deathKnight.FrostFeverDisease.ExpiresAt()
				bloodPlagueExpireTime := deathKnight.BloodPlagueDisease.ExpiresAt()

				runeWaitTime := 0 * time.Second
				if deathKnight.CurrentFrostRunes(sim) == 0 {
					if deathKnight.CurrentUnholyRunes(sim) == 0 {
						runeWaitTime = core.MaxDuration(deathKnight.FrostRuneReadyAt(sim), deathKnight.UnholyRuneReadyAt(sim))
					} else {
						runeWaitTime = deathKnight.FrostRuneReadyAt(sim)
					}
				} else {
					if deathKnight.CurrentUnholyRunes(sim) == 0 {
						runeWaitTime = deathKnight.UnholyRuneReadyAt(sim)
					}
				}

				waitTime := core.MinDuration(frostFeverExpireTime, bloodPlagueExpireTime)
				if runeWaitTime != 0 {
					waitTime = core.MinDuration(waitTime, runeWaitTime)
				}

				if waitTime > sim.CurrentTime {
					deathKnight.WaitUntil(sim, waitTime)
				}
			}
		*/
	}
}
