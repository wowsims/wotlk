package dps

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterDpsDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_Deathknight{},
		proto.Spec_SpecDeathknight,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewDpsDeathknight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Deathknight)
			if !ok {
				panic("Invalid spec value for Deathknight!")
			}
			player.Spec = playerSpec
		},
	)
}

type DpsDeathknight struct {
	*deathknight.Deathknight

	sr SharedRotation
	br BloodRotation
	fr FrostRotation
	ur UnholyRotation

	CustomRotation *common.CustomRotation

	Rotation *proto.Deathknight_Rotation
}

func NewDpsDeathknight(character core.Character, player *proto.Player) *DpsDeathknight {
	dk := player.GetDeathknight()

	dpsDk := &DpsDeathknight{
		Deathknight: deathknight.NewDeathknight(character, deathknight.DeathknightInputs{
			StartingRunicPower:  dk.Options.StartingRunicPower,
			PrecastGhoulFrenzy:  dk.Options.PrecastGhoulFrenzy,
			PrecastHornOfWinter: dk.Options.PrecastHornOfWinter,
			PetUptime:           dk.Options.PetUptime,
			DrwPestiApply:       dk.Options.DrwPestiApply,
			BloodOpener:         dk.Rotation.BloodOpener,
			IsDps:               true,

			RefreshHornOfWinter: dk.Rotation.RefreshHornOfWinter,
			ArmyOfTheDeadType:   dk.Rotation.ArmyOfTheDead,
			StartingPresence:    dk.Rotation.StartingPresence,
			UseAMS:              dk.Rotation.UseAms,
			AvgAMSSuccessRate:   dk.Rotation.AvgAmsSuccessRate,
			AvgAMSHit:           dk.Rotation.AvgAmsHit,
		}, player.TalentsString, dk.Rotation.PreNerfedGargoyle),
		Rotation: dk.Rotation,
	}

	dpsDk.Inputs.UnholyFrenzyTarget = dk.Options.UnholyFrenzyTarget

	dpsDk.EnableAutoAttacks(dpsDk, core.AutoAttackOptions{
		MainHand:       dpsDk.WeaponFromMainHand(dpsDk.DefaultMeleeCritMultiplier()),
		OffHand:        dpsDk.WeaponFromOffHand(dpsDk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if dpsDk.RuneStrike.CanCast(sim, nil) {
				return dpsDk.RuneStrike
			} else {
				return nil
			}
		},
	})

	if dpsDk.Talents.SummonGargoyle && dpsDk.Rotation.UseGargoyle && dpsDk.Rotation.EnableWeaponSwap {
		dpsDk.EnableItemSwap(dpsDk.Rotation.WeaponSwap, dpsDk.DefaultMeleeCritMultiplier(), dpsDk.DefaultMeleeCritMultiplier(), 0)
	}

	dpsDk.br.dk = dpsDk
	dpsDk.sr.dk = dpsDk
	dpsDk.ur.dk = dpsDk

	return dpsDk
}

func (dk *DpsDeathknight) FrostPointsInBlood() int32 {
	return dk.Talents.Butchery + dk.Talents.Subversion + dk.Talents.BladeBarrier + dk.Talents.DarkConviction
}

func (dk *DpsDeathknight) FrostPointsInUnholy() int32 {
	return dk.Talents.ViciousStrikes + dk.Talents.Virulence + dk.Talents.Epidemic + dk.Talents.RavenousDead + dk.Talents.Necrosis + dk.Talents.BloodCakedBlade
}

