package protection

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (prot *ProtectionPaladin) OnGCDReady(sim *core.Simulation) {

	var success bool

	if !success {
		waitTime := time.Second * 5
		prot.Metrics.MarkOOM(&prot.Unit, waitTime)
		prot.WaitUntil(sim, sim.CurrentTime+waitTime)
	}
}

func (prot *ProtectionPaladin) nextCDAt(sim *core.Simulation) time.Duration {
	nextCDAt := core.MinDuration(prot.HolyShield.ReadyAt(), prot.JudgementOfWisdom.ReadyAt())
	nextCDAt = core.MinDuration(nextCDAt, prot.Consecration.ReadyAt())
	return nextCDAt
}

func (prot *ProtectionPaladin) shouldExorcism(sim *core.Simulation) bool {
	return prot.Rotation.UseExorcism &&
		prot.CanExorcism(prot.CurrentTarget) &&
		prot.Exorcism.IsReady(sim) &&
		prot.CurrentMana() > prot.MaxMana()*0.4
}
