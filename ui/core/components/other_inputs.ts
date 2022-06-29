import { BooleanPicker } from '/tbc/core/components/boolean_picker.js';
import { EnumPicker, EnumPickerConfig } from '/tbc/core/components/enum_picker.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { StrengthOfEarthType } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js';
import { Party } from '/tbc/core/party.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { Target } from '/tbc/core/target.js';
import { Encounter } from '/tbc/core/encounter.js';
import { Raid } from '/tbc/core/raid.js';
import { SimUI } from '/tbc/core/sim_ui.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { emptyRaidTarget } from '/tbc/core/proto_utils/utils.js';

export function makeShow1hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: [
			'show-1h-weapons-selector',
		],
		label: '1H',
		changedEvent: (sim: Sim) => sim.show1hWeaponsChangeEmitter,
		getValue: (sim: Sim) => sim.getShow1hWeapons(),
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			sim.setShow1hWeapons(eventID, newValue);
		},
	});
}

export function makeShow2hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: [
			'show-2h-weapons-selector',
		],
		label: '2H',
		changedEvent: (sim: Sim) => sim.show2hWeaponsChangeEmitter,
		getValue: (sim: Sim) => sim.getShow2hWeapons(),
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			sim.setShow2hWeapons(eventID, newValue);
		},
	});
}

export function makeShowMatchingGemsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: [
			'show-matching-gems-selector',
		],
		label: 'Match Socket',
		changedEvent: (sim: Sim) => sim.showMatchingGemsChangeEmitter,
		getValue: (sim: Sim) => sim.getShowMatchingGems(),
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			sim.setShowMatchingGems(eventID, newValue);
		},
	});
}

export function makePhaseSelector(parent: HTMLElement, sim: Sim): EnumPicker<Sim> {
	return new EnumPicker<Sim>(parent, sim, {
		extraCssClasses: [
			'phase-selector',
		],
		values: [
			{ name: 'Phase 1', value: 1 },
			{ name: 'Phase 2', value: 2 },
			{ name: 'Phase 3', value: 3 },
			{ name: 'Phase 4', value: 4 },
			{ name: 'Phase 5', value: 5 },
		],
		changedEvent: (sim: Sim) => sim.phaseChangeEmitter,
		getValue: (sim: Sim) => sim.getPhase(),
		setValue: (eventID: EventID, sim: Sim, newValue: number) => {
			sim.setPhase(eventID, newValue);
		},
	});
}

export const StartingPotion = {
	type: 'enum' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'starting-potion-picker',
		],
		label: 'Starting Potion',
		labelTooltip: 'If set, this potion will be used instead of the default potion for the first few uses.',
		values: [
			{ name: 'None', value: Potions.UnknownPotion },
			{ name: 'Destruction', value: Potions.DestructionPotion },
			{ name: 'Haste', value: Potions.HastePotion },
			{ name: 'Super Mana', value: Potions.SuperManaPotion },
			{ name: 'Fel Mana', value: Potions.FelManaPotion },
		],
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes().startingPotion,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newConsumes = player.getConsumes();
			newConsumes.startingPotion = newValue;
			player.setConsumes(eventID, newConsumes);
		},
	},
};

export const NumStartingPotions = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'num-starting-potions-picker',
		],
		label: '# to use',
		labelTooltip: 'The number of starting potions to use before going back to the default potion.',
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes().numStartingPotions,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newConsumes = player.getConsumes();
			newConsumes.numStartingPotions = newValue;
			player.setConsumes(eventID, newConsumes);
		},
		enableWhen: (player: Player<any>) => player.getConsumes().startingPotion != Potions.UnknownPotion,
	},
};

export const StartingConjured = {
	type: 'enum' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'starting-conjured-picker',
		],
		label: 'Starting Conjured',
		labelTooltip: 'If set, this conjured will be used instead of the default conjured for the first few uses.',
		values: [
			{ name: 'None', value: Conjured.ConjuredUnknown },
			{ name: 'Dark Rune', value: Conjured.ConjuredDarkRune },
			{ name: 'Flame Cap', value: Conjured.ConjuredFlameCap },
			{ name: 'Mana Gem', value: Conjured.ConjuredMageManaEmerald },
			{ name: 'Thistle Tea', value: Conjured.ConjuredRogueThistleTea },
		],
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes().startingConjured,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newConsumes = player.getConsumes();
			newConsumes.startingConjured = newValue;
			player.setConsumes(eventID, newConsumes);
		},
	},
};

