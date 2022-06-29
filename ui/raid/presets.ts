import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Raid as RaidProto } from '/tbc/core/proto/api.js';
import { Party as PartyProto } from '/tbc/core/proto/api.js';
import { Class } from '/tbc/core/proto/common.js';
import { Consumes } from '/tbc/core/proto/common.js';
import { Drums } from '/tbc/core/proto/common.js';
import { Encounter as EncounterProto } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto/common.js';
import { SpecOptions } from '/tbc/core/proto_utils/utils.js';
import { SpecRotation } from '/tbc/core/proto_utils/utils.js';
import { playerToSpec } from '/tbc/core/proto_utils/utils.js';
import { specIconsLarge } from '/tbc/core/proto_utils/utils.js';
import { specNames } from '/tbc/core/proto_utils/utils.js';
import { talentTreeIcons } from '/tbc/core/proto_utils/utils.js';
import { NO_TARGET } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import { BuffBot } from './buff_bot.js';

import * as BalanceDruidPresets from '/tbc/balance_druid/presets.js';
import * as FeralDruidPresets from '/tbc/feral_druid/presets.js';
import * as FeralTankDruidPresets from '/tbc/feral_tank_druid/presets.js';
import * as ElementalShamanPresets from '/tbc/elemental_shaman/presets.js';
import * as EnhancementShamanPresets from '/tbc/enhancement_shaman/presets.js';
import * as HunterPresets from '/tbc/hunter/presets.js';
import * as MagePresets from '/tbc/mage/presets.js';
import * as RoguePresets from '/tbc/rogue/presets.js';
import * as RetributionPaladinPresets from '/tbc/retribution_paladin/presets.js';
import * as ProtectionPaladinPresets from '/tbc/protection_paladin/presets.js';
import * as ShadowPriestPresets from '/tbc/shadow_priest/presets.js';
import * as SmitePriestPresets from '/tbc/smite_priest/presets.js';
import * as WarriorPresets from '/tbc/warrior/presets.js';
import * as ProtectionWarriorPresets from '/tbc/protection_warrior/presets.js';
import * as WarlockPresets from '/tbc/warlock/presets.js';


import { BalanceDruidSimUI } from '/tbc/balance_druid/sim.js';
import { FeralDruidSimUI } from '/tbc/feral_druid/sim.js';
import { FeralTankDruidSimUI } from '/tbc/feral_tank_druid/sim.js';
import { EnhancementShamanSimUI } from '/tbc/enhancement_shaman/sim.js';
import { ElementalShamanSimUI } from '/tbc/elemental_shaman/sim.js';
import { HunterSimUI } from '/tbc/hunter/sim.js';
import { MageSimUI } from '/tbc/mage/sim.js';
import { RogueSimUI } from '/tbc/rogue/sim.js';
import { RetributionPaladinSimUI } from '/tbc/retribution_paladin/sim.js';
import { ProtectionPaladinSimUI } from '/tbc/protection_paladin/sim.js';
import { ShadowPriestSimUI } from '/tbc/shadow_priest/sim.js';
import { SmitePriestSimUI } from '/tbc/smite_priest/sim.js';
import { WarriorSimUI } from '/tbc/warrior/sim.js';
import { ProtectionWarriorSimUI } from '/tbc/protection_warrior/sim.js';
import { WarlockSimUI } from '/tbc/warlock/sim.js';

export const specSimFactories: Partial<Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>>> = {
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecFeralTankDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralTankDruidSimUI(parentElem, player),
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecHunter]: (parentElem: HTMLElement, player: Player<any>) => new HunterSimUI(parentElem, player),
	[Spec.SpecMage]: (parentElem: HTMLElement, player: Player<any>) => new MageSimUI(parentElem, player),
	[Spec.SpecRogue]: (parentElem: HTMLElement, player: Player<any>) => new RogueSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
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
	talents: string,
	specOptions: SpecOptions<SpecType>,
	consumes: Consumes,

	defaultName: string,
	defaultFactionRaces: Record<Faction, Race>,
	defaultGear: Record<Faction, Record<number, EquipmentSpec>>,

	tooltip: string,
	iconUrl: string,
}

