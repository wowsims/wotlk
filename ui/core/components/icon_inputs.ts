import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Alchohol } from '/wotlk/core/proto/common.js';
import { BattleElixir } from '/wotlk/core/proto/common.js';
import { Explosive } from '/wotlk/core/proto/common.js';
import { Flask } from '/wotlk/core/proto/common.js';
import { Food } from '/wotlk/core/proto/common.js';
import { GuardianElixir } from '/wotlk/core/proto/common.js';
import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Conjured } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';

import { PetFood } from '/wotlk/core/proto/common.js';
import { Potions } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js';
import { WeaponImbue } from '/wotlk/core/proto/common.js';
import { IndividualSimIconPickerConfig } from '/wotlk/core/individual_sim_ui.js';
import { Party } from '/wotlk/core/party.js';
import { Player } from '/wotlk/core/player.js';
import { Raid } from '/wotlk/core/raid.js';
import { Sim } from '/wotlk/core/sim.js';
import { Target } from '/wotlk/core/target.js';
import { Encounter } from '/wotlk/core/encounter.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import { ExclusivityTag } from '/wotlk/core/individual_sim_ui.js';
import { IconPickerConfig } from './icon_picker.js';
import { IconEnumPicker, IconEnumPickerConfig, IconEnumValueConfig } from './icon_enum_picker.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_WeaveType as WeaveType,
} from '/wotlk/core/proto/hunter.js';

// Keep each section in alphabetical order.

// Raid Buffs
export const ArcaneBrilliance = makeBooleanRaidBuffInput(ActionId.fromSpellId(27127), 'arcaneBrilliance');
export const DivineSpirit = makeBooleanRaidBuffInput(ActionId.fromSpellId(48073), 'divineSpirit');
export const GiftOfTheWild = makeTristateRaidBuffInput(ActionId.fromSpellId(26991), ActionId.fromSpellId(17051), 'giftOfTheWild');
export const Thorns = makeTristateRaidBuffInput(ActionId.fromSpellId(26992), ActionId.fromSpellId(16840), 'thorns');
export const PowerWordFortitude = makeTristateRaidBuffInput(ActionId.fromSpellId(25389), ActionId.fromSpellId(14767), 'powerWordFortitude');
export const ShadowProtection = makeBooleanRaidBuffInput(ActionId.fromSpellId(39374), 'shadowProtection');
export const FerociousInspiration = makeBooleanRaidBuffInput(ActionId.fromSpellId(34460), 'ferociousInspiration');
export const Bloodlust = makeBooleanRaidBuffInput(ActionId.fromSpellId(2825), 'bloodlust');
export const CommandingShout = makeTristateRaidBuffInput(ActionId.fromSpellId(469), ActionId.fromSpellId(12861), 'commandingShout');
export const DevotionAura = makeTristateRaidBuffInput(ActionId.fromSpellId(27149), ActionId.fromSpellId(20142), 'devotionAura');
export const LeaderOfThePack = makeTristateRaidBuffInput(ActionId.fromSpellId(17007), ActionId.fromItemId(32387), 'leaderOfThePack');
export const ManaSpringTotem = makeTristateRaidBuffInput(ActionId.fromSpellId(25570), ActionId.fromSpellId(16206), 'manaSpringTotem');
export const MoonkinAura = makeTristateRaidBuffInput(ActionId.fromSpellId(24907), ActionId.fromSpellId(48396), 'moonkinAura');
export const RetributionAura = makeTristateRaidBuffInput(ActionId.fromSpellId(27150), ActionId.fromSpellId(20092), 'retributionAura');
export const TotemOfWrath = makeBooleanRaidBuffInput(ActionId.fromSpellId(30706), 'totemOfWrath');
export const TrueshotAura = makeBooleanRaidBuffInput(ActionId.fromSpellId(27066), 'trueshotAura');
export const WrathOfAirTotem = makeBooleanRaidBuffInput(ActionId.fromSpellId(3738), 'wrathOfAirTotem');
export const BloodPact = makeTristateRaidBuffInput(ActionId.fromSpellId(27268), ActionId.fromSpellId(18696), 'bloodPact');
export const UnleashedRage = makeBooleanRaidBuffInput(ActionId.fromSpellId(30811), 'unleashedRage');

