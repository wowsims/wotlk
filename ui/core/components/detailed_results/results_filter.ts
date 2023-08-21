import { UnitReference, UnitReference_Type as UnitType } from '../../proto/common.js';
import { SimResult, SimResultFilter } from '../../proto_utils/sim_result.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { UnitPicker, UnitValueConfig, UnitValue } from '../../components/unit_picker.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

const ALL_UNITS = -1;

interface FilterData {
	player: number,
	target: number,
};

export class ResultsFilter extends ResultComponent {
	private readonly currentFilter: FilterData;

	readonly changeEmitter: TypedEvent<void>;

	private readonly playerFilter: UnitPicker<FilterData>;
	private readonly targetFilter: UnitPicker<FilterData>;

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
				'd-none',
			],
			changedEvent: (_filterData: FilterData) => this.changeEmitter,
			sourceToValue: (src: UnitReference|undefined) => this.refToValue(src),
			valueToSource: (val: UnitValue) => val.value,
			getValue: (filterData: FilterData) => this.numToRef(filterData.player, true),
			setValue: (eventID: EventID, filterData: FilterData, newValue: UnitReference|undefined) => this.setPlayer(eventID, this.refToNum(newValue)),
			values: [],
		});

		this.targetFilter = new UnitPicker(this.rootElem, this.currentFilter, {
			extraCssClasses: [
				'target-filter-root',
				'd-none',
			],
			changedEvent: (_filterData: FilterData) => this.changeEmitter,
			sourceToValue: (src: UnitReference|undefined) => this.refToValue(src),
			valueToSource: (val: UnitValue) => val.value,
			getValue: (filterData: FilterData) => this.numToRef(filterData.target, false),
			setValue: (eventID: EventID, filterData: FilterData, newValue: UnitReference|undefined) => this.setTarget(eventID, this.refToNum(newValue)),
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
		this.playerFilter.rootElem.classList.remove('d-none');
		this.targetFilter.rootElem.classList.remove('d-none');
	}

	setPlayer(eventID: EventID, newPlayer: number | null) {
		this.currentFilter.player = (newPlayer === null) ? ALL_UNITS : newPlayer;
		this.changeEmitter.emit(eventID);
	}

	setTarget(eventID: EventID, newTarget: number | null) {
		this.currentFilter.target = (newTarget === null) ? ALL_UNITS : newTarget;
		this.changeEmitter.emit(eventID);
	}

	private refToValue(ref: UnitReference|undefined): UnitValue {
		if (!ref || ref.type == UnitType.Unknown) {
			return {
				value: ref,
			};
		} else if (ref.type == UnitType.AllPlayers) {
			return {
				iconUrl: '',
				text: 'All Players',
				value: ref,
			};
		} else if (ref.type == UnitType.AllTargets) {
			return {
				iconUrl: '',
				text: 'All Targets',
				value: ref,
			};
		} else if (this.hasLastSimResult()) {
			const simResult = this.getLastSimResult();
			const unit = ref.type == UnitType.Player
				? simResult.result.getPlayerWithRaidIndex(ref.index)
				: ref.type == UnitType.Target 
					? simResult.result.getTargetWithEncounterIndex(ref.index)
					: null;

			if (unit) {
				return {
					iconUrl: unit.iconUrl || '',
					text: unit.label,
					color: unit.classColor || '',
					value: ref,
				};
			}
		}

		return {
			value: ref,
		};
	}

	private refToNum(ref: UnitReference|undefined): number {
		return (!ref || ref.type == UnitType.AllPlayers || ref.type == UnitType.AllTargets) ? ALL_UNITS : ref.index;
	}

	private numToRef(idx: number, isPlayer: boolean): UnitReference {
		if (isPlayer) {
			return idx == ALL_UNITS
				? UnitReference.create({type: UnitType.AllPlayers})
				: UnitReference.create({type: UnitType.Player, index: idx});
		} else {
			return idx == ALL_UNITS
				? UnitReference.create({type: UnitType.AllTargets})
				: UnitReference.create({type: UnitType.Target, index: idx});
		}
	}

	private getUnitOptions(eventID: EventID, simResult: SimResult, isPlayer: boolean): Array<UnitValueConfig> {
		const allUnitsOption = UnitReference.create({type: isPlayer ? UnitType.AllPlayers : UnitType.AllTargets});

		const unitOptions = (isPlayer ? simResult.getPlayers() : simResult.getTargets())
			.map(unit => UnitReference.create({type: isPlayer ? UnitType.Player : UnitType.Target, index: unit.index}));

		const options = [allUnitsOption].concat(unitOptions);

		const curRef = this.numToRef(isPlayer ? this.currentFilter.player : this.currentFilter.target, isPlayer);
		const hasSameOption = options.find(option => UnitReference.equals(option, curRef)) != null;
		if (!hasSameOption) {
			if (isPlayer) {
				this.currentFilter.player = ALL_UNITS;
			} else {
				this.currentFilter.target = ALL_UNITS;
			}
			this.changeEmitter.emit(eventID);
		}

		return options.map(o => {
			return {
				value: this.refToValue(o),
			};
		});
	}
}