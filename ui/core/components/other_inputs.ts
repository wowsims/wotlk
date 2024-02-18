import {BooleanPicker} from '../components/boolean_picker.js';
import {EnumPicker} from '../components/enum_picker.js';
import {ItemSlot, UnitReference} from '../proto/common.js';
import {Player} from '../player.js';
import {Sim} from '../sim.js';
import {EventID} from '../typed_event.js';
import {emptyUnitReference} from '../proto_utils/utils.js';

export function makeShow1hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: ['show-1h-weapons-selector', 'mb-0'],
		label: '1H',
		inline: true,
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().oneHandedWeapons,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.oneHandedWeapons = newValue;
			sim.setFilters(eventID, filters);
		},
	});
}

export function makeShow2hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: ['show-2h-weapons-selector', 'mb-0'],
		label: '2H',
		inline: true,
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().twoHandedWeapons,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.twoHandedWeapons = newValue;
			sim.setFilters(eventID, filters);
		},
	});
}

export function makeShowMatchingGemsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: ['show-matching-gems-selector', 'input-inline', 'mb-0'],
		label: 'Match Socket',
		inline: true,
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().matchingGemsOnly,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.matchingGemsOnly = newValue;
			sim.setFilters(eventID, filters);
		},
	});
}

export function makeShowEPValuesSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: ['show-ep-values-selector', 'input-inline', 'mb-0'],
		label: 'Show EP',
		inline: true,
		changedEvent: (sim: Sim) => sim.showEPValuesChangeEmitter,
		getValue: (sim: Sim) => sim.getShowEPValues(),
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			sim.setShowEPValues(eventID, newValue);
		},
	});
}

export function makePhaseSelector(parent: HTMLElement, sim: Sim): EnumPicker<Sim> {
	return new EnumPicker<Sim>(parent, sim, {
		extraCssClasses: ['phase-selector'],
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

export const ReactionTime = {
	type: 'number' as const,
	label: 'Reaction Time',
	labelTooltip: 'Reaction time of the player, in milliseconds. Used with certain APL values (such as \'Aura Is Active With Reaction Time\').',
	changedEvent: (player: Player<any>) => player.miscOptionsChangeEmitter,
	getValue: (player: Player<any>) => player.getReactionTime(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setReactionTime(eventID, newValue);
	},
};

export const ChannelClipDelay = {
	type: 'number' as const,
	label: 'Channel Clip Delay',
	labelTooltip: 'Clip delay following channeled spells, in milliseconds. This delay occurs following any full or partial channel ending after the GCD becomes available, due to the player not being able to queue the next spell.',
	changedEvent: (player: Player<any>) => player.miscOptionsChangeEmitter,
	getValue: (player: Player<any>) => player.getChannelClipDelay(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setChannelClipDelay(eventID, newValue);
	},
};

export const InFrontOfTarget = {
	type: 'boolean' as const,
	label: 'In Front of Target',
	labelTooltip: 'Stand in front of the target, causing Blocks and Parries to be included in the attack table.',
	changedEvent: (player: Player<any>) => player.inFrontOfTargetChangeEmitter,
	getValue: (player: Player<any>) => player.getInFrontOfTarget(),
	setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
		player.setInFrontOfTarget(eventID, newValue);
	},
};

export const DistanceFromTarget = {
	type: 'number' as const,
	label: 'Distance From Target',
	labelTooltip: 'Distance from targets, in yards. Used to calculate travel time for certain spells.',
	changedEvent: (player: Player<any>) => player.distanceFromTargetChangeEmitter,
	getValue: (player: Player<any>) => player.getDistanceFromTarget(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setDistanceFromTarget(eventID, newValue);
	},
};

export const nibelungAverageCasts =  {
	type: 'number' as const,
	label: "Nibelung's Valkyr Survival (in # of casts)",
	labelTooltip: 'Number of casts of Nibelung\'s summoned Valkyrs get out before they die (max 16)',
	changedEvent: (player: Player<any>) => player.changeEmitter,
	getValue: (player: Player<any>) => player.getNibelungAverageCasts(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setNibelungAverageCastsSet(eventID, true);
		player.setNibelungAverageCasts(eventID, newValue);
	},
	showWhen: (player: Player<any>) => [49992, 50648].includes(player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.id || 0)
}

export const TankAssignment = {
	type: 'enum' as const,
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
	getValue: (player: Player<any>) => (player.getRaid()?.getTanks() || []).findIndex(tank => UnitReference.equals(tank, player.makeUnitReference())),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const newTanks = [];
		if (newValue != -1) {
			for (let i = 0; i < newValue; i++) {
				newTanks.push(emptyUnitReference());
			}
			newTanks.push(player.makeUnitReference());
		}
		player.getRaid()!.setTanks(eventID, newTanks);
	},
};

export const IncomingHps = {
	type: 'number' as const,
	label: 'Incoming HPS',
	labelTooltip: `
		<p>Average amount of healing received per second. Used for calculating chance of death.</p>
		<p>If set to 0, defaults to 17.5% of the primary target's base DPS.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().hps,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.hps = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const HealingCadence = {
	type: 'number' as const,
	float: true,
	label: 'Healing Cadence',
	labelTooltip: `
		<p>How often the incoming heal 'ticks', in seconds. Generally, longer durations favor Effective Hit Points (EHP) for minimizing Chance of Death, while shorter durations favor avoidance.</p>
		<p>Example: if Incoming HPS is set to 1000 and this is set to 1s, then every 1s a heal will be received for 1000. If this is instead set to 2s, then every 2s a heal will be recieved for 2000.</p>
		<p>If set to 0, defaults to 1.5 times the primary target's base swing timer, and half that for dual wielding targets.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().cadenceSeconds,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.cadenceSeconds = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const HealingCadenceVariation = {
	type: 'number' as const,
	float: true,
	label: 'Cadence +/-',
	labelTooltip: `
		<p>Magnitude of random variation in healing intervals, in seconds.</p>
		<p>Example: if Healing Cadence is set to 1s with 0.5s variation, then the interval between successive heals will vary uniformly between 0.5 and 1.5s. If the variation is instead set to 2s, then 50% of healing intervals will fall between 0s and 1s, and the other 50% will fall between 1s and 3s.</p>
		<p>The amount of healing per 'tick' is automatically scaled up or down based on the randomized time since the last tick, so as to keep HPS constant.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().cadenceVariation,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.cadenceVariation = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const BurstWindow = {
	type: 'number' as const,
	float: false,
	label: 'TMI Burst Window',
	labelTooltip: `
		<p>Size in whole seconds of the burst window for calculating TMI. It is important to use a consistent setting when comparing this metric.</p>
		<p>Default is 6 seconds. If set to 0, TMI calculations are disabled.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().burstWindow,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.burstWindow = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const HpPercentForDefensives = {
	type: 'number' as const,
	float: true,
	label: 'HP % for Defensive CDs',
	labelTooltip: `
		<p>% of Maximum Health, below which defensive cooldowns are allowed to be used.</p>
		<p>If set to 0, this restriction is disabled.</p>
	`,
	changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
	getValue: (player: Player<any>) => player.getSimpleCooldowns().hpPercentForDefensives * 100,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const cooldowns = player.getSimpleCooldowns();
		cooldowns.hpPercentForDefensives = newValue / 100;
		player.setSimpleCooldowns(eventID, cooldowns);
	},
};

export const InspirationUptime = {
	type: 'number' as const,
	float: true,
	label: 'Inspiration % Uptime',
	labelTooltip: `
		<p>% average of Encounter Duration, during which you have the Inspiration buff.</p>
		<p>If set to 0, the buff isn't applied.</p>
	`,
	changedEvent: (player: Player<any>) => player.healingModelChangeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().inspirationUptime * 100,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.inspirationUptime = newValue / 100;
		player.setHealingModel(eventID, healingModel);
	},
};