// Party Buffs
export const AtieshMage = makeMultistatePartyBuffInput(ActionId.fromSpellId(28142), 5, 'atieshMage');
export const AtieshWarlock = makeMultistatePartyBuffInput(ActionId.fromSpellId(28143), 5, 'atieshWarlock');
export const BraidedEterniumChain = makeBooleanPartyBuffInput(ActionId.fromSpellId(31025), 'braidedEterniumChain');
export const ChainOfTheTwilightOwl = makeBooleanPartyBuffInput(ActionId.fromSpellId(31035), 'chainOfTheTwilightOwl');
export const HeroicPresence = makeBooleanPartyBuffInput(ActionId.fromSpellId(6562), 'heroicPresence');
export const EyeOfTheNight = makeBooleanPartyBuffInput(ActionId.fromSpellId(31033), 'eyeOfTheNight');
export const ManaTideTotem = makeMultistatePartyBuffInput(ActionId.fromSpellId(16190), 5, 'manaTideTotems');

// TODO: drum buff icons

// Individual Buffs
export const BlessingOfKings = makeBooleanIndividualBuffInput(ActionId.fromSpellId(25898), 'blessingOfKings');
export const BlessingOfMight = makeTristateIndividualBuffInput(ActionId.fromSpellId(27140), ActionId.fromSpellId(20048), 'blessingOfMight');
export const BlessingOfSanctuary = makeBooleanIndividualBuffInput(ActionId.fromSpellId(27169), 'blessingOfSanctuary');
export const BlessingOfWisdom = makeTristateIndividualBuffInput(ActionId.fromSpellId(27143), ActionId.fromSpellId(20245), 'blessingOfWisdom');
export const Innervate = makeMultistateIndividualBuffInput(ActionId.fromSpellId(29166), 11, 'innervates');
export const PowerInfusion = makeMultistateIndividualBuffInput(ActionId.fromSpellId(10060), 11, 'powerInfusions');
export const Replenishment = makeBooleanIndividualBuffInput(ActionId.fromSpellId(57669), 'replenishment');

// Debuffs
export const BloodFrenzy = makeBooleanDebuffInput(ActionId.fromSpellId(29859), 'bloodFrenzy');
export const HuntersMark = makeTristateDebuffInput(ActionId.fromSpellId(14325), ActionId.fromSpellId(19425), 'huntersMark');
export const ImprovedShadowbolt = makeBooleanDebuffInput(ActionId.fromSpellId(17803), 'improvedShadowBolt');
export const ImprovedScorch = makeBooleanDebuffInput(ActionId.fromSpellId(12873), 'improvedScorch');
export const WintersChill = makeBooleanDebuffInput(ActionId.fromSpellId(28593), 'wintersChill');
export const JudgementOfWisdom = makeBooleanDebuffInput(ActionId.fromSpellId(53408), 'judgementOfWisdom');
export const JudgementOfLight = makeBooleanDebuffInput(ActionId.fromSpellId(27163), 'judgementOfLight');
export const Mangle = makeBooleanDebuffInput(ActionId.fromSpellId(33876), 'mangle');
export const Misery = makeBooleanDebuffInput(ActionId.fromSpellId(33198), 'misery');
export const CurseOfElements = makeBooleanDebuffInput(ActionId.fromSpellId(47865), 'curseOfElements');
export const EbonPlagueBringer = makeBooleanDebuffInput(ActionId.fromSpellId(51161), 'ebonPlaguebringer');
export const EarthAndMoon = makeBooleanDebuffInput(ActionId.fromSpellId(48511), 'earthAndMoon');
export const CurseOfWeakness = makeBooleanDebuffInput(ActionId.fromSpellId(27226), 'curseOfWeakness');
export const FaerieFire = makeTristateDebuffInput(ActionId.fromSpellId(770), ActionId.fromSpellId(33602), 'faerieFire');
export const ExposeArmor = makeTristateDebuffInput(ActionId.fromSpellId(26866), ActionId.fromSpellId(14169), 'exposeArmor');
export const SunderArmor = makeBooleanDebuffInput(ActionId.fromSpellId(25225), 'sunderArmor');
export const GiftOfArthas = makeBooleanDebuffInput(ActionId.fromSpellId(11374), 'giftOfArthas');
export const DemoralizingRoar = makeTristateDebuffInput(ActionId.fromSpellId(26998), ActionId.fromSpellId(16862), 'demoralizingRoar');
export const DemoralizingShout = makeTristateDebuffInput(ActionId.fromSpellId(25203), ActionId.fromSpellId(12879), 'demoralizingShout');
export const Screech = makeBooleanDebuffInput(ActionId.fromSpellId(27051), 'screech');
export const ThunderClap = makeTristateDebuffInput(ActionId.fromSpellId(25264), ActionId.fromSpellId(12666), 'thunderClap');
export const InsectSwarm = makeBooleanDebuffInput(ActionId.fromSpellId(27013), 'insectSwarm');
export const ScorpidSting = makeBooleanDebuffInput(ActionId.fromSpellId(3043), 'scorpidSting');