// Configuration necessary for creating new BuffBots.
export interface BuffBotSettings {
	// The value of this field must never change, to preserve local storage data.
	buffBotId: string,

	// Set this to true to remove a buff bot option after launching a real sim.
	// This will allow users with saved settings to properly load the buffbot but
	// also remove the buffbot as an option from the UI.
	deprecated?: boolean,

	spec: Spec,
	name: string,
	tooltip: string,
	iconUrl: string,

	// Callback to apply buffs from this buff bot.
	modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => void,
}

export const playerPresets: Array<PresetSpecSettings<any>> = [
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
				1: BalanceDruidPresets.P1_ALLIANCE_PRESET.gear,
				2: BalanceDruidPresets.P2_ALLIANCE_PRESET.gear,
				3: BalanceDruidPresets.P3_PRESET.gear,
				4: BalanceDruidPresets.P4_PRESET.gear,
				5: BalanceDruidPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: BalanceDruidPresets.P1_HORDE_PRESET.gear,
				2: BalanceDruidPresets.P2_HORDE_PRESET.gear,
				3: BalanceDruidPresets.P3_PRESET.gear,
				4: BalanceDruidPresets.P4_PRESET.gear,
				5: BalanceDruidPresets.P5_PRESET.gear,
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
				2: FeralDruidPresets.P2_PRESET.gear,
				3: FeralDruidPresets.P3_PRESET.gear,
				4: FeralDruidPresets.P4_PRESET.gear,
				5: FeralDruidPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: FeralDruidPresets.P1_PRESET.gear,
				2: FeralDruidPresets.P2_PRESET.gear,
				3: FeralDruidPresets.P3_PRESET.gear,
				4: FeralDruidPresets.P4_PRESET.gear,
				5: FeralDruidPresets.P5_PRESET.gear,
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
				2: FeralTankDruidPresets.P2_PRESET.gear,
				3: FeralTankDruidPresets.P3_PRESET.gear,
				4: FeralTankDruidPresets.P4_PRESET.gear,
				5: FeralTankDruidPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: FeralTankDruidPresets.P1_PRESET.gear,
				2: FeralTankDruidPresets.P2_PRESET.gear,
				3: FeralTankDruidPresets.P3_PRESET.gear,
				4: FeralTankDruidPresets.P4_PRESET.gear,
				5: FeralTankDruidPresets.P5_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecFeralTankDruid],
		iconUrl: specIconsLarge[Spec.SpecFeralTankDruid],
	},
	{
		spec: Spec.SpecHunter,
		rotation: HunterPresets.DefaultRotation,
		talents: HunterPresets.BeastMasteryTalents.data,
		specOptions: HunterPresets.DefaultOptions,
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
				1: HunterPresets.P1_BM_PRESET.gear,
				2: HunterPresets.P2_BM_PRESET.gear,
				3: HunterPresets.P3_BM_PRESET.gear,
				4: HunterPresets.P4_BM_PRESET.gear,
				5: HunterPresets.P5_BM_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.P1_BM_PRESET.gear,
				2: HunterPresets.P2_BM_PRESET.gear,
				3: HunterPresets.P3_BM_PRESET.gear,
				4: HunterPresets.P4_BM_PRESET.gear,
				5: HunterPresets.P5_BM_PRESET.gear,
			},
		},
		tooltip: 'BM Hunter',
		iconUrl: talentTreeIcons[Class.ClassHunter][0],
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
				1: HunterPresets.P1_SV_PRESET.gear,
				2: HunterPresets.P2_SV_PRESET.gear,
				3: HunterPresets.P3_SV_PRESET.gear,
				4: HunterPresets.P4_SV_PRESET.gear,
				5: HunterPresets.P5_SV_PRESET.gear,
			},
			[Faction.Horde]: {
				1: HunterPresets.P1_SV_PRESET.gear,
				2: HunterPresets.P2_SV_PRESET.gear,
				3: HunterPresets.P3_SV_PRESET.gear,
				4: HunterPresets.P4_SV_PRESET.gear,
				5: HunterPresets.P5_SV_PRESET.gear,
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
			[Faction.Horde]: Race.RaceTroll10,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.P1_ARCANE_PRESET.gear,
				2: MagePresets.P2_ARCANE_PRESET.gear,
				3: MagePresets.P3_ARCANE_PRESET.gear,
				4: MagePresets.P4_ARCANE_PRESET.gear,
				5: MagePresets.P5_ARCANE_PRESET.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.P1_ARCANE_PRESET.gear,
				2: MagePresets.P2_ARCANE_PRESET.gear,
				3: MagePresets.P3_ARCANE_PRESET.gear,
				4: MagePresets.P4_ARCANE_PRESET.gear,
				5: MagePresets.P5_ARCANE_PRESET.gear,
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
			[Faction.Horde]: Race.RaceTroll10,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.P1_FIRE_PRESET.gear,
				2: MagePresets.P2_FIRE_PRESET.gear,
				3: MagePresets.P3_FIRE_PRESET.gear,
				4: MagePresets.P4_FIRE_PRESET.gear,
				5: MagePresets.P5_FIRE_PRESET.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.P1_FIRE_PRESET.gear,
				2: MagePresets.P2_FIRE_PRESET.gear,
				3: MagePresets.P3_FIRE_PRESET.gear,
				4: MagePresets.P4_FIRE_PRESET.gear,
				5: MagePresets.P5_FIRE_PRESET.gear,
			},
		},
		tooltip: 'Fire Mage',
		iconUrl: talentTreeIcons[Class.ClassMage][1],
	},
	{
		spec: Spec.SpecMage,
		rotation: MagePresets.DefaultFrostRotation,
		talents: MagePresets.DeepFrostTalents.data,
		specOptions: MagePresets.DefaultFrostOptions,
		consumes: MagePresets.DefaultFrostConsumes,
		defaultName: 'Frost Mage',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceGnome,
			[Faction.Horde]: Race.RaceTroll10,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: MagePresets.P1_FROST_PRESET.gear,
				2: MagePresets.P2_FROST_PRESET.gear,
				3: MagePresets.P3_FROST_PRESET.gear,
				4: MagePresets.P4_FROST_PRESET.gear,
				5: MagePresets.P5_FROST_PRESET.gear,
			},
			[Faction.Horde]: {
				1: MagePresets.P1_FROST_PRESET.gear,
				2: MagePresets.P2_FROST_PRESET.gear,
				3: MagePresets.P3_FROST_PRESET.gear,
				4: MagePresets.P4_FROST_PRESET.gear,
				5: MagePresets.P5_FROST_PRESET.gear,
			},
		},
		tooltip: 'Frost Mage',
		iconUrl: talentTreeIcons[Class.ClassMage][2],
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
				1: RoguePresets.P1_PRESET.gear,
				2: RoguePresets.P2_PRESET.gear,
				3: RoguePresets.P3_PRESET.gear,
				4: RoguePresets.P4_PRESET.gear,
				5: RoguePresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.P1_PRESET.gear,
				2: RoguePresets.P2_PRESET.gear,
				3: RoguePresets.P3_PRESET.gear,
				4: RoguePresets.P4_PRESET.gear,
				5: RoguePresets.P5_PRESET.gear,
			},
		},
		tooltip: 'Combat Rogue',
		iconUrl: specIconsLarge[Spec.SpecRogue],
	},
	{
		spec: Spec.SpecRogue,
		rotation: RoguePresets.DefaultRotation,
		talents: RoguePresets.HemoTalents.data,
		specOptions: RoguePresets.DefaultOptions,
		consumes: RoguePresets.DefaultConsumes,
		defaultName: 'Hemo Rogue',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceOrc,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: RoguePresets.P1_PRESET.gear,
				2: RoguePresets.P2_PRESET.gear,
				3: RoguePresets.P3_PRESET.gear,
				4: RoguePresets.P4_PRESET.gear,
				5: RoguePresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: RoguePresets.P1_PRESET.gear,
				2: RoguePresets.P2_PRESET.gear,
				3: RoguePresets.P3_PRESET.gear,
				4: RoguePresets.P4_PRESET.gear,
				5: RoguePresets.P5_PRESET.gear,
			},
		},
		tooltip: 'Hemo Rogue',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
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
				2: ElementalShamanPresets.P2_PRESET.gear,
				3: ElementalShamanPresets.P3_PRESET.gear,
				4: ElementalShamanPresets.P4_PRESET.gear,
				5: ElementalShamanPresets.P5_ALLIANCE_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ElementalShamanPresets.P1_PRESET.gear,
				2: ElementalShamanPresets.P2_PRESET.gear,
				3: ElementalShamanPresets.P3_PRESET.gear,
				4: ElementalShamanPresets.P4_PRESET.gear,
				5: ElementalShamanPresets.P5_HORDE_PRESET.gear,
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
				2: EnhancementShamanPresets.P2_PRESET.gear,
				3: EnhancementShamanPresets.P3_PRESET.gear,
				4: EnhancementShamanPresets.P4_PRESET.gear,
				5: EnhancementShamanPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: EnhancementShamanPresets.P1_PRESET.gear,
				2: EnhancementShamanPresets.P2_PRESET.gear,
				3: EnhancementShamanPresets.P3_PRESET.gear,
				4: EnhancementShamanPresets.P4_PRESET.gear,
				5: EnhancementShamanPresets.P5_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecEnhancementShaman],
		iconUrl: specIconsLarge[Spec.SpecEnhancementShaman],
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
				2: ShadowPriestPresets.P2_PRESET.gear,
				3: ShadowPriestPresets.P3_PRESET.gear,
				4: ShadowPriestPresets.P4_PRESET.gear,
				5: ShadowPriestPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ShadowPriestPresets.P1_PRESET.gear,
				2: ShadowPriestPresets.P2_PRESET.gear,
				3: ShadowPriestPresets.P3_PRESET.gear,
				4: ShadowPriestPresets.P4_PRESET.gear,
				5: ShadowPriestPresets.P5_PRESET.gear,
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
				2: SmitePriestPresets.P2_PRESET.gear,
				3: SmitePriestPresets.P3_PRESET.gear,
				4: SmitePriestPresets.P4_PRESET.gear,
				5: SmitePriestPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: SmitePriestPresets.P1_PRESET.gear,
				2: SmitePriestPresets.P2_PRESET.gear,
				3: SmitePriestPresets.P3_PRESET.gear,
				4: SmitePriestPresets.P4_PRESET.gear,
				5: SmitePriestPresets.P5_PRESET.gear,
			},
		},
		tooltip: specNames[Spec.SpecSmitePriest],
		iconUrl: specIconsLarge[Spec.SpecSmitePriest],
	},
	{
		spec: Spec.SpecWarrior,
		rotation: WarriorPresets.ArmsRotation,
		talents: WarriorPresets.ArmsSlamTalents.data,
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
				2: WarriorPresets.P2_ARMS_PRESET.gear,
				3: WarriorPresets.P3_ARMS_PRESET.gear,
				4: WarriorPresets.P4_ARMS_PRESET.gear,
				5: WarriorPresets.P5_ARMS_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarriorPresets.P1_ARMS_PRESET.gear,
				2: WarriorPresets.P2_ARMS_PRESET.gear,
				3: WarriorPresets.P3_ARMS_PRESET.gear,
				4: WarriorPresets.P4_ARMS_PRESET.gear,
				5: WarriorPresets.P5_ARMS_PRESET.gear,
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
				2: WarriorPresets.P2_FURY_PRESET.gear,
				3: WarriorPresets.P3_FURY_PRESET.gear,
				4: WarriorPresets.P4_FURY_PRESET.gear,
				5: WarriorPresets.P5_FURY_PRESET.gear,
			},
			[Faction.Horde]: {
				1: WarriorPresets.P1_FURY_PRESET.gear,
				2: WarriorPresets.P2_FURY_PRESET.gear,
				3: WarriorPresets.P3_FURY_PRESET.gear,
				4: WarriorPresets.P4_FURY_PRESET.gear,
				5: WarriorPresets.P5_FURY_PRESET.gear,
			},
		},
		tooltip: 'Fury Warrior',
		iconUrl: talentTreeIcons[Class.ClassWarrior][1],
	},
	{
		spec: Spec.SpecProtectionWarrior,
		rotation: ProtectionWarriorPresets.DefaultRotation,
		talents: ProtectionWarriorPresets.ImpaleProtTalents.data,
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
				2: ProtectionWarriorPresets.P2_BALANCED_PRESET.gear,
				3: ProtectionWarriorPresets.P3_BALANCED_PRESET.gear,
				4: ProtectionWarriorPresets.P4_BALANCED_PRESET.gear,
				5: ProtectionWarriorPresets.P5_BALANCED_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ProtectionWarriorPresets.P1_BALANCED_PRESET.gear,
				2: ProtectionWarriorPresets.P2_BALANCED_PRESET.gear,
				3: ProtectionWarriorPresets.P3_BALANCED_PRESET.gear,
				4: ProtectionWarriorPresets.P4_BALANCED_PRESET.gear,
				5: ProtectionWarriorPresets.P5_BALANCED_PRESET.gear,
			},
		},
		tooltip: 'Protection Warrior',
		iconUrl: talentTreeIcons[Class.ClassWarrior][2],
	},
	{
		spec: Spec.SpecRetributionPaladin,
		rotation: RetributionPaladinPresets.DefaultRotation,
		talents: RetributionPaladinPresets.RetKingsPaladinTalents.data,
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
		spec: Spec.SpecProtectionPaladin,
		rotation: ProtectionPaladinPresets.DefaultRotation,
		talents: ProtectionPaladinPresets.SanctityTalents.data,
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
				2: ProtectionPaladinPresets.P2_PRESET.gear,
				3: ProtectionPaladinPresets.P3_PRESET.gear,
				4: ProtectionPaladinPresets.P4_PRESET.gear,
				5: ProtectionPaladinPresets.P5_PRESET.gear,
			},
			[Faction.Horde]: {
				1: ProtectionPaladinPresets.P1_PRESET.gear,
				2: ProtectionPaladinPresets.P2_PRESET.gear,
				3: ProtectionPaladinPresets.P3_PRESET.gear,
				4: ProtectionPaladinPresets.P4_PRESET.gear,
				5: ProtectionPaladinPresets.P5_PRESET.gear,
			},
		},
		tooltip: 'Protection Paladin',
		iconUrl: talentTreeIcons[Class.ClassPaladin][1],
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.DefaultRotation,
		talents: WarlockPresets.DestructionTalents.data,
		specOptions: WarlockPresets.DefaultOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Destro Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_DESTRO.gear,
				2: WarlockPresets.P2_DESTRO.gear,
				3: WarlockPresets.P3_DESTRO.gear,
				4: WarlockPresets.P4_DESTRO.gear,
				5: WarlockPresets.P5_DESTRO.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_DESTRO.gear,
				2: WarlockPresets.P2_DESTRO.gear,
				3: WarlockPresets.P3_DESTRO.gear,
				4: WarlockPresets.P4_DESTRO.gear,
				5: WarlockPresets.P5_DESTRO.gear,
			},
		},
		tooltip: 'Destruction Warlock: defaults to casting Curse of Doom.',
		iconUrl: talentTreeIcons[Class.ClassWarlock][2],
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.AfflictionRotation,
		talents: WarlockPresets.AfflicationTalents.data,
		specOptions: WarlockPresets.AfflictionOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Aff Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_DESTRO.gear,
				2: WarlockPresets.P2_DESTRO.gear,
				3: WarlockPresets.P3_DESTRO.gear,
				4: WarlockPresets.P4_DESTRO.gear,
				5: WarlockPresets.P5_DESTRO.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_DESTRO.gear,
				2: WarlockPresets.P2_DESTRO.gear,
				3: WarlockPresets.P3_DESTRO.gear,
				4: WarlockPresets.P4_DESTRO.gear,
				5: WarlockPresets.P5_DESTRO.gear,
			},
		},
		tooltip: 'Affliction Warlock: by default casts CoE with Malediction',
		iconUrl: talentTreeIcons[Class.ClassWarlock][0],
	},
	{
		spec: Spec.SpecWarlock,
		rotation: WarlockPresets.DemonologyRotation,
		talents: WarlockPresets.DemonologistTalents.data,
		specOptions: WarlockPresets.DemonologyOptions,
		consumes: WarlockPresets.DefaultConsumes,
		defaultName: 'Demo Warlock',
		defaultFactionRaces: {
			[Faction.Unknown]: Race.RaceUnknown,
			[Faction.Alliance]: Race.RaceHuman,
			[Faction.Horde]: Race.RaceBloodElf,
		},
		defaultGear: {
			[Faction.Unknown]: {},
			[Faction.Alliance]: {
				1: WarlockPresets.P1_DESTRO.gear,
				2: WarlockPresets.P2_DESTRO.gear,
				3: WarlockPresets.P3_DESTRO.gear,
				4: WarlockPresets.P4_DESTRO.gear,
				5: WarlockPresets.P5_DESTRO.gear,
			},
			[Faction.Horde]: {
				1: WarlockPresets.P1_DESTRO.gear,
				2: WarlockPresets.P2_DESTRO.gear,
				3: WarlockPresets.P3_DESTRO.gear,
				4: WarlockPresets.P4_DESTRO.gear,
				5: WarlockPresets.P5_DESTRO.gear,
			},
		},
		tooltip: 'Demonology Warlock',
		iconUrl: talentTreeIcons[Class.ClassWarlock][1],
	},
];

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];

