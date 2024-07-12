import { ActionMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';
import { bucket } from '../../utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class MeleeMetricsTable extends MetricsTable<ActionMetrics> {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'melee-metrics-root';
		super(config, [
			MetricsTable.nameCellConfig((metric: ActionMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
				};
			}),
			{
				name: 'DPS',
				tooltip: '伤害 / 战斗时间',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				getDisplayString: (metric: ActionMetrics) => metric.dps.toFixed(1),
			},
			{
				name: '平均施法',
				tooltip: '伤害 / 施法次数',
				getValue: (metric: ActionMetrics) => metric.avgCast,
				getDisplayString: (metric: ActionMetrics) => metric.avgCast.toFixed(1),
			},
			{
				name: '平均命中',
				tooltip: '伤害 / (命中 + 暴击 + 刮擦 + 格挡)',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => metric.avgHit.toFixed(1),
			},
			{
				name: 'TPS',
				tooltip: '仇恨 / 战斗时间',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.tps,
				getDisplayString: (metric: ActionMetrics) => metric.tps.toFixed(1),
			},
			{
				name: '平均施法威胁',
				tooltip: '仇恨 / 施法次数',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgCastThreat,
				getDisplayString: (metric: ActionMetrics) => metric.avgCastThreat.toFixed(1),
			},
			{
				name: '平均命中威胁',
				tooltip: '仇恨 / (命中 + 暴击 + 刮擦 + 格挡)',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgHitThreat,
				getDisplayString: (metric: ActionMetrics) => metric.avgHitThreat.toFixed(1),
			},
			{
				name: '施法次数',
				tooltip: '施法次数',
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => metric.casts.toFixed(1),
			},
			{
				name: '命中次数',
				tooltip: '命中 + 暴击 + 刮擦 + 格挡',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				getDisplayString: (metric: ActionMetrics) => metric.landedHits.toFixed(1),
			},
			{
				name: '未命中率',
				tooltip: '未命中 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.missPercent,
				getDisplayString: (metric: ActionMetrics) => metric.missPercent.toFixed(2) + '%',
			},
			{
				name: '躲闪率',
				tooltip: '躲闪 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.dodgePercent,
				getDisplayString: (metric: ActionMetrics) => metric.dodgePercent.toFixed(2) + '%',
			},
			{
				name: '招架率',
				tooltip: '招架 / 攻击次数',
				columnClass: 'in-front-of-target',
				getValue: (metric: ActionMetrics) => metric.parryPercent,
				getDisplayString: (metric: ActionMetrics) => metric.parryPercent.toFixed(2) + '%',
			},
			{
				name: '格挡率',
				tooltip: '格挡 / 攻击次数',
				columnClass: 'in-front-of-target',
				getValue: (metric: ActionMetrics) => metric.blockPercent,
				getDisplayString: (metric: ActionMetrics) => metric.blockPercent.toFixed(2) + '%',
			},
			{
				name: '偏斜率',
				tooltip: '偏斜 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.glancePercent,
				getDisplayString: (metric: ActionMetrics) => metric.glancePercent.toFixed(2) + '%',
			},
			{
				name: '暴击率',
				tooltip: '暴击 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => metric.critPercent.toFixed(2) + '%',
			},
		]);

	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ActionMetrics>> {
		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		if (player.inFrontOfTarget) {
			this.rootElem.classList.remove('hide-in-front-of-target');
		} else {
			this.rootElem.classList.add('hide-in-front-of-target');
		}

		const actions = player.getMeleeActions().map(action => action.forTarget(resultData.filter));
		const actionGroups = ActionMetrics.groupById(actions);

		const petsByName = bucket(player.pets, pet => pet.name);
		const petGroups = Object.values(petsByName).map(pets => ActionMetrics.joinById(pets.map(pet => pet.getMeleeActions().map(action => action.forTarget(resultData.filter))).flat(), true));

		return actionGroups.concat(petGroups);
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