// Consumes
export const SuperSapper = makeBooleanConsumeInput(ActionId.fromItemId(23827), 'superSapper', [], onSetExplosives);
export const GoblinSapper = makeBooleanConsumeInput(ActionId.fromItemId(10646), 'goblinSapper', [], onSetExplosives);

export const KiblersBits = makeEnumValueConsumeInput(ActionId.fromItemId(33874), 'petFood', PetFood.PetFoodKiblersBits, ['Pet Food']);

export const ScrollOfAgilityV = makeEnumValueConsumeInput(ActionId.fromItemId(27498), 'scrollOfAgility', 5);
export const ScrollOfSpiritV = makeEnumValueConsumeInput(ActionId.fromItemId(27501), 'scrollOfSpirit', 5, ['Spirit']);
export const ScrollOfStrengthV = makeEnumValueConsumeInput(ActionId.fromItemId(27503), 'scrollOfStrength', 5);
export const ScrollOfProtectionV = makeEnumValueConsumeInput(ActionId.fromItemId(27500), 'scrollOfProtection', 5);

export const PetScrollOfAgilityV = makeEnumValueConsumeInput(ActionId.fromItemId(27498), 'petScrollOfAgility', 5);
export const PetScrollOfStrengthV = makeEnumValueConsumeInput(ActionId.fromItemId(27503), 'petScrollOfStrength', 5);

function makeBooleanRaidBuffInput(id: ActionId, buffsFieldName: keyof RaidBuffs, exclusivityTags?: Array<ExclusivityTag>): IndividualSimIconPickerConfig<Raid, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (raid: Raid) => raid.buffsChangeEmitter,
		getValue: (raid: Raid) => raid.getBuffs()[buffsFieldName] as boolean,
		setValue: (eventID: EventID, raid: Raid, newValue: boolean) => {
			const newBuffs = raid.getBuffs();
			(newBuffs[buffsFieldName] as boolean) = newValue;
			raid.setBuffs(eventID, newBuffs);
		},
	}
}

function makeTristateRaidBuffInput(id: ActionId, impId: ActionId, buffsFieldName: keyof RaidBuffs, exclusivityTags?: Array<ExclusivityTag>): IndividualSimIconPickerConfig<Raid, number> {
	return {
		id: id,
		states: 3,
		improvedId: impId,
		exclusivityTags: exclusivityTags,
		changedEvent: (raid: Raid) => raid.buffsChangeEmitter,
		getValue: (raid: Raid) => raid.getBuffs()[buffsFieldName] as number,
		setValue: (eventID: EventID, raid: Raid, newValue: number) => {
			const newBuffs = raid.getBuffs();
			(newBuffs[buffsFieldName] as number) = newValue;
			raid.setBuffs(eventID, newBuffs);
		},
	}
}

function makeBooleanPartyBuffInput(id: ActionId, buffsFieldName: keyof PartyBuffs, exclusivityTags?: Array<ExclusivityTag>): IndividualSimIconPickerConfig<Party, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs()[buffsFieldName] as boolean,
		setValue: (eventID: EventID, party: Party, newValue: boolean) => {
			const newBuffs = party.getBuffs();
			(newBuffs[buffsFieldName] as boolean) = newValue;
			party.setBuffs(eventID, newBuffs);
		},
	}
}

function makeTristatePartyBuffInput(id: ActionId, impId: ActionId, buffsFieldName: keyof PartyBuffs): IndividualSimIconPickerConfig<Party, number> {
	return {
		id: id,
		states: 3,
		improvedId: impId,
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs()[buffsFieldName] as number,
		setValue: (eventID: EventID, party: Party, newValue: number) => {
			const newBuffs = party.getBuffs();
			(newBuffs[buffsFieldName] as number) = newValue;
			party.setBuffs(eventID, newBuffs);
		},
	}
}

function makeMultistatePartyBuffInput(id: ActionId, numStates: number, buffsFieldName: keyof PartyBuffs): IndividualSimIconPickerConfig<Party, number> {
	return {
		id: id,
		states: numStates,
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs()[buffsFieldName] as number,
		setValue: (eventID: EventID, party: Party, newValue: number) => {
			const newBuffs = party.getBuffs();
			(newBuffs[buffsFieldName] as number) = newValue;
			party.setBuffs(eventID, newBuffs);
		},
	}
}

