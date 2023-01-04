import { IndividualSimUI, OtherDefaults } from '../core/individual_sim_ui.js';
import { Raid as RaidProto } from '../core/proto/api.js';
import { Party as PartyProto } from '../core/proto/api.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';

import { Encounter as EncounterProto } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Race } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { SpecOptions } from '../core/proto_utils/utils.js';
import { SpecRotation } from '../core/proto_utils/utils.js';
import { playerToSpec } from '../core/proto_utils/utils.js';
import { specIconsLarge } from '../core/proto_utils/utils.js';
import { specNames } from '../core/proto_utils/utils.js';
import { talentTreeIcons } from '../core/proto_utils/utils.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { Player } from '../core/player.js';

import * as TankDeathknightPresets from '../tank_deathknight/presets.js';
import * as DeathknightPresets from '../deathknight/presets.js';
import * as BalanceDruidPresets from '../balance_druid/presets.js';
import * as FeralDruidPresets from '../feral_druid/presets.js';
import * as FeralTankDruidPresets from '../feral_tank_druid/presets.js';
import * as RestorationDruidPresets from '../restoration_druid/presets.js';
import * as ElementalShamanPresets from '../elemental_shaman/presets.js';
import * as EnhancementShamanPresets from '../enhancement_shaman/presets.js';
import * as RestorationShamanPresets from '../restoration_shaman/presets.js';
import * as HunterPresets from '../hunter/presets.js';
import * as MagePresets from '../mage/presets.js';
import * as RoguePresets from '../rogue/presets.js';
import * as HolyPaladinPresets from '../holy_paladin/presets.js';
import * as ProtectionPaladinPresets from '../protection_paladin/presets.js';
import * as RetributionPaladinPresets from '../retribution_paladin/presets.js';
import * as HealingPriestPresets from '../healing_priest/presets.js';
import * as ShadowPriestPresets from '../shadow_priest/presets.js';
import * as SmitePriestPresets from '../smite_priest/presets.js';
import * as WarriorPresets from '../warrior/presets.js';
import * as ProtectionWarriorPresets from '../protection_warrior/presets.js';
import * as WarlockPresets from '../warlock/presets.js';

import { TankDeathknightSimUI } from '../tank_deathknight/sim.js';
import { DeathknightSimUI } from '../deathknight/sim.js';
import { BalanceDruidSimUI } from '../balance_druid/sim.js';
import { FeralDruidSimUI } from '../feral_druid/sim.js';
import { FeralTankDruidSimUI } from '../feral_tank_druid/sim.js';
import { RestorationDruidSimUI } from '../restoration_druid/sim.js';
import { ElementalShamanSimUI } from '../elemental_shaman/sim.js';
import { EnhancementShamanSimUI } from '../enhancement_shaman/sim.js';
import { RestorationShamanSimUI } from '../restoration_shaman/sim.js';
import { HunterSimUI } from '../hunter/sim.js';
import { MageSimUI } from '../mage/sim.js';
import { RogueSimUI } from '../rogue/sim.js';
import { HolyPaladinSimUI } from '../holy_paladin/sim.js';
import { ProtectionPaladinSimUI } from '../protection_paladin/sim.js';
import { RetributionPaladinSimUI } from '../retribution_paladin/sim.js';
import { HealingPriestSimUI } from '../healing_priest/sim.js';
import { ShadowPriestSimUI } from '../shadow_priest/sim.js';
import { SmitePriestSimUI } from '../smite_priest/sim.js';
import { WarriorSimUI } from '../warrior/sim.js';
import { ProtectionWarriorSimUI } from '../protection_warrior/sim.js';
import { WarlockSimUI } from '../warlock/sim.js';

