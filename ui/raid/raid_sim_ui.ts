import { BooleanPicker } from "../core/components/boolean_picker.js";
import { DetailedResults } from "../core/components/detailed_results.js";
import { EncounterPicker } from "../core/components/encounter_picker.js";
import { LogRunner } from "../core/components/log_runner.js";
import { addRaidSimAction, RaidSimResultsManager, ReferenceData } from "../core/components/raid_sim_action.js";
import { SavedDataManager } from "../core/components/saved_data_manager.js";
import { SettingsMenu } from "../core/components/settings_menu.js";

import * as Tooltips from "../core/constants/tooltips.js";
import { Encounter } from "../core/encounter.js";
import { Player } from "../core/player.js";
import { Raid as RaidProto } from "../core/proto/api.js";
import { Class, Encounter as EncounterProto, Faction, Stat, TristateEffect } from "../core/proto/common.js";
import { Blessings } from "../core/proto/paladin.js";
import { BlessingsAssignments, BuffBot as BuffBotProto, RaidSimSettings, SavedEncounter, SavedRaid } from "../core/proto/ui.js";
import { playerToSpec } from "../core/proto_utils/utils.js";
import { Raid } from "../core/raid.js";
import { Sim } from "../core/sim.js";
import { SimUI } from "../core/sim_ui.js";
import { LaunchStatus, raidSimLaunched } from '../core/launched_sims.js';
import { EventID, TypedEvent } from "../core/typed_event.js";

import { AssignmentsPicker } from "./assignments_picker.js";
import { BlessingsPicker } from "./blessings_picker.js";
import { BuffBot } from "./buff_bot.js";
import { newRaidExporters, newRaidImporters } from "./import_export.js";
import { implementedSpecs } from "./presets.js";
import { RaidPicker } from "./raid_picker.js";
import { TanksPicker } from "./tanks_picker.js";

declare var Muuri: any;
declare var tippy: any;
declare var pako: any;

export interface RaidSimConfig {
	knownIssues?: Array<string>,
}

const extraKnownIssues: Array<string> = [
	//'We\'re still missing implementations for many specs. If you\'d like to help us out, check out our <a href="https://github.com/wowsims/wotlk">Github project</a> or <a href="https://discord.gg/jJMPr9JWwx">join our discord</a>!',
];

export class RaidSimUI extends SimUI {
	private readonly config: RaidSimConfig;
	private raidSimResultsManager: RaidSimResultsManager | null = null;
	public raidPicker: RaidPicker | null = null;
	private blessingsPicker: BlessingsPicker | null = null;

	// Emits when the raid comp changes. Includes changes to buff bots.
	readonly compChangeEmitter = new TypedEvent<void>();
	readonly changeEmitter = new TypedEvent<void>();

	readonly referenceChangeEmitter = new TypedEvent<void>();

	private settingsMuuri: any;