function makeEnumValuePartyBuffInput(id: ActionId, buffsFieldName: keyof PartyBuffs, enumValue: number, exclusivityTags?: Array<ExclusivityTag>): IndividualSimIconPickerConfig<Party, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs()[buffsFieldName] == enumValue,
		setValue: (eventID: EventID, party: Party, newValue: boolean) => {
			const newBuffs = party.getBuffs();
			(newBuffs[buffsFieldName] as number) = newValue ? enumValue : 0;
			party.setBuffs(eventID, newBuffs);
		},
	}
}

function makeBooleanIndividualBuffInput(id: ActionId, buffsFieldName: keyof IndividualBuffs, exclusivityTags?: Array<ExclusivityTag>): IndividualSimIconPickerConfig<Player<any>, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (player: Player<any>) => player.buffsChangeEmitter,
		getValue: (player: Player<any>) => player.getBuffs()[buffsFieldName] as boolean,
		setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
			const newBuffs = player.getBuffs();
			(newBuffs[buffsFieldName] as boolean) = newValue;
			player.setBuffs(eventID, newBuffs);
		},
	}
}

function makeTristateIndividualBuffInput(id: ActionId, impId: ActionId, buffsFieldName: keyof IndividualBuffs): IndividualSimIconPickerConfig<Player<any>, number> {
	return {
		id: id,
		states: 3,
		improvedId: impId,
		changedEvent: (player: Player<any>) => player.buffsChangeEmitter,
		getValue: (player: Player<any>) => player.getBuffs()[buffsFieldName] as number,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newBuffs = player.getBuffs();
			(newBuffs[buffsFieldName] as number) = newValue;
			player.setBuffs(eventID, newBuffs);
		},
	}
}

function makeMultistateIndividualBuffInput(id: ActionId, numStates: number, buffsFieldName: keyof IndividualBuffs): IndividualSimIconPickerConfig<Player<any>, number> {
	return {
		id: id,
		states: numStates,
		changedEvent: (player: Player<any>) => player.buffsChangeEmitter,
		getValue: (player: Player<any>) => player.getBuffs()[buffsFieldName] as number,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newBuffs = player.getBuffs();
			(newBuffs[buffsFieldName] as number) = newValue;
			player.setBuffs(eventID, newBuffs);
		},
	}
}

function makeBooleanDebuffInput(id: ActionId, debuffsFieldName: keyof Debuffs, exclusivityTags?: Array<ExclusivityTag>): IndividualSimIconPickerConfig<Raid, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (raid: Raid) => raid.debuffsChangeEmitter,
		getValue: (raid: Raid) => raid.getDebuffs()[debuffsFieldName] as boolean,
		setValue: (eventID: EventID, raid: Raid, newValue: boolean) => {
			const newDebuffs = raid.getDebuffs();
			(newDebuffs[debuffsFieldName] as boolean) = newValue;
			raid.setDebuffs(eventID, newDebuffs);
		},
	}
}

function makeTristateDebuffInput(id: ActionId, impId: ActionId, debuffsFieldName: keyof Debuffs): IndividualSimIconPickerConfig<Raid, number> {
	return {
		id: id,
		states: 3,
		improvedId: impId,
		changedEvent: (raid: Raid) => raid.debuffsChangeEmitter,
		getValue: (raid: Raid) => raid.getDebuffs()[debuffsFieldName] as number,
		setValue: (eventID: EventID, raid: Raid, newValue: number) => {
			const newDebuffs = raid.getDebuffs();
			(newDebuffs[debuffsFieldName] as number) = newValue;
			raid.setDebuffs(eventID, newDebuffs);
		},
	}
}

function makeBooleanConsumeInput(id: ActionId, consumesFieldName: keyof Consumes, exclusivityTags?: Array<ExclusivityTag>, onSet?: (eventID: EventID, player: Player<any>, newValue: boolean) => void): IndividualSimIconPickerConfig<Player<any>, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes()[consumesFieldName] as boolean,
		setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
			const newConsumes = player.getConsumes();
			(newConsumes[consumesFieldName] as boolean) = newValue;
			TypedEvent.freezeAllAndDo(() => {
				player.setConsumes(eventID, newConsumes);
				if (onSet) {
					onSet(eventID, player, newValue);
				}
			});
		},
	}
}