export const specSimFactories: Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>> = {
	[Spec.SpecTankDeathknight]: (parentElem: HTMLElement, player: Player<any>) => new TankDeathknightSimUI(parentElem, player),
	[Spec.SpecDeathknight]: (parentElem: HTMLElement, player: Player<any>) => new DeathknightSimUI(parentElem, player),
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecFeralTankDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralTankDruidSimUI(parentElem, player),
	[Spec.SpecRestorationDruid]: (parentElem: HTMLElement, player: Player<any>) => new RestorationDruidSimUI(parentElem, player),
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecRestorationShaman]: (parentElem: HTMLElement, player: Player<any>) => new RestorationShamanSimUI(parentElem, player),
	[Spec.SpecHunter]: (parentElem: HTMLElement, player: Player<any>) => new HunterSimUI(parentElem, player),
	[Spec.SpecMage]: (parentElem: HTMLElement, player: Player<any>) => new MageSimUI(parentElem, player),
	[Spec.SpecRogue]: (parentElem: HTMLElement, player: Player<any>) => new RogueSimUI(parentElem, player),
	[Spec.SpecHolyPaladin]: (parentElem: HTMLElement, player: Player<any>) => new HolyPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	[Spec.SpecHealingPriest]: (parentElem: HTMLElement, player: Player<any>) => new HealingPriestSimUI(parentElem, player),
	[Spec.SpecShadowPriest]: (parentElem: HTMLElement, player: Player<any>) => new ShadowPriestSimUI(parentElem, player),
	[Spec.SpecSmitePriest]: (parentElem: HTMLElement, player: Player<any>) => new SmitePriestSimUI(parentElem, player),
	[Spec.SpecWarrior]: (parentElem: HTMLElement, player: Player<any>) => new WarriorSimUI(parentElem, player),
	[Spec.SpecProtectionWarrior]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionWarriorSimUI(parentElem, player),
	[Spec.SpecWarlock]: (parentElem: HTMLElement, player: Player<any>) => new WarlockSimUI(parentElem, player),
};

// Configuration necessary for creating new players.
export interface PresetSpecSettings<SpecType extends Spec> {
	spec: Spec,
	rotation: SpecRotation<SpecType>,
	talents: SavedTalents,
	specOptions: SpecOptions<SpecType>,
	consumes: Consumes,

	defaultName: string,
	defaultFactionRaces: Record<Faction, Race>,
	defaultGear: Record<Faction, Record<number, EquipmentSpec>>,
	otherDefaults?: OtherDefaults,

	tooltip: string,
	iconUrl: string,
}

