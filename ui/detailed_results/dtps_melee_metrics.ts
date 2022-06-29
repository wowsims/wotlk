import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { ActionMetrics, SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';

import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare var $: any;
declare var tippy: any;

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
				tooltip: 'Damage / Encounter Duration',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.dps,
				getDisplayString: (metric: ActionMetrics) => metric.dps.toFixed(1),
			},
			{
				name: 'Avg Cast',
				tooltip: 'Damage / Casts',
				getValue: (metric: ActionMetrics) => metric.avgCast,
				getDisplayString: (metric: ActionMetrics) => metric.avgCast.toFixed(1),
			},
			{
				name: 'Avg Hit',
				tooltip: 'Damage / (Hits + Crits + Glances + Blocks)',
				getValue: (metric: ActionMetrics) => metric.avgHit,
				getDisplayString: (metric: ActionMetrics) => metric.avgHit.toFixed(1),
			},
			{
				name: 'Casts',
				tooltip: 'Casts',
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => metric.casts.toFixed(1),
			},
			{
				name: 'Hits',
				tooltip: 'Hits + Crits + Glances + Blocks',
				getValue: (metric: ActionMetrics) => metric.landedHits,
				getDisplayString: (metric: ActionMetrics) => metric.landedHits.toFixed(1),
			},
			{
				name: 'Miss %',
				tooltip: 'Misses / Swings',
				getValue: (metric: ActionMetrics) => metric.missPercent,
				getDisplayString: (metric: ActionMetrics) => metric.missPercent.toFixed(2) + '%',
			},
			{
				name: 'Dodge %',
				tooltip: 'Dodges / Swings',
				getValue: (metric: ActionMetrics) => metric.dodgePercent,
				getDisplayString: (metric: ActionMetrics) => metric.dodgePercent.toFixed(2) + '%',
			},
			{
				name: 'Parry %',
				tooltip: 'Parries / Swings',
				getValue: (metric: ActionMetrics) => metric.parryPercent,
				getDisplayString: (metric: ActionMetrics) => metric.parryPercent.toFixed(2) + '%',
			},
			{
				name: 'Block %',
				tooltip: 'Blocks / Swings',
				getValue: (metric: ActionMetrics) => metric.blockPercent,
				getDisplayString: (metric: ActionMetrics) => metric.blockPercent.toFixed(2) + '%',
			},
			{
				name: 'Crit %',
				tooltip: 'Crits / Swings',
				getValue: (metric: ActionMetrics) => metric.critPercent,
				getDisplayString: (metric: ActionMetrics) => metric.critPercent.toFixed(2) + '%',
			},
			{
				name: 'Crush %',
				tooltip: 'Crushes / Swings',
				getValue: (metric: ActionMetrics) => metric.crushPercent,
				getDisplayString: (metric: ActionMetrics) => metric.crushPercent.toFixed(2) + '%',
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