func (dk *DpsDeathknight) SetupRotations() {
	bl, fr, uh := deathknight.PointsInTalents(dk.Talents)

	if dk.Rotation.AutoRotation {
		if uh > fr && uh > bl {
			// Unholy
			dk.Rotation.BtGhoulFrenzy = false
			dk.Rotation.UseEmpowerRuneWeapon = true
			dk.Rotation.HoldErwArmy = true
			dk.Rotation.UseGargoyle = true
			dk.Inputs.ArmyOfTheDeadType = proto.Deathknight_Rotation_AsMajorCd
			dk.Rotation.BloodTap = proto.Deathknight_Rotation_GhoulFrenzy
			dk.Rotation.FirstDisease = proto.Deathknight_Rotation_FrostFever
			dk.Rotation.StartingPresence = proto.Deathknight_Rotation_Unholy
			dk.Rotation.BlPresence = proto.Deathknight_Rotation_Blood
			dk.Rotation.Presence = proto.Deathknight_Rotation_Blood
			dk.Rotation.GargoylePresence = proto.Deathknight_Rotation_Unholy

			mh := dk.GetMHWeapon()
			oh := dk.GetOHWeapon()

			if mh != nil && oh != nil {
				// DW
				dk.Rotation.BloodRuneFiller = proto.Deathknight_Rotation_BloodBoil
				dk.Rotation.UseDeathAndDecay = true
			} else {
				// 2h
				if dk.Env.GetNumTargets() > 1 {
					dk.Rotation.BloodRuneFiller = proto.Deathknight_Rotation_BloodBoil
					dk.Rotation.UseDeathAndDecay = true
				} else {
					dk.Rotation.BloodRuneFiller = proto.Deathknight_Rotation_BloodStrike
					dk.Rotation.UseDeathAndDecay = false
				}
			}
			// Always use DnD if you have the glyph.
			if dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathAndDecay) {
				dk.Rotation.UseDeathAndDecay = true
			}
		} else if fr > uh && fr > bl {
			// Frost rotations here.
		} else if bl > fr && bl > uh {
			// Blood rotations here.

			// AotD not good as Major CD in blood due to DRW confclits
			if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_AsMajorCd {
				dk.Inputs.ArmyOfTheDeadType = proto.Deathknight_Rotation_PreCast
				dk.Rotation.HoldErwArmy = false
			}
		} else {
			// some weird spec where two trees are equal...
		}
	}

	dk.RotationSequence.Clear()

	dk.Inputs.FuStrike = deathknight.FuStrike_Obliterate

	dk.CustomRotation = dk.makeCustomRotation()
	if dk.CustomRotation == nil || dk.Rotation.FrostRotationType == proto.Deathknight_Rotation_SingleTarget {
		dk.Rotation.FrostRotationType = proto.Deathknight_Rotation_SingleTarget
		if fr > uh && fr > bl {
			// AotD as major CD doesnt work well with frost
			if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_AsMajorCd {
				dk.Inputs.ArmyOfTheDeadType = proto.Deathknight_Rotation_PreCast
				dk.Rotation.HoldErwArmy = false
			}
			if bl > uh {
				if dk.Rotation.DesyncRotation {
					dk.setupFrostSubBloodDesyncOpener()
				} else if dk.Rotation.UseEmpowerRuneWeapon {
					dk.setupFrostSubBloodERWOpener()
				} else {
					dk.setupFrostSubBloodNoERWOpener()
				}
			} else {
				dk.Rotation.FrostRotationType = proto.Deathknight_Rotation_SingleTarget
				if dk.Rotation.UseEmpowerRuneWeapon {
					dk.setupFrostSubUnholyERWOpener()
				} else {
					// TODO you can't unh sub without ERW in the opener...yet
					dk.Rotation.UseEmpowerRuneWeapon = true
					dk.setupFrostSubUnholyERWOpener()
				}
			}
		} else if uh > fr && uh > bl {
			dk.setupUnholyRotations()
		} else if bl > fr && bl > uh {
			if dk.Inputs.ArmyOfTheDeadType == proto.Deathknight_Rotation_AsMajorCd {
				dk.Inputs.ArmyOfTheDeadType = proto.Deathknight_Rotation_PreCast
				dk.Rotation.HoldErwArmy = false
			}
			dk.setupBloodRotations()
		}
	} else {
		dk.setupCustomRotations()
	}
}

func (dk *DpsDeathknight) GetDeathknight() *deathknight.Deathknight {
	return dk.Deathknight
}

func (dk *DpsDeathknight) Initialize() {
	dk.Deathknight.Initialize()

	if dk.Talents.DancingRuneWeapon {
		dk.br.drwSnapshot = core.NewSnapshotManager(dk.GetCharacter())
		dk.setupDrwProcTrackers()
	}

	if dk.Talents.SummonGargoyle {
		dk.ur.gargoyleSnapshot = core.NewSnapshotManager(dk.GetCharacter())
		dk.setupGargProcTrackers()
	}

	dk.sr.Initialize(dk)
	dk.br.Initialize(dk)
	dk.fr.Initialize(dk)
	dk.ur.Initialize(dk)
}

func (dk *DpsDeathknight) setupGargProcTrackers() {
	snapshotManager := dk.ur.gargoyleSnapshot

	// Don't need to wait for haste snapshots anymore
	if dk.Rotation.PreNerfedGargoyle {
		snapshotManager.AddProc(40211, "Potion of Speed", true)
		snapshotManager.AddProc(54999, "Hyperspeed Acceleration", true)
		snapshotManager.AddProc(26297, "Berserking (Troll)", true)
		snapshotManager.AddProc(33697, "Blood Fury", true)

		snapshotManager.AddProc(55379, "Thundering Skyflare Diamond Proc", false)
		snapshotManager.AddProc(59626, "Black Magic Proc", false)
		snapshotManager.AddProc(53344, "Rune Of The Fallen Crusader Proc", false)

		snapshotManager.AddProc(37390, "Meteorite Whetstone Proc", false)
		snapshotManager.AddProc(39229, "Embrace of the Spider Proc", false)
		snapshotManager.AddProc(44308, "Signet of Edward the Odd Proc", false)
		snapshotManager.AddProc(43573, "Tears of Bitter Anguish Proc", false)
		snapshotManager.AddProc(45609, "Comet's Trail Proc", false)
		snapshotManager.AddProc(45866, "Elemental Focus Stone Proc", false)

		snapshotManager.AddProc(53344, "Rune Of The Fallen Crusader Proc", false)
	} else {
		fcEnchantId := int32(3368)
		// Only worth snapshotting if both are on (might want to re-visit this after P2)
		if mh, oh := dk.Character.GetMHWeapon(), dk.Character.GetOHWeapon(); mh != nil && oh != nil && mh.Enchant.EffectID == fcEnchantId && oh.Enchant.EffectID == fcEnchantId {
			snapshotManager.AddProc(53344, "Rune Of The Fallen Crusader Proc", false)
		}
	}

	snapshotManager.AddProc(42987, "DMC Greatness Strength Proc", false)

	snapshotManager.AddProc(47115, "Deaths Verdict Strength Proc", false)
	snapshotManager.AddProc(47131, "Deaths Verdict H Strength Proc", false)
	snapshotManager.AddProc(47303, "Deaths Choice Strength Proc", false)
	snapshotManager.AddProc(47464, "Deaths Choice H Strength Proc", false)

	snapshotManager.AddProc(71484, "Deathbringer's Will Strength Proc", false)
	snapshotManager.AddProc(71492, "Deathbringer's Will Haste Proc", false)
	snapshotManager.AddProc(71561, "Deathbringer's Will H Strength Proc", false)
	snapshotManager.AddProc(71560, "Deathbringer's Will H Haste Proc", false)

	snapshotManager.AddProc(40684, "Mirror of Truth Proc", false)
	snapshotManager.AddProc(40767, "Sonic Booster Proc", false)
	snapshotManager.AddProc(44914, "Anvil of Titans Proc", false)
	snapshotManager.AddProc(45286, "Pyrite Infuser Proc", false)
	snapshotManager.AddProc(45522, "Blood of the Old God Proc", false)
	snapshotManager.AddProc(47214, "Banner of Victory Proc", false)
	snapshotManager.AddProc(49074, "Coren's Chromium Coaster Proc", false)
	snapshotManager.AddProc(50342, "Whispering Fanged Skull Proc", false)
	snapshotManager.AddProc(50343, "Whispering Fanged Skull H Proc", false)
	snapshotManager.AddProc(50401, "Ashen Band of Unmatched Vengeance Proc", false)
	snapshotManager.AddProc(50402, "Ashen Band of Endless Vengeance Proc", false)
	snapshotManager.AddProc(52571, "Ashen Band of Unmatched Might Proc", false)
	snapshotManager.AddProc(52572, "Ashen Band of Endless Might Proc", false)
	snapshotManager.AddProc(54569, "Sharpened Twilight Scale Proc", false)
	snapshotManager.AddProc(54590, "Sharpened Twilight Scale H Proc", false)
}

func (dk *DpsDeathknight) setupGargoyleCooldowns() {
	dk.ur.gargoyleSnapshot.ClearMajorCooldowns()

	// hyperspeed accelerators
	dk.gargoyleHasteCooldownSync(core.ActionID{SpellID: 54758}, false)

	// berserking (troll)
	dk.gargoyleHasteCooldownSync(core.ActionID{SpellID: 26297}, false)

	// blood fury (orc)
	dk.gargoyleAPCooldownSync(core.ActionID{SpellID: 33697}, false)

	// potion of speed
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 40211}, true)

	// active ap trinkets
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 35937}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 36871}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 37166}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 37556}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 37557}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 38080}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 38081}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 38761}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 39257}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 45263}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 46086}, false)
	dk.gargoyleAPCooldownSync(core.ActionID{ItemID: 47734}, false)

	// active haste trinkets
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 36972}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 37558}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 37560}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 37562}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 38070}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 38258}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 38259}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 38764}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 40531}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 43836}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 45466}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 46088}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 48722}, false)
	dk.gargoyleHasteCooldownSync(core.ActionID{ItemID: 50260}, false)
}

