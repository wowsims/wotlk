import { ResourceType } from '../../proto/api.js';
import { resourceNames } from '../../proto_utils/names.js';
import { ResourceMetrics, SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';
import { orderedResourceTypes } from '../../proto_utils/utils.js';
import { ColumnSortType, MetricsTable } from './metrics_table.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class ResourceMetricsTable extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'resource-metrics-root';
		super(config);

		orderedResourceTypes.forEach(resourceType => {
			const containerElem = document.createElement('div');
			containerElem.classList.add('resource-metrics-table-container', 'hide');
			containerElem.innerHTML = `<span class="resource-metrics-table-title">${resourceNames.get(resourceType)}</span>`;
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
				name: '施法次数',
				tooltip: '施法次数',
				getValue: (metric: ResourceMetrics) => metric.events,
				getDisplayString: (metric: ResourceMetrics) => metric.events.toFixed(1),
			},
			{
				name: '获取',
				tooltip: '获取',
				sort: ColumnSortType.Descending,
				getValue: (metric: ResourceMetrics) => metric.gain,
				getDisplayString: (metric: ResourceMetrics) => metric.gain.toFixed(1),
			},
			{
				name: '每秒获取',
				tooltip: '每秒获取',
				getValue: (metric: ResourceMetrics) => metric.gainPerSecond,
				getDisplayString: (metric: ResourceMetrics) => metric.gainPerSecond.toFixed(1),
			},
			{
				name: '平均获取',
				tooltip: '每次事件获取',
				getValue: (metric: ResourceMetrics) => metric.avgGain,
				getDisplayString: (metric: ResourceMetrics) => metric.avgGain.toFixed(1),
			},
			{
				name: '浪费获取',
				tooltip: '由于资源上限而浪费的获取。',
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
