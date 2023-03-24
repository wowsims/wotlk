package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/deathknight"
)

type FrostRotation struct {
	oblitCount int32

	// CDS
	hyperSpeedMCD           *core.MajorCooldown
	stoneformMCD            *core.MajorCooldown
	bloodfuryMCD            *core.MajorCooldown
	berserkingMCD           *core.MajorCooldown
	potionOfSpeedMCD        *core.MajorCooldown
	indestructiblePotionMCD *core.MajorCooldown
	potionUsed              bool

	onUseTrinkets []*core.MajorCooldown

	oblitRPRegen float64
}

func (fr *FrostRotation) Initialize(dk *DpsDeathknight) {
	fr.oblitRPRegen = core.TernaryFloat64(dk.HasSetBonus(deathknight.ItemSetScourgeborneBattlegear, 4), 25.0, 20.0)
	fr.onUseTrinkets = make([]*core.MajorCooldown, 0)
}

func (fr *FrostRotation) Reset(sim *core.Simulation) {
	fr.oblitCount = 0

	fr.hyperSpeedMCD = nil
	fr.stoneformMCD = nil
	fr.bloodfuryMCD = nil
	fr.berserkingMCD = nil
	fr.potionOfSpeedMCD = nil
	fr.indestructiblePotionMCD = nil
	fr.onUseTrinkets = nil
	fr.potionUsed = false
}

func (dk *DpsDeathknight) addOnUseTrinketCooldown(actionID core.ActionID) {
	if majorCd := dk.Character.GetMajorCooldown(actionID); majorCd != nil {
		majorCd.Disable()
		dk.fr.onUseTrinkets = append(dk.fr.onUseTrinkets, majorCd)
	}
}

func (dk *DpsDeathknight) getMajorCooldown(actionID core.ActionID) *core.MajorCooldown {
	if majorCd := dk.Character.GetMajorCooldown(actionID); majorCd != nil {
		majorCd.Disable()
		return majorCd
	}
	return nil
}

func (dk *DpsDeathknight) setupUnbreakableArmorCooldowns() {
	fr := &dk.fr

	// hyperspeed accelerators
	fr.hyperSpeedMCD = dk.getMajorCooldown(core.ActionID{SpellID: 54758})

	// stoneform (dwarf)
	fr.stoneformMCD = dk.getMajorCooldown(core.ActionID{SpellID: 20594})

	// bloodfury (orc)
	fr.bloodfuryMCD = dk.getMajorCooldown(core.ActionID{SpellID: 33697})

	// berserking (troll)
	fr.berserkingMCD = dk.getMajorCooldown(core.ActionID{SpellID: 26297})

	// potion of speed
	fr.potionOfSpeedMCD = dk.getMajorCooldown(core.ActionID{ItemID: 40211})

	// indestructible potion
	fr.indestructiblePotionMCD = dk.getMajorCooldown(core.ActionID{ItemID: 40093})

	// On use trinkets
	dk.addOnUseTrinketCooldown(core.ActionID{ItemID: 40531}) // Mark of nogganon
	dk.addOnUseTrinketCooldown(core.ActionID{ItemID: 37166}) // Sphere of Red Dragon's Blood
	dk.addOnUseTrinketCooldown(core.ActionID{ItemID: 37723}) // Incisor Fragment
	dk.addOnUseTrinketCooldown(core.ActionID{ItemID: 39257}) // Loatheb's Shadow
	dk.addOnUseTrinketCooldown(core.ActionID{ItemID: 44014}) // Fezzik's Pocketwatch
}

func (dk *DpsDeathknight) castMajorCooldown(mcd *core.MajorCooldown, sim *core.Simulation, target *core.Unit) {
	if mcd != nil {
		if mcd.Spell.IsReady(sim) && (dk.GCD.IsReady(sim) || mcd.Spell.DefaultCast.GCD == 0) {
			mcd.Spell.Cast(sim, target)
		}
	}
}

func (dk *DpsDeathknight) castMajorCooldownConditional(mcd *core.MajorCooldown, conditionalMCDs []*core.MajorCooldown, sim *core.Simulation, target *core.Unit) {
	if mcd != nil && !dk.fr.potionUsed {
		if mcd.Spell.IsReady(sim) && dk.GCD.IsReady(sim) {
			allReady := true
			for i := range conditionalMCDs {
				if conditionalMCDs[i] != nil && allReady {
					spell := conditionalMCDs[i].Spell
					allReady = allReady && spell.CD.IsReady(sim)
					//// TODO: Find a way to get the racial aura durations
					if !allReady {
						if dk.Env.Encounter.Duration < time.Duration(float64(spell.CD.ReadyAt())*1.25) {
							allReady = true
						}
					}
				}
			}

			if allReady {
				mcd.Spell.Cast(sim, target)

				dk.fr.potionUsed = true
			}
		}
	}
}

func (dk *DpsDeathknight) castAllMajorCooldowns(sim *core.Simulation) {
	fr := &dk.fr
	target := dk.CurrentTarget

	racialConditionals := []*core.MajorCooldown{fr.bloodfuryMCD, fr.berserkingMCD, fr.stoneformMCD}

	dk.castMajorCooldownConditional(fr.potionOfSpeedMCD, racialConditionals, sim, target)
	dk.castMajorCooldownConditional(fr.indestructiblePotionMCD, racialConditionals, sim, target)
	dk.castMajorCooldown(fr.hyperSpeedMCD, sim, target)
	dk.castMajorCooldown(fr.stoneformMCD, sim, target)
	dk.castMajorCooldown(fr.bloodfuryMCD, sim, target)
	dk.castMajorCooldown(fr.berserkingMCD, sim, target)

	for _, trinket := range fr.onUseTrinkets {
		dk.castMajorCooldown(trinket, sim, target)
	}
}

func (dk *DpsDeathknight) RotationActionCallback_UA_Frost(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.UnbreakableArmor != nil {
		casted := dk.UnbreakableArmor.Cast(sim, target)

		if casted {
			dk.castAllMajorCooldowns(sim)
			s.ConditionalAdvance(casted)
			return sim.CurrentTime
		} else {
			s.ConditionalAdvance(casted)
			return -1
		}
	} else {
		casted := dk.BloodStrike.Cast(sim, target)
		if casted {
			dk.castAllMajorCooldowns(sim)
			s.ConditionalAdvance(casted)
			return sim.CurrentTime
		} else {
			s.ConditionalAdvance(casted)
			return -1
		}
	}
}

func (dk *DpsDeathknight) RotationActionCallback_Frost_FS_HB(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	if dk.RimeAura.IsActive() && dk.Talents.HowlingBlast {
		dk.HowlingBlast.Cast(sim, target)
	} else if dk.Talents.FrostStrike {
		dk.FrostStrike.Cast(sim, target)
	}

	s.Advance()
	return -1
}

func (dk *DpsDeathknight) RotationActionCallback_Frost_Pesti_ERW(sim *core.Simulation, target *core.Unit, s *deathknight.Sequence) time.Duration {
	// Casts Pesti then ERW in the same GCD (rather than chaining them sequentially which will cause ERW to be delayed by pesti GCD)
	// This is a DPS increase since it allows rune grace to start as soon as possible
	casted := dk.Pestilence.Cast(sim, target)
	advance := casted && dk.LastOutcome.Matches(core.OutcomeLanded)

	if advance {
		casted = dk.EmpowerRuneWeapon.Cast(sim, target)
		advance = casted && advance
		s.ConditionalAdvance(advance)
	}
	return -1
}
