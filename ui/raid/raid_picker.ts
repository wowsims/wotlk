import { CloseButton } from '../core/components/close_button.js';
import { Component } from '../core/components/component.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { makePhaseSelector } from '../core/components/other_inputs.js';
import { Raid } from '../core/raid.js';
import { MAX_PARTY_SIZE } from '../core/party.js';
import { Party } from '../core/party.js';
import { Player } from '../core/player.js';
import { Player as PlayerProto } from '../core/proto/api.js';
import { Encounter as EncounterProto } from '../core/proto/common.js';
import { Raid as RaidProto } from '../core/proto/api.js';
import { Party as PartyProto } from '../core/proto/api.js';
import { Class } from '../core/proto/common.js';
import { Profession } from '../core/proto/common.js';
import { Race } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { playerToSpec, specNames } from '../core/proto_utils/utils.js';
import { classColors } from '../core/proto_utils/utils.js';
import { isTankSpec } from '../core/proto_utils/utils.js';
import { specToClass } from '../core/proto_utils/utils.js';
import { newRaidTarget } from '../core/proto_utils/utils.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { camelToSnakeCase } from '../core/utils.js';
import { formatDeltaTextElem } from '../core/utils.js';
import { getEnumValues } from '../core/utils.js';
import { hexToRgba } from '../core/utils.js';

import { RaidSimUI } from './raid_sim_ui.js';
import { playerPresets, specSimFactories } from './presets.js';

import { BalanceDruid_Options as BalanceDruidOptions } from '../core/proto/druid.js';
import { Mage_Options as MageOptions } from '../core/proto/mage.js';
import { SmitePriest_Options as SmitePriestOptions } from '../core/proto/priest.js';
import { MessageType } from '@protobuf-ts/runtime';
import { BaseModal } from '../core/components/base_modal.js';

declare var tippy: any;
declare var $: any;

const NEW_PLAYER: number = -1;

enum DragType {
	None,
	New,
	Move,
	Swap,
	Copy,
}