function makeEnumValueConsumeInput(id: ActionId, consumesFieldName: keyof Consumes, enumValue: number, exclusivityTags?: Array<ExclusivityTag>, onSet?: (eventID: EventID, player: Player<any>, newValue: boolean) => void, showWhen?: (player: Player<any>) => boolean): IndividualSimIconPickerConfig<Player<any>, boolean> {
	return {
		id: id,
		states: 2,
		exclusivityTags: exclusivityTags,
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes()[consumesFieldName] == enumValue,
		setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
			const newConsumes = player.getConsumes();
			(newConsumes[consumesFieldName] as number) = newValue ? enumValue : 0;
			TypedEvent.freezeAllAndDo(() => {
				player.setConsumes(eventID, newConsumes);
				if (onSet) {
					onSet(eventID, player, newValue);
				}
			});
		},
		showWhen: showWhen,
	}
}

//////////////////////////////////////////////////////////////////////
// Custom buffs that don't fit into any of the helper functions above.
//////////////////////////////////////////////////////////////////////

export const StrengthOfEarthTotem = {
	id: ActionId.fromSpellId(25528),
	states: 3,
	improvedId: ActionId.fromSpellId(52456),
	changedEvent: (raid: Raid) => raid.buffsChangeEmitter,
	getValue: (raid: Raid) => raid.getBuffs().strengthOfEarthTotem,
	setValue: (eventID: EventID, raid: Raid, newValue: number) => {
		const newBuffs = raid.getBuffs();
		newBuffs.strengthOfEarthTotem = newValue;
		raid.setBuffs(eventID, newBuffs);
	},
};

export const WindfuryTotem = {
	id: ActionId.fromSpellId(25587),
	states: 3,
	improvedId: ActionId.fromSpellId(29193),
	changedEvent: (raid: Raid) => raid.buffsChangeEmitter,
	getValue: (raid: Raid) => raid.getBuffs().windfuryTotem,
	setValue: (eventID: EventID, raid: Raid, newValue: number) => {
		const newBuffs = raid.getBuffs();
		newBuffs.windfuryTotem = newValue;
		raid.setBuffs(eventID, newBuffs);
	},
};

export const BattleShout = {
	id: ActionId.fromSpellId(2048),
	states: 3,
	improvedId: ActionId.fromSpellId(12861),
	improvedId2: ActionId.fromItemId(30446),
	changedEvent: (raid: Raid) => raid.buffsChangeEmitter,
	getValue: (raid: Raid) => raid.getBuffs().battleShout,
	setValue: (eventID: EventID, raid: Raid, newValue: number) => {
		const newBuffs = raid.getBuffs();
		newBuffs.battleShout = newValue;
		raid.setBuffs(eventID, newBuffs);
	},
};

export const makePotionsInput = makeConsumeInputFactory('defaultPotion', [
	{ actionId: ActionId.fromItemId(33447), value: Potions.RunicHealingPotion },
	{ actionId: ActionId.fromItemId(33448), value: Potions.RunicManaPotion },
	{ actionId: ActionId.fromItemId(40093), value: Potions.IndestructiblePotion },
	{ actionId: ActionId.fromItemId(40211), value: Potions.PotionOfSpeed },
	{ actionId: ActionId.fromItemId(40212), value: Potions.PotionOfWildMagic },

	{ actionId: ActionId.fromItemId(22839), value: Potions.DestructionPotion },
	{ actionId: ActionId.fromItemId(22838), value: Potions.HastePotion },
	{ actionId: ActionId.fromItemId(13442), value: Potions.MightyRagePotion },
	{ actionId: ActionId.fromItemId(22832), value: Potions.SuperManaPotion },
	{ actionId: ActionId.fromItemId(31677), value: Potions.FelManaPotion },
	{ actionId: ActionId.fromItemId(22828), value: Potions.InsaneStrengthPotion },
	{ actionId: ActionId.fromItemId(22849), value: Potions.IronshieldPotion },
	{ actionId: ActionId.fromItemId(22837), value: Potions.HeroicPotion },
] as Array<IconEnumValueConfig<Player<any>, Potions>>);

export const makeConjuredInput = makeConsumeInputFactory('defaultConjured', [
	{ actionId: ActionId.fromItemId(12662), value: Conjured.ConjuredDarkRune },
	{ actionId: ActionId.fromItemId(22788), value: Conjured.ConjuredFlameCap },
	{ actionId: ActionId.fromItemId(22105), value: Conjured.ConjuredHealthstone },
	{ actionId: ActionId.fromItemId(22044), value: Conjured.ConjuredMageManaEmerald },
	{ actionId: ActionId.fromItemId(7676), value: Conjured.ConjuredRogueThistleTea },
] as Array<IconEnumValueConfig<Player<any>, Conjured>>);

