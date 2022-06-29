import { ActionMetrics, UnitMetrics, SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare var $: any;
declare var tippy: any;

export enum ColumnSortType {
	None,
	Ascending,
	Descending,
}

export interface MetricsColumnConfig<T> {
	name: string,
	tooltip?: string,
	headerCellClass?: string,
	columnClass?: string,
	sort?: ColumnSortType,

	getValue?: (metric: T) => number,

	// Either getDisplayString or fillCell must be specified.
	getDisplayString?: (metric: T) => string,
	fillCell?: (metric: T, cellElem: HTMLElement, rowElem: HTMLElement) => void,
};

export abstract class MetricsTable<T> extends ResultComponent {
	private readonly columnConfigs: Array<MetricsColumnConfig<T>>;

	protected readonly tableElem: HTMLElement;
	protected readonly bodyElem: HTMLElement;

	readonly onUpdate = new TypedEvent<void>('MetricsTableUpdate');

	constructor(config: ResultComponentConfig, columnConfigs: Array<MetricsColumnConfig<T>>) {
		super(config);
		this.columnConfigs = columnConfigs;

		this.rootElem.innerHTML = `
		<table class="metrics-table tablesorter">
			<thead class="metrics-table-header">
				<tr class="metrics-table-header-row"></tr>
			</thead>
			<tbody class="metrics-table-body">
			</tbody>
		</table>
		`;

		this.tableElem = this.rootElem.getElementsByClassName('metrics-table')[0] as HTMLTableSectionElement;
		this.bodyElem = this.rootElem.getElementsByClassName('metrics-table-body')[0] as HTMLElement;

		const headerRowElem = this.rootElem.getElementsByClassName('metrics-table-header-row')[0] as HTMLElement;
		this.columnConfigs.forEach(columnConfig => {
			const headerCell = document.createElement('th');
			headerCell.classList.add('metrics-table-header-cell');
			if (columnConfig.headerCellClass) {
				headerCell.classList.add(columnConfig.headerCellClass);
			}
			if (columnConfig.columnClass) {
				headerCell.classList.add(columnConfig.columnClass);
			}
			headerCell.innerHTML = `<span>${columnConfig.name}</span>`;
			if (columnConfig.tooltip) {
				tippy(headerCell, {
					'content': columnConfig.tooltip,
					'allowHTML': true,
				});
			}
			headerRowElem.appendChild(headerCell);
		});

		const sortList = this.columnConfigs
			.map((config, i) => [i, config.sort == ColumnSortType.Ascending ? 0 : 1])
			.filter(sortData => this.columnConfigs[sortData[0]].sort);
		$(this.tableElem).tablesorter({
			sortList: sortList,
			cssChildRow: 'child-metric',
		});
	}

	protected sortMetrics(metrics: Array<T>) {
		this.columnConfigs.filter(config => config.sort).forEach(config => {
			if (!config.getValue) {
				throw new Error('Can\' apply group sorting without getValue');
			}
			if (config.sort == ColumnSortType.Ascending) {
				metrics.sort((a, b) => config.getValue!(a) - config.getValue!(b));
			} else {
				metrics.sort((a, b) => config.getValue!(b) - config.getValue!(a));
			}
		});
	}

	private addRow(metric: T): HTMLElement {
		const rowElem = document.createElement('tr');
		this.bodyElem.appendChild(rowElem);

		this.columnConfigs.forEach(columnConfig => {
			const cellElem = document.createElement('td');
			if (columnConfig.columnClass) {
				cellElem.classList.add(columnConfig.columnClass);
			}
			if (columnConfig.fillCell) {
				columnConfig.fillCell(metric, cellElem, rowElem);
			} else if (columnConfig.getDisplayString) {
				cellElem.textContent = columnConfig.getDisplayString(metric);
			} else {
				throw new Error('Metrics column config does not provide content function: ' + columnConfig.name);
			}
			rowElem.appendChild(cellElem);
		});

		this.customizeRowElem(metric, rowElem);
		return rowElem;
	}

	private addGroup(metrics: Array<T>) {
		if (metrics.length == 0) {
			return;
		}

		if (metrics.length == 1 && this.shouldCollapse(metrics[0])) {
			this.addRow(metrics[0]);
			return;
		}

		// Manually sort because tablesorter doesn't let us apply sorting to child rows.
		this.sortMetrics(metrics);

		const mergedMetrics = this.mergeMetrics(metrics);
		const parentRow = this.addRow(mergedMetrics);
		const childRows = metrics.map(metric => this.addRow(metric));
		childRows.forEach(childRow => childRow.classList.add('child-metric'));

		let expand = true;
		parentRow.classList.add('parent-metric', 'expand');
		parentRow.addEventListener('click', event => {
			expand = !expand;
			if (expand) {
				childRows.forEach(row => row.classList.remove('hide'));
				parentRow.classList.add('expand');
			} else {
				childRows.forEach(row => row.classList.add('hide'));
				parentRow.classList.remove('expand');
			}
		});
	}

	onSimResult(resultData: SimResultData) {
		this.bodyElem.textContent = '';
		const groupedMetrics = this.getGroupedMetrics(resultData).filter(group => group.length > 0);
		if (groupedMetrics.length == 0) {
			this.rootElem.classList.add('hide');
			this.onUpdate.emit(resultData.eventID);
			return;
		} else {
			this.rootElem.classList.remove('hide');
		}

		groupedMetrics.forEach(group => this.addGroup(group));
		$(this.tableElem).trigger('update');
		this.onUpdate.emit(resultData.eventID);
	}

	// Whether a single-element group should have its parent row removed.
	// Override this to add custom behavior.
	protected shouldCollapse(metric: T): boolean {
		return true;
	}

	// Override this to customize rowElem after it has been populated.
	protected customizeRowElem(metric: T, rowElem: HTMLElement) { }

	// Override this to provide custom merge behavior.
	protected mergeMetrics(metrics: Array<T>): T {
		return metrics[0];
	}

	// Returns grouped metrics to display.
	abstract getGroupedMetrics(resultData: SimResultData): Array<Array<T>>;

	static nameCellConfig<T>(getData: (metric: T) => { name: string, actionId: ActionId }): MetricsColumnConfig<T> {
		return {
			name: 'Name',
			fillCell: (metric: T, cellElem: HTMLElement, rowElem: HTMLElement) => {
				const data = getData(metric);
				cellElem.innerHTML = `
				<a class="metrics-action-icon"></a>
				<span class="metrics-action-name">${data.name}</span>
				<span class="expand-toggle fa fa-caret-down"></span>
				<span class="expand-toggle fa fa-caret-right"></span>
				`;

				const iconElem = cellElem.getElementsByClassName('metrics-action-icon')[0] as HTMLAnchorElement;
				data.actionId.setBackgroundAndHref(iconElem);
			},
		};
	}

	static playerNameCellConfig(): MetricsColumnConfig<UnitMetrics> {
		return {
			name: 'Name',
			fillCell: (player: UnitMetrics, cellElem: HTMLElement, rowElem: HTMLElement) => {
				cellElem.innerHTML = `
				<img class="metrics-action-icon" src="${player.iconUrl}"></img>
				<span class="metrics-action-name" style="color:${player.classColor}">${player.label}</span>
				`;
			},
		};
	}
}