func (dk *DpsDeathknight) gargoyleAPCooldownSync(actionID core.ActionID, isPotion bool) {
	if majorCd := dk.Character.GetMajorCooldown(actionID); majorCd != nil {

		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			if dk.ur.activatingGargoyle {
				return true
			}
			if dk.SummonGargoyle.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration && !isPotion {
				return true
			}
			if dk.SummonGargoyle.CD.ReadyAt() > sim.Duration {
				return true
			}

			return false
		}

		dk.ur.gargoyleSnapshot.AddMajorCooldown(majorCd)
	}
}

func (dk *DpsDeathknight) gargoyleHasteCooldownSync(actionID core.ActionID, isPotion bool) {
	if majorCd := dk.Character.GetMajorCooldown(actionID); majorCd != nil {

		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			if !dk.Rotation.PreNerfedGargoyle {
				aura := dk.GetAura("Summon Gargoyle")

				if aura != nil && aura.IsActive() {
					return true
				}
				if dk.SummonGargoyle.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration-10*time.Second && !isPotion {
					return true
				}
				if dk.SummonGargoyle.CD.ReadyAt() > sim.Duration {
					return true
				}

				return false
			} else {
				if dk.ur.activatingGargoyle {
					return true
				}
				if dk.SummonGargoyle.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration-10*time.Second && !isPotion {
					return true
				}
				if dk.SummonGargoyle.CD.ReadyAt() > sim.Duration {
					return true
				}
			}

			return false
		}

		dk.ur.gargoyleSnapshot.AddMajorCooldown(majorCd)
	}
}

func (dk *DpsDeathknight) setupDrwProcTrackers() {
	snapshotManager := dk.br.drwSnapshot

	snapshotManager.AddProc(40211, "Potion of Speed", true)
	snapshotManager.AddProc(54999, "Hyperspeed Acceleration", true)
	snapshotManager.AddProc(26297, "Berserking (Troll)", true)
	snapshotManager.AddProc(33697, "Blood Fury", true)

	snapshotManager.AddProc(55379, "Thundering Skyflare Diamond Proc", false)
	snapshotManager.AddProc(59626, "Black Magic Proc", false)
	snapshotManager.AddProc(53344, "Rune Of The Fallen Crusader Proc", false)

	snapshotManager.AddProc(37390, "Meteorite Whetstone Proc", false)
	snapshotManager.AddProc(39229, "Embrace of the Spider Proc", false)
	snapshotManager.AddProc(44308, "Signet of Edward the Odd Proc", false)
	snapshotManager.AddProc(43573, "Tears of Bitter Anguish Proc", false)
	snapshotManager.AddProc(45609, "Comet's Trail Proc", false)
	snapshotManager.AddProc(45866, "Elemental Focus Stone Proc", false)

	snapshotManager.AddProc(53344, "Rune Of The Fallen Crusader Proc", false)

	snapshotManager.AddProc(42987, "DMC Greatness Strength Proc", false)

	snapshotManager.AddProc(47115, "Deaths Verdict Strength Proc", false)
	snapshotManager.AddProc(47131, "Deaths Verdict H Strength Proc", false)
	snapshotManager.AddProc(47303, "Deaths Choice Strength Proc", false)
	snapshotManager.AddProc(47464, "Deaths Choice H Strength Proc", false)

	snapshotManager.AddProc(71484, "Deathbringer's Will Strength Proc", false)
	snapshotManager.AddProc(71492, "Deathbringer's Will Haste Proc", false)
	snapshotManager.AddProc(71491, "Deathbringer's Will Crit Proc", false)
	snapshotManager.AddProc(71561, "Deathbringer's Will H Strength Proc", false)
	snapshotManager.AddProc(71560, "Deathbringer's Will H Haste Proc", false)
	snapshotManager.AddProc(71559, "Deathbringer's Will H Crit Proc", false)

	snapshotManager.AddProc(40684, "Mirror of Truth Proc", false)
	snapshotManager.AddProc(40767, "Sonic Booster Proc", false)
	snapshotManager.AddProc(44914, "Anvil of Titans Proc", false)
	snapshotManager.AddProc(45286, "Pyrite Infuser Proc", false)
	snapshotManager.AddProc(45522, "Blood of the Old God Proc", false)
	snapshotManager.AddProc(47214, "Banner of Victory Proc", false)
	snapshotManager.AddProc(49074, "Coren's Chromium Coaster Proc", false)
	snapshotManager.AddProc(50342, "Whispering Fanged Skull Proc", false)
	snapshotManager.AddProc(50343, "Whispering Fanged Skull H Proc", false)
	snapshotManager.AddProc(50401, "Ashen Band of Unmatched Vengeance Proc", false)
	snapshotManager.AddProc(50402, "Ashen Band of Endless Vengeance Proc", false)
	snapshotManager.AddProc(52571, "Ashen Band of Unmatched Might Proc", false)
	snapshotManager.AddProc(52572, "Ashen Band of Endless Might Proc", false)
	snapshotManager.AddProc(54569, "Sharpened Twilight Scale Proc", false)
	snapshotManager.AddProc(54590, "Sharpened Twilight Scale H Proc", false)

	//snapshotManager.AddProc(40256, "Grim Toll Proc", false)
	//snapshotManager.AddProc(45931, "Mjolnir Runestone Proc", false)
	snapshotManager.AddProc(46038, "Dark Matter Proc", false)
	snapshotManager.AddProc(50198, "Needle-Encrusted Scorpion Proc", false)
}

