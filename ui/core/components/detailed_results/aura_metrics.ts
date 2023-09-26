import { ActionId } from '../../proto_utils/action_id';
import { AuraMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result';

import { ColumnSortType, MetricsTable } from './metrics_table';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component';

export class AuraMetricsTable extends MetricsTable<AuraMetrics> {
	private readonly useDebuffs: boolean;

	constructor(config: ResultComponentConfig, useDebuffs: boolean) {
		if (useDebuffs) {
			config.rootCssClass = 'debuff-metrics-root';
		} else {
			config.rootCssClass = 'buff-metrics-root';
		}
		super(config, [
			MetricsTable.nameCellConfig((metric: AuraMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
				};
			}),
			{
				name: 'Procs',
				tooltip: 'Procs',
				getValue: (metric: AuraMetrics) => metric.averageProcs,
				getDisplayString: (metric: AuraMetrics) => metric.averageProcs.toFixed(2),
			},
			{
				name: 'PPM',
				tooltip: 'Procs Per Minute',
				getValue: (metric: AuraMetrics) => metric.ppm,
				getDisplayString: (metric: AuraMetrics) => metric.ppm.toFixed(2),
			},
			{
				name: 'Uptime',
				tooltip: 'Uptime / Encounter Duration',
				sort: ColumnSortType.Descending,
				getValue: (metric: AuraMetrics) => metric.uptimePercent,
				getDisplayString: (metric: AuraMetrics) => metric.uptimePercent.toFixed(2) + '%',
			},
		]);
		this.useDebuffs = useDebuffs;
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<AuraMetrics>> {
		if (this.useDebuffs) {
			return AuraMetrics.groupById(resultData.result.getDebuffMetrics(resultData.filter));
		} else {
			const players = resultData.result.getPlayers(resultData.filter);
			if (players.length != 1) {
				return [];
			}
			const player = players[0];

			const auras = player.auras;
			const actionGroups = AuraMetrics.groupById(auras);
			const petGroups = player.pets.map(pet => pet.auras);

			return actionGroups.concat(petGroups);
		}
	}

	mergeMetrics(metrics: Array<AuraMetrics>): AuraMetrics {
		return AuraMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}

	shouldCollapse(metric: AuraMetrics): boolean {
		return !metric.unit?.isPet;
	}
}
