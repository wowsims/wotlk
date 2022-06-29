import { SimResult, SimResultFilter, UnitMetrics } from '/tbc/core/proto_utils/sim_result.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { EnumPicker } from '/tbc/core/components/enum_picker.js';
import { Input } from '/tbc/core/components/input.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

const ALL_UNITS = -1;

interface FilterData {
	player: number,
	target: number,
};

export class ResultsFilter extends ResultComponent {
	private readonly currentFilter: FilterData;

	readonly changeEmitter: TypedEvent<void>;

	private readonly playerFilter: PlayerFilter;
	private readonly targetFilter: TargetFilter;

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'results-filter-root';
		super(config);
		this.currentFilter = {
			player: ALL_UNITS,
			target: ALL_UNITS,
		};
		this.changeEmitter = new TypedEvent<void>();

		this.playerFilter = new PlayerFilter(this.rootElem, this.currentFilter);
		this.playerFilter.changeEmitter.on(eventID => this.changeEmitter.emit(eventID));

		this.targetFilter = new TargetFilter(this.rootElem, this.currentFilter);
		this.targetFilter.changeEmitter.on(eventID => this.changeEmitter.emit(eventID));
	}

	getFilter(): SimResultFilter {
		return {
			player: this.currentFilter.player == ALL_UNITS ? null : this.currentFilter.player,
			target: this.currentFilter.target == ALL_UNITS ? null : this.currentFilter.target,
		};
	}

	onSimResult(resultData: SimResultData) {
		this.playerFilter.setOptions(resultData.eventID, resultData.result);
		this.targetFilter.setOptions(resultData.eventID, resultData.result);
	}

	setPlayer(eventID: EventID, newPlayer: number | null) {
		this.currentFilter.player = (newPlayer === null) ? ALL_UNITS : newPlayer;
		this.playerFilter.changeEmitter.emit(eventID);
	}

	setTarget(eventID: EventID, newTarget: number | null) {
		this.currentFilter.target = (newTarget === null) ? ALL_UNITS : newTarget;
		this.targetFilter.changeEmitter.emit(eventID);
	}
}

interface UnitFilterOption {
	iconUrl: string,
	text: string,
	color: string,
	value: number,
};

// Dropdown menu for filtering by player.
abstract class UnitGroupFilter extends Input<FilterData, number> {
	private readonly filterData: FilterData;
	readonly changeEmitter: TypedEvent<void>;

	private allUnitsOption: UnitFilterOption;
	private currentOptions: Array<UnitFilterOption>;

	private readonly buttonElem: HTMLElement;
	private readonly dropdownElem: HTMLElement;

	constructor(parent: HTMLElement, filterData: FilterData, allUnitsLabel: string) {
		const changeEmitter = new TypedEvent<void>();
		super(parent, 'unit-filter-root', filterData, {
			extraCssClasses: [
				'dropdown-root',
			],
			changedEvent: (filterData: FilterData) => changeEmitter,
			getValue: (filterData: FilterData) => this.getFilterDataValue(filterData),
			setValue: (eventID: EventID, filterData: FilterData, newValue: number) => this.setFilterDataValue(filterData, newValue),
		});
		this.filterData = filterData;
		this.changeEmitter = changeEmitter;

		this.allUnitsOption = {
			iconUrl: '',
			text: allUnitsLabel,
			color: 'black',
			value: ALL_UNITS,
		};
		this.currentOptions = [this.allUnitsOption];

		this.rootElem.innerHTML = `
			<div class="dropdown-button unit-filter-button"></div>
			<div class="dropdown-panel unit-filter-dropdown"></div>
    `;

		this.buttonElem = this.rootElem.getElementsByClassName('unit-filter-button')[0] as HTMLElement;
		this.dropdownElem = this.rootElem.getElementsByClassName('unit-filter-dropdown')[0] as HTMLElement;

		this.buttonElem.addEventListener('click', event => {
			event.preventDefault();
		});

		this.init();
	}

	abstract getFilterDataValue(filterData: FilterData): number;
	abstract setFilterDataValue(filterData: FilterData, newValue: number): void;
	abstract getAllUnits(simResult: SimResult): Array<UnitMetrics>;

	setOptions(eventID: EventID, simResult: SimResult) {
		this.currentOptions = [this.allUnitsOption].concat(this.getAllUnits(simResult).map(unit => {
			return {
				iconUrl: unit.iconUrl || '',
				text: unit.label,
				color: unit.classColor || 'black',
				value: unit.index,
			};
		}));

		const hasSameOption = this.currentOptions.find(option => option.value == this.getInputValue()) != null;
		if (!hasSameOption) {
			this.setFilterDataValue(this.filterData, this.allUnitsOption.value);
			this.changeEmitter.emit(eventID);
		}

		this.dropdownElem.innerHTML = '';
		this.currentOptions.forEach(option => this.dropdownElem.appendChild(this.makeOption(option)));
	}

	private makeOption(data: UnitFilterOption): HTMLElement {
		const option = this.makeOptionElem(data);

		option.addEventListener('click', event => {
			event.preventDefault();
			this.setFilterDataValue(this.filterData, data.value);
			this.changeEmitter.emit(TypedEvent.nextEventID());
		});

		return option;
	}

	private makeOptionElem(data: UnitFilterOption): HTMLElement {
		const optionContainer = document.createElement('div');
		optionContainer.classList.add('dropdown-option-container');

		const option = document.createElement('div');
		option.classList.add('dropdown-option', 'unit-filter-option');
		optionContainer.appendChild(option);

		if (data.color) {
			option.style.backgroundColor = data.color;
		}

		if (data.iconUrl) {
			const icon = document.createElement('img');
			icon.src = data.iconUrl;
			icon.classList.add('unit-filter-icon');
			option.appendChild(icon);
		}

		if (data.text) {
			const label = document.createElement('span');
			label.textContent = data.text;
			label.classList.add('unit-filter-label');
			option.appendChild(label);
		}

		return optionContainer;
	}

	getInputElem(): HTMLElement {
		return this.buttonElem;
	}

	getInputValue(): number {
		return this.getFilterDataValue(this.filterData);
	}

	setInputValue(newValue: number) {
		this.setFilterDataValue(this.filterData, newValue);

		const optionData = this.currentOptions.find(optionData => optionData.value == newValue);
		if (!optionData) {
			return;
		}

		this.buttonElem.innerHTML = '';
		this.buttonElem.appendChild(this.makeOptionElem(optionData));
	}
}

class PlayerFilter extends UnitGroupFilter {
	constructor(parent: HTMLElement, filterData: FilterData) {
		super(parent, filterData, 'All Players');
		this.rootElem.classList.add('player-filter-root');
	}

	getFilterDataValue(filterData: FilterData): number {
		return filterData.player;
	}
	setFilterDataValue(filterData: FilterData, newValue: number): void {
		filterData.player = newValue;
	}
	getAllUnits(simResult: SimResult): Array<UnitMetrics> {
		return simResult.getPlayers();
	}
}

class TargetFilter extends UnitGroupFilter {
	constructor(parent: HTMLElement, filterData: FilterData) {
		super(parent, filterData, 'All Targets');
		this.rootElem.classList.add('target-filter-root');
	}

	getFilterDataValue(filterData: FilterData): number {
		return filterData.target;
	}
	setFilterDataValue(filterData: FilterData, newValue: number): void {
		filterData.target = newValue;
	}
	getAllUnits(simResult: SimResult): Array<UnitMetrics> {
		return simResult.getTargets();
	}
}
