import {BooleanPicker} from '../components/boolean_picker.js';
import {EnumPicker} from '../components/enum_picker.js';
import {Player} from '../player.js';
import {ItemSlot, UnitReference} from '../proto/common.js';
import {emptyUnitReference} from '../proto_utils/utils.js';
import {Sim} from '../sim.js';
import {EventID} from '../typed_event.js';

export function makeShow1hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		extraCssClasses: ['show-1h-weapons-selector', 'mb-0'],
		label: '单手',
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
		label: '双手',
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
		label: '符合孔位颜色',
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
		label: '显示装备权重',
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
			{ name: 'P1阶段', value: 1 },
			{ name: 'P2阶段', value: 2 },
			{ name: 'P3阶段', value: 3 },
			{ name: 'P4阶段', value: 4 },
			{ name: 'P5阶段', value: 5 },
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
	label: '反应时间',
	labelTooltip: '玩家的反应时间，以毫秒为单位。用于某些 APL 值（例如“光环在反应时间内激活”）。',
	changedEvent: (player: Player<any>) => player.miscOptionsChangeEmitter,
	getValue: (player: Player<any>) => player.getReactionTime(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setReactionTime(eventID, newValue);
	},
};

export const ChannelClipDelay = {
	type: 'number' as const,
	label: '通道剪辑延迟',
	labelTooltip: '引导法术后的剪辑延迟，以毫秒为单位。由于玩家无法在GCD可用后排队下一个法术，因此在任何完整或部分引导结束后会发生此延迟。',
	changedEvent: (player: Player<any>) => player.miscOptionsChangeEmitter,
	getValue: (player: Player<any>) => player.getChannelClipDelay(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setChannelClipDelay(eventID, newValue);
	},
};

export const InFrontOfTarget = {
	type: 'boolean' as const,
	label: '正面战斗(无法打背)',
	labelTooltip: '站在目标前方，使格挡和招架包含在攻击表中。',
	changedEvent: (player: Player<any>) => player.inFrontOfTargetChangeEmitter,
	getValue: (player: Player<any>) => player.getInFrontOfTarget(),
	setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
		player.setInFrontOfTarget(eventID, newValue);
	},
};

export const DistanceFromTarget = {
	type: 'number' as const,
	label: '距离目标的距离',
	labelTooltip: '距离目标的距离，以码为单位。用于计算某些法术的飞行时间。',
	changedEvent: (player: Player<any>) => player.distanceFromTargetChangeEmitter,
	getValue: (player: Player<any>) => player.getDistanceFromTarget(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setDistanceFromTarget(eventID, newValue);
	},
};

export const nibelungAverageCasts =  {
	type: 'number' as const,
	label: "尼伯龙之瓦基里生存（以施法次数计）",
	labelTooltip: '尼伯龙之召唤的瓦基里在死亡前能施放的次数（最多16次）',
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
	label: '坦克分配',
	labelTooltip: '确定哪些怪物将由谁坦克。大多数怪物默认会攻击主坦克，但在预设的多目标战斗中，这并不总是如此。',
	values: [
		{ name: '无', value: -1 },
		{ name: '主坦克', value: 0 },
		{ name: '坦克2', value: 1 },
		{ name: '坦克3', value: 2 },
		{ name: '坦克4', value: 3 },
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
	label: 'HPS',
	labelTooltip: `
		<p>每秒接收的平均治疗量。用于计算死亡概率。</p>
		<p>如果设置为 0，则默认为主要目标基础 DPS 的 17.5%。</p>
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
	label: '治疗节奏',
	labelTooltip: `
		<p>输入治疗“跳动”的频率，以秒为单位。一般来说，较长的时间段有利于有效生命值 (EHP) 来最大限度减少死亡概率，而较短的时间段则有利于闪避。</p>
		<p>例如：如果输入 HPS 设置为 1000 而此项设置为 1s，则每 1s 将接收 1000 的治疗。如果此项设置为 2s，则每 2s 将接收 2000 的治疗。</p>
		<p>如果设置为 0，则默认为主要目标基础攻击间隔的 1.5 倍，双持目标则为其一半。</p>
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
	label: '节奏 +/-',
	labelTooltip: `
		<p>治疗间隔中的随机变化幅度，以秒为单位。</p>
		<p>例如：如果治疗节奏设置为 1s，变化幅度为 0.5s，则连续治疗之间的间隔将均匀变化在 0.5 到 1.5s 之间。如果变化幅度设置为 2s，则 50% 的治疗间隔将在 0s 到 1s 之间，另 50% 的治疗间隔将在 1s 到 3s 之间。</p>
		<p>每“跳”的治疗量会根据上一次跳动以来的随机时间自动调整，以保持 HPS 恒定。</p>
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
	label: 'TMI 爆发窗口',
	labelTooltip: `
		<p>用于计算 TMI 的爆发窗口大小，以整秒为单位。在比较此指标时，使用一致的设置非常重要。</p>
		<p>默认是 6 秒。如果设置为 0，则禁用 TMI 计算。</p>
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
	label: '防御性技能使用的HP百分比',
	labelTooltip: `
		<p>防御性技能允许使用时的最大生命值百分比。</p>
		<p>如果设置为 0，则禁用此限制。</p>
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
	label: '灵感持续时间百分比',
	labelTooltip: `
		<p>战斗期间获得灵感增益的平均持续时间百分比。</p>
		<p>如果设置为 0，则不应用此增益。</p>
	`,
	changedEvent: (player: Player<any>) => player.healingModelChangeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().inspirationUptime * 100,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.inspirationUptime = newValue / 100;
		player.setHealingModel(eventID, healingModel);
	},
};