export const makeFlasksInput = makeConsumeInputFactory('flask', [
	{ actionId: ActionId.fromItemId(46376), value: Flask.FlaskOfTheFrostWyrm },
	{ actionId: ActionId.fromItemId(46377), value: Flask.FlaskOfEndlessRage },
	{ actionId: ActionId.fromItemId(46378), value: Flask.FlaskOfPureMojo },
	{ actionId: ActionId.fromItemId(46379), value: Flask.FlaskOfStoneblood },
	{ actionId: ActionId.fromItemId(40079), value: Flask.LesserFlaskOfToughness },
	{ actionId: ActionId.fromItemId(44939), value: Flask.LesserFlaskOfResistance },
	{ actionId: ActionId.fromItemId(22861), value: Flask.FlaskOfBlindingLight },
	{ actionId: ActionId.fromItemId(22853), value: Flask.FlaskOfMightyRestoration },
	{ actionId: ActionId.fromItemId(22866), value: Flask.FlaskOfPureDeath },
	{ actionId: ActionId.fromItemId(22854), value: Flask.FlaskOfRelentlessAssault },
	{ actionId: ActionId.fromItemId(13512), value: Flask.FlaskOfSupremePower },
	{ actionId: ActionId.fromItemId(22851), value: Flask.FlaskOfFortification },
	{ actionId: ActionId.fromItemId(33208), value: Flask.FlaskOfChromaticWonder },
] as Array<IconEnumValueConfig<Player<any>, Flask>>, (eventID: EventID, player: Player<any>, newValue: Flask) => {
	if (newValue) {
		const newConsumes = player.getConsumes();
		newConsumes.battleElixir = BattleElixir.BattleElixirUnknown;
		newConsumes.guardianElixir = GuardianElixir.GuardianElixirUnknown;
		player.setConsumes(eventID, newConsumes);
	}
});

export const makeBattleElixirsInput = makeConsumeInputFactory('battleElixir', [
	{ actionId: ActionId.fromItemId(44325), value: BattleElixir.ElixirOfAccuracy },
	{ actionId: ActionId.fromItemId(44330), value: BattleElixir.ElixirOfArmorPiercing },
	{ actionId: ActionId.fromItemId(44327), value: BattleElixir.ElixirOfDeadlyStrikes },
	{ actionId: ActionId.fromItemId(44329), value: BattleElixir.ElixirOfExpertise },
	{ actionId: ActionId.fromItemId(44331), value: BattleElixir.ElixirOfLightningSpeed },
	{ actionId: ActionId.fromItemId(39666), value: BattleElixir.ElixirOfMightyAgility },
	{ actionId: ActionId.fromItemId(40073), value: BattleElixir.ElixirOfMightyStrength },
	{ actionId: ActionId.fromItemId(40076), value: BattleElixir.GurusElixir },
	{ actionId: ActionId.fromItemId(40070), value: BattleElixir.SpellpowerElixir },
	{ actionId: ActionId.fromItemId(40068), value: BattleElixir.WrathElixir },
	{ actionId: ActionId.fromItemId(28103), value: BattleElixir.AdeptsElixir },
	{ actionId: ActionId.fromItemId(9224), value: BattleElixir.ElixirOfDemonslaying },
	{ actionId: ActionId.fromItemId(22831), value: BattleElixir.ElixirOfMajorAgility },
	{ actionId: ActionId.fromItemId(22833), value: BattleElixir.ElixirOfMajorFirePower },
	{ actionId: ActionId.fromItemId(22827), value: BattleElixir.ElixirOfMajorFrostPower },
	{ actionId: ActionId.fromItemId(22835), value: BattleElixir.ElixirOfMajorShadowPower },
	{ actionId: ActionId.fromItemId(22824), value: BattleElixir.ElixirOfMajorStrength },
	{ actionId: ActionId.fromItemId(28104), value: BattleElixir.ElixirOfMastery },
	{ actionId: ActionId.fromItemId(13452), value: BattleElixir.ElixirOfTheMongoose },
	{ actionId: ActionId.fromItemId(31679), value: BattleElixir.FelStrengthElixir },
	{ actionId: ActionId.fromItemId(27155), value: BattleElixir.GreaterArcaneElixir },
] as Array<IconEnumValueConfig<Player<any>, BattleElixir>>, (eventID: EventID, player: Player<any>, newValue: BattleElixir) => {
	if (newValue) {
		const newConsumes = player.getConsumes();
		newConsumes.flask = Flask.FlaskUnknown;
		player.setConsumes(eventID, newConsumes);
	}
});

