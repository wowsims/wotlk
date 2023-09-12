import { Component } from '../core/components/component.js';
import { UnitReferencePicker } from '../core/components/raid_target_picker.js';

import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import { Class, UnitReference, Spec } from '../core/proto/common.js';
import { emptyUnitReference } from '../core/proto_utils/utils.js';

import { RaidSimUI } from './raid_sim_ui.js';
import { PriestTalents } from 'ui/core/proto/priest.js';
import { DeathknightTalents } from 'ui/core/proto/deathknight.js';

export class AssignmentsPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly innervatesPicker: InnervatesPicker;
	private readonly powerInfusionsPicker: PowerInfusionsPicker;
	private readonly tricksOfTheTradesPicker: TricksOfTheTradesPicker;
	private readonly unholyFrenzyPicker: UnholyFrenzyPicker;
	private readonly focusMagicsPicker: FocusMagicsPicker;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'assignments-picker-root');
		this.raidSimUI = raidSimUI;

		this.innervatesPicker = new InnervatesPicker(this.rootElem, raidSimUI);
		this.powerInfusionsPicker = new PowerInfusionsPicker(this.rootElem, raidSimUI);
		this.tricksOfTheTradesPicker = new TricksOfTheTradesPicker(this.rootElem, raidSimUI);
		this.unholyFrenzyPicker = new UnholyFrenzyPicker(this.rootElem, raidSimUI);
		this.focusMagicsPicker = new FocusMagicsPicker(this.rootElem, raidSimUI);
	}
}

interface AssignmentTargetPicker {
	player: Player<any>,
	targetPicker: UnitReferencePicker<Player<any>>,
	targetPlayer: Player<any> | null;
};

abstract class AssignedBuffPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly playersContainer: HTMLElement;

	private targetPickers: Array<AssignmentTargetPicker>;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'assigned-buff-picker-root');
		this.raidSimUI = raidSimUI;
		this.targetPickers = [];

		this.playersContainer = document.createElement('div');
		this.playersContainer.classList.add('assigned-buff-container');
		this.rootElem.appendChild(this.playersContainer);

		this.raidSimUI.changeEmitter.on(eventID => this.update());
		this.update();
	}

	private update() {
		this.playersContainer.innerHTML = `
			<label class="assignmented-buff-label form-label">${this.getTitle()}</label>
		`

		const sourcePlayers = this.getSourcePlayers();
		if (sourcePlayers.length == 0)
			this.rootElem.classList.add('hide');
		else
			this.rootElem.classList.remove('hide');

		this.targetPickers = sourcePlayers.map((sourcePlayer, sourcePlayerIndex) => {
			const row = document.createElement('div');
			row.classList.add('assigned-buff-player', 'input-inline');
			this.playersContainer.appendChild(row);

			let sourceElem = document.createElement('div');
			sourceElem.classList.add('raid-target-picker-root');
			sourceElem.appendChild(
				UnitReferencePicker.makeOptionElem({ player: sourcePlayer, isDropdown: false })
			);
			row.appendChild(sourceElem);

			const arrow = document.createElement('i');
			arrow.classList.add('assigned-buff-arrow', 'fa', 'fa-arrow-right');
			row.appendChild(arrow);

			const raidTargetPicker: UnitReferencePicker<Player<any>> | null = new UnitReferencePicker<Player<any>>(row, this.raidSimUI.sim.raid, sourcePlayer, {
				extraCssClasses: ['assigned-buff-target-picker'],
				noTargetLabel: 'Unassigned',
				compChangeEmitter: this.raidSimUI.sim.raid.compChangeEmitter,

				changedEvent: (player: Player<any>) => player.specOptionsChangeEmitter,
				getValue: (player: Player<any>) => this.getPlayerValue(player),
				setValue: (eventID: EventID, player: Player<any>, newValue: UnitReference) => this.setPlayerValue(eventID, player, newValue),
			});

			const targetPickerData = {
				player: sourcePlayer,
				targetPicker: raidTargetPicker!,
				targetPlayer: this.raidSimUI.sim.raid.getPlayerFromUnitReference(raidTargetPicker!.getInputValue()),
			};

			raidTargetPicker!.changeEmitter.on(eventID => {
				targetPickerData.targetPlayer = this.raidSimUI.sim.raid.getPlayerFromUnitReference(raidTargetPicker!.getInputValue());
			});

			return targetPickerData;
		});
	}

	abstract getTitle(): string;
	abstract getSourcePlayers(): Array<Player<any>>;

	abstract getPlayerValue(player: Player<any>): UnitReference;
	abstract setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference): void;
}

class InnervatesPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Innervate';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassDruid));
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecBalanceDruid>).getSpecOptions().innervateTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecBalanceDruid>).getSpecOptions();
		newOptions.innervateTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

class PowerInfusionsPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Power Infusion';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassPriest) && (player.getTalents() as PriestTalents).powerInfusion);
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecSmitePriest>).getSpecOptions().powerInfusionTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecSmitePriest>).getSpecOptions();
		newOptions.powerInfusionTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

class TricksOfTheTradesPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Tricks of the Trade';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassRogue));
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecRogue>).getSpecOptions().tricksOfTheTradeTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecRogue>).getSpecOptions();
		newOptions.tricksOfTheTradeTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

class UnholyFrenzyPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Unholy Frenzy';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassDeathknight) && (player.getTalents() as DeathknightTalents).hysteria);
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecDeathknight>).getSpecOptions().unholyFrenzyTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecDeathknight>).getSpecOptions();
		newOptions.unholyFrenzyTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

class FocusMagicsPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Focus Magic';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassMage));
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecMage>).getSpecOptions().focusMagicTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecMage>).getSpecOptions();
		newOptions.focusMagicTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}
