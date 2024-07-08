import { ActionMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class DtpsMeleeMetricsTable extends MetricsTable<ActionMetrics> {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'dtps-melee-metrics-root';
		super(config, [
			MetricsTable.nameCellConfig((metric: ActionMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
				};
			}),
			{
				name: 'DPS',
				tooltip: '伤害 / 战斗时长',
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
				tooltip: '伤害 / (命中 + 暴击 + 擦过 + 格挡)',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => metric.avgHit.toFixed(1),
			},
			{
				name: '施法次数',
				tooltip: '施法次数',
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => metric.casts.toFixed(1),
			},
			{
				name: '命中次数',
				tooltip: '命中 + 暴击 + 擦过 + 格挡',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				getDisplayString: (metric: ActionMetrics) => metric.landedHits.toFixed(1),
			},
			{
				name: '未命中 %',
				tooltip: '未命中 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.missPercent,
				getDisplayString: (metric: ActionMetrics) => metric.missPercent.toFixed(2) + '%',
			},
			{
				name: '闪避 %',
				tooltip: '闪避 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.dodgePercent,
				getDisplayString: (metric: ActionMetrics) => metric.dodgePercent.toFixed(2) + '%',
			},
			{
				name: '招架 %',
				tooltip: '招架 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.parryPercent,
				getDisplayString: (metric: ActionMetrics) => metric.parryPercent.toFixed(2) + '%',
			},
			{
				name: '格挡 %',
				tooltip: '格挡 / 攻击次数',
				getValue: (metric: ActionMetrics) => metric.blockPercent,
				getDisplayString: (metric: ActionMetrics) => metric.blockPercent.toFixed(2) + '%',
			},
			{
				name: '暴击 %',
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

		const targets = resultData.result.getTargets(resultData.filter);
		const targetActions = targets.map(target => target.getMeleeActions().map(action => action.forTarget(resultData.filter))).flat();
		const actionGroups = ActionMetrics.groupById(targetActions);

		return actionGroups;
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		// TODO: Use NPC ID here instead of pet ID.
		return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}
}