	constructor(parentElem: HTMLElement, config: RaidSimConfig) {
		super(parentElem, new Sim(), {
			spec: null,
			launchStatus: raidSimLaunched ? LaunchStatus.Launched : LaunchStatus.Unlaunched,
			knownIssues: (config.knownIssues || []).concat(extraKnownIssues),
		});
		this.rootElem.classList.add('raid-sim-ui');

		this.config = config;
		this.settingsMuuri = null;

		this.sim.raid.compChangeEmitter.on(eventID => this.compChangeEmitter.emit(eventID));
		this.sim.setModifyRaidProto(raidProto => this.modifyRaidProto(raidProto));

		[
			this.compChangeEmitter,
			this.sim.changeEmitter,
		].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));

		this.sim.waitForInit().then(() => this.loadSettings());

		this.addSidebarComponents();
		this.addTopbarComponents();
		this.addRaidTab();
		this.addSettingsTab();
		this.addDetailedResultsTab();
		this.addLogTab();

		this.changeEmitter.on(() => this.recomputeSettingsLayout());
	}

	private loadSettings() {
		const initEventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			let loadedSettings = false;

			const savedSettings = window.localStorage.getItem(this.getSettingsStorageKey());
			if (savedSettings != null) {
				try {
					const settings = RaidSimSettings.fromJsonString(savedSettings);
					this.fromProto(initEventID, settings);
					loadedSettings = true;
				} catch (e) {
					console.warn('Failed to parse saved settings: ' + e);
				}
			}

			if (!loadedSettings) {
				this.applyDefaults(initEventID);
			}

			// This needs to go last so it doesn't re-store things as they are initialized.
			this.changeEmitter.on(eventID => {
				const jsonStr = RaidSimSettings.toJsonString(this.toProto());
				window.localStorage.setItem(this.getSettingsStorageKey(), jsonStr);
			});
		});
	}

	private addSidebarComponents() {
		this.raidSimResultsManager = addRaidSimAction(this);
		this.raidSimResultsManager.changeEmitter.on(eventID => this.referenceChangeEmitter.emit(eventID));
	}

	private addTopbarComponents() {
		this.addImportLink(newRaidImporters(this));
		this.addExportLink(newRaidExporters(this));
	}

	private addRaidTab() {
		this.addTab('RAID', 'raid-tab', `
			<div class="raid-picker">
			</div>
			<div class="saved-raids-div">
				<div class="saved-raids-manager">
				</div>
			</div>
		`);

		this.raidPicker = new RaidPicker(this.rootElem.getElementsByClassName('raid-picker')[0] as HTMLElement, this);

		const savedRaidManager = new SavedDataManager<RaidSimUI, SavedRaid>(this.rootElem.getElementsByClassName('saved-raids-manager')[0] as HTMLElement, this, {
			label: 'Raid',
			storageKey: this.getSavedRaidStorageKey(),
			getData: (raidSimUI: RaidSimUI) => SavedRaid.create({
				raid: this.sim.raid.toProto(),
				buffBots: this.getBuffBots().map(b => b.toProto()),
				blessings: this.blessingsPicker!.getAssignments(),
				faction: this.sim.getFaction(),
				phase: this.sim.getPhase(),
			}),
			setData: (eventID: EventID, raidSimUI: RaidSimUI, newRaid: SavedRaid) => {
				TypedEvent.freezeAllAndDo(() => {
					this.sim.raid.fromProto(eventID, newRaid.raid || RaidProto.create());
					this.raidPicker!.setBuffBots(eventID, newRaid.buffBots);
					this.blessingsPicker!.setAssignments(eventID, newRaid.blessings || BlessingsAssignments.create());
					if (newRaid.faction) this.sim.setFaction(eventID, newRaid.faction);
					if (newRaid.phase) this.sim.setPhase(eventID, newRaid.phase);
				});
			},
			changeEmitters: [this.changeEmitter, this.sim.changeEmitter],
			equals: (a: SavedRaid, b: SavedRaid) => {
				return SavedRaid.equals(a, b);
			},
			toJson: (a: SavedRaid) => SavedRaid.toJson(a),
			fromJson: (obj: any) => SavedRaid.fromJson(obj),
		});
		this.sim.waitForInit().then(() => {
			savedRaidManager.loadUserData();
		});
	}

	private addSettingsTab() {
		this.addTab('SETTINGS', 'raid-settings-tab', `
			<div class="raid-settings-sections">
				<div class="settings-section-container raid-settings-section-container">
					<fieldset class="settings-section raid-encounter-section">
						<legend>Encounter</legend>
					</fieldset>
				</div>
				<div class="settings-section-container blessings-section-container">
					<fieldset class="settings-section blessings-section">
						<legend>Blessings</legend>
					</fieldset>
				</div>
				<div class="settings-section-container assignments-section-container">
				</div>
				<div class="settings-section-container tanks-section-container">
				</div>
				<div class="settings-section-container raid-settings-section-container">
					<fieldset class="settings-section other-options-section">
						<legend>Other Options</legend>
					</fieldset>
				</div>
			</div>
			<div class="settings-bottom-bar">
				<div class="saved-encounter-manager">
				</div>
			</div>
		`);

		const encounterSectionElem = this.rootElem.getElementsByClassName('raid-encounter-section')[0] as HTMLElement;
		new EncounterPicker(encounterSectionElem, this.sim.encounter, {
			showExecuteProportion: true,
		}, this);
		const savedEncounterManager = new SavedDataManager<Encounter, SavedEncounter>(this.rootElem.getElementsByClassName('saved-encounter-manager')[0] as HTMLElement, this.sim.encounter, {
			label: 'Encounter',
			storageKey: this.getSavedEncounterStorageKey(),
			getData: (encounter: Encounter) => SavedEncounter.create({ encounter: encounter.toProto() }),
			setData: (eventID: EventID, encounter: Encounter, newEncounter: SavedEncounter) => encounter.fromProto(eventID, newEncounter.encounter!),
			changeEmitters: [this.sim.encounter.changeEmitter],
			equals: (a: SavedEncounter, b: SavedEncounter) => SavedEncounter.equals(a, b),
			toJson: (a: SavedEncounter) => SavedEncounter.toJson(a),
			fromJson: (obj: any) => SavedEncounter.fromJson(obj),
		});
		this.sim.waitForInit().then(() => {
			savedEncounterManager.loadUserData();
		});

		const blessingsSection = this.rootElem.getElementsByClassName('blessings-section')[0] as HTMLElement;
		this.blessingsPicker = new BlessingsPicker(blessingsSection, this);
		this.blessingsPicker.changeEmitter.on(eventID => this.changeEmitter.emit(eventID));
		tippy(blessingsSection, {
			content: Tooltips.BLESSINGS_SECTION,
			allowHTML: true,
			placement: 'left',
		});

		const assignmentsPicker = new AssignmentsPicker(this.rootElem.getElementsByClassName('assignments-section-container')[0] as HTMLElement, this);
		const tanksPicker = new TanksPicker(this.rootElem.getElementsByClassName('tanks-section-container')[0] as HTMLElement, this);

		const otherOptionsSectionElem = this.rootElem.getElementsByClassName('other-options-section')[0] as HTMLElement;
		//new BooleanPicker(otherOptionsSectionElem, this.sim.raid, {
		//	label: 'Stagger Stormstrikes',
		//	labelTooltip: 'When there are multiple Enhancement Shaman in the raid, causes them to coordinate their Stormstrike casts for optimal SS charge usage.',
		//	changedEvent: (raid: Raid) => raid.staggerStormstrikesChangeEmitter,
		//	getValue: (raid: Raid) => raid.getStaggerStormstrikes(),
		//	setValue: (eventID: EventID, raid: Raid, newValue: boolean) => {
		//		raid.setStaggerStormstrikes(eventID, newValue);
		//	},
		//});

		// Init Muuri layout only when settings tab is clicked, because it needs the elements
		// to be shown so it can calculate sizes.
		(this.rootElem.getElementsByClassName('raid-settings-tab-tab')[0] as HTMLElement)!.addEventListener('click', event => {
			if (this.settingsMuuri == null) {
				setTimeout(() => {
					this.settingsMuuri = new Muuri('.raid-settings-sections');
				}, 200); // Magic amount of time before Muuri init seems to work
			}

			setTimeout(() => {
				this.recomputeSettingsLayout();
			}, 200);
		});
	}

	private addDetailedResultsTab() {
		this.addTab('DETAILED RESULTS', 'detailed-results-tab', `
			<div class="detailed-results">
			</div>
		`);

		const detailedResults = new DetailedResults(this.rootElem.getElementsByClassName('detailed-results')[0] as HTMLElement, this, this.raidSimResultsManager!);
	}

	private addLogTab() {
		this.addTab('LOG', 'log-tab', `
			<div class="log-runner">
			</div>
		`);

		const logRunner = new LogRunner(this.rootElem.getElementsByClassName('log-runner')[0] as HTMLElement, this);
	}

	private recomputeSettingsLayout() {
		if (this.settingsMuuri) {
			//this.settingsMuuri.refreshItems();
		}
		window.dispatchEvent(new Event('resize'));
	}

	private modifyRaidProto(raidProto: RaidProto) {
		// Invoke all the buff bot callbacks.
		this.getBuffBots().forEach(buffBot => {
			const partyProto = raidProto.parties[buffBot.getPartyIndex()];
			if (!partyProto) {
				throw new Error('No party proto for party index: ' + buffBot.getPartyIndex());
			}
			buffBot.settings.modifyRaidProto(buffBot, raidProto, partyProto);
		});

		// Apply blessings.
		const numPaladins = this.getClassCount(Class.ClassPaladin);
		const blessingsAssignments = this.blessingsPicker!.getAssignments();
		implementedSpecs.forEach(spec => {
			const playerProtos = raidProto.parties
				.map(party => party.players.filter(player => player.class != Class.ClassUnknown && playerToSpec(player) == spec))
				.flat();

			blessingsAssignments.paladins.forEach((paladin, i) => {
				if (i >= numPaladins) {
					return;
				}

				if (paladin.blessings[spec] == Blessings.BlessingOfKings) {
					playerProtos.forEach(playerProto => playerProto.buffs!.blessingOfKings = true);
				} else if (paladin.blessings[spec] == Blessings.BlessingOfMight) {
					playerProtos.forEach(playerProto => playerProto.buffs!.blessingOfMight = TristateEffect.TristateEffectImproved);
				} else if (paladin.blessings[spec] == Blessings.BlessingOfWisdom) {
					playerProtos.forEach(playerProto => playerProto.buffs!.blessingOfWisdom = TristateEffect.TristateEffectImproved);
				} else if (paladin.blessings[spec] == Blessings.BlessingOfSanctuary) {
					playerProtos.forEach(playerProto => playerProto.buffs!.blessingOfSanctuary = true);
				}
			});
		});
	}

	getCurrentData(): ReferenceData | null {
		if (this.raidSimResultsManager) {
			return this.raidSimResultsManager.getCurrentData();
		} else {
			return null;
		}
	}

	getReferenceData(): ReferenceData | null {
		if (this.raidSimResultsManager) {
			return this.raidSimResultsManager.getReferenceData();
		} else {
			return null;
		}
	}

	getClassCount(playerClass: Class): number {
		return this.sim.raid.getClassCount(playerClass)
			+ this.getBuffBots()
				.filter(buffBot => buffBot.getClass() == playerClass).length;
	}

	getBuffBots(): Array<BuffBot> {
		return this.raidPicker!.getBuffBots();
	}

	setBuffBots(eventID: EventID, buffBotProtos: BuffBotProto[]): void {
		this.raidPicker!.setBuffBots(eventID, buffBotProtos);
	}

	clearBuffBots(eventID: EventID): void {
		this.raidPicker!.setBuffBots(eventID, []);
	}

	getPlayersAndBuffBots(): Array<Player<any> | BuffBot | null> {
		const players = this.sim.raid.getPlayers();
		const buffBots = this.getBuffBots();

		const playersAndBuffBots: Array<Player<any> | BuffBot | null> = players.slice();
		buffBots.forEach(buffBot => {
			playersAndBuffBots[buffBot.getRaidIndex()] = buffBot;
		});

		return playersAndBuffBots;
	}

	applyDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			this.sim.raid.fromProto(eventID, RaidProto.create());
			this.sim.encounter.applyDefaults(eventID);
			this.sim.applyDefaults(eventID, true, true);
			this.sim.setShowDamageMetrics(eventID, true);
		});
	}

	toProto(): RaidSimSettings {
		return RaidSimSettings.create({
			settings: this.sim.toProto(),
			raid: this.sim.raid.toProto(true),
			buffBots: this.getBuffBots().map(b => b.toProto()),
			blessings: this.blessingsPicker!.getAssignments(),
			encounter: this.sim.encounter.toProto(),
		});
	}

	toLink(): string {
		const proto = this.toProto();
		// When sharing links, people generally don't intend to share settings.
		proto.settings = undefined;

		const protoBytes = RaidSimSettings.toBinary(proto);
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		const encoded = btoa(String.fromCharCode(...deflated));

		const linkUrl = new URL(window.location.href);
		linkUrl.hash = encoded;
		return linkUrl.toString();
	}

	fromProto(eventID: EventID, settings: RaidSimSettings) {
		TypedEvent.freezeAllAndDo(() => {
			if (settings.settings) {
				this.sim.fromProto(eventID, settings.settings);
			}
			this.sim.raid.fromProto(eventID, settings.raid || RaidProto.create());
			this.sim.encounter.fromProto(eventID, settings.encounter || EncounterProto.create());
			this.raidPicker!.setBuffBots(eventID, settings.buffBots);
			this.blessingsPicker!.setAssignments(eventID, settings.blessings || BlessingsAssignments.create());
		});
	}

	clearRaid(eventID: EventID) {
		this.sim.raid.clear(eventID);
		this.clearBuffBots(eventID);
	}

	// Returns the actual key to use for local storage, based on the given key part and the site context.
	getStorageKey(keyPart: string): string {
		return '__wotlk_raid__' + keyPart;
	}

	getSavedRaidStorageKey(): string {
		return this.getStorageKey('__savedRaid__');
	}
}