export const playerPresets: Array<PresetSpecSettings<any>> = [
	{
		spec: Spec.SpecTankDeathknight,
		rotation: TankDeathknightPresets.DefaultRotation,
		talents: TankDeathknightPresets.BloodTalents.data,
		specOptions: TankDeathknightPresets.DefaultOptions,
		consumes: TankDeathknightPresets.DefaultConsumes,
		defaultName: 'Blood DK',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: TankDeathknightPresets.P1_BLOOD_BIS_PRESET.gear,
			},
			[Faction.Horde]: {
				1: TankDeathknightPresets.P1_BLOOD_BIS_PRESET.gear,
			},
		},
		tooltip: 'Blood Tank Death Knight',
		iconUrl: talentTreeIcons[Class.ClassDeathknight][0],
	},
	{
		spec: Spec.SpecDeathknight,
		rotation: DeathknightPresets.DefaultFrostRotation,
		talents: DeathknightPresets.FrostTalents.data,
		specOptions: DeathknightPresets.DefaultFrostOptions,
		consumes: DeathknightPresets.DefaultConsumes,
		defaultName: 'Frost DK',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: DeathknightPresets.P1_FROST_BIS_PRESET.gear,
			},
			[Faction.Horde]: {
				1: DeathknightPresets.P1_FROST_BIS_PRESET.gear,
			},
		},
		otherDefaults: DeathknightPresets.OtherDefaults,
		tooltip: 'Frost Death Knight',
		iconUrl: talentTreeIcons[Class.ClassDeathknight][1],
	},
	{
		spec: Spec.SpecDeathknight,
		rotation: DeathknightPresets.DefaultUnholyRotation,
		talents: DeathknightPresets.UnholyDualWieldTalents.data,
		specOptions: DeathknightPresets.DefaultUnholyOptions,
		consumes: DeathknightPresets.DefaultConsumes,
		defaultName: 'DW Unholy DK',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: DeathknightPresets.P1_UNHOLY_DW_BIS_PRESET.gear,
			},
			[Faction.Horde]: {
				1: DeathknightPresets.P1_UNHOLY_DW_BIS_PRESET.gear,
			},
		},
		otherDefaults: DeathknightPresets.OtherDefaults,
		tooltip: 'Dual Wield Unholy DK',
		iconUrl: talentTreeIcons[Class.ClassDeathknight][2],
	},
	//{
	//	spec: Spec.SpecDeathknight,
	//	rotation: DeathknightPresets.DefaultBloodRotation,
	//	talents: DeathknightPresets.BloodTalents.data,
	//	specOptions: DeathknightPresets.DefaultBloodOptions,
	//	consumes: DeathknightPresets.DefaultConsumes,
	//	defaultName: 'Blood Dps DK',
	//	defaultFactionRaces: {
	//		[Faction.Unknown]: Race.RaceUnknown,
	//		[Faction.Alliance]: Race.RaceHuman,
	//		[Faction.Horde]: Race.RaceTroll,
	//	},
	//	defaultGear: {
	//		[Faction.Unknown]: {},
	//		[Faction.Alliance]: {
	//			1: DeathknightPresets.P1_UNHOLY_2H_BIS_PRESET.gear,
	//		},
	//		[Faction.Horde]: {
	//			1: DeathknightPresets.P1_UNHOLY_2H_BIS_PRESET.gear,
	//		},
	//	},
	//	otherDefaults: DeathknightPresets.OtherDefaults,
	//	tooltip: 'Blood Dps DK',
	//	iconUrl: talentTreeIcons[Class.ClassDeathknight][0],
	//},
	{
		spec: Spec.SpecBalanceDruid,
		rotation: BalanceDruidPresets.DefaultRotation,
		talents: BalanceDruidPresets.StandardTalents.data,
		specOptions: BalanceDruidPresets.DefaultOptions,
		consumes: BalanceDruidPresets.DefaultConsumes,
		defaultName: 'Balance Druid',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: BalanceDruidPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: BalanceDruidPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecBalanceDruid],
		iconUrl: specIconsLarge[Spec.SpecBalanceDruid],
	},
	{
		spec: Spec.SpecFeralDruid,
		rotation: FeralDruidPresets.DefaultRotation,
		talents: FeralDruidPresets.StandardTalents.data,
		specOptions: FeralDruidPresets.DefaultOptions,
		consumes: FeralDruidPresets.DefaultConsumes,
		defaultName: 'Cat Druid',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: FeralDruidPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: FeralDruidPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecFeralDruid],
		iconUrl: specIconsLarge[Spec.SpecFeralDruid],
	},
	{
		spec: Spec.SpecFeralTankDruid,
		rotation: FeralTankDruidPresets.DefaultRotation,
		talents: FeralTankDruidPresets.StandardTalents.data,
		specOptions: FeralTankDruidPresets.DefaultOptions,
		consumes: FeralTankDruidPresets.DefaultConsumes,
		defaultName: 'Bear Druid',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: FeralTankDruidPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: FeralTankDruidPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecFeralTankDruid],
		iconUrl: specIconsLarge[Spec.SpecFeralTankDruid],
	},
	{
		spec: Spec.SpecRestorationDruid,
		rotation: RestorationDruidPresets.DefaultRotation,
		talents: RestorationDruidPresets.CelestialFocusTalents.data,
		specOptions: RestorationDruidPresets.DefaultOptions,
		consumes: RestorationDruidPresets.DefaultConsumes,
		defaultName: 'Restoration Druid',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceTauren,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RestorationDruidPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RestorationDruidPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecRestorationDruid],
		iconUrl: specIconsLarge[Spec.SpecRestorationDruid],
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.BeastMasteryTalents.data,
		specOptions: HunterPresets.BMDefaultOptions,
		consumes: HunterPresets.DefaultConsumes,
		defaultName: 'BM Hunter',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HunterPresets.MM_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.MM_P1_PRESET.gear,
			},
		},
		tooltip: 'BM Hunter',
		iconUrl: talentTreeIcons[Class.ClassHunter][0],
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.MarksmanTalents.data,
		specOptions: HunterPresets.DefaultOptions,
		consumes: HunterPresets.DefaultConsumes,
		defaultName: 'MM Hunter',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HunterPresets.MM_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.MM_P1_PRESET.gear,
			},
		},
		tooltip: 'MM Hunter',
		iconUrl: talentTreeIcons[Class.ClassHunter][1],
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.SurvivalTalents.data,
		specOptions: HunterPresets.DefaultOptions,
		consumes: HunterPresets.DefaultConsumes,
		defaultName: 'SV Hunter',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceNightElf,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HunterPresets.SV_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.SV_P1_PRESET.gear,
			},
		},
		tooltip: 'SV Hunter',
		iconUrl: talentTreeIcons[Class.ClassHunter][2],
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultArcaneRotation,
		talents: MagePresets.ArcaneTalents.data,
		specOptions: MagePresets.DefaultArcaneOptions,
		consumes: MagePresets.DefaultArcaneConsumes,
		defaultName: 'Arcane Mage',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.ARCANE_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.ARCANE_P1_PRESET.gear,
			},
		},
		tooltip: 'Arcane Mage',
		iconUrl: talentTreeIcons[Class.ClassMage][0],
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultFireRotation,
		talents: MagePresets.FireTalents.data,
		specOptions: MagePresets.DefaultFireOptions,
		consumes: MagePresets.DefaultFireConsumes,
		defaultName: 'Fire Mage',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.FIRE_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.FIRE_P1_PRESET.gear,
			},
		},
		tooltip: 'Fire Mage',
		iconUrl: talentTreeIcons[Class.ClassMage][1],
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultFrostRotation,
		talents: MagePresets.FrostTalents.data,
		specOptions: MagePresets.DefaultFrostOptions,
		consumes: MagePresets.DefaultFrostConsumes,
		defaultName: 'Frost Mage',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.FROST_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.FROST_P1_PRESET.gear,
			},
		},
		tooltip: 'Frost Mage',
		iconUrl: talentTreeIcons[Class.ClassMage][2],
	},
	{
		spec: Spec.SpecRogue,
		rotation: RoguePresets.DefaultRotation,
		talents: RoguePresets.AssassinationTalents.data,
		specOptions: RoguePresets.DefaultOptions,
		consumes: RoguePresets.DefaultConsumes,
		defaultName: 'Assassination Rogue',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RoguePresets.P1_PRESET_ASSASSINATION.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.P1_PRESET_ASSASSINATION.gear,
			},
		},
		tooltip: 'Assassination Rogue',
		iconUrl: talentTreeIcons[Class.ClassRogue][0],
	},
	{
		spec: Spec.SpecRogue,
		rotation: RoguePresets.DefaultRotation,
		talents: RoguePresets.CombatTalents.data,
		specOptions: RoguePresets.DefaultOptions,
		consumes: RoguePresets.DefaultConsumes,
		defaultName: 'Combat Rogue',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RoguePresets.P1_PRESET_COMBAT.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.P1_PRESET_COMBAT.gear,
			},
		},
		tooltip: 'Combat Rogue',
		iconUrl: talentTreeIcons[Class.ClassRogue][1],
	},
	{
		spec: Spec.SpecElementalShaman,
		rotation: ElementalShamanPresets.DefaultRotation,
		talents: ElementalShamanPresets.StandardTalents.data,
		specOptions: ElementalShamanPresets.DefaultOptions,
		consumes: ElementalShamanPresets.DefaultConsumes,
		defaultName: 'Ele Shaman',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDraenei,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ElementalShamanPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ElementalShamanPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecElementalShaman],
		iconUrl: specIconsLarge[Spec.SpecElementalShaman],
	},
	{
		spec: Spec.SpecEnhancementShaman,
		rotation: EnhancementShamanPresets.DefaultRotation,
		talents: EnhancementShamanPresets.StandardTalents.data,
		specOptions: EnhancementShamanPresets.DefaultOptions,
		consumes: EnhancementShamanPresets.DefaultConsumes,
		defaultName: 'Enh Shaman',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDraenei,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: EnhancementShamanPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: EnhancementShamanPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecEnhancementShaman],
		iconUrl: specIconsLarge[Spec.SpecEnhancementShaman],
	},
	{
		spec: Spec.SpecRestorationShaman,
		rotation: RestorationShamanPresets.DefaultRotation,
		talents: RestorationShamanPresets.RaidHealingTalents.data,
		specOptions: RestorationShamanPresets.DefaultOptions,
		consumes: RestorationShamanPresets.DefaultConsumes,
		defaultName: 'Resto Shaman',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDraenei,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RestorationShamanPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RestorationShamanPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecRestorationShaman],
		iconUrl: specIconsLarge[Spec.SpecRestorationShaman],
	},
	{
		spec: Spec.SpecHealingPriest,
		rotation: HealingPriestPresets.DiscDefaultRotation,
		talents: HealingPriestPresets.DiscTalents.data,
		specOptions: HealingPriestPresets.DefaultOptions,
		consumes: HealingPriestPresets.DefaultConsumes,
		defaultName: 'Disc Priest',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HealingPriestPresets.DISC_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HealingPriestPresets.DISC_P1_PRESET.gear,
			},
		},
		tooltip: 'Disc Priest',
		iconUrl: talentTreeIcons[Class.ClassPriest][0],
	},
	{
		spec: Spec.SpecHealingPriest,
		rotation: HealingPriestPresets.HolyDefaultRotation,
		talents: HealingPriestPresets.HolyTalents.data,
		specOptions: HealingPriestPresets.DefaultOptions,
		consumes: HealingPriestPresets.DefaultConsumes,
		defaultName: 'Holy Priest',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HealingPriestPresets.HOLY_P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HealingPriestPresets.HOLY_P1_PRESET.gear,
			},
		},
		tooltip: 'Holy Priest',
		iconUrl: talentTreeIcons[Class.ClassPriest][1],
	},
	{
		spec: Spec.SpecShadowPriest,
		rotation: ShadowPriestPresets.DefaultRotation,
		talents: ShadowPriestPresets.StandardTalents.data,
		specOptions: ShadowPriestPresets.DefaultOptions,
		consumes: ShadowPriestPresets.DefaultConsumes,
		defaultName: 'Shadow Priest',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ShadowPriestPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ShadowPriestPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecShadowPriest],
		iconUrl: specIconsLarge[Spec.SpecShadowPriest],
	},
	{
		spec: Spec.SpecSmitePriest,
		rotation: SmitePriestPresets.DefaultRotation,
		talents: SmitePriestPresets.StandardTalents.data,
		specOptions: SmitePriestPresets.DefaultOptions,
		consumes: SmitePriestPresets.DefaultConsumes,
		defaultName: 'Smite Priest',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceDwarf,
			[Faction.Horde]: Race.RaceUndead,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: SmitePriestPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: SmitePriestPresets.P1_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecSmitePriest],
		iconUrl: specIconsLarge[Spec.SpecSmitePriest],
	},
	{
		spec: Spec.SpecWarrior,
		rotation: WarriorPresets.ArmsRotation,
		talents: WarriorPresets.ArmsTalents.data,
		specOptions: WarriorPresets.DefaultOptions,
		consumes: WarriorPresets.DefaultConsumes,
		defaultName: 'Arms Warrior',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarriorPresets.P1_ARMS_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarriorPresets.P1_ARMS_PRESET.gear,
			},
		},
		tooltip: 'Arms Warrior',
		iconUrl: talentTreeIcons[Class.ClassWarrior][0],
	},
	{
		spec: Spec.SpecWarrior,
		rotation: WarriorPresets.DefaultRotation,
		talents: WarriorPresets.FuryTalents.data,
		specOptions: WarriorPresets.DefaultOptions,
		consumes: WarriorPresets.DefaultConsumes,
		defaultName: 'Fury Warrior',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarriorPresets.P1_FURY_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarriorPresets.P1_FURY_PRESET.gear,
			},
		},
		tooltip: 'Fury Warrior',
		iconUrl: talentTreeIcons[Class.ClassWarrior][1],
	},
	{
		spec: Spec.SpecProtectionWarrior,
		rotation: ProtectionWarriorPresets.DefaultRotation,
		talents: ProtectionWarriorPresets.StandardTalents.data,
		specOptions: ProtectionWarriorPresets.DefaultOptions,
		consumes: ProtectionWarriorPresets.DefaultConsumes,
		defaultName: 'Prot Warrior',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ProtectionWarriorPresets.P1_BALANCED_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ProtectionWarriorPresets.P1_BALANCED_PRESET.gear,
			},
		},
		tooltip: 'Protection Warrior',
		iconUrl: talentTreeIcons[Class.ClassWarrior][2],
	},
	{
		spec: Spec.SpecHolyPaladin,
		rotation: HolyPaladinPresets.DefaultRotation,
		talents: HolyPaladinPresets.StandardTalents.data,
		specOptions: HolyPaladinPresets.DefaultOptions,
		consumes: HolyPaladinPresets.DefaultConsumes,
		defaultName: 'Holy Paladin',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: HolyPaladinPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HolyPaladinPresets.P1_PRESET.gear,
			},
		},
		tooltip: 'Holy Paladin',
		iconUrl: talentTreeIcons[Class.ClassPaladin][0],
	},
	{
		spec: Spec.SpecProtectionPaladin,
		rotation: ProtectionPaladinPresets.DefaultRotation,
		talents: ProtectionPaladinPresets.GenericAoeTalents.data,
		specOptions: ProtectionPaladinPresets.DefaultOptions,
		consumes: ProtectionPaladinPresets.DefaultConsumes,
		defaultName: 'Prot Paladin',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: ProtectionPaladinPresets.P1_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ProtectionPaladinPresets.P1_PRESET.gear,
			},
		},
		tooltip: 'Protection Paladin',
		iconUrl: talentTreeIcons[Class.ClassPaladin][1],
	},
	{
		spec: Spec.SpecRetributionPaladin,
		rotation: RetributionPaladinPresets.DefaultRotation,
		talents: RetributionPaladinPresets.AuraMasteryTalents.data,
		specOptions: RetributionPaladinPresets.DefaultOptions,
		consumes: RetributionPaladinPresets.DefaultConsumes,
		defaultName: 'Ret Paladin',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RetributionPaladinPresets.P1_PRESET.gear,
				2: RetributionPaladinPresets.P2_PRESET.gear,
				3: RetributionPaladinPresets.P3_PRESET.gear,
				4: RetributionPaladinPresets.P4_PRESET.gear,
				5: RetributionPaladinPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RetributionPaladinPresets.P1_PRESET.gear,
				2: RetributionPaladinPresets.P2_PRESET.gear,
				3: RetributionPaladinPresets.P3_PRESET.gear,
				4: RetributionPaladinPresets.P4_PRESET.gear,
				5: RetributionPaladinPresets.P5_PRESET.gear,
			},
		},
		tooltip: 'Ret Paladin',
		iconUrl: talentTreeIcons[Class.ClassPaladin][2],
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.AfflictionRotation,
		talents: WarlockPresets.AfflictionTalents.data,
		specOptions: WarlockPresets.AfflictionOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Affli Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {
				1: WarlockPresets.P1_Preset_Affliction.gear,
				2: WarlockPresets.P1_PreBiS_11.gear,
			},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_Preset_Affliction.gear,
				2: WarlockPresets.P1_PreBiS_11.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_Preset_Affliction.gear,
				2: WarlockPresets.P1_PreBiS_11.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Affliction Warlock',
		iconUrl: talentTreeIcons[Class.ClassWarlock][0],
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.DemonologyRotation,
		talents: WarlockPresets.DemonologyTalents.data,
		specOptions: WarlockPresets.DemonologyOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Demo Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {
				1: WarlockPresets.P1_Preset_Demo_Destro.gear,
				2: WarlockPresets.P1_PreBiS_14.gear,
			},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_Preset_Demo_Destro.gear,
				2: WarlockPresets.P1_PreBiS_14.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_Preset_Demo_Destro.gear,
				2: WarlockPresets.P1_PreBiS_14.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Demonology Warlock',
		iconUrl: talentTreeIcons[Class.ClassWarlock][1],
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.DestructionRotation,
		talents: WarlockPresets.DestructionTalents.data,
		specOptions: WarlockPresets.DestructionOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Destro Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {
				1: WarlockPresets.P1_Preset_Demo_Destro.gear,
				2: WarlockPresets.P1_PreBiS_14.gear,
			},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_Preset_Demo_Destro.gear,
				2: WarlockPresets.P1_PreBiS_14.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_Preset_Demo_Destro.gear,
				2: WarlockPresets.P1_PreBiS_14.gear,
			},
		},
		otherDefaults: WarlockPresets.OtherDefaults,
		tooltip: 'Destruction Warlock',
		iconUrl: talentTreeIcons[Class.ClassWarlock][2],
	},
];

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];