export class RaidPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly raid: Raid;
	readonly partyPickers: Array<PartyPicker>;
	readonly newPlayerPicker: NewPlayerPicker;

	// Hold data about the player being dragged while the drag is happening.
	currentDragPlayer: Player<any> | null = null;
	currentDragPlayerFromIndex: number = NEW_PLAYER;
	currentDragType: DragType = DragType.New;

	// Hold data about the party being dragged while the drag is happening.
	currentDragParty: PartyPicker | null = null;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI) {
		super(parent, 'raid-picker-root');
		this.raidSimUI = raidSimUI;
		this.raid = raidSimUI.sim.raid;

		this.rootElem.innerHTML = `
			<div class="current-raid-viewer">
				<div class="raid-controls">
				</div>
				<div class="parties-container">
				</div>
			</div>
			<div class="new-player-picker">
			</div>
		`;

		const raidControls = this.rootElem.getElementsByClassName('raid-controls')[0] as HTMLDivElement;
		const activePartiesSelector = new EnumPicker<Raid>(raidControls, this.raidSimUI.sim.raid, {
			label: 'Raid Size',
			extraCssClasses: ['input-inline'],
			labelTooltip: 'Number of players participating in the sim.',
			values: [
				{ name: '5', value: 1 },
				{ name: '10', value: 2 },
				{ name: '25', value: 5 },
				{ name: '40', value: 8 },
			],
			changedEvent: (raid: Raid) => raid.numActivePartiesChangeEmitter,
			getValue: (raid: Raid) => raid.getNumActiveParties(),
			setValue: (eventID: EventID, raid: Raid, newValue: number) => {
				raid.setNumActiveParties(eventID, newValue);
			},
		});

		const partiesContainer = this.rootElem.getElementsByClassName('parties-container')[0] as HTMLDivElement;
		this.partyPickers = this.raid.getParties().map((party, i) => new PartyPicker(partiesContainer, party, i, this));

		const updateActiveParties = () => {
			this.partyPickers.forEach(partyPicker => {
				if (partyPicker.index < this.raidSimUI.sim.raid.getNumActiveParties()) {
					partyPicker.rootElem.classList.add('active');
				} else {
					partyPicker.rootElem.classList.remove('active');
				}
			});
		};
		this.raidSimUI.sim.raid.numActivePartiesChangeEmitter.on(updateActiveParties);
		updateActiveParties();

		const newPlayerPickerRoot = this.rootElem.getElementsByClassName('new-player-picker')[0] as HTMLDivElement;
		this.newPlayerPicker = new NewPlayerPicker(newPlayerPickerRoot, this);

		this.rootElem.ondragend = event => {
			// Uncomment to remove player when dropped 'off' the raid.
			//if (this.currentDragPlayerFromIndex != NEW_PLAYER) {
			//	const playerPicker = this.getPlayerPicker(this.currentDragPlayerFromIndex);
			//	playerPicker.setPlayer(null, null, DragType.None);
			//}

			this.clearDragPlayer();
			this.clearDragParty();
		};
	}

	getCurrentFaction(): Faction {
		return this.raid.sim.getFaction();
	}

	getCurrentPhase(): number {
		return this.raid.sim.getPhase();
	}

	getPlayerPicker(raidIndex: number): PlayerPicker {
		return this.partyPickers[Math.floor(raidIndex / MAX_PARTY_SIZE)].playerPickers[raidIndex % MAX_PARTY_SIZE];
	}

	getPlayerPickers(): Array<PlayerPicker> {
		return [...new Array(25).keys()].map(i => this.getPlayerPicker(i));
	}

	setDragPlayer(player: Player<any>, fromIndex: number, type: DragType) {
		this.clearDragPlayer();

		this.currentDragPlayer = player;
		this.currentDragPlayerFromIndex = fromIndex;
		this.currentDragType = type;

		if (fromIndex != NEW_PLAYER) {
			const playerPicker = this.getPlayerPicker(fromIndex);
			playerPicker.rootElem.classList.add('dragfrom');
		}
	}

	clearDragPlayer() {
		if (this.currentDragPlayerFromIndex != NEW_PLAYER) {
			const playerPicker = this.getPlayerPicker(this.currentDragPlayerFromIndex);
			playerPicker.rootElem.classList.remove('dragfrom');
		}

		this.currentDragPlayer = null;
		this.currentDragPlayerFromIndex = NEW_PLAYER;
		this.currentDragType = DragType.New;
	}

	setDragParty(party: PartyPicker) {
		this.currentDragParty = party;
		party.rootElem.classList.add('dragfrom');
	}
	clearDragParty() {
		if (this.currentDragParty) {
			this.currentDragParty.rootElem.classList.remove('dragfrom');
			this.currentDragParty = null;
		}
	}
}

export class PartyPicker extends Component {
	readonly party: Party;
	readonly index: number;
	readonly raidPicker: RaidPicker;
	readonly playerPickers: Array<PlayerPicker>;