export const makeGuardianElixirsInput = makeConsumeInputFactory('guardianElixir', [
	{ actionId: ActionId.fromItemId(44328), value: GuardianElixir.ElixirOfMightyDefense },
	{ actionId: ActionId.fromItemId(40078), value: GuardianElixir.ElixirOfMightyFortitude },
	{ actionId: ActionId.fromItemId(40109), value: GuardianElixir.ElixirOfMightyMageblood },
	{ actionId: ActionId.fromItemId(44332), value: GuardianElixir.ElixirOfMightyThoughts },
	{ actionId: ActionId.fromItemId(40097), value: GuardianElixir.ElixirOfProtection },
	{ actionId: ActionId.fromItemId(40072), value: GuardianElixir.ElixirOfSpirit },
	{ actionId: ActionId.fromItemId(9088), value: GuardianElixir.GiftOfArthas },
	{ actionId: ActionId.fromItemId(32067), value: GuardianElixir.ElixirOfDraenicWisdom },
	{ actionId: ActionId.fromItemId(32068), value: GuardianElixir.ElixirOfIronskin },
	{ actionId: ActionId.fromItemId(22834), value: GuardianElixir.ElixirOfMajorDefense },
	{ actionId: ActionId.fromItemId(32062), value: GuardianElixir.ElixirOfMajorFortitude },
	{ actionId: ActionId.fromItemId(22840), value: GuardianElixir.ElixirOfMajorMageblood },
] as Array<IconEnumValueConfig<Player<any>, GuardianElixir>>, (eventID: EventID, player: Player<any>, newValue: GuardianElixir) => {
	if (newValue) {
		const newConsumes = player.getConsumes();
		newConsumes.flask = Flask.FlaskUnknown;
		player.setConsumes(eventID, newConsumes);
	}
});

export const makeFoodInput = makeConsumeInputFactory('food', [
	{ actionId: ActionId.fromItemId(43015), value: Food.FoodFishFeast },
	{ actionId: ActionId.fromItemId(34753), value: Food.FoodGreatFeast },
	{ actionId: ActionId.fromItemId(42999), value: Food.FoodBlackenedDragonfin },
	{ actionId: ActionId.fromItemId(42995), value: Food.FoodHeartyRhino },
	{ actionId: ActionId.fromItemId(34754), value: Food.FoodMegaMammothMeal },
	{ actionId: ActionId.fromItemId(34756), value: Food.FoodSpicedWormBurger },
	{ actionId: ActionId.fromItemId(42994), value: Food.FoodRhinoliciousWormsteak },
	{ actionId: ActionId.fromItemId(34769), value: Food.FoodImperialMantaSteak },
	{ actionId: ActionId.fromItemId(42996), value: Food.FoodSnapperExtreme },
	{ actionId: ActionId.fromItemId(34758), value: Food.FoodMightyRhinoDogs },
	{ actionId: ActionId.fromItemId(34767), value: Food.FoodFirecrackerSalmon },
	{ actionId: ActionId.fromItemId(42998), value: Food.FoodCuttlesteak },
	{ actionId: ActionId.fromItemId(43000), value: Food.FoodDragonfinFilet },

	{ actionId: ActionId.fromItemId(27657), value: Food.FoodBlackenedBasilisk },
	{ actionId: ActionId.fromItemId(27664), value: Food.FoodGrilledMudfish },
	{ actionId: ActionId.fromItemId(27655), value: Food.FoodRavagerDog },
	{ actionId: ActionId.fromItemId(27658), value: Food.FoodRoastedClefthoof },
	{ actionId: ActionId.fromItemId(33872), value: Food.FoodSpicyHotTalbuk },
	{ actionId: ActionId.fromItemId(33825), value: Food.FoodSkullfishSoup },
	{ actionId: ActionId.fromItemId(33052), value: Food.FoodFishermansFeast },
] as Array<IconEnumValueConfig<Player<any>, Food>>);

export const makeAlcoholInput = makeConsumeInputFactory('alchohol', [
	{ actionId: ActionId.fromItemId(18284), value: Alchohol.AlchoholKreegsStoutBeatdown },
] as Array<IconEnumValueConfig<Player<any>, Alchohol>>);

export const makePetFoodInput = makeConsumeInputFactory('petFood', [
	{ actionId: ActionId.fromItemId(33874), value: PetFood.PetFoodKiblersBits },
] as Array<IconEnumValueConfig<Player<any>, PetFood>>);

