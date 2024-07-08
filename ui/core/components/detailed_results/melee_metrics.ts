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
				name: '施法次数',
				tooltip: '施法次数',
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => metric.casts.toFixed(1),
			},
			{
				name: 'CPM',
				tooltip: '施法次数 / (战斗时长 / 60 秒)',
				getValue: (metric: ActionMetrics) => metric.castsPerMinute,
				getDisplayString: (metric: ActionMetrics) => metric.castsPerMinute.toFixed(1),
			},
			{
				name: '施法时间',
				tooltip: '平均施法时间（秒）',
				getValue: (metric: ActionMetrics) => metric.avgCastTimeMs,
				getDisplayString: (metric: ActionMetrics) => (metric.avgCastTimeMs / 1000).toFixed(2),
			},
			{
				name: '每法力治疗量',
				tooltip: '治疗 / 法力',
				getValue: (metric: ActionMetrics) => metric.hpm,
				getDisplayString: (metric: ActionMetrics) => metric.hpm.toFixed(1),
			},
			{
				name: '每施法时间治疗量',
				tooltip: '治疗 / 平均施法时间',
				getValue: (metric: ActionMetrics) => metric.healingThroughput,
				getDisplayString: (metric: ActionMetrics) => metric.healingThroughput.toFixed(1),
			},
			{
				name: 'HPS',
				tooltip: '治疗 / 战斗时长',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.hps,
				getDisplayString: (metric: ActionMetrics) => metric.hps.toFixed(1),
			},
			{
				name: '平均施法治疗量',
				tooltip: '治疗 / 施法次数',
				getValue: (metric: ActionMetrics) => metric.avgCastHealing,
				getDisplayString: (metric: ActionMetrics) => metric.avgCastHealing.toFixed(1),
			},
			{
				name: 'TPS',
				tooltip: '仇恨 / 战斗时长',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.tps,
				getDisplayString: (metric: ActionMetrics) => metric.tps.toFixed(1),
			},
			{
				name: '平均施法仇恨',
				tooltip: '仇恨 / 施法次数',
				columnClass: 'threat-metrics',
				getValue: (metric: ActionMetrics) => metric.avgCastThreat,
				getDisplayString: (metric: ActionMetrics) => metric.avgCastThreat.toFixed(1),
			},
			{
				name: '暴击 %',
				tooltip: '暴击 / 命中',
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