export const NumStartingConjured = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'num-starting-conjureds-picker',
		],
		label: '# to use',
		labelTooltip: 'The number of starting conjured items to use before going back to the default conjured.',
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes().numStartingConjured,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newConsumes = player.getConsumes();
			newConsumes.numStartingConjured = newValue;
			player.setConsumes(eventID, newConsumes);
		},
		enableWhen: (player: Player<any>) => player.getConsumes().startingConjured != Conjured.ConjuredUnknown,
	},
};

export const ShadowPriestDPS = {
	type: 'number' as const,
	cssClass: 'shadow-priest-dps-picker',
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'shadow-priest-dps-picker',
			'within-raid-sim-hide',
		],
		label: 'Shadow Priest DPS',
		changedEvent: (player: Player<any>) => player.buffsChangeEmitter,
		getValue: (player: Player<any>) => player.getBuffs().shadowPriestDps,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const buffs = player.getBuffs();
			buffs.shadowPriestDps = newValue;
			player.setBuffs(eventID, buffs);
		},
	},
};

export const ISBUptime = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.sim.raid,
	config: {
		extraCssClasses: [
			'isb-uptime-picker',
			'within-raid-sim-hide',
		],
		label: 'Improved Shadowbolt Uptime %',
		labelTooltip: "Uptime for the Improved Shadowbolt debuff, applied by 1 or more warlocks in your raid.",
		changedEvent: (raid: Raid) => raid.debuffsChangeEmitter,
		getValue: (raid: Raid) => Math.round(raid.getDebuffs().isbUptime * 100),
		setValue: (eventID: EventID, raid: Raid, newValue: number) => {
			const newDebuffs = raid.getDebuffs();
			newDebuffs.isbUptime = newValue / 100;
			raid.setDebuffs(eventID, newDebuffs);
		},
	},
};

export const InspirationUptime = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'inspiration-uptime-picker',
		],
		label: 'Inspiration Uptime %',
		labelTooltip: "Uptime for the Inspiration or Ancestral Fortitude (+25% armor) buffs.",
		changedEvent: (player: Player<any>) => player.buffsChangeEmitter,
		getValue: (player: Player<any>) => Math.round(player.getBuffs().inspirationUptime * 100),
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newBuffs = player.getBuffs();
			newBuffs.inspirationUptime = newValue / 100;
			player.setBuffs(eventID, newBuffs);
		},
	},
};

export const ExposeWeaknessUptime = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.sim.raid,
	config: {
		extraCssClasses: [
			'expose-weakness-uptime-picker',
			'within-raid-sim-hide',
		],
		label: 'Expose Weakness Uptime %',
		labelTooltip: 'Uptime for the Expose Weakness debuff, applied by 1 or more Survival hunters in your raid.',
		changedEvent: (raid: Raid) => raid.debuffsChangeEmitter,
		getValue: (raid: Raid) => Math.round(raid.getDebuffs().exposeWeaknessUptime * 100),
		setValue: (eventID: EventID, raid: Raid, newValue: number) => {
			const newDebuffs = raid.getDebuffs();
			newDebuffs.exposeWeaknessUptime = newValue / 100;
			raid.setDebuffs(eventID, newDebuffs);
		},
	},
};

export const ExposeWeaknessHunterAgility = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.sim.raid,
	config: {
		extraCssClasses: [
			'expose-weakness-hunter-agility-picker',
			'within-raid-sim-hide',
		],
		label: 'EW Hunter Agility',
		labelTooltip: 'The amount of agility on the Expose Weakness hunter.',
		changedEvent: (raid: Raid) => raid.debuffsChangeEmitter,
		getValue: (raid: Raid) => Math.round(raid.getDebuffs().exposeWeaknessHunterAgility),
		setValue: (eventID: EventID, raid: Raid, newValue: number) => {
			const newDebuffs = raid.getDebuffs();
			newDebuffs.exposeWeaknessHunterAgility = newValue;
			raid.setDebuffs(eventID, newDebuffs);
		},
	},
};

export const SnapshotImprovedStrengthOfEarthTotem = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player.getParty()!,
	config: {
		extraCssClasses: [
			'snapshot-improved-strength-of-earth-totem-picker',
			'within-raid-sim-hide',
		],
		label: 'Snapshot Imp Strength of Earth',
		labelTooltip: 'An enhancement shaman in your party is snapshotting their improved Strength of Earth totem bonus from T4 2pc (+12 Strength) for the first 1:50s of the fight.',
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs().snapshotImprovedStrengthOfEarthTotem,
		setValue: (eventID: EventID, party: Party, newValue: boolean) => {
			const buffs = party.getBuffs();
			buffs.snapshotImprovedStrengthOfEarthTotem = newValue;
			party.setBuffs(eventID, buffs);
		},
		enableWhen: (party: Party) => party.getBuffs().strengthOfEarthTotem == StrengthOfEarthType.Basic || party.getBuffs().strengthOfEarthTotem == StrengthOfEarthType.EnhancingTotems,
	},
};

export const SnapshotImprovedWrathOfAirTotem = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player.getParty()!,
	config: {
		extraCssClasses: [
			'snapshot-improved-wrath-of-air-totem-picker',
			'within-raid-sim-hide',
		],
		label: 'Snapshot Imp Wrath of Air',
		labelTooltip: 'An elemental shaman in your party is snapshotting their improved wrath of air totem bonus from T4 2pc (+20 spell power) for the first 1:50s of the fight.',
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs().snapshotImprovedWrathOfAirTotem,
		setValue: (eventID: EventID, party: Party, newValue: boolean) => {
			const buffs = party.getBuffs();
			buffs.snapshotImprovedWrathOfAirTotem = newValue;
			party.setBuffs(eventID, buffs);
		},
		enableWhen: (party: Party) => party.getBuffs().wrathOfAirTotem == TristateEffect.TristateEffectRegular,
	},
};

export const SnapshotBsSolarianSapphire = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player.getParty()!,
	config: {
		extraCssClasses: [
			'snapshot-bs-solarian-sapphire-picker',
			'within-raid-sim-hide',
		],
		label: 'Snapshot BS Solarian\'s Sapphire',
		labelTooltip: 'A Warrior in your party is snapshotting their Battle Shout before combat, using the bonus from Solarian\'s Sapphire (+70 attack power) for the first 1:50s of the fight.',
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs().snapshotBsSolarianSapphire,
		setValue: (eventID: EventID, party: Party, newValue: boolean) => {
			const buffs = party.getBuffs();
			buffs.snapshotBsSolarianSapphire = newValue;
			party.setBuffs(eventID, buffs);
		},
		enableWhen: (party: Party) => party.getBuffs().battleShout > 0 && !party.getBuffs().bsSolarianSapphire,
	},
};

export const SnapshotBsT2 = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player.getParty()!,
	config: {
		extraCssClasses: [
			'snapshot-bs-t2-picker',
			'within-raid-sim-hide',
		],
		label: 'Snapshot BS T2',
		labelTooltip: 'A Warrior in your party is snapshotting their Battle Shout before combat, using the bonus from T2 3pc (+30 attack power) for the first 1:50s of the fight.',
		changedEvent: (party: Party) => party.buffsChangeEmitter,
		getValue: (party: Party) => party.getBuffs().snapshotBsT2,
		setValue: (eventID: EventID, party: Party, newValue: boolean) => {
			const buffs = party.getBuffs();
			buffs.snapshotBsT2 = newValue;
			party.setBuffs(eventID, buffs);
		},
		enableWhen: (party: Party) => party.getBuffs().battleShout > 0,
	},
};

