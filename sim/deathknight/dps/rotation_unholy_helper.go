package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type WeaponSwapType int32

const (
	WeaponSwap_None WeaponSwapType = iota
	WeaponSwap_BlackMagic
	WeaponSwap_Berserking
	WeaponSwap_FallenCrusader
)

type SigilType int32

const (
	Sigil_Other SigilType = iota
	Sigil_Virulence
	Sigil_HangedMan
)

type UnholyRotation struct {
	dk *DpsDeathknight

	syncTimeFF time.Duration

	gargoyleSnapshot   *core.SnapshotManager
	activatingGargoyle bool
	gargoyleMaxDelay   time.Duration
	gargoyleMinTime    time.Duration

	mhSwap    WeaponSwapType
	mhSwapped bool

	ohSwap    WeaponSwapType
	ohSwapped bool

	bmIcd time.Duration

	sigil       SigilType
	unholyMight bool
}

func (ur *UnholyRotation) Reset(sim *core.Simulation) {
	ur.syncTimeFF = 0
	ur.activatingGargoyle = false
	ur.gargoyleMaxDelay = -1

	ur.mhSwapped = false
	ur.ohSwapped = false
	ur.bmIcd = -1

	if ur.dk.Talents.SummonGargoyle {
		gargMcd := ur.dk.getMajorCooldown(ur.dk.SummonGargoyle.ActionID)
		if gargMcd != nil {
			timings := gargMcd.GetTimings()
			if len(timings) > 0 {
				ur.gargoyleMinTime = timings[0]
			}
		}

		ur.gargoyleSnapshot.ResetProcTrackers()
	}
}

func (ur *UnholyRotation) Initialize(dk *DpsDeathknight) {
	if dk.Talents.SummonGargoyle && dk.Rotation.UseGargoyle {
		dk.setupWeaponSwap()
	}

	// Init Sigil of Virulence Rotation
	if dk.Talents.ScourgeStrike && dk.Equip[core.ItemSlotRanged].ID == 47673 {
		ur.sigil = Sigil_Virulence
	}

	// Init T9 2P Proc
	if dk.HasSetBonus(deathknight.ItemSetThassariansBattlegear, 2) {
		ur.unholyMight = true
	}
}

func (dk *DpsDeathknight) getFirstDiseaseAction() deathknight.RotationAction {
	if dk.sr.ffFirst {
		return dk.RotationActionCallback_IT
	}
	return dk.RotationActionCallback_PS
}

func (dk *DpsDeathknight) getSecondDiseaseAction() deathknight.RotationAction {
	if dk.sr.ffFirst {
		return dk.RotationActionCallback_PS
	}
	return dk.RotationActionCallback_IT
}

func (dk *DpsDeathknight) uhBloodRuneAction(isFirst bool) deathknight.RotationAction {
	if isFirst {
		if dk.Env.GetNumTargets() > 1 {
			return dk.RotationActionCallback_Pesti
		} else {
			return dk.RotationActionCallback_BS
		}
	} else {
		return dk.RotationActionCallback_BS
	}
}

func (dk *DpsDeathknight) bonusProcRotationChecks(sim *core.Simulation) (bool, bool) {
	// If we have sigil of virulence
	// Higher prio SS then Dnd when gargoyle is ready
	prioSs := false
	if dk.Talents.ScourgeStrike && dk.ur.sigil == Sigil_Virulence && (dk.SummonGargoyle.IsReady(sim) || dk.SummonGargoyle.CD.TimeToReady(sim) < 10*time.Second) {
		virulenceAura := dk.GetAura("Sigil of Virulence Proc")
		prioSs = !virulenceAura.IsActive() || virulenceAura.RemainingDuration(sim) < 10*time.Second
	}

	// If we have T9 2P we prio BS over BB for refreshing the buff
	prioBs := false
	if dk.ur.unholyMight {
		unholyMightAura := dk.GetAura("Unholy Might Proc")
		prioBs = unholyMightAura.StartedAt() == 0 || unholyMightAura.StartedAt() < sim.CurrentTime-45*time.Second
	}
	return prioSs, prioBs
}

