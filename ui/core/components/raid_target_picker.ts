import { Input, InputConfig } from '../components/input.js';
import { Player } from '../player.js';
import { Raid } from '../raid.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { UnitReference } from '../proto/common.js';
import { emptyUnitReference, cssClassForClass } from '../proto_utils/utils.js';

export interface UnitReferencePickerConfig<ModObject> extends InputConfig<ModObject, UnitReference> {
	noTargetLabel: string,
	compChangeEmitter: TypedEvent<void>,
}

interface OptionElemOptions {
	isDropdown?: boolean,
	player: Player<any> | null,
}

// Dropdown menu for selecting a player.
export class UnitReferencePicker<ModObject> extends Input<ModObject, UnitReference> {
	private readonly config: UnitReferencePickerConfig<ModObject>;
	private readonly raid: Raid;

	private curPlayer: Player<any> | null;
	private curUnitReference: UnitReference;

	private currentOptions: Array<OptionElemOptions>;

	private readonly buttonElem: HTMLElement;
	private readonly dropdownElem: HTMLElement;

	constructor(parent: HTMLElement, raid: Raid, modObj: ModObject, config: UnitReferencePickerConfig<ModObject>) {
		super(parent, 'raid-target-picker-root', modObj, config);
		this.rootElem.classList.add('dropdown');
		this.config = config;
		this.raid = raid;
		this.curPlayer = this.raid.getPlayerFromUnitReference(config.getValue(modObj));
		this.curUnitReference = this.getInputValue();

		this.rootElem.innerHTML = `
			<a
				class="raid-target-picker-button"
				href="javascript:void(0)"
				role="button"
				data-bs-toggle="dropdown"
			></a>
			<div class="dropdown-menu"></div>
    `;

		this.buttonElem = this.rootElem.querySelector('.raid-target-picker-button') as HTMLElement;
		this.dropdownElem = this.rootElem.querySelector('.dropdown-menu') as HTMLElement;

		this.buttonElem.addEventListener('click', event => event.preventDefault());

		this.currentOptions = [];
		this.updateOptions(TypedEvent.nextEventID());
		config.compChangeEmitter.on(eventID => this.updateOptions(eventID));

		this.init();
	}

	private makeTargetOptions(): Array<OptionElemOptions> {
		const unassignedOption = { player: null, isDropdown: true }
		const playerOptions = this.raid.getPlayers().filter(player => player != null).map(player => {
			return { player: player, isDropdown: true }
		});

		return [unassignedOption, ...playerOptions]
	}

	private updateOptions(eventID: EventID) {
		this.currentOptions = this.makeTargetOptions();

		this.dropdownElem.innerHTML = '';
		this.currentOptions.forEach(option => this.dropdownElem.appendChild(this.makeOption(option)));

		const prevUnitReference = this.curUnitReference;
		this.curUnitReference = this.getInputValue();
		if (!UnitReference.equals(prevUnitReference, this.curUnitReference)) {
			this.inputChanged(eventID);
		} else {
			this.setInputValue(this.curUnitReference);
		}
	}

	private makeOption(data: OptionElemOptions): HTMLElement {
		const option = UnitReferencePicker.makeOptionElem(data);

		option.addEventListener('click', event => {
			event.preventDefault();
			this.curPlayer = data.player;
			this.curUnitReference = this.getInputValue();
			this.inputChanged(TypedEvent.nextEventID());
		});

		return option;
	}

	getInputElem(): HTMLElement {
		return this.buttonElem;
	}

	getInputValue(): UnitReference {
		if (this.curPlayer) {
			return this.curPlayer.makeUnitReference();
		} else {
			return emptyUnitReference();
		}
	}

	setInputValue(newValue: UnitReference) {
		this.curUnitReference = UnitReference.clone(newValue);
		this.curPlayer = this.raid.getPlayerFromUnitReference(this.curUnitReference);

		const optionData = this.currentOptions.find(optionData => optionData.player == this.curPlayer);

		if (optionData)
			this.buttonElem.innerHTML = UnitReferencePicker.makeOptionElem({ player: optionData.player }).outerHTML;
	}

	static makeOptionElem(data: OptionElemOptions): HTMLElement {
		const classCssClass = data.player ? cssClassForClass(data.player.getClass()) : '';
		let playerFragment = document.createElement('fragment');

		playerFragment.innerHTML = `
			<div class="player ${classCssClass ? `bg-${classCssClass}-dampened` : ''}">
				<div class="player-label">
					${data.player ? `<img class="player-icon" src="${data.player.getSpecIcon()}" draggable="false"/>` : ''}
					<div class="player-details">
						<span class="player-name ${classCssClass ? `text-${classCssClass}` : ''}">
							${data.player ? data.player.getName() : 'Unassigned'}
						</span>
					</div>
				</div>
			</div>
		`

		if (data.isDropdown) {
			playerFragment.innerHTML = `
				<a class="dropdown-option" href="javascript:void(0) role="button">${playerFragment.innerHTML}</a>
			`
		}

		return playerFragment.children[0] as HTMLElement;
	}
}
