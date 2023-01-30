import { BooleanPicker } from '../components/boolean_picker.js';
import { EnumPicker, EnumPickerConfig } from '../components/enum_picker.js';
import { Conjured } from '../proto/common.js';
import { RaidTarget } from '../proto/common.js';
import { TristateEffect } from '../proto/common.js';
import { Party } from '../party.js';
import { Player } from '../player.js';
import { Sim } from '../sim.js';
import { Target } from '../target.js';
import { Encounter } from '../encounter.js';
import { Raid } from '../raid.js';
import { SimUI } from '../sim_ui.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { emptyRaidTarget } from '../proto_utils/utils.js';

export function makeShow1hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: [
			'show-1h-weapons-selector', 'mb-0'
		],
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
		extraCssClasses: [
			'show-2h-weapons-selector', 'mb-0'
		],
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
		extraCssClasses: [
			'show-matching-gems-selector',
		],
		label: 'Match Socket',
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().matchingGemsOnly,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.matchingGemsOnly = newValue;
			sim.setFilters(eventID, filters);
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
	getValue: (player: Player<any>) => (player.getRaid()?.getTanks() || []).findIndex(tank => RaidTarget.equals(tank, player.makeRaidTarget())),
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
};

export const IncomingHps = {
	type: 'number' as const,
	label: 'Incoming HPS',
	labelTooltip: `
		<p>Average amount of healing received per second. Used for calculating chance of death.</p>
		<p>If set to 0, defaults to 150% of DTPS.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().hps,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.hps = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => RaidTarget.equals(tank, player.makeRaidTarget())) != null,
};

export const HealingCadence = {
	type: 'number' as const,
	float: true,
	label: 'Healing Cadence',
	labelTooltip: `
		<p>How often the incoming heal 'ticks', in seconds. Generally, longer durations favor Effective Hit Points (EHP) for minimizing Chance of Death, while shorter durations favor avoidance.</p>
		<p>Example: if Incoming HPS is set to 1000 and this is set to 1s, then every 1s a heal will be received for 1000. If this is instead set to 2s, then every 2s a heal will be recieved for 2000.</p>
		<p>If set to 0, defaults to 2.0 seconds.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().cadenceSeconds,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.cadenceSeconds = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => RaidTarget.equals(tank, player.makeRaidTarget())) != null,
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
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => RaidTarget.equals(tank, player.makeRaidTarget())) != null,
};

export const HpPercentForDefensives = {
	type: 'number' as const,
	float: true,
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