func (dk *DpsDeathknight) weaponSwapCheck(sim *core.Simulation) bool {
	if !dk.ItemSwap.IsEnabled() {
		return false
	}

	// Swap if gargoyle will still be on CD for full ICD or if gargoyle is already active
	shouldSwapBm := dk.ur.bmIcd < sim.CurrentTime && (dk.SummonGargoyle.CD.TimeToReady(sim) > 45*time.Second || dk.HasActiveAura("Summon Gargoyle"))
	shouldSwapBackFromBm := dk.HasActiveAura("Black Magic Proc") // || dk.GetAura("Rune Of The Fallen Crusader Proc").RemainingDuration(sim) < 5*time.Second

	if dk.ur.mhSwap == WeaponSwap_BlackMagic {
		if !dk.ur.mhSwapped && shouldSwapBm {
			// Swap to BM
			if sim.Log != nil {
				sim.Log("Swapping MH to BM")
			}
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, true)
			dk.ur.mhSwapped = true
		} else if dk.ur.mhSwapped && shouldSwapBackFromBm {
			// Swap to Normal set and set BM Icd tracker
			if sim.Log != nil {
				sim.Log("Swapping MH to Normal")
			}
			dk.ur.bmIcd = dk.GetAura("Black Magic Proc").ExpiresAt() + 35*time.Second
			dk.ur.mhSwapped = false
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, true)
		}
	}

	if dk.ur.ohSwap == WeaponSwap_BlackMagic {
		if !dk.ur.ohSwapped && shouldSwapBm {
			// Swap to BM
			if sim.Log != nil {
				sim.Log("Swapping OH to BM")
			}
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, true)
			dk.ur.ohSwapped = true
		} else if dk.ur.ohSwapped && shouldSwapBackFromBm {
			// Swap to Normal set and set BM Icd tracker
			if sim.Log != nil {
				sim.Log("Swapping OH to Normal")
			}
			dk.ur.bmIcd = dk.GetAura("Black Magic Proc").ExpiresAt() + 35*time.Second
			dk.ur.ohSwapped = false
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, true)
		}
	}

	shouldSwapBerserking := dk.HasActiveAura("Rune Of The Fallen Crusader Proc") &&
		dk.GetAura("Rune Of The Fallen Crusader Proc").RemainingDuration(sim) > time.Second*10

	shouldSwapBackfromBerserking := false //dk.GetAura("Rune Of The Fallen Crusader Proc").RemainingDuration(sim) < 5*time.Second

	if dk.ur.mhSwap == WeaponSwap_Berserking {
		if !dk.ur.mhSwapped && !dk.HasActiveAura("Berserking MH Proc") && shouldSwapBerserking {
			// Swap to Berserking
			if sim.Log != nil {
				sim.Log("Swapping MH to Berserking")
			}
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, true)
			dk.ur.mhSwapped = true
		} else if dk.ur.mhSwapped && (dk.HasActiveAura("Berserking MH Proc") || shouldSwapBackfromBerserking) {
			// Swap to Normal set
			if sim.Log != nil {
				sim.Log("Swapping MH to Normal")
			}
			dk.ur.mhSwapped = false
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand}, true)
		}
	}

	if dk.ur.ohSwap == WeaponSwap_Berserking {
		if !dk.ur.ohSwapped && !dk.HasActiveAura("Berserking OH Proc") && shouldSwapBerserking {
			// Swap to Berserking
			if sim.Log != nil {
				sim.Log("Swapping OH to Berserking")
			}
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, true)
			dk.ur.ohSwapped = true
		} else if dk.ur.ohSwapped && (dk.HasActiveAura("Berserking OH Proc") || shouldSwapBackfromBerserking) {
			// Swap to Normal set
			if sim.Log != nil {
				sim.Log("Swapping OH to Normal")
			}
			dk.ur.ohSwapped = false
			dk.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand}, true)
		}
	}

	return false
}

func (dk *DpsDeathknight) desolationAuraCheck(sim *core.Simulation) bool {
	return !dk.DesolationAura.IsActive() || dk.DesolationAura.RemainingDuration(sim) < 10*time.Second ||
		dk.Rotation.BloodRuneFiller == proto.Deathknight_Rotation_BloodStrike
}

func (dk *DpsDeathknight) uhDiseaseCheck(sim *core.Simulation, target *core.Unit, spell *core.Spell, costRunes bool, casts int) bool {
	return dk.shDiseaseCheck(sim, target, spell, costRunes, casts, 0)
}

func (dk *DpsDeathknight) uhSpreadDiseases(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) bool {
	if dk.uhDiseaseCheck(sim, target, dk.Pestilence, true, 1) {
		casted := dk.Pestilence.Cast(sim, target)
		landed := dk.LastOutcome.Matches(core.OutcomeLanded)

		// Reset flags on succesfull cast
		dk.sr.recastedFF = !(casted && landed)
		dk.sr.recastedBP = !(casted && landed)
		return casted
	} else {
		dk.uhRecastDiseasesSequence(sim)
		return true
	}
}

// Simpler but somehow more effective for overall dps dnd check
func (dk *DpsDeathknight) uhShouldWaitForDnD(sim *core.Simulation, blood bool, frost bool, unholy bool) bool {
	if dk.Talents.ImprovedUnholyPresence > 0 {
		return dk.DeathAndDecay.IsReady(sim) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1))
	} else {
		return !(!(dk.DeathAndDecay.CD.IsReady(sim) || dk.DeathAndDecay.CD.TimeToReady(sim) <= 4*time.Second) || ((!blood || dk.CurrentBloodRunes() > 1) && (!frost || dk.CurrentFrostRunes() > 1) && (!unholy || dk.CurrentUnholyRunes() > 1)))
	}
}

func (dk *DpsDeathknight) uhGhoulFrenzyCheck(sim *core.Simulation, target *core.Unit) bool {
	if !dk.Ghoul.Pet.IsEnabled() {
		return false
	}

	// If no Ghoul Frenzy Aura or duration less then 10 seconds we try recasting
	if !dk.GhoulFrenzyAura.IsActive() || dk.GhoulFrenzyAura.RemainingDuration(sim) < 10*time.Second {
		// Use Ghoul Frenzy with a Blood Tap and Blood rune if all blood runes are on CD and Garg wont come off cd in less then a minute.
		if (dk.Rotation.BloodTap == proto.Deathknight_Rotation_GhoulFrenzy || dk.Rotation.BtGhoulFrenzy) && dk.BloodTap.CanCast(sim, nil) && dk.GhoulFrenzy.IsReady(sim) && dk.CurrentBloodRunes() == 0 && dk.CurrentUnholyRunes() == 0 {
			if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 1) {
				dk.uhGhoulFrenzySequence(sim, true)
				return true
			} else {
				dk.uhRecastDiseasesSequence(sim)
				return true
			}
		} else if !dk.Rotation.BtGhoulFrenzy && dk.GhoulFrenzy.CanCast(sim, nil) && dk.IcyTouch.CanCast(sim, nil) {
			if dk.uhGargoyleCheck(sim, target, dk.SpellGCD()*2+50*time.Millisecond) {
				dk.uhAfterGargoyleSequence(sim)
				return true
			}
			// Use Ghoul Frenzy with an Unholy Rune and sync the frost rune with Icy Touch
			if dk.uhDiseaseCheck(sim, target, dk.GhoulFrenzy, true, 5) && dk.uhDiseaseCheck(sim, target, dk.IcyTouch, true, 5) {
				// TODO: This can spend runes that should be spent on DnD fix it!
				dk.uhGhoulFrenzySequence(sim, false)
				return true
			} else {
				dk.uhRecastDiseasesSequence(sim)
				return true
			}
		}
	}
	return false
}

func (dk *DpsDeathknight) uhBloodTap(sim *core.Simulation, target *core.Unit) bool {
	if !dk.GCD.IsReady(sim) {
		return false
	}

	if dk.Rotation.BloodTap != proto.Deathknight_Rotation_GhoulFrenzy && dk.BloodTap.IsReady(sim) && dk.CurrentBloodRunes() == 0 {
		switch dk.Rotation.BloodTap {
		case proto.Deathknight_Rotation_IcyTouch:
			if dk.CurrentFrostRunes() == 0 {
				dk.BloodTap.Cast(sim, dk.CurrentTarget)
				dk.IcyTouch.Cast(sim, target)
				return true
			}
		case proto.Deathknight_Rotation_BloodStrikeBT:
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodStrike.Cast(sim, target)
			return true
		case proto.Deathknight_Rotation_BloodBoilBT:
			dk.BloodTap.Cast(sim, dk.CurrentTarget)
			dk.BloodBoil.Cast(sim, target)
			return true
		}
	}

	return false
}

func (dk *DpsDeathknight) uhEmpoweredRuneWeapon(sim *core.Simulation, target *core.Unit) bool {
	if !dk.Rotation.UseEmpowerRuneWeapon {
		return false
	}

	if !dk.EmpowerRuneWeapon.IsReady(sim) {
		return false
	}

	// Save ERW for best Army Snapshot after Garg
	if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_AsMajorCd && dk.Rotation.HoldErwArmy && dk.ArmyOfTheDead.IsReady(sim) {
		return false
	}

	if dk.CurrentBloodRunes() > 0 || dk.CurrentFrostRunes() > 0 || dk.CurrentUnholyRunes() > 0 || dk.CurrentDeathRunes() > 0 {
		return false
	}

	timeToNextRune := dk.AnyRuneReadyAt(sim) - sim.CurrentTime
	if timeToNextRune < 2*time.Second {
		return false
	}

	dk.EmpowerRuneWeapon.Cast(sim, target)
	return true
}

func (dk *DpsDeathknight) uhMindFreeze(sim *core.Simulation, target *core.Unit) bool {
	if dk.Talents.EndlessWinter == 2 && dk.SummonGargoyle.IsReady(sim) {
		if dk.MindFreezeSpell.IsReady(sim) {
			dk.MindFreezeSpell.Cast(sim, target)
			return true
		}
	}
	return false
}

// Save up Runic Power for Summon Gargoyle - Allow casts above 100 rp or garg CD > 5 sec
func (dk *DpsDeathknight) uhDeathCoilCheck(sim *core.Simulation) bool {
	return !dk.Talents.SummonGargoyle || !(dk.SummonGargoyle.IsReady(sim) || dk.SummonGargoyle.CD.TimeToReady(sim) < 5*time.Second) || sim.CurrentTime < dk.ur.gargoyleMinTime-5*time.Second || dk.CurrentRunicPower() >= 100 || !dk.Rotation.UseGargoyle
}

