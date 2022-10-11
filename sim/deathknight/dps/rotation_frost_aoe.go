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
			dk.CustomRotation.Cast(sim)
		}
		return -1
	})
}

func (dk *DpsDeathknight) makeCustomRotation() *common.CustomRotation {
	return common.NewCustomRotation(dk.Rotation.FrostCustomRotation, dk.GetCharacter(), map[int32]common.CustomSpell{
		int32(proto.Deathknight_Rotation_CustomIcyTouch): common.CustomSpell{
			RuneSpell: dk.IcyTouch,
			Condition: func(sim *core.Simulation) bool {
				return !dk.FrostFeverDisease[dk.CurrentTarget.Index].IsActive() && dk.IcyTouch.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomPlagueStrike): common.CustomSpell{
			RuneSpell: dk.PlagueStrike,
			Condition: func(sim *core.Simulation) bool {
				return !dk.BloodPlagueDisease[dk.CurrentTarget.Index].IsActive() && dk.PlagueStrike.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomPestilence): common.CustomSpell{
			RuneSpell: dk.Pestilence,
			Condition: func(sim *core.Simulation) bool {
				if !dk.Pestilence.CanCast(sim) {
					return false
				}

				if dk.FrostFeverDisease[dk.CurrentTarget.Index].ExpiresAt()-sim.CurrentTime <= 2*time.Second && dk.BloodPlagueDisease[dk.CurrentTarget.Index].ExpiresAt()-sim.CurrentTime <= 2*time.Second {
					return true
				}

				needToSpread := false
				numHits := dk.Env.GetNumTargets()
				for i := int32(0); i < numHits; i++ {
					target := &dk.Env.GetTarget(i).Unit
					diseases := dk.FrostFeverDisease[target.Index].IsActive() && dk.BloodPlagueDisease[target.Index].IsActive()

					if !diseases {
						needToSpread = true
					}
				}

				return needToSpread
			},
		},
		int32(proto.Deathknight_Rotation_CustomObliterate): common.CustomSpell{
			RuneSpell: dk.Obliterate,
			Condition: func(sim *core.Simulation) bool {
				return dk.Obliterate.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomHowlingBlast): common.CustomSpell{
			RuneSpell: dk.HowlingBlast,
			Condition: func(sim *core.Simulation) bool {
				return dk.HowlingBlast.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomHowlingBlastRime): common.CustomSpell{
			RuneSpell: dk.HowlingBlast,
			Condition: func(sim *core.Simulation) bool {
				return dk.HowlingBlast.CanCast(sim) && dk.Rime()
			},
		},
		int32(proto.Deathknight_Rotation_CustomBloodBoil): common.CustomSpell{
			RuneSpell: dk.BloodBoil,
			Condition: func(sim *core.Simulation) bool {
				return dk.BloodBoil.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomBloodStrike): common.CustomSpell{
			RuneSpell: dk.BloodStrike,
			Condition: func(sim *core.Simulation) bool {
				return dk.BloodStrike.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomDeathAndDecay): common.CustomSpell{
			RuneSpell: dk.DeathAndDecay,
			Condition: func(sim *core.Simulation) bool {
				return dk.DeathAndDecay.CanCast(sim)
			},
		},
		int32(proto.Deathknight_Rotation_CustomHornOfWinter): common.CustomSpell{
			RuneSpell: dk.HornOfWinter,
			Condition: func(sim *core.Simulation) bool {
				return dk.HornOfWinter.CanCast(sim)
			},
		},
	})
}
