import { Component } from '../core/components/component.js';
import { Input, InputConfig } from '../core/components/input.js';
import { RaidTargetPicker } from '../core/components/raid_target_picker.js';
import { Player } from '../core/player.js';
import { Raid } from '../core/raid.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { Class } from '../core/proto/common.js';
import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { getEnumValues } from '../core/utils.js';
import { wait } from '../core/utils.js';
import { emptyRaidTarget } from '../core/proto_utils/utils.js';

import { BuffBot } from './buff_bot.js';
import { RaidSimUI } from './raid_sim_ui.js';

declare var tippy: any;

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
	playerOrBot: Player<any> | BuffBot,
	targetPicker: RaidTargetPicker<Player<any>> | RaidTargetPicker<BuffBot>,
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

		this.playersContainer = document.createElement('fieldset');
		this.playersContainer.classList.add('assigned-buff-players-container', 'settings-section');
		this.rootElem.appendChild(this.playersContainer);

		this.update();
		this.raidSimUI.changeEmitter.on(eventID => {
			this.update();
		});
	}

	private update() {
		this.playersContainer.innerHTML = `
			<legend>${this.getTitle()}</legend>
		`;

		const sourcePlayers = this.getSourcePlayers();
		if (sourcePlayers.length == 0) {
			this.rootElem.style.display = 'none';
		} else {
			this.rootElem.style.display = 'initial';
		}

		this.targetPickers = sourcePlayers.map((sourcePlayer, sourcePlayerIndex) => {
			const row = document.createElement('div');
			row.classList.add('assigned-buff-player');
			this.playersContainer.appendChild(row);

			const sourceElem = RaidTargetPicker.makeOptionElem({
				iconUrl: sourcePlayer instanceof Player ? sourcePlayer.getTalentTreeIcon() : sourcePlayer.settings.iconUrl,
				text: sourcePlayer.getLabel(),
				color: sourcePlayer.getClassColor(),
				isDropdown: false,
			});
			sourceElem.classList.add('raid-target-picker-root');
			row.appendChild(sourceElem);

			const arrow = document.createElement('span');
			arrow.classList.add('assigned-buff-arrow', 'fa', 'fa-arrow-right');
			row.appendChild(arrow);

			let raidTargetPicker: RaidTargetPicker<Player<any>> | RaidTargetPicker<BuffBot> | null = null;
			if (sourcePlayer instanceof Player) {
				raidTargetPicker = new RaidTargetPicker<Player<any>>(row, this.raidSimUI.sim.raid, sourcePlayer, {
					extraCssClasses: [
						'assigned-buff-target-picker',
					],
					noTargetLabel: 'Unassigned',
					compChangeEmitter: this.raidSimUI.sim.raid.compChangeEmitter,

					changedEvent: (player: Player<any>) => player.specOptionsChangeEmitter,
					getValue: (player: Player<any>) => this.getPlayerValue(player),
					setValue: (eventID: EventID, player: Player<any>, newValue: RaidTarget) => this.setPlayerValue(eventID, player, newValue),
				});
			} else {
				raidTargetPicker = new RaidTargetPicker<BuffBot>(row, this.raidSimUI.sim.raid, sourcePlayer, {
					extraCssClasses: [
						'assigned-buff-target-picker',
					],
					noTargetLabel: 'Unassigned',
					compChangeEmitter: this.raidSimUI.sim.raid.compChangeEmitter,

					changedEvent: (buffBot: BuffBot) => buffBot.changeEmitter,
					getValue: (buffBot: BuffBot) => this.getBuffBotValue(buffBot),
					setValue: (eventID: EventID, buffBot: BuffBot, newValue: RaidTarget) => this.setBuffBotValue(eventID, buffBot, newValue),
				});
			}

			const targetPickerData = {
				playerOrBot: sourcePlayer,
				targetPicker: raidTargetPicker!,
				targetPlayer: this.raidSimUI.sim.raid.getPlayerFromRaidTarget(raidTargetPicker!.getInputValue()),
			};

			raidTargetPicker!.changeEmitter.on(eventID => {
				targetPickerData.targetPlayer = this.raidSimUI.sim.raid.getPlayerFromRaidTarget(raidTargetPicker!.getInputValue());
			});

			return targetPickerData;
		});
	}

	abstract getTitle(): string;
	abstract getSourcePlayers(): Array<Player<any> | BuffBot>;

	abstract getPlayerValue(player: Player<any>): RaidTarget;
	abstract setPlayerValue(eventID: EventID, player: Player<any>, newValue: RaidTarget): void;

	abstract getBuffBotValue(buffBot: BuffBot): RaidTarget;
	abstract setBuffBotValue(eventID: EventID, buffBot: BuffBot, newValue: RaidTarget): void;
}

class InnervatesPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Innervate';
	}

	getSourcePlayers(): Array<Player<any> | BuffBot> {
		return this.raidSimUI.getPlayersAndBuffBots().filter(playerOrBot => playerOrBot?.getClass() == Class.ClassDruid) as Array<Player<any> | BuffBot>;
	}

	getPlayerValue(player: Player<any>): RaidTarget {
		return (player as Player<Spec.SpecBalanceDruid>).getSpecOptions().innervateTarget || emptyRaidTarget();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: RaidTarget) {
		const newOptions = (player as Player<Spec.SpecBalanceDruid>).getSpecOptions();
		newOptions.innervateTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}

	getBuffBotValue(buffBot: BuffBot): RaidTarget {
		return buffBot.getInnervateAssignment();
	}

	setBuffBotValue(eventID: EventID, buffBot: BuffBot, newValue: RaidTarget) {
		buffBot.setInnervateAssignment(eventID, newValue);
	}
}

class PowerInfusionsPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Power Infusion';
	}

	getSourcePlayers(): Array<Player<any> | BuffBot> {
		return this.raidSimUI.getPlayersAndBuffBots()
			.filter(playerOrBot => playerOrBot?.getClass() == Class.ClassPriest)
			.filter(playerOrBot => {
				if (playerOrBot instanceof BuffBot) {
					return playerOrBot.settings.buffBotId == 'Divine Spirit Priest';
				} else {
					const player = playerOrBot as Player<any>;
					if (!(player as Player<Spec.SpecSmitePriest>).getTalents().powerInfusion) {
						return false;
					}
					// Don't include shadow priests even if they have the talent, because they
					// don't have a raid target option for this.
					return player.spec != Spec.SpecShadowPriest;
				}
			}) as Array<Player<any> | BuffBot>;
	}

	getPlayerValue(player: Player<any>): RaidTarget {
		return (player as Player<Spec.SpecSmitePriest>).getSpecOptions().powerInfusionTarget || emptyRaidTarget();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: RaidTarget) {
		const newOptions = (player as Player<Spec.SpecSmitePriest>).getSpecOptions();
		newOptions.powerInfusionTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}

	getBuffBotValue(buffBot: BuffBot): RaidTarget {
		return buffBot.getPowerInfusionAssignment();
	}

	setBuffBotValue(eventID: EventID, buffBot: BuffBot, newValue: RaidTarget) {
		buffBot.setPowerInfusionAssignment(eventID, newValue);
	}
}

class TricksOfTheTradesPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Tricks of the Trade';
	}

	getSourcePlayers(): Array<Player<any> | BuffBot> {
		return this.raidSimUI.getPlayersAndBuffBots().filter(playerOrBot => playerOrBot?.getClass() == Class.ClassRogue) as Array<Player<any> | BuffBot>;
	}

	getPlayerValue(player: Player<any>): RaidTarget {
		return (player as Player<Spec.SpecRogue>).getSpecOptions().tricksOfTheTradeTarget || emptyRaidTarget();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: RaidTarget) {
		const newOptions = (player as Player<Spec.SpecRogue>).getSpecOptions();
		newOptions.tricksOfTheTradeTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}

	getBuffBotValue(buffBot: BuffBot): RaidTarget {
		return buffBot.getTricksOfTheTradeAssignment();
	}

	setBuffBotValue(eventID: EventID, buffBot: BuffBot, newValue: RaidTarget) {
		buffBot.setTricksOfTheTradeAssignment(eventID, newValue);
	}
}

class UnholyFrenzyPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Unholy Frenzy';
	}

	getSourcePlayers(): Array<Player<any> | BuffBot> {
		return this.raidSimUI.getPlayersAndBuffBots()
			.filter(playerOrBot => playerOrBot?.getClass() == Class.ClassDeathknight)
			.filter(playerOrBot => {
				if (playerOrBot instanceof BuffBot) {
					return playerOrBot.settings.buffBotId == 'Blood DK Tank';
				} else {
					const player = playerOrBot as Player<Spec.SpecDeathknight>;
					return player.getTalents().hysteria;
				}
			}) as Array<Player<any> | BuffBot>;
	}

	getPlayerValue(player: Player<any>): RaidTarget {
		return (player as Player<Spec.SpecDeathknight>).getSpecOptions().unholyFrenzyTarget || emptyRaidTarget();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: RaidTarget) {
		const newOptions = (player as Player<Spec.SpecDeathknight>).getSpecOptions();
		newOptions.unholyFrenzyTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}

	getBuffBotValue(buffBot: BuffBot): RaidTarget {
		return buffBot.getUnholyFrenzyAssignment();
	}

	setBuffBotValue(eventID: EventID, buffBot: BuffBot, newValue: RaidTarget) {
		buffBot.setUnholyFrenzyAssignment(eventID, newValue);
	}
}

class FocusMagicsPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Focus Magic';
	}

	getSourcePlayers(): Array<Player<any> | BuffBot> {
		return this.raidSimUI.getPlayersAndBuffBots().filter(playerOrBot => playerOrBot?.getClass() == Class.ClassMage) as Array<Player<any> | BuffBot>;
	}

	getPlayerValue(player: Player<any>): RaidTarget {
		return (player as Player<Spec.SpecMage>).getSpecOptions().focusMagicTarget || emptyRaidTarget();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: RaidTarget) {
		const newOptions = (player as Player<Spec.SpecMage>).getSpecOptions();
		newOptions.focusMagicTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}

	getBuffBotValue(buffBot: BuffBot): RaidTarget {
		throw new Error('Unimplemented focus magic buff bot');
	}

	setBuffBotValue(eventID: EventID, buffBot: BuffBot, newValue: RaidTarget) {
		throw new Error('Unimplemented focus magic buff bot');
	}
}
