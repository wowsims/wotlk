package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

type ProcTracker struct {
	id          int32
	aura        *core.Aura
	didActivate bool
	isActive    bool
	expiresAt   time.Duration
}

func (warlock *Warlock) addProc(id int32, label string, isActive bool) bool {
	if !warlock.HasAura(label) {
		return false
	}
	warlock.procTrackers = append(warlock.procTrackers, &ProcTracker{
		id:          id,
		didActivate: false,
		isActive:    isActive,
		expiresAt:   -1,
		aura:        warlock.GetAura(label),
	})
	return true
}

func (warlock *Warlock) resetProcTrackers() {
	for _, procTracker := range warlock.procTrackers {
		procTracker.didActivate = false
		procTracker.expiresAt = -1
	}
}

func (warlock *Warlock) initProcTrackers() {
	warlock.procTrackers = make([]*ProcTracker, 0)

	warlock.addProc(40211, "Potion of Speed", true)
	warlock.addProc(47197, "Eradication", true)
	warlock.addProc(10060, "Power Infusion", true)
	warlock.addProc(54999, "Hyperspeed Acceleration", true)
	warlock.addProc(26297, "Berserking (Troll)", true)
	warlock.addProc(33697, "Blood Fury", true)

	warlock.addProc(40255, "Dying Curse Proc", false)
	warlock.addProc(55379, "Illustration of the Dragon Soul Proc", false)
	warlock.addProc(45518, "Flare of the Heavens Proc", false)

	warlock.addProc(45466, "Scale of Fates Proc", false)
	warlock.addProc(39229, "Embrace of the Spider Proc", false)
}

func (warlock *Warlock) setupDSCooldowns() {
	warlock.majorCds = make([]*core.MajorCooldown, 0)

	// hyperspeed accelerators
	warlock.DSCooldownSync(core.ActionID{SpellID: 54758}, false)

	// berserking (troll)
	warlock.DSCooldownSync(core.ActionID{SpellID: 26297}, false)

	// blood fury (orc)
	warlock.DSCooldownSync(core.ActionID{SpellID: 33697}, false)

	// potion of speed
	warlock.DSCooldownSync(core.ActionID{ItemID: 40211}, true)

	// Power Infusion
	warlock.DSCooldownSync(core.ActionID{SpellID: 10060}, false)

	// Eradication
	warlock.DSCooldownSync(core.ActionID{SpellID: 47197}, false)

	// active sp trinkets
	warlock.DSCooldownSync(core.ActionID{ItemID: 40255}, false)
	warlock.DSCooldownSync(core.ActionID{ItemID: 40432}, false)
	warlock.DSCooldownSync(core.ActionID{ItemID: 45518}, false)

	// active haste trinkets
	warlock.DSCooldownSync(core.ActionID{ItemID: 39229}, false)
	warlock.DSCooldownSync(core.ActionID{ItemID: 36972}, false)
}

func (warlock *Warlock) DSCooldownSync(actionID core.ActionID, isPotion bool) {
	if majorCd := warlock.Character.GetMajorCooldown(actionID); majorCd != nil {
		warlock.majorCds = append(warlock.majorCds, majorCd)
	}
}

func logMessage(sim *core.Simulation, message string) {
	if sim.Log != nil {
		sim.Log(message)
	}
}

func (warlock *Warlock) DSProcCheck(sim *core.Simulation, castTime time.Duration) bool {
	for _, procTracker := range warlock.procTrackers {
		if !procTracker.didActivate && procTracker.aura.IsActive() {
			procTracker.didActivate = true
			procTracker.expiresAt = procTracker.aura.ExpiresAt()
		}

		// A proc is about to drop
		if procTracker.didActivate && procTracker.expiresAt <= sim.CurrentTime+castTime {
			logMessage(sim, "Proc dropping "+procTracker.aura.Label)
			return false
		}
	}

	for _, procTracker := range warlock.procTrackers {
		if !procTracker.didActivate && !procTracker.isActive {
			logMessage(sim, "Waiting on procs..")
			return true
		}
	}

	return false
}