	constructor(parent: HTMLElement, party: Party, index: number, raidPicker: RaidPicker) {
		super(parent, 'party-picker-root');
		this.party = party;
		this.index = index;
		this.raidPicker = raidPicker;

		this.rootElem.innerHTML = `
			<div class="party-header">
				<span class="party-label" draggable="true">Group ${index + 1}</span>
				<div class="party-results">
					<span class="party-results-dps"></span>
					<span class="party-results-reference-delta"></span>
				</div>
			</div>
			<div class="players-container">
			</div>
		`;

		const playersContainer = this.rootElem.getElementsByClassName('players-container')[0] as HTMLDivElement;
		this.playerPickers = [...Array(MAX_PARTY_SIZE).keys()].map(i => new PlayerPicker(playersContainer, this, i));

		const dpsResultElem = this.rootElem.getElementsByClassName('party-results-dps')[0] as HTMLElement;
		const referenceDeltaElem = this.rootElem.getElementsByClassName('party-results-reference-delta')[0] as HTMLElement;

		this.raidPicker.raidSimUI.referenceChangeEmitter.on(() => {
			const currentData = this.raidPicker.raidSimUI.getCurrentData();
			const referenceData = this.raidPicker.raidSimUI.getReferenceData();

			const partyDps = currentData?.simResult.raidMetrics.parties[this.index]?.dps.avg || 0;
			const referenceDps = referenceData?.simResult.raidMetrics.parties[this.index]?.dps.avg || 0;

			if (partyDps == 0 && referenceDps == 0) {
				dpsResultElem.textContent = '';
				referenceDeltaElem.textContent = '';
				return;
			}

			dpsResultElem.textContent = partyDps.toFixed(1);

			if (!referenceData) {
				referenceDeltaElem.textContent = '';
				return;
			}

			formatDeltaTextElem(referenceDeltaElem, referenceDps, partyDps, 1);
		});

		const dragStart = (event: DragEvent, type: DragType) => {
			event.dataTransfer!.dropEffect = 'move';
			event.dataTransfer!.effectAllowed = 'all';
			this.raidPicker.setDragParty(this);
		};
		const labelElem = this.rootElem.getElementsByClassName('party-label')[0] as HTMLElement;
		labelElem.ondragstart = event => {
			dragStart(event, DragType.Swap);
		};

		let dragEnterCounter = 0;
		this.rootElem.ondragenter = event => {
			event.preventDefault();
			if (!this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter++;
			this.rootElem.classList.add('dragto');
		};
		this.rootElem.ondragleave = event => {
			event.preventDefault();
			if (!this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter--;
			if (dragEnterCounter <= 0) {
				this.rootElem.classList.remove('dragto');
			}
		};
		this.rootElem.ondragover = event => {
			event.preventDefault();
		};
		this.rootElem.ondrop = event => {
			if (!this.raidPicker.currentDragParty) {
				return;
			}

			event.preventDefault();
			dragEnterCounter = 0;
			this.rootElem.classList.remove('dragto');

			const eventID = TypedEvent.nextEventID();
			TypedEvent.freezeAllAndDo(() => {
				const srcPartyPicker = this.raidPicker.currentDragParty!;

				for (let i = 0; i < MAX_PARTY_SIZE; i++) {
					const srcPlayerPicker = srcPartyPicker.playerPickers[i]!;
					const dstPlayerPicker = this.playerPickers[i]!;

					const srcPlayer = srcPlayerPicker.player;
					const dstPlayer = dstPlayerPicker.player;

					srcPlayerPicker.setPlayer(eventID, dstPlayer, DragType.Swap);
					dstPlayerPicker.setPlayer(eventID, srcPlayer, DragType.Swap);
				}
			});

			this.raidPicker.clearDragParty();
		};
	}
}

export class PlayerPicker extends Component {
	// Index of this player within its party (0-4).
	readonly index: number;

	// Index of this player within the whole raid (0-24).
	readonly raidIndex: number;

	player: Player<any> | null;

	readonly partyPicker: PartyPicker;
	readonly raidPicker: RaidPicker;

	private readonly labelElem: HTMLElement;
	private readonly iconElem: HTMLImageElement;
	private readonly nameElem: HTMLSpanElement;
	private readonly resultsElem: HTMLElement;
	private readonly dpsResultElem: HTMLElement;
	private readonly referenceDeltaElem: HTMLElement;

	constructor(parent: HTMLElement, partyPicker: PartyPicker, index: number) {
		super(parent, 'player-picker-root');
		this.index = index;
		this.raidIndex = partyPicker.index * MAX_PARTY_SIZE + index;
		this.player = null;
		this.partyPicker = partyPicker;
		this.raidPicker = partyPicker.raidPicker;

		this.partyPicker.party.compChangeEmitter.on(eventID => {
			const newPlayer = this.partyPicker.party.getPlayer(this.index);
			if (newPlayer != this.player) {
				this.setPlayer(eventID, newPlayer, DragType.None);
			}
		});

		this.rootElem.innerHTML = `
			<div class="player-label">
				<img class="player-icon"></img>
				<span class="player-name" contenteditable></span>
			</div>
			<div class="player-spacer">
			</div>
			<div class="player-options">
				<span class="player-edit fa fa-edit"></span>
				<span class="player-copy fa fa-copy" draggable="true"></span>
				<span class="player-delete fa fa-times"></span>
			</div>
			<div class="player-results">
				<span class="player-results-dps"></span>
				<span class="player-results-reference-delta"></span>
			</div>
		`;

		this.labelElem = this.rootElem.getElementsByClassName('player-label')[0] as HTMLElement;
		this.iconElem = this.rootElem.getElementsByClassName('player-icon')[0] as HTMLImageElement;
		this.nameElem = this.rootElem.getElementsByClassName('player-name')[0] as HTMLSpanElement;
		this.resultsElem = this.rootElem.getElementsByClassName('player-results')[0] as HTMLElement;
		this.dpsResultElem = this.rootElem.getElementsByClassName('player-results-dps')[0] as HTMLElement;
		this.referenceDeltaElem = this.rootElem.getElementsByClassName('player-results-reference-delta')[0] as HTMLElement;

		this.nameElem.addEventListener('input', event => {
			if (this.player) {
				this.player.setName(TypedEvent.nextEventID(), this.nameElem.textContent || '');
			}
		});

		const maxLength = 15;
		this.nameElem.addEventListener('keydown', event => {
			// 9 is tab, 13 is enter
			if (event.keyCode == 9 || event.keyCode == 13) {
				event.preventDefault();
				const realPlayerPickers = this.raidPicker.getPlayerPickers().filter(pp => pp.player != null);
				const indexOfThis = realPlayerPickers.indexOf(this);
				if (indexOfThis != -1 && realPlayerPickers.length > indexOfThis + 1) {
					realPlayerPickers[indexOfThis + 1].nameElem.focus();
				} else {
					this.nameElem.blur();
				}
			}

			// escape
			if (event.keyCode == 27) {
				this.nameElem.blur();
			}

			// 8 is backspace, 46 is delete, 
			if ((event.keyCode != 8 && event.keyCode != 46) && (this.nameElem.textContent?.length || 0) >= maxLength) {
				event.preventDefault();
			}
		});

		const emptyName = 'Unnamed';
		this.nameElem.addEventListener('focusin', event => {
			const selection = window.getSelection();
			if (selection) {
				const range = document.createRange();
				range.selectNodeContents(this.nameElem);
				selection.removeAllRanges();
				selection.addRange(range);
			}
		});
		this.nameElem.addEventListener('focusout', event => {
			if (!this.nameElem.textContent) {
				this.nameElem.textContent = emptyName;
				if (this.player) {
					this.player.setName(TypedEvent.nextEventID(), emptyName);
				}
			}
		});

		const dragStart = (event: DragEvent, type: DragType) => {
			if (this.player == null) {
				event.preventDefault();
				return;
			}
			event.dataTransfer!.dropEffect = 'move';
			event.dataTransfer!.effectAllowed = 'all';
			event.dataTransfer!.setDragImage(this.iconElem, 30, 30);
			if (this.player) {
				var playerDataProto = this.player.toProto(true);
				event.dataTransfer!.setData("text/plain", btoa(String.fromCharCode(...PlayerProto.toBinary(playerDataProto))));
			}
			this.raidPicker.setDragPlayer(this.player, this.raidIndex, type);
		};

		this.labelElem.ondragstart = event => {
			dragStart(event, DragType.Swap);
		};
		this.resultsElem.ondragstart = event => {
			dragStart(event, DragType.Swap);
		};

		const copyElem = this.rootElem.getElementsByClassName('player-copy')[0] as HTMLSpanElement;
		tippy(copyElem, {
			'content': 'Drag to Copy',
			'allowHTML': true,
		});
		copyElem.ondragstart = event => {
			dragStart(event, DragType.Copy);
		};

		const deleteElem = this.rootElem.getElementsByClassName('player-delete')[0] as HTMLSpanElement;
		tippy(deleteElem, {
			'content': 'Click to Delete',
			'allowHTML': true,
		});
		deleteElem.addEventListener('click', event => {
			this.setPlayer(TypedEvent.nextEventID(), null, DragType.None);
		});

		let dragEnterCounter = 0;
		this.rootElem.ondragenter = event => {
			event.preventDefault();
			if (this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter++;
			this.rootElem.classList.add('dragto');
		};
		this.rootElem.ondragleave = event => {
			event.preventDefault();
			if (this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter--;
			if (dragEnterCounter <= 0) {
				this.rootElem.classList.remove('dragto');
			}
		};
		this.rootElem.ondragover = event => {
			event.preventDefault();
		};
		this.rootElem.ondrop = event => {
			if (this.raidPicker.currentDragParty) {
				return;
			}
			var dropData = event.dataTransfer!.getData("text/plain");

			event.preventDefault();
			dragEnterCounter = 0;
			this.rootElem.classList.remove('dragto');

			const eventID = TypedEvent.nextEventID();
			TypedEvent.freezeAllAndDo(() => {
				if (this.raidPicker.currentDragPlayer == null && dropData.length == 0) {
					return;
				}

				if (this.raidPicker.currentDragPlayerFromIndex == this.raidIndex) {
					this.raidPicker.clearDragPlayer();
					return;
				}

				const dragType = this.raidPicker.currentDragType;

				if (this.raidPicker.currentDragPlayerFromIndex != NEW_PLAYER) {
					const fromPlayerPicker = this.raidPicker.getPlayerPicker(this.raidPicker.currentDragPlayerFromIndex);
					if (dragType == DragType.Swap) {
						fromPlayerPicker.setPlayer(eventID, this.player, dragType);
					} else if (dragType == DragType.Move) {
						fromPlayerPicker.setPlayer(eventID, null, dragType);
					}
				} else if (this.raidPicker.currentDragPlayer == null) {
					// This would be a copy from another window.
					const binary = atob(dropData);
					const bytes = new Uint8Array(binary.length);
					for (let i = 0; i < bytes.length; i++) {
						bytes[i] = binary.charCodeAt(i);
					}
					const playerProto = PlayerProto.fromBinary(bytes);

					var localPlayer = new Player(playerToSpec(playerProto), this.raidPicker.raidSimUI.sim);
					localPlayer.fromProto(eventID, playerProto);
					this.raidPicker.currentDragPlayer = localPlayer;
				}

				if (dragType == DragType.Copy) {
					this.setPlayer(eventID, this.raidPicker.currentDragPlayer!.clone(eventID), dragType);
				} else {
					this.setPlayer(eventID, this.raidPicker.currentDragPlayer, dragType);
				}

				this.raidPicker.clearDragPlayer();
			});
		};

		const editElem = this.rootElem.getElementsByClassName('player-edit')[0] as HTMLSpanElement;
		tippy(editElem, {
			'content': 'Edit',
			'allowHTML': true,
		});
		editElem.addEventListener('click', event => {
			if (this.player) {
				new PlayerEditorModal(this.player);
			}
		});

		this.raidPicker.raidSimUI.referenceChangeEmitter.on(() => {
			const currentData = this.raidPicker.raidSimUI.getCurrentData();
			const referenceData = this.raidPicker.raidSimUI.getReferenceData();

			const playerDps = currentData?.simResult.getPlayerWithRaidIndex(this.raidIndex)?.dps.avg || 0;
			const referenceDps = referenceData?.simResult.getPlayerWithRaidIndex(this.raidIndex)?.dps.avg || 0;

			if (playerDps == 0 && referenceDps == 0) {
				this.dpsResultElem.textContent = '';
				this.referenceDeltaElem.textContent = '';
				return;
			}

			this.dpsResultElem.textContent = playerDps.toFixed(1);

			if (!referenceData) {
				this.referenceDeltaElem.textContent = '';
				return;
			}

			formatDeltaTextElem(this.referenceDeltaElem, referenceDps, playerDps, 1);
		});

		this.update();
	}

	setPlayer(eventID: EventID, newPlayer: Player<any> | null, dragType: DragType) {
		if (newPlayer == this.player) {
			return;
		}

		this.dpsResultElem.textContent = '';
		this.referenceDeltaElem.textContent = '';

		TypedEvent.freezeAllAndDo(() => {
			this.player = newPlayer;
			if (newPlayer) {
				this.partyPicker.party.setPlayer(eventID, this.index, newPlayer);

				if (dragType == DragType.New) {
					applyNewPlayerAssignments(eventID, newPlayer, this.raidPicker.raid);
				}
			} else {
				this.partyPicker.party.setPlayer(eventID, this.index, newPlayer);
				this.partyPicker.party.compChangeEmitter.emit(eventID);
			}
		});

		this.update();
	}

	private update() {
		if (this.player == null) {
			this.rootElem.classList.add('empty');
			this.rootElem.style.backgroundColor = 'black';
			this.labelElem.setAttribute('draggable', 'false');
			this.resultsElem.setAttribute('draggable', 'false');
			this.nameElem.textContent = '';
			this.nameElem.removeAttribute('contenteditable');
		} else {
			this.rootElem.classList.remove('empty');
			this.rootElem.style.backgroundColor = this.player.getClassColor();
			this.labelElem.setAttribute('draggable', 'true');
			this.resultsElem.setAttribute('draggable', 'true');
			this.nameElem.textContent = this.player.getName();
			this.nameElem.setAttribute('contenteditable', '');
			this.iconElem.src = this.player.getTalentTreeIcon();
		}
	}
}

class PlayerEditorModal extends BaseModal {
	constructor(player: Player<any>) {
		super('player-editor-modal', {
			closeButton: {fixed: true, text: false},
			header: false
		});

		this.rootElem.id = 'playerEditorModal';
		this.body.insertAdjacentHTML('beforeend', `
			<div class="player-editor within-raid-sim"></div>
		`);

		const editorRoot = this.rootElem.getElementsByClassName('player-editor')[0] as HTMLElement;
		const individualSim = specSimFactories[player.spec]!(editorRoot, player);
	}
}

class NewPlayerPicker extends Component {
	readonly raidPicker: RaidPicker;

	constructor(parent: HTMLElement, raidPicker: RaidPicker) {
		super(parent, 'new-player-picker-root');
		this.raidPicker = raidPicker;

		this.rootElem.innerHTML = `
			<div class="new-player-picker-controls">
				<div class="faction-selector"></div>
				<div class="phase-selector"></div>
			</div>
			<div class="presets-container"></div>
		`;

		const factionSelector = new EnumPicker<NewPlayerPicker>(this.rootElem.getElementsByClassName('faction-selector')[0] as HTMLElement, this, {
			label: 'Faction',
			labelTooltip: 'Default faction for newly-created players.',
			values: [
				{ name: 'Alliance', value: Faction.Alliance },
				{ name: 'Horde', value: Faction.Horde },
			],
			changedEvent: (picker: NewPlayerPicker) => this.raidPicker.raid.sim.factionChangeEmitter,
			getValue: (picker: NewPlayerPicker) => this.raidPicker.raid.sim.getFaction(),
			setValue: (eventID: EventID, picker: NewPlayerPicker, newValue: Faction) => {
				this.raidPicker.raid.sim.setFaction(eventID, newValue);
			},
		});

		const phaseSelector = new EnumPicker<NewPlayerPicker>(this.rootElem.getElementsByClassName('phase-selector')[0] as HTMLElement, this, {
			label: 'Phase',
			labelTooltip: 'Newly-created players will start with approximate BIS gear from this phase.',
			values: [
				{ name: '1', value: 1 },
				// Presets aren't filled for most roles so disable these options for now.
				//{ name: '2', value: 2 },
				//{ name: '3', value: 3 },
				//{ name: '4', value: 4 },
				//{ name: '5', value: 5 },
			],
			changedEvent: (picker: NewPlayerPicker) => this.raidPicker.raid.sim.phaseChangeEmitter,
			getValue: (picker: NewPlayerPicker) => this.raidPicker.raid.sim.getPhase(),
			setValue: (eventID: EventID, picker: NewPlayerPicker, newValue: number) => {
				this.raidPicker.raid.sim.setPhase(eventID, newValue);
			},
		});

		const presetsContainer = this.rootElem.getElementsByClassName('presets-container')[0] as HTMLElement;
		getEnumValues(Class).forEach(wowClass => {
			if (wowClass == Class.ClassUnknown) {
				return;
			}

			const matchingPresets = playerPresets.filter(preset => specToClass[preset.spec] == wowClass);
			if (matchingPresets.length == 0) {
				return;
			}

			const classPresetsContainer = document.createElement('div');
			classPresetsContainer.classList.add('class-presets-container');
			presetsContainer.appendChild(classPresetsContainer);
			classPresetsContainer.style.backgroundColor = hexToRgba(classColors[wowClass as Class], 0.5);

			matchingPresets.forEach(matchingPreset => {
				const presetElem = document.createElement('div');
				presetElem.classList.add('preset-picker');
				classPresetsContainer.appendChild(presetElem);

				const presetIconElem = document.createElement('img');
				presetIconElem.classList.add('preset-picker-icon');
				presetElem.appendChild(presetIconElem);
				presetIconElem.src = matchingPreset.iconUrl;
				tippy(presetIconElem, {
					'content': matchingPreset.tooltip,
					'allowHTML': true,
				});

				presetElem.setAttribute('draggable', 'true');
				presetElem.ondragstart = event => {
					const eventID = TypedEvent.nextEventID();
					TypedEvent.freezeAllAndDo(() => {
						const dragImage = new Image();
						dragImage.src = matchingPreset.iconUrl;
						event.dataTransfer!.setDragImage(dragImage, 30, 30);
						event.dataTransfer!.setData("text/plain", "");
						event.dataTransfer!.dropEffect = 'copy';


						const newPlayer = new Player(matchingPreset.spec, this.raidPicker.raid.sim);
						newPlayer.applySharedDefaults(eventID);
						newPlayer.setRace(eventID, matchingPreset.defaultFactionRaces[this.raidPicker.getCurrentFaction()]);
						newPlayer.setRotation(eventID, matchingPreset.rotation);
						newPlayer.setTalentsString(eventID, matchingPreset.talents.talentsString);
						newPlayer.setGlyphs(eventID, matchingPreset.talents.glyphs || Glyphs.create());
						newPlayer.setSpecOptions(eventID, matchingPreset.specOptions);
						newPlayer.setConsumes(eventID, matchingPreset.consumes);
						newPlayer.setName(eventID, matchingPreset.defaultName);
						newPlayer.setProfession1(eventID, matchingPreset.otherDefaults?.profession1 || Profession.Engineering);
						newPlayer.setProfession2(eventID, matchingPreset.otherDefaults?.profession2 || Profession.Jewelcrafting);
						newPlayer.setDistanceFromTarget(eventID, matchingPreset.otherDefaults?.distanceFromTarget || 0);

						// Need to wait because the gear might not be loaded yet.
						this.raidPicker.raid.sim.waitForInit().then(() => {
							newPlayer.setGear(
								eventID,
								this.raidPicker.raid.sim.db.lookupEquipmentSpec(
									matchingPreset.defaultGear[this.raidPicker.getCurrentFaction()][this.raidPicker.getCurrentPhase()]));
						});

						this.raidPicker.setDragPlayer(newPlayer, NEW_PLAYER, DragType.New);
					});
				};
			});
		});
	}
}

function applyNewPlayerAssignments(eventID: EventID, newPlayer: Player<any>, raid: Raid) {
	if (isTankSpec(newPlayer.spec)) {
		const tanks = raid.getTanks();
		const emptyIdx = tanks.findIndex(tank => raid.getPlayerFromRaidTarget(tank) == null);
		if (emptyIdx == -1) {
			if (tanks.length < 3) {
				raid.setTanks(eventID, tanks.concat([newPlayer.makeRaidTarget()]));
			}
		} else {
			tanks[emptyIdx] = newPlayer.makeRaidTarget();
			raid.setTanks(eventID, tanks);
		}
	}

	// Spec-specific assignments. For most cases, default to buffing self.
	if (newPlayer.spec == Spec.SpecBalanceDruid) {
		const newOptions = newPlayer.getSpecOptions() as BalanceDruidOptions;
		newOptions.innervateTarget = newRaidTarget(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	} else if (newPlayer.spec == Spec.SpecSmitePriest) {
		const newOptions = newPlayer.getSpecOptions() as SmitePriestOptions;
		newOptions.powerInfusionTarget = newRaidTarget(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	} else if (newPlayer.spec == Spec.SpecMage) {
		const newOptions = newPlayer.getSpecOptions() as MageOptions;
		newOptions.focusMagicTarget = newRaidTarget(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	}
}