func (dk *DpsDeathknight) setupDrwCooldowns() {
	dk.br.drwSnapshot.ClearMajorCooldowns()

	// Unholy Frenzy
	dk.drwCooldownSync(core.ActionID{SpellID: 49016, Tag: dk.Index}, false)

	// hyperspeed accelerators
	dk.drwCooldownSync(core.ActionID{SpellID: 54758}, false)

	// berserking (troll)
	dk.drwCooldownSync(core.ActionID{SpellID: 26297}, false)

	// blood fury (orc)
	dk.drwCooldownSync(core.ActionID{SpellID: 33697}, false)

	// potion of speed
	dk.drwCooldownSync(core.ActionID{ItemID: 40211}, true)

	// active ap trinkets
	dk.drwCooldownSync(core.ActionID{ItemID: 35937}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 36871}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 37166}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 37556}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 37557}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38080}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38081}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38761}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 39257}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 45263}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 46086}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 47734}, false)

	// active haste trinkets
	dk.drwCooldownSync(core.ActionID{ItemID: 36972}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 37558}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 37560}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 37562}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38070}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38258}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38259}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 38764}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 40531}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 43836}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 45466}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 46088}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 48722}, false)
	dk.drwCooldownSync(core.ActionID{ItemID: 50260}, false)
}

func (dk *DpsDeathknight) drwCooldownSync(actionID core.ActionID, isPotion bool) {
	if majorCd := dk.Character.GetMajorCooldown(actionID); majorCd != nil {

		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			if character != &dk.Character {
				return true
			}
			// Opener use everything
			if sim.CurrentTime < 2*time.Second {
				return true
			}
			// If the fight is long enough for Unholy Frenzy we use potion with it
			if isPotion && dk.br.activatingDrw && sim.Duration > 200*time.Second {
				if dk.UnholyFrenzy.IsReady(sim) || dk.UnholyFrenzyAura.IsActive() {
					return true
				}
				return false
			}
			if dk.br.activatingDrw {
				return true
			}
			if dk.DancingRuneWeapon.CD.TimeToReady(sim) > majorCd.Spell.CD.Duration && !isPotion {
				return true
			}
			if dk.DancingRuneWeapon.CD.ReadyAt() > sim.Duration {
				return true
			}

			return false
		}

		dk.br.drwSnapshot.AddMajorCooldown(majorCd)
	}
}

func (dk *DpsDeathknight) Reset(sim *core.Simulation) {
	dk.Deathknight.Reset(sim)

	dk.sr.Reset(sim)
	dk.br.Reset(sim)
	dk.fr.Reset(sim)
	dk.ur.Reset(sim)

	dk.SetupRotations()

	dk.Presence = deathknight.UnsetPresence

	b, f, u := deathknight.PointsInTalents(dk.Talents)

	if f > u && f > b {
		if dk.Rotation.Presence == proto.Deathknight_Rotation_Blood {
			dk.ChangePresence(sim, deathknight.BloodPresence)
		} else if dk.Rotation.Presence == proto.Deathknight_Rotation_Frost {
			dk.ChangePresence(sim, deathknight.FrostPresence)
		} else if dk.Rotation.Presence == proto.Deathknight_Rotation_Unholy {
			dk.ChangePresence(sim, deathknight.UnholyPresence)
		}
	}

	if u > f && u > b {
		if dk.Rotation.StartingPresence == proto.Deathknight_Rotation_Unholy {
			dk.ChangePresence(sim, deathknight.UnholyPresence)
		} else if dk.Talents.SummonGargoyle {
			dk.ChangePresence(sim, deathknight.BloodPresence)
		}
	}

	if b > f && b > u {
		dk.ChangePresence(sim, deathknight.BloodPresence)
	}
}
