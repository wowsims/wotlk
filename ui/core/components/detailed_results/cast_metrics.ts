import { ActionMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';

import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class CastMetricsTable extends MetricsTable<ActionMetrics> {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'cast-metrics-root';
		super(config, [
			MetricsTable.nameCellConfig((metric: ActionMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
				};
			}),
			{
				name: 'Casts',
				tooltip: 'Casts',
				sort: ColumnSortType.Descending,
				getValue: (metric: ActionMetrics) => metric.casts,
				getDisplayString: (metric: ActionMetrics) => metric.casts.toFixed(1),
			},
			{
				name: 'CPM',
				tooltip: 'Casts / (Encounter Duration / 60 Seconds)',
				getValue: (metric: ActionMetrics) => metric.castsPerMinute,
				getDisplayString: (metric: ActionMetrics) => metric.castsPerMinute.toFixed(1),
			},
		]);
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ActionMetrics>> {
		//const actionMetrics = resultData.result.getActionMetrics(resultData.filter);
		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		const actions = player.actions.filter(action => action.casts != 0).map(action => action.forTarget(resultData.filter));
		const actionGroups = ActionMetrics.groupById(actions);
		const petGroups = player.pets.map(pet => pet.actions.filter(action => action.casts != 0).map(action => action.forTarget(resultData.filter)));

		return actionGroups.concat(petGroups);
	}

	mergeMetrics(metrics: Array<ActionMetrics>): ActionMetrics {
		return ActionMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}

	shouldCollapse(metric: ActionMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