export const InFrontOfTarget = {
	type: 'boolean' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'in-front-of-target-picker',
		],
		label: 'In Front of Target',
		labelTooltip: 'Stand in front of the target, causing Blocks and Parries to be included in the attack table.',
		changedEvent: (player: Player<any>) => player.inFrontOfTargetChangeEmitter,
		getValue: (player: Player<any>) => player.getInFrontOfTarget(),
		setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
			player.setInFrontOfTarget(eventID, newValue);
		},
	},
};

export const TankAssignment = {
	type: 'enum' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'tank-selector',
			'threat-metrics',
			'within-raid-sim-hide',
		],
		label: 'Tank Assignment',
		labelTooltip: 'Determines which mobs will be tanked. Most mobs default to targeting the Main Tank, but in preset multi-target encounters this is not always true.',
		values: [
			{ name: 'None', value: -1 },
			{ name: 'Main Tank', value: 0 },
			{ name: 'Tank 2', value: 1 },
			{ name: 'Tank 3', value: 2 },
			{ name: 'Tank 4', value: 3 },
		],
		changedEvent: (player: Player<any>) => player.getRaid()!.tanksChangeEmitter,
		getValue: (player: Player<any>) => player.getRaid()!.getTanks().findIndex(tank => RaidTarget.equals(tank, player.makeRaidTarget())),
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const newTanks = [];
			if (newValue != -1) {
				for (let i = 0; i < newValue; i++) {
					newTanks.push(emptyRaidTarget());
				}
				newTanks.push(player.makeRaidTarget());
			}
			player.getRaid()!.setTanks(eventID, newTanks);
		},
	},
};

export const IncomingHps = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		extraCssClasses: [
			'incoming-hps-picker',
		],
		label: 'Incoming HPS',
		labelTooltip: `
			<p>Average amount of healing received per second. Used for calculating chance of death.</p>
			<p>If set to 0, defaults to 125% of DTPS.</p>
		`,
		changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
		getValue: (player: Player<any>) => player.getHealingModel().hps,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const healingModel = player.getHealingModel();
			healingModel.hps = newValue;
			player.setHealingModel(eventID, healingModel);
		},
		enableWhen: (player: Player<any>) => player.getRaid()!.getTanks().find(tank => RaidTarget.equals(tank, player.makeRaidTarget())) != null,
	},
};

export const HealingCadence = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		float: true,
		extraCssClasses: [
			'healing-cadence-picker',
		],
		label: 'Healing Cadence',
		labelTooltip: `
			<p>How often the incoming heal 'ticks', in seconds. Generally, longer durations favor Effective Hit Points (EHP) for minimizing Chance of Death, while shorter durations favor avoidance.</p>
			<p>Example: if Incoming HPS is set to 1000 and this is set to 1s, then every 1s a heal will be received for 1000. If this is instead set to 2s, then every 2s a heal will be recieved for 2000.</p>
			<p>If set to 0, defaults to 2.5 seconds.</p>
		`,
		changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
		getValue: (player: Player<any>) => player.getHealingModel().cadenceSeconds,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const healingModel = player.getHealingModel();
			healingModel.cadenceSeconds = newValue;
			player.setHealingModel(eventID, healingModel);
		},
		enableWhen: (player: Player<any>) => player.getRaid()!.getTanks().find(tank => RaidTarget.equals(tank, player.makeRaidTarget())) != null,
	},
};

export const HpPercentForDefensives = {
	type: 'number' as const,
	getModObject: (simUI: IndividualSimUI<any>) => simUI.player,
	config: {
		float: true,
		extraCssClasses: [
			'hp-percent-for-defensives-picker',
		],
		label: 'HP % for Defensive CDs',
		labelTooltip: `
			<p>% of Maximum Health, below which defensive cooldowns are allowed to be used.</p>
			<p>If set to 0, this restriction is disabled.</p>
		`,
		changedEvent: (player: Player<any>) => player.cooldownsChangeEmitter,
		getValue: (player: Player<any>) => player.getCooldowns().hpPercentForDefensives * 100,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			const cooldowns = player.getCooldowns();
			cooldowns.hpPercentForDefensives = newValue / 100;
			player.setCooldowns(eventID, cooldowns);
		},
	},
};