export const buffBotPresets: Array<BuffBotSettings> = [
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Bear',
		deprecated: true,
		spec: Spec.SpecBalanceDruid,
		name: 'Bear',
		tooltip: 'Bear: Adds Gift of the Wild, an Innervate, Faerie Fire, and Leader of the Pack.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_racial_bearform.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			raidProto.buffs!.giftOfTheWild = Math.max(raidProto.buffs!.giftOfTheWild, TristateEffect.TristateEffectRegular);
			raidProto.buffs!.thorns = Math.max(raidProto.buffs!.thorns, TristateEffect.TristateEffectRegular);
			raidProto.debuffs!.faerieFire = Math.max(raidProto.debuffs!.faerieFire, TristateEffect.TristateEffectRegular);
			partyProto.buffs!.leaderOfThePack = Math.max(partyProto.buffs!.leaderOfThePack, TristateEffect.TristateEffectRegular);

			const innervateIndex = buffBot.getInnervateAssignment().targetIndex;
			if (innervateIndex != NO_TARGET) {
				const partyIndex = Math.floor(innervateIndex / 5);
				const playerIndex = innervateIndex % 5;
				const playerProto = raidProto.parties[partyIndex].players[playerIndex];
				if (playerProto.buffs) {
					playerProto.buffs.innervates++;
				}
			}
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Resto Druid',
		spec: Spec.SpecBalanceDruid,
		name: 'Resto Druid',
		tooltip: 'Resto Druid: Adds Improved Gift of the Wild, and an Innervate.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingtouch.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			raidProto.buffs!.giftOfTheWild = TristateEffect.TristateEffectImproved;
			raidProto.buffs!.thorns = Math.max(raidProto.buffs!.thorns, TristateEffect.TristateEffectRegular);

			const innervateIndex = buffBot.getInnervateAssignment().targetIndex;
			if (innervateIndex != NO_TARGET) {
				const partyIndex = Math.floor(innervateIndex / 5);
				const playerIndex = innervateIndex % 5;
				const playerProto = raidProto.parties[partyIndex].players[playerIndex];
				if (playerProto.buffs) {
					playerProto.buffs.innervates++;
				}
			}
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Dreamstate',
		spec: Spec.SpecBalanceDruid,
		name: 'Dreamstate',
		tooltip: 'Dreamstate: Adds Improved Gift of the Wild, an Innervate, and Improved Faerie Fire.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_faeriefire.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			raidProto.buffs!.giftOfTheWild = TristateEffect.TristateEffectImproved;
			raidProto.buffs!.thorns = TristateEffect.TristateEffectImproved;
			raidProto.debuffs!.faerieFire = TristateEffect.TristateEffectImproved;

			const innervateIndex = buffBot.getInnervateAssignment().targetIndex;
			if (innervateIndex != NO_TARGET) {
				const partyIndex = Math.floor(innervateIndex / 5);
				const playerIndex = innervateIndex % 5;
				const playerProto = raidProto.parties[partyIndex].players[playerIndex];
				if (playerProto.buffs) {
					playerProto.buffs.innervates++;
				}
			}
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Mage',
		deprecated: true,
		spec: Spec.SpecMage,
		name: 'Mage',
		tooltip: 'Mage: Adds Arcane Brilliance.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_arcaneintellect.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			raidProto.buffs!.arcaneBrilliance = true;
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Paladin',
		spec: Spec.SpecRetributionPaladin,
		name: 'Holy Paladin',
		tooltip: 'Holy Paladin: Adds a set of blessings.',
		iconUrl: talentTreeIcons[Class.ClassPaladin][0],
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			// Do nothing, blessings are handled elswhere.
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'JoW Paladin',
		spec: Spec.SpecRetributionPaladin,
		name: 'JoW Paladin',
		tooltip: 'JoW Paladin: Adds a set of blessings and Judgement of Wisdom.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousnessaura.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			// Blessings are handled elswhere.
			raidProto.debuffs!.judgementOfWisdom = true;
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'JoC Paladin',
		spec: Spec.SpecRetributionPaladin,
		deprecated: true,
		name: 'JoC Paladin',
		tooltip: 'JoC Paladin: Adds a set of blessings and Improved Judgement of the Crusader (+3% crit).',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_holysmite.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			// Blessings are handled elswhere.
			raidProto.debuffs!.improvedSealOfTheCrusader = true;
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Holy Priest',
		spec: Spec.SpecShadowPriest,
		name: 'Holy Priest',
		tooltip: 'Holy Priest: Adds Improved PW Fortitude and Shadow Protection.',
		iconUrl: talentTreeIcons[Class.ClassPriest][1],
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			raidProto.buffs!.shadowProtection = true;
			raidProto.buffs!.powerWordFortitude = TristateEffect.TristateEffectImproved;
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Divine Spirit Priest',
		spec: Spec.SpecShadowPriest,
		name: 'Disc Priest',
		tooltip: 'Disc Priest: Adds Improved PW Fort, Shadow Protection, Improved Divine Spirit and a Power Infusion.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_powerinfusion.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			raidProto.buffs!.shadowProtection = true;
			raidProto.buffs!.powerWordFortitude = TristateEffect.TristateEffectImproved;
			raidProto.buffs!.divineSpirit = TristateEffect.TristateEffectImproved;

			const powerInfusionIndex = buffBot.getPowerInfusionAssignment().targetIndex;
			if (powerInfusionIndex != NO_TARGET) {
				const partyIndex = Math.floor(powerInfusionIndex / 5);
				const playerIndex = powerInfusionIndex % 5;
				const playerProto = raidProto.parties[partyIndex].players[playerIndex];
				if (playerProto.buffs) {
					playerProto.buffs.powerInfusions++;
				}
			}
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Resto Shaman',
		spec: Spec.SpecElementalShaman,
		name: 'Resto Shaman',
		tooltip: 'Resto Shaman: Adds Bloodlust, Mana Spring Totem, Mana Tide Totem, Strength of Earth Totem, and Drums of Battle. Chooses air totem based on party composition.',
		iconUrl: talentTreeIcons[Class.ClassShaman][2],
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			partyProto.buffs!.drums = Drums.DrumsOfBattle;
			partyProto.buffs!.bloodlust++;
			partyProto.buffs!.manaSpringTotem = TristateEffect.TristateEffectImproved;
			partyProto.buffs!.manaTideTotems++;

			// Choose which air totem to drop based on party composition.
			const woaSpecs = [
				Spec.SpecBalanceDruid,
				Spec.SpecMage,
				Spec.SpecShadowPriest,
				Spec.SpecSmitePriest,
				Spec.SpecProtectionPaladin,
				Spec.SpecEnhancementShaman,
				Spec.SpecElementalShaman,
				Spec.SpecWarlock,
			];
			const wfSpecs = [
				Spec.SpecRetributionPaladin,
				Spec.SpecRogue,
				Spec.SpecWarrior,
				Spec.SpecProtectionWarrior,
			];
			const goaSpecs = [
				Spec.SpecFeralDruid,
				Spec.SpecFeralTankDruid,
				Spec.SpecHunter,
			];
			const [woaVotes, wfVotes, goaVotes] = [woaSpecs, wfSpecs, goaSpecs]
				.map(specs => partyProto.players
					.filter(player => player.class != Class.ClassUnknown)
					.map(player => playerToSpec(player))
					.filter(playerSpec => specs.includes(playerSpec))
					.length);

			if (woaVotes >= wfVotes && woaVotes >= goaVotes) {
				partyProto.buffs!.wrathOfAirTotem = Math.max(partyProto.buffs!.wrathOfAirTotem, TristateEffect.TristateEffectRegular);
			} else if (wfVotes >= goaVotes) {
				partyProto.buffs!.windfuryTotemRank = 5;
			} else {
				partyProto.buffs!.graceOfAirTotem = Math.max(partyProto.buffs!.graceOfAirTotem, TristateEffect.TristateEffectRegular);
			}
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'CoE Warlock',
		spec: Spec.SpecWarlock,
		deprecated: true,
		name: 'CoE Warlock',
		tooltip: 'CoE Warlock: Adds Curse of Elements (regular). Also adds +20% uptime to ISB.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_chilltouch.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			const debuffs = raidProto.debuffs!;
			debuffs.curseOfElements = Math.max(debuffs.curseOfElements, TristateEffect.TristateEffectRegular);
			debuffs.isbUptime = Math.min(1.0, debuffs.isbUptime + 0.2);
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Malediction Warlock',
		spec: Spec.SpecWarlock,
		name: 'Aff Warlock',
		deprecated: true,
		tooltip: 'Afflication Warlock: Adds Curse of Elements (improved). Also adds +20% uptime to ISB.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_curseofachimonde.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			const debuffs = raidProto.debuffs!;
			debuffs.curseOfElements = TristateEffect.TristateEffectImproved;
			debuffs.isbUptime = Math.min(1.0, debuffs.isbUptime + 0.2);
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'CoR Warlock',
		spec: Spec.SpecWarlock,
		deprecated: true,
		name: 'CoR Warlock',
		tooltip: 'CoR Warlock: Adds Curse of Recklessness. Also adds +20% uptime to ISB.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unholystrength.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			const debuffs = raidProto.debuffs!;
			debuffs.curseOfRecklessness = true;
			debuffs.isbUptime = Math.min(1.0, debuffs.isbUptime + 0.2);
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Arms Warrior',
		deprecated: true,
		spec: Spec.SpecWarrior,
		name: 'Arms Warrior',
		tooltip: 'Arms Warrior: Adds Sunder Armor, Blood Frenzy, and Improved Battle Shout.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/ability_warrior_savageblow.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			partyProto.buffs!.battleShout = TristateEffect.TristateEffectImproved;
			const debuffs = raidProto.debuffs!;
			debuffs.sunderArmor = true;
			debuffs.bloodFrenzy = true;
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Fury Warrior',
		deprecated: true,
		spec: Spec.SpecWarrior,
		name: 'Fury Warrior',
		tooltip: 'Fury Warrior: Adds Sunder Armor and Improved Battle Shout.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/ability_warrior_innerrage.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			partyProto.buffs!.battleShout = TristateEffect.TristateEffectImproved;
			const debuffs = raidProto.debuffs!;
			debuffs.sunderArmor = true;
		},
	},
	{
		// The value of this field must never change, to preserve local storage data.
		buffBotId: 'Prot Warrior',
		deprecated: true,
		spec: Spec.SpecWarrior,
		name: 'Prot Warrior',
		tooltip: 'Prot Warrior: Adds Sunder Armor.',
		iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/inv_shield_06.jpg',
		modifyRaidProto: (buffBot: BuffBot, raidProto: RaidProto, partyProto: PartyProto) => {
			const debuffs = raidProto.debuffs!;
			debuffs.sunderArmor = true;
		},
	},
];
