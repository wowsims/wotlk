package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func (dk *DpsDeathknight) setupCustomRotations() {
	dk.RotationSequence.NewAction(func(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
		if dk.CustomRotation != nil {
			if !dk.CustomRotation.Cast(sim) {
				return -1
			}
		} else {
			dk.LastCast = nil
		}

		if dk.LastCast == dk.EmpowerRuneWeapon.Spell || dk.LastCast == dk.BloodTap.Spell || dk.LastCast == dk.UnbreakableArmor.Spell {
			return sim.CurrentTime
		}

		return -1
	})
}

func (dk *DpsDeathknight) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(dk.Rotation.FrostCustomRotation, dk.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.Deathknight_Rotation_CustomIcyTouch): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.IcyTouch.CurCast.Cost
				return dk.IcyTouch.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return !dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() && dk.IcyTouch.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomPlagueStrike): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.PlagueStrike.CurCast.Cost
				return dk.PlagueStrike.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return !dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() && dk.PlagueStrike.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomPestilence): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.Pestilence.CurCast.Cost
				return dk.Pestilence.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				if !dk.Pestilence.CanCast(sim) {
					return false
				}

				ff := dk.FrostFeverDisease[dk.CurrentTarget.Index].ExpiresAt() - sim.CurrentTime
				bp := dk.BloodPlagueDisease[dk.CurrentTarget.Index].ExpiresAt() - sim.CurrentTime
				ffHalfDuration := time.Duration(0.5 * float64(dk.FrostFeverDisease[dk.CurrentTarget.Index].Duration))
				bpHalfDuration := time.Duration(0.5 * float64(dk.BloodPlagueDisease[dk.CurrentTarget.Index].Duration))
				if ff <= 2*time.Second && bp <= 2*time.Second && sim.GetRemainingDuration() >= ffHalfDuration && sim.GetRemainingDuration() >= bpHalfDuration {
					return true
				}

				numHits := dk.Env.GetNumTargets()
				numDiseased := numHits
				for i := int32(0); i < numHits; i++ {
					target := &dk.Env.GetTarget(i).Unit
					diseases := dk.FrostFeverDisease[target.Index].IsActive() && dk.BloodPlagueDisease[target.Index].IsActive()

					if !diseases {
						numDiseased--
					}
				}

				return float64(numDiseased)/float64(numHits) <= 0.5
			},
		},
		int32(proto.Deathknight_Rotation_CustomObliterate): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.Obliterate.CurCast.Cost
				if dk.Deathchill != nil && dk.Deathchill.IsReady(sim) {
					dk.Deathchill.Cast(sim, target)
				}
				return dk.Obliterate.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.Obliterate.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomHowlingBlast): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.HowlingBlast.CurCast.Cost
				return dk.HowlingBlast.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.HowlingBlast.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomHowlingBlastRime): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.HowlingBlast.CurCast.Cost
				return dk.HowlingBlast.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.HowlingBlast.CanCast(sim) && dk.Rime()
			},
		},
		int32(proto.Deathknight_Rotation_CustomBloodBoil): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.BloodBoil.CurCast.Cost
				return dk.BloodBoil.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.BloodBoil.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomBloodStrike): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.BloodStrike.CurCast.Cost
				return dk.BloodStrike.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.BloodStrike.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomDeathAndDecay): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.DeathAndDecay.CurCast.Cost
				return dk.DeathAndDecay.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.DeathAndDecay.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomHornOfWinter): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.HornOfWinter.CurCast.Cost
				return dk.HornOfWinter.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.HornOfWinter.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomUnbreakableArmor): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.UnbreakableArmor.CurCast.Cost
				return dk.UnbreakableArmor.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.UnbreakableArmor.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomBloodTap): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.BloodTap.CurCast.Cost
				return dk.BloodTap.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.BloodTap.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomEmpoweredRuneWeapon): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.EmpowerRuneWeapon.CurCast.Cost
				return dk.EmpowerRuneWeapon.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.EmpowerRuneWeapon.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomFrostStrike): {
			Action: func(sim *core.Simulation, target *core.Unit) (bool, float64) {
				cost := dk.FrostStrike.CurCast.Cost
				return dk.FrostStrike.Cast(sim, target), cost
			},
			Condition: func(sim *core.Simulation) bool {
				return dk.FrostStrike.CanCast(sim)
			},
		},
	})
}