function onSetExplosives(eventID: EventID, player: Player<any>, newValue: Explosive | boolean) {
	if (newValue) {
		const playerConsumes = player.getConsumes();
		player.setConsumes(eventID, playerConsumes);
	}
};

export const FillerExplosiveInput = makeConsumeInput('fillerExplosive', [
	{ actionId: ActionId.fromItemId(23736), value: Explosive.ExplosiveFelIronBomb },
	{ actionId: ActionId.fromItemId(23737), value: Explosive.ExplosiveAdamantiteGrenade },
	{ actionId: ActionId.fromItemId(23841), value: Explosive.ExplosiveGnomishFlameTurret },
	{ actionId: ActionId.fromItemId(13180), value: Explosive.ExplosiveHolyWater },
] as Array<IconEnumValueConfig<Player<any>, Explosive>>, onSetExplosives);

export function makeWeaponImbueInput(isMainHand: boolean, options: Array<WeaponImbue>): IconEnumPickerConfig<Player<any>, WeaponImbue> {
	const allOptions = [
		{ actionId: ActionId.fromItemId(18262), value: WeaponImbue.WeaponImbueElementalSharpeningStone },
		{ actionId: ActionId.fromItemId(20749), value: WeaponImbue.WeaponImbueBrilliantWizardOil },
		{ actionId: ActionId.fromItemId(22522), value: WeaponImbue.WeaponImbueSuperiorWizardOil },
		{ actionId: ActionId.fromItemId(34539), value: WeaponImbue.WeaponImbueRighteousWeaponCoating },
		{
			actionId: ActionId.fromItemId(23529), value: WeaponImbue.WeaponImbueAdamantiteSharpeningStone,
			showWhen: (player: Player<any>) => !(isMainHand ? player.getGear().hasBluntMHWeapon() : player.getGear().hasBluntOHWeapon()),
		},
		{
			actionId: ActionId.fromItemId(28421), value: WeaponImbue.WeaponImbueAdamantiteWeightstone,
			showWhen: (player: Player<any>) => (isMainHand ? player.getGear().hasBluntMHWeapon() : player.getGear().hasBluntOHWeapon()),
		},
		{ actionId: ActionId.fromSpellId(27186), value: WeaponImbue.WeaponImbueRogueDeadlyPoison },
		{ actionId: ActionId.fromSpellId(26891), value: WeaponImbue.WeaponImbueRogueInstantPoison },
		{ actionId: ActionId.fromSpellId(25505), value: WeaponImbue.WeaponImbueShamanWindfury },
		{ actionId: ActionId.fromSpellId(58790), value: WeaponImbue.WeaponImbueShamanFlametongue },
		{ actionId: ActionId.fromSpellId(25500), value: WeaponImbue.WeaponImbueShamanFrostbrand },
		{ actionId: ActionId.fromSpellId(10399), value: WeaponImbue.WeaponImbueShamanRockbiter },
	];
	if (isMainHand) {
		const config = makeConsumeInputFactory('mainHandImbue', allOptions)(options);
		config.enableWhen = (player: Player<any>) => !player.getParty()
			|| (player.spec == Spec.SpecHunter && (player.getRotation() as HunterRotation).weave == WeaveType.WeaveNone);
		config.changedEvent = (player: Player<any>) => TypedEvent.onAny([player.getRaid()?.changeEmitter || player.consumesChangeEmitter]);
		return config;
	} else {
		return makeConsumeInputFactory('offHandImbue', allOptions)(options);
	}
}

function makeConsumeInputFactory<T extends number>(consumesFieldName: keyof Consumes, allOptions: Array<IconEnumValueConfig<Player<any>, T>>, onSet?: (eventID: EventID, player: Player<any>, newValue: T) => void): (options: Array<T>) => IconEnumPickerConfig<Player<any>, T> {
	return (options: Array<T>) => {
		return {
			numColumns: 1,
			values: [
				{ color: 'grey', value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>,
			].concat(options.map(option => allOptions.find(allOption => allOption.value == option)!)),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
			getValue: (player: Player<any>) => player.getConsumes()[consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();
				(newConsumes[consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (onSet) {
						onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}

function makeConsumeInput<T extends number>(consumesFieldName: keyof Consumes, allOptions: Array<IconEnumValueConfig<Player<any>, T>>, onSet?: (eventID: EventID, player: Player<any>, newValue: T) => void): IconEnumPickerConfig<Player<any>, T> {
	const factory = makeConsumeInputFactory(consumesFieldName, allOptions, onSet);
	return factory(allOptions.map(option => option.value));
}
