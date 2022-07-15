package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
)

func OutcomeEitherWeaponHitOrCrit(mhOutcome core.HitOutcome, ohOutcome core.HitOutcome) bool {
	return mhOutcome == core.OutcomeHit || mhOutcome == core.OutcomeCrit || ohOutcome == core.OutcomeHit || ohOutcome == core.OutcomeCrit
}