// Combined checks for casting gargoyle sequence & going back to blood presence after
func (dk *DpsDeathknight) uhGargoyleCheck(sim *core.Simulation, target *core.Unit, castTime time.Duration) bool {
	if !dk.Rotation.UseGargoyle {
		return false
	}

	if dk.uhGargoyleCanCast(sim, castTime) {
		if !dk.PresenceMatches(deathknight.UnholyPresence) && (dk.Rotation.PreNerfedGargoyle || dk.Rotation.GargoylePresence == proto.Deathknight_Rotation_Unholy) {
			if dk.CurrentUnholyRunes() == 0 {
				if dk.BloodTap.IsReady(sim) {
					dk.BloodTap.Cast(sim, dk.CurrentTarget)
				} else {
					return false
				}
			}
			dk.UnholyPresence.Cast(sim, dk.CurrentTarget)
		}

		dk.ur.activatingGargoyle = true
		dk.OnGargoyleStartFirstCast = func() {
			dk.ur.gargoyleSnapshot.ActivateMajorCooldowns(sim)
		}
		dk.ur.gargoyleSnapshot.ActivateMajorCooldowns(sim)
		dk.ur.activatingGargoyle = false

		if dk.SummonGargoyle.Cast(sim, target) {
			dk.UpdateMajorCooldowns()
			dk.ur.gargoyleSnapshot.ResetProcTrackers()
			dk.ur.gargoyleMaxDelay = -1
			return true
		}
	}

	// Go back to Blood Presence after Gargoyle
	if !dk.Rotation.PreNerfedGargoyle && !dk.SummonGargoyle.IsReady(sim) && dk.Rotation.Presence == proto.Deathknight_Rotation_Blood && dk.Rotation.GargoylePresence == proto.Deathknight_Rotation_Unholy && dk.PresenceMatches(deathknight.UnholyPresence) && !dk.HasActiveAura("Summon Gargoyle") {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		return dk.BloodPresence.Cast(sim, target)
	}

	// Go back to Unholy Presence after Gargoyle
	if !dk.Rotation.PreNerfedGargoyle && !dk.SummonGargoyle.IsReady(sim) && dk.Rotation.Presence == proto.Deathknight_Rotation_Unholy && dk.Rotation.GargoylePresence == proto.Deathknight_Rotation_Blood && dk.PresenceMatches(deathknight.BloodPresence) && !dk.HasActiveAura("Summon Gargoyle") {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		return dk.UnholyPresence.Cast(sim, target)
	}

	// Do not switch presences if gargoyle is still up if it's nerfed gargoyle
	if !dk.Rotation.PreNerfedGargoyle && dk.HasActiveAura("Summon Gargoyle") {
		return false
	}

	// Go back to Blood Presence after Bloodlust
	if dk.Rotation.Presence == proto.Deathknight_Rotation_Blood && dk.Rotation.BlPresence == proto.Deathknight_Rotation_Unholy && dk.PresenceMatches(deathknight.UnholyPresence) && !dk.HasActiveAuraWithTag("Bloodlust") {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		return dk.BloodPresence.Cast(sim, target)
	}

	// Go back to Unholy Presence after Bloodlust
	if dk.Rotation.Presence == proto.Deathknight_Rotation_Unholy && dk.Rotation.BlPresence == proto.Deathknight_Rotation_Blood && dk.PresenceMatches(deathknight.BloodPresence) && !dk.HasActiveAuraWithTag("Bloodlust") {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		return dk.UnholyPresence.Cast(sim, target)
	}

	// Go back to Blood Presence after gargoyle cast
	if dk.Rotation.BlPresence == proto.Deathknight_Rotation_Blood && dk.PresenceMatches(deathknight.UnholyPresence) && !dk.SummonGargoyle.IsReady(sim) && dk.HasActiveAuraWithTag("Bloodlust") {
		if dk.BloodTapAura.IsActive() {
			dk.BloodTapAura.Deactivate(sim)
		}
		if dk.BloodPresence.Cast(sim, target) {
			return true
		}
	}
	return false
}

func (dk *DpsDeathknight) uhGargoyleCanCast(sim *core.Simulation, castTime time.Duration) bool {
	if !dk.SummonGargoyle.IsReady(sim) {
		return false
	}
	if sim.CurrentTime < dk.ur.gargoyleMinTime {
		return false
	}
	if !dk.CastCostPossible(sim, 60.0, 0, 0, 0) {
		return false
	}
	// Setup max delay possible
	if dk.ur.gargoyleMaxDelay == -1 {
		gargCd := dk.SummonGargoyle.CD.Duration
		timeLeft := sim.GetRemainingDuration()
		for timeLeft > gargCd {
			timeLeft = timeLeft - (gargCd + 2*time.Second)
		}
		dk.ur.gargoyleMaxDelay = sim.CurrentTime + timeLeft - 2*time.Second
	}
	// Cast it if holding will result in less total Gargoyles for the encounter
	if sim.CurrentTime > dk.ur.gargoyleMaxDelay {
		return true
	}
	// Cast it if holding will take from its duration
	if sim.GetRemainingDuration() < 32*time.Second {
		return true
	}
	if !dk.PresenceMatches(deathknight.UnholyPresence) && (!dk.BloodTap.CanCast(sim, nil) && dk.CurrentUnholyRunes() == 0) {
		return false
	}
	if !dk.ur.gargoyleSnapshot.CanSnapShot(sim, castTime) {
		return false
	}

	return true
}
