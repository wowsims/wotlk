import { SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { UnitPicker, UnitValueConfig } from '../../components/unit_picker.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

const ALL_UNITS = -1;

interface FilterData {
	player: number,
	target: number,
};

export class ResultsFilter extends ResultComponent {
	private readonly currentFilter: FilterData;

	readonly changeEmitter: TypedEvent<void>;

	private readonly playerFilter: UnitPicker<FilterData, number>;
	private readonly targetFilter: UnitPicker<FilterData, number>;

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'results-filter-root';
		super(config);
		this.currentFilter = {
			player: ALL_UNITS,
			target: ALL_UNITS,
		};
		this.changeEmitter = new TypedEvent<void>();

		this.playerFilter = new UnitPicker(this.rootElem, this.currentFilter, {
			extraCssClasses: [
				'player-filter-root',
			],
			changedEvent: (filterData: FilterData) => this.changeEmitter,
			getValue: (filterData: FilterData) => filterData.player,
			setValue: (eventID: EventID, filterData: FilterData, newValue: number) => this.setPlayer(eventID, newValue),
			equals: (a, b) => a == b,
			values: [],
		});

		this.targetFilter = new UnitPicker(this.rootElem, this.currentFilter, {
			extraCssClasses: [
				'target-filter-root',
			],
			changedEvent: (filterData: FilterData) => this.changeEmitter,
			getValue: (filterData: FilterData) => filterData.target,
			setValue: (eventID: EventID, filterData: FilterData, newValue: number) => this.setTarget(eventID, newValue),
			equals: (a, b) => a == b,
			values: [],
		});
	}

	getFilter(): SimResultFilter {
		return {
			player: this.currentFilter.player == ALL_UNITS ? null : this.currentFilter.player,
			target: this.currentFilter.target == ALL_UNITS ? null : this.currentFilter.target,
		};
	}

	onSimResult(resultData: SimResultData) {
		this.playerFilter.setOptions(this.getUnitOptions(resultData.eventID, resultData.result, true));
		this.targetFilter.setOptions(this.getUnitOptions(resultData.eventID, resultData.result, false));
	}

	setPlayer(eventID: EventID, newPlayer: number | null) {
		this.currentFilter.player = (newPlayer === null) ? ALL_UNITS : newPlayer;
		this.changeEmitter.emit(eventID);
	}

	setTarget(eventID: EventID, newTarget: number | null) {
		this.currentFilter.target = (newTarget === null) ? ALL_UNITS : newTarget;
		this.changeEmitter.emit(eventID);
	}

	private getUnitOptions(eventID: EventID, simResult: SimResult, isPlayer: boolean): Array<UnitValueConfig<number>> {
		const allUnitsOption = {
			iconUrl: '',
			text: isPlayer ? 'All Players' : 'All Targets',
			color: 'black',
			value: ALL_UNITS,
		};

		const unitOptions = (isPlayer ? simResult.getPlayers() : simResult.getTargets()).map(unit => {
			return {
				iconUrl: unit.iconUrl || '',
				text: unit.label,
				color: unit.classColor || 'black',
				value: unit.unitIndex,
			};
		});

		const options = [allUnitsOption].concat(unitOptions);

		const curValue = isPlayer ? this.currentFilter.player : this.currentFilter.target;
		const hasSameOption = options.find(option => option.value == curValue) != null;
		if (!hasSameOption) {
			if (isPlayer) {
				this.currentFilter.player = ALL_UNITS;
			} else {
				this.currentFilter.target = ALL_UNITS;
			}
			this.changeEmitter.emit(eventID);
		}

		return options;
	}
}