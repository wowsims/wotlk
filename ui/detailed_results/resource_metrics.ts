import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { ResourceMetrics, SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';
import { ResourceType } from '/tbc/core/proto/api.js';
import { resourceNames } from '/tbc/core/proto_utils/names.js';
import { orderedResourceTypes } from '/tbc/core/proto_utils/utils.js';
import { getEnumValues } from '/tbc/core/utils.js';

import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare var $: any;
declare var tippy: any;

export class ResourceMetricsTable extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'resource-metrics-root';
		super(config);

		orderedResourceTypes.forEach(resourceType => {
			const containerElem = document.createElement('div');
			containerElem.classList.add('resource-metrics-table-container', 'hide');
			containerElem.innerHTML = `<span class="resource-metrics-table-title">${resourceNames[resourceType]}</span>`;
			this.rootElem.appendChild(containerElem);

			const childConfig = config;
			childConfig.parent = containerElem;
			const table = new TypedResourceMetricsTable(childConfig, resourceType);
			table.onUpdate.on(() => {
				if (table.rootElem.classList.contains('hide')) {
					containerElem.classList.add('hide');
				} else {
					containerElem.classList.remove('hide');
				}
			});
		});
	}

	onSimResult(resultData: SimResultData) {
	}
}

export class TypedResourceMetricsTable extends MetricsTable<ResourceMetrics> {
	readonly resourceType: ResourceType;

	constructor(config: ResultComponentConfig, resourceType: ResourceType) {
		config.rootCssClass = 'resource-metrics-table-root';
		super(config, [
			MetricsTable.nameCellConfig((metric: ResourceMetrics) => {
				return {
					name: metric.name,
					actionId: metric.actionId,
				};
			}),
			{
				name: 'Casts',
				tooltip: 'Casts',
				getValue: (metric: ResourceMetrics) => metric.events,
				getDisplayString: (metric: ResourceMetrics) => metric.events.toFixed(1),
			},
			{
				name: 'Gain',
				tooltip: 'Gain',
				sort: ColumnSortType.Descending,
				getValue: (metric: ResourceMetrics) => metric.gain,
				getDisplayString: (metric: ResourceMetrics) => metric.gain.toFixed(1),
			},
			{
				name: 'Gain / s',
				tooltip: 'Gain / Second',
				getValue: (metric: ResourceMetrics) => metric.gainPerSecond,
				getDisplayString: (metric: ResourceMetrics) => metric.gainPerSecond.toFixed(1),
			},
			{
				name: 'Avg Gain',
				tooltip: 'Gain / Event',
				getValue: (metric: ResourceMetrics) => metric.avgGain,
				getDisplayString: (metric: ResourceMetrics) => metric.avgGain.toFixed(1),
			},
			{
				name: 'Wasted Gain',
				tooltip: 'Gain that was wasted because of resource cap.',
				getValue: (metric: ResourceMetrics) => metric.wastedGain,
				getDisplayString: (metric: ResourceMetrics) => metric.wastedGain.toFixed(1),
			},
		]);
		this.resourceType = resourceType;
	}

	getGroupedMetrics(resultData: SimResultData): Array<Array<ResourceMetrics>> {
		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length != 1) {
			return [];
		}
		const player = players[0];

		const resources = player.getResourceMetrics(this.resourceType);
		const resourceGroups = ResourceMetrics.groupById(resources);
		return resourceGroups;
	}

	mergeMetrics(metrics: Array<ResourceMetrics>): ResourceMetrics {
		return ResourceMetrics.merge(metrics, true, metrics[0].unit?.petActionId || undefined);
	}
}
