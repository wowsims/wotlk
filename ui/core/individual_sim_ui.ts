import { aplLaunchStatuses, LaunchStatus, simLaunchStatuses } from './launched_sims';
import { Player } from './player';
import { SimUI, SimWarning } from './sim_ui';
import { EventID, TypedEvent } from './typed_event';

import { CharacterStats, StatMods } from './components/character_stats';
import { ContentBlock } from './components/content_block';
import { EmbeddedDetailedResults } from './components/detailed_results';
import { EncounterPickerConfig } from './components/encounter_picker';
import { EnumPicker } from './components/enum_picker';
import { IconEnumPicker } from './components/icon_enum_picker';
import { LogRunner } from './components/log_runner';
import { addRaidSimAction, RaidSimResultsManager } from './components/raid_sim_action';
import { SavedDataConfig, SavedDataManager } from './components/saved_data_manager';
import { addStatWeightsAction } from './components/stat_weights_action';

import { BulkTab } from './components/individual_sim_ui/bulk_tab';
import { GearTab } from './components/individual_sim_ui/gear_tab';
import { SettingsTab } from './components/individual_sim_ui/settings_tab';
import { RotationTab } from './components/individual_sim_ui/rotation_tab';

import {
	Class,
	Consumes,
	Debuffs,
	Encounter as EncounterProto,
	EquipmentSpec,
	Glyphs,
	HandType,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	Profession,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
} from './proto/common';

import { IndividualSimSettings, SavedRotation, SavedTalents } from './proto/ui';
import { StatWeightsResult } from './proto/api';

import { Gear } from './proto_utils/gear';
import { getMetaGemConditionDescription } from './proto_utils/gems';
import { professionNames } from './proto_utils/names';
import { Stats } from './proto_utils/stats';
import {
	getTalentPoints,
	isHealingSpec,
	isTankSpec,
	SpecOptions,
	SpecRotation,
	specToEligibleRaces,
	specToLocalStorageKey,
} from './proto_utils/utils';

import { HunterPetTalentsPicker, makePetTypeInputConfig } from './talents/hunter_pet';
import { newGlyphsPicker, newTalentsPicker } from './talents/factory';

import * as Exporters from './components/exporters';
import * as Importers from './components/importers';
import * as IconInputs from './components/icon_inputs';
import * as InputHelpers from './components/input_helpers';
import * as Mechanics from './constants/mechanics';
import * as Tooltips from './constants/tooltips';

declare var pako: any;

const SAVED_GEAR_STORAGE_KEY = '__savedGear__';
const SAVED_ROTATION_STORAGE_KEY = '__savedRotation__';
const SAVED_SETTINGS_STORAGE_KEY = '__savedSettings__';
const SAVED_TALENTS_STORAGE_KEY = '__savedTalents__';

export type InputConfig<ModObject> = (
	InputHelpers.TypedBooleanPickerConfig<ModObject> |
	InputHelpers.TypedNumberPickerConfig<ModObject> |
	InputHelpers.TypedEnumPickerConfig<ModObject> |
	InputHelpers.TypedCustomRotationPickerConfig<any, any> |
	InputHelpers.TypedItemSwapPickerConfig<any, any>
);

export interface InputSection {
	tooltip?: string,
	inputs: Array<InputConfig<Player<any>>>,
}

export interface OtherDefaults {
	profession1?: Profession,
	profession2?: Profession,
	distanceFromTarget?: number,
}

export interface IndividualSimUIConfig<SpecType extends Spec> {
	// Additional css class to add to the root element.
	cssClass: string,
	// Used to generate schemed components. E.g. 'shaman', 'druid', 'raid'
	cssScheme: string,

	knownIssues?: Array<string>;
	warnings?: Array<(simUI: IndividualSimUI<SpecType>) => SimWarning>,

	epStats: Array<Stat>;
	epPseudoStats?: Array<PseudoStat>;
	epReferenceStat: Stat;
	displayStats: Array<Stat>;
	modifyDisplayStats?: (player: Player<SpecType>) => StatMods,

	defaults: {
		gear: EquipmentSpec,
		epWeights: Stats,
		consumes: Consumes,
		rotation: SpecRotation<SpecType>,
		talents: SavedTalents,
		specOptions: SpecOptions<SpecType>,

		raidBuffs: RaidBuffs,
		partyBuffs: PartyBuffs,
		individualBuffs: IndividualBuffs,

		debuffs: Debuffs,

		other?: OtherDefaults,
	},

	playerInputs?: InputSection,
	playerIconInputs: Array<IconInputs.IconInputConfig<Player<SpecType>, any>>,
	petConsumeInputs?: Array<IconInputs.IconInputConfig<Player<SpecType>, any>>,
	rotationInputs: InputSection;
	rotationIconInputs?: Array<IconInputs.IconInputConfig<Player<any>, any>>;
	includeBuffDebuffInputs: Array<any>,
	excludeBuffDebuffInputs: Array<any>,
	otherInputs: InputSection;

	// For when extra sections are needed (e.g. Shaman totems)
	customSections?: Array<(parentElem: HTMLElement, simUI: IndividualSimUI<SpecType>) => ContentBlock>,

	encounterPicker: EncounterPickerConfig,

	presets: {
		gear: Array<PresetGear>,
		talents: Array<SavedDataConfig<Player<any>, SavedTalents>>,
		rotations?: Array<PresetRotation>,
	},
}

export interface GearAndStats {
	gear: Gear,
	bonusStats?: Stats,
}

export interface PresetGear {
	name: string;
	gear: EquipmentSpec;
	tooltip?: string;
	enableWhen?: (obj: Player<any>) => boolean;
}

export interface PresetRotation {
	name: string;
	rotation: SavedRotation;
	tooltip?: string;
	enableWhen?: (obj: Player<any>) => boolean;
}

export interface Settings {
	raidBuffs: RaidBuffs,
	partyBuffs: PartyBuffs,
	individualBuffs: IndividualBuffs,
	consumes: Consumes,
	race: Race,
	professions?: Array<Profession>;
}

// Extended shared UI for all individual player sims.
export abstract class IndividualSimUI<SpecType extends Spec> extends SimUI {
	readonly player: Player<SpecType>;
	readonly individualConfig: IndividualSimUIConfig<SpecType>;

	private raidSimResultsManager: RaidSimResultsManager | null;

	prevEpIterations: number;
	prevEpSimResult: StatWeightsResult | null;
	dpsRefStat?: Stat;
	healRefStat?: Stat;
	tankRefStat?: Stat;

	readonly bt: BulkTab;

	constructor(parentElem: HTMLElement, player: Player<SpecType>, config: IndividualSimUIConfig<SpecType>) {
		super(parentElem, player.sim, {
			cssClass: config.cssClass,
			cssScheme: config.cssScheme,
			spec: player.spec,
			knownIssues: config.knownIssues,
			launchStatus: simLaunchStatuses[player.spec],
			noticeText: aplLaunchStatuses[player.spec] == LaunchStatus.Alpha || aplLaunchStatuses[player.spec] == LaunchStatus.Beta ? 'Rotation settings have been moved to the \'Rotation\' tab, where experimental APL options are also available. Try them out!' : undefined,
		});
		this.rootElem.classList.add('individual-sim-ui');
		this.player = player;
		this.individualConfig = config;
		this.raidSimResultsManager = null;
		this.prevEpIterations = 0;
		this.prevEpSimResult = null;

		this.addWarning({
			updateOn: this.player.gearChangeEmitter,
			getContent: () => {
				if (!this.player.getGear().hasInactiveMetaGem(this.player.isBlacksmithing())) {
					return '';
				}

				const metaGem = this.player.getGear().getMetaGem()!;
				return `Meta gem disabled (${metaGem.name}): ${getMetaGemConditionDescription(metaGem)}`;
			},
		});
		this.addWarning({
			updateOn: TypedEvent.onAny([this.player.gearChangeEmitter, this.player.professionChangeEmitter]),
			getContent: () => {
				const failedProfReqs = this.player.getGear().getFailedProfessionRequirements(this.player.getProfessions());
				if (failedProfReqs.length == 0) {
					return '';
				}

				return failedProfReqs.map(fpr => `${fpr.name} requires ${professionNames[fpr.requiredProfession]}, but it is not selected.`);
			},
		});
		this.addWarning({
			updateOn: this.player.gearChangeEmitter,
			getContent: () => {
				const jcGems = this.player.getGear().getJCGems(this.player.isBlacksmithing());
				if (jcGems.length <= 3) {
					return '';
				}

				return `Only 3 Jewelcrafting Gems are allowed, but ${jcGems.length} are equipped.`;
			},
		});
		this.addWarning({
			updateOn: this.player.talentsChangeEmitter,
			getContent: () => {
				const talentPoints = getTalentPoints(this.player.getTalentsString());

				if (talentPoints == 0) {
					// Just return here, so we don't show a warning during page load.
					return '';
				} else if (talentPoints < Mechanics.MAX_TALENT_POINTS) {
					return 'Unspent talent points.';
				} else if (talentPoints > Mechanics.MAX_TALENT_POINTS) {
					return 'More than maximum talent points spent.';
				} else {
					return '';
				}
			},
		});
		this.addWarning({
			updateOn: TypedEvent.onAny([this.player.gearChangeEmitter, this.player.talentsChangeEmitter]),
			getContent: () => {
				if (!this.player.canDualWield2H() &&
					(this.player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType == HandType.HandTypeTwoHand &&
						this.player.getEquippedItem(ItemSlot.ItemSlotOffHand) != null ||
						this.player.getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.handType == HandType.HandTypeTwoHand)) {
					return "Dual wielding two-handed weapon(s) without Titan's Grip spec."
				} else {
					return '';
				}
			},
		});
		(config.warnings || []).forEach(warning => this.addWarning(warning(this)));

		if (!this.isWithinRaidSim) {
			// This needs to go before all the UI components so that gear loading is the
			// first callback invoked from waitForInit().
			this.sim.waitForInit().then(() => {
				this.loadSettings();

				if (isHealingSpec(this.player.spec)) {
					alert(Tooltips.HEALING_SIM_DISCLAIMER);
				}
			});
		}

		this.addSidebarComponents();
		this.addGearTab();
		this.bt = this.addBulkTab();
		this.addSettingsTab();
		this.addTalentsTab();
		if (aplLaunchStatuses[this.player.spec] != LaunchStatus.Unlaunched) {
			this.addRotationTab();
		}

		if (!this.isWithinRaidSim) {
			this.addDetailedResultsTab();
			this.addLogTab();
		}

		this.addTopbarComponents();
	}

	private loadSettings() {
		const initEventID = TypedEvent.nextEventID();
		TypedEvent.freezeAllAndDo(() => {
			let loadedSettings = false;

			let hash = window.location.hash;
			if (hash.length > 1) {
				// Remove leading '#'
				hash = hash.substring(1);
				try {
					const binary = atob(hash);
					const bytes = new Uint8Array(binary.length);
					for (let i = 0; i < bytes.length; i++) {
						bytes[i] = binary.charCodeAt(i);
					}

					const settingsBytes = pako.inflate(bytes);
					const settings = IndividualSimSettings.fromBinary(settingsBytes);
					this.fromProto(initEventID, settings);
					loadedSettings = true;
				} catch (e) {
					console.warn('Failed to parse settings from window hash: ' + e);
				}
			}
			window.location.hash = '';

			const savedSettings = window.localStorage.getItem(this.getSettingsStorageKey());
			if (!loadedSettings && savedSettings != null) {
				try {
					const settings = IndividualSimSettings.fromJsonString(savedSettings);
					this.fromProto(initEventID, settings);
					loadedSettings = true;
				} catch (e) {
					console.warn('Failed to parse saved settings: ' + e);
				}
			}

			if (!loadedSettings) {
				this.applyDefaults(initEventID);
			}
			this.player.setName(initEventID, 'Player');

			// This needs to go last so it doesn't re-store things as they are initialized.
			this.changeEmitter.on(eventID => {
				const jsonStr = IndividualSimSettings.toJsonString(this.toProto());
				window.localStorage.setItem(this.getSettingsStorageKey(), jsonStr);
			});
		});
	}

	private addSidebarComponents() {
		this.raidSimResultsManager = addRaidSimAction(this);
		addStatWeightsAction(this, this.individualConfig.epStats, this.individualConfig.epPseudoStats, this.individualConfig.epReferenceStat);

		const characterStats = new CharacterStats(
			this.rootElem.getElementsByClassName('sim-sidebar-footer')[0] as HTMLElement,
			this.player,
			this.individualConfig.displayStats,
			this.individualConfig.modifyDisplayStats);
	}

	private addGearTab() {
		let gearTab = new GearTab(this.simTabContentsContainer, this);
		gearTab.rootElem.classList.add('active', 'show');
	}

	private addBulkTab(): BulkTab {
		let bulkTab = new BulkTab(this.simTabContentsContainer, this);
		bulkTab.navLink.hidden = !this.sim.getShowExperimental()
		this.sim.showExperimentalChangeEmitter.on(() => {
			bulkTab.navLink.hidden = !this.sim.getShowExperimental();
		});
		return bulkTab;
	}

	private addSettingsTab() {
		new SettingsTab(this.simTabContentsContainer, this);
	}

	private addTalentsTab() {
		this.addTab('Talents', 'talents-tab', `
			<div class="talents-content tab-pane-content-container">
				<div class="talents-tab-content tab-panel-left">
					<div class="player-pet-toggle hide"></div>
					<div class="talents-picker"></div>
					<div class="glyphs-picker">
						<span>Glyphs</span>
					</div>
					<div class="pet-talents-picker hide"></div>
				</div>
				<div class="saved-talents-manager tab-panel-right"></div>
			</div>
		`);

		const talentsPicker = newTalentsPicker(this.rootElem.getElementsByClassName('talents-picker')[0] as HTMLElement, this.player);
		const glyphsPicker = newGlyphsPicker(this.rootElem.getElementsByClassName('glyphs-picker')[0] as HTMLElement, this.player);

		this.rootElem.querySelector('#talents-tab-tab')?.classList.add('sim-tab');

		const savedTalentsManager = new SavedDataManager<Player<any>, SavedTalents>(
			this.rootElem.getElementsByClassName('saved-talents-manager')[0] as HTMLElement, this, this.player, {
			label: 'Talents',
			header: { title: 'Saved Talents' },
			storageKey: this.getSavedTalentsStorageKey(),
			getData: (player: Player<any>) => SavedTalents.create({
				talentsString: player.getTalentsString(),
				glyphs: player.getGlyphs(),
			}),
			setData: (eventID: EventID, player: Player<any>, newTalents: SavedTalents) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setTalentsString(eventID, newTalents.talentsString);
					player.setGlyphs(eventID, newTalents.glyphs || Glyphs.create());
				});
			},
			changeEmitters: [this.player.talentsChangeEmitter, this.player.glyphsChangeEmitter],
			equals: (a: SavedTalents, b: SavedTalents) => SavedTalents.equals(a, b),
			toJson: (a: SavedTalents) => SavedTalents.toJson(a),
			fromJson: (obj: any) => SavedTalents.fromJson(obj),
		}
		);

		this.sim.waitForInit().then(() => {
			savedTalentsManager.loadUserData();
			this.individualConfig.presets.talents.forEach(config => {
				config.isPreset = true;
				savedTalentsManager.addSavedData({
					name: config.name,
					isPreset: true,
					data: config.data,
				});
			});

			if (this.player.getClass() == Class.ClassHunter) {
				const petTalentsPicker = new HunterPetTalentsPicker(
					this.rootElem.getElementsByClassName('pet-talents-picker')[0] as HTMLElement, this, this.player as Player<Spec.SpecHunter>
				);

				let curShown = 0;
				const updateToggle = () => {
					this.rootElem.querySelector('.talents-picker')?.classList.toggle('hide');
					this.rootElem.querySelector('.glyphs-picker')?.classList.toggle('hide');
					this.rootElem.querySelector('.pet-talents-picker')?.classList.toggle('hide');
				}

				const toggleContainer = this.rootElem.getElementsByClassName('player-pet-toggle')[0] as HTMLElement;
				toggleContainer.classList.remove('hide');
				const playerPetToggle = new EnumPicker(toggleContainer, this, {
					values: [
						{ name: 'Player', value: 0 },
						{ name: 'Pet', value: 1 },
					],
					changedEvent: sim => new TypedEvent(),
					getValue: sim => curShown,
					setValue: (eventID, sim, newValue) => {
						curShown = newValue;
						updateToggle();
					},
				});
				const petTypeToggle = new IconEnumPicker(toggleContainer, this.player as Player<Spec.SpecHunter>, makePetTypeInputConfig(false));
			}
		});
	}

	private addRotationTab() {
		new RotationTab(this.simTabContentsContainer, this);
	}

	private addDetailedResultsTab() {
		this.addTab('Results', 'detailed-results-tab', `
			<div class="detailed-results">
			</div>
		`);

		const detailedResults = new EmbeddedDetailedResults(this.rootElem.getElementsByClassName('detailed-results')[0] as HTMLElement, this, this.raidSimResultsManager!);
	}

	private addLogTab() {
		this.addTab('Log', 'log-tab', `
			<div class="log-runner">
			</div>
		`);

		const logRunner = new LogRunner(this.rootElem.getElementsByClassName('log-runner')[0] as HTMLElement, this);
	}

	private addTopbarComponents() {
		this.simHeader.addImportLink('JSON', _parent => new Importers.IndividualJsonImporter(this.rootElem, this), true);
		this.simHeader.addImportLink('80U', _parent => new Importers.Individual80UImporter(this.rootElem, this), true);
		this.simHeader.addImportLink('WoWHead', _parent => new Importers.IndividualWowheadGearPlannerImporter(this.rootElem, this), false);
		this.simHeader.addImportLink('Addon', _parent => new Importers.IndividualAddonImporter(this.rootElem, this), true);

		this.simHeader.addExportLink('Link', _parent => new Exporters.IndividualLinkExporter(this.rootElem, this), false);
		this.simHeader.addExportLink('JSON', _parent => new Exporters.IndividualJsonExporter(this.rootElem, this), true);
		this.simHeader.addExportLink('WoWHead', _parent => new Exporters.IndividualWowheadGearPlannerExporter(this.rootElem, this), false);
		this.simHeader.addExportLink('80U EP', _parent => new Exporters.Individual80UEPExporter(this.rootElem, this), false);
		this.simHeader.addExportLink('Pawn EP', _parent => new Exporters.IndividualPawnEPExporter(this.rootElem, this), false);
		this.simHeader.addExportLink("CLI", _parent => new Exporters.IndividualCLIExporter(this.rootElem, this), true);
	}

	applyDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			const tankSpec = isTankSpec(this.player.spec);
			const healingSpec = isHealingSpec(this.player.spec);

			//Special case for Totem of Wrath keeps buff and debuff sync'd
			const towEnabled = this.individualConfig.defaults.raidBuffs.totemOfWrath || this.individualConfig.defaults.debuffs.totemOfWrath
			this.individualConfig.defaults.raidBuffs.totemOfWrath = towEnabled;
			this.individualConfig.defaults.debuffs.totemOfWrath = towEnabled;

			this.player.applySharedDefaults(eventID);
			this.player.setRace(eventID, specToEligibleRaces[this.player.spec][0]);
			this.player.setGear(eventID, this.sim.db.lookupEquipmentSpec(this.individualConfig.defaults.gear));
			this.player.setConsumes(eventID, this.individualConfig.defaults.consumes);
			this.player.setRotation(eventID, this.individualConfig.defaults.rotation);
			this.player.setTalentsString(eventID, this.individualConfig.defaults.talents.talentsString);
			this.player.setGlyphs(eventID, this.individualConfig.defaults.talents.glyphs || Glyphs.create());
			this.player.setSpecOptions(eventID, this.individualConfig.defaults.specOptions);
			this.player.setBuffs(eventID, this.individualConfig.defaults.individualBuffs);
			this.player.getParty()!.setBuffs(eventID, this.individualConfig.defaults.partyBuffs);
			this.player.getRaid()!.setBuffs(eventID, this.individualConfig.defaults.raidBuffs);
			this.player.setEpWeights(eventID, this.individualConfig.defaults.epWeights);
			const defaultRatios = this.player.getDefaultEpRatios(tankSpec, healingSpec)
			this.player.setEpRatios(eventID, defaultRatios);
			this.player.setProfession1(eventID, this.individualConfig.defaults.other?.profession1 || Profession.Engineering);
			this.player.setProfession2(eventID, this.individualConfig.defaults.other?.profession2 || Profession.Jewelcrafting);
			this.player.setDistanceFromTarget(eventID, this.individualConfig.defaults.other?.distanceFromTarget || 0);

			if (this.isWithinRaidSim) {
				this.sim.raid.setTargetDummies(eventID, 0);
			} else {
				this.sim.raid.setTargetDummies(eventID, healingSpec ? 9 : 0);
				this.sim.encounter.applyDefaults(eventID);
				this.sim.raid.setDebuffs(eventID, this.individualConfig.defaults.debuffs);
				this.sim.applyDefaults(eventID, tankSpec, healingSpec);

				if (tankSpec) {
					this.sim.raid.setTanks(eventID, [this.player.makeUnitReference()]);
				} else {
					this.sim.raid.setTanks(eventID, []);
				}
			}
		});
	}

	getSavedGearStorageKey(): string {
		return this.getStorageKey(SAVED_GEAR_STORAGE_KEY);
	}

	getSavedRotationStorageKey(): string {
		return this.getStorageKey(SAVED_ROTATION_STORAGE_KEY);
	}

	getSavedSettingsStorageKey(): string {
		return this.getStorageKey(SAVED_SETTINGS_STORAGE_KEY);
	}

	getSavedTalentsStorageKey(): string {
		return this.getStorageKey(SAVED_TALENTS_STORAGE_KEY);
	}

	// Returns the actual key to use for local storage, based on the given key part and the site context.
	getStorageKey(keyPart: string): string {
		// Local storage is shared by all sites under the same domain, so we need to use
		// different keys for each spec site.
		return specToLocalStorageKey[this.player.spec] + keyPart;
	}

	toProto(): IndividualSimSettings {
		return IndividualSimSettings.create({
			settings: this.sim.toProto(),
			player: this.player.toProto(true),
			raidBuffs: this.sim.raid.getBuffs(),
			debuffs: this.sim.raid.getDebuffs(),
			tanks: this.sim.raid.getTanks(),
			partyBuffs: this.player.getParty()?.getBuffs() || PartyBuffs.create(),
			encounter: this.sim.encounter.toProto(),
			epWeightsStats: this.player.getEpWeights().toProto(),
			epRatios: this.player.getEpRatios(),
			targetDummies: this.sim.raid.getTargetDummies(),
		});
	}

	toLink(): string {
		const proto = this.toProto();
		// When sharing links, people generally don't intend to share settings/ep weights.
		proto.settings = undefined;
		proto.epWeights = [];

		const protoBytes = IndividualSimSettings.toBinary(proto);
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		const encoded = btoa(String.fromCharCode(...deflated));

		const linkUrl = new URL(window.location.href);
		linkUrl.hash = encoded;
		return linkUrl.toString();
	}

	fromProto(eventID: EventID, settings: IndividualSimSettings) {
		TypedEvent.freezeAllAndDo(() => {
			if (!settings.player) {
				return;
			}
			this.player.fromProto(eventID, settings.player);
			if (settings.epWeights?.length > 0) {
				this.player.setEpWeights(eventID, new Stats(settings.epWeights));
			} else if (settings.epWeightsStats) {
				this.player.setEpWeights(eventID, Stats.fromProto(settings.epWeightsStats));
			} else {
				this.player.setEpWeights(eventID, this.individualConfig.defaults.epWeights);
			}

			const tankSpec = isTankSpec(this.player.spec);
			const healingSpec = isHealingSpec(this.player.spec);
			const defaultRatios = this.player.getDefaultEpRatios(tankSpec, healingSpec);
			if (settings.epRatios) {
				const missingRatios = new Array<number>(defaultRatios.length - settings.epRatios.length).fill(0);
				this.player.setEpRatios(eventID, settings.epRatios.concat(missingRatios));
			} else {
				this.player.setEpRatios(eventID, defaultRatios);
			}

			this.sim.raid.setBuffs(eventID, settings.raidBuffs || RaidBuffs.create());
			this.sim.raid.setDebuffs(eventID, settings.debuffs || Debuffs.create());
			this.sim.raid.setTanks(eventID, settings.tanks || []);
			this.sim.raid.setTargetDummies(eventID, settings.targetDummies);
			const party = this.player.getParty();
			if (party) {
				party.setBuffs(eventID, settings.partyBuffs || PartyBuffs.create());
			}

			this.sim.encounter.fromProto(eventID, settings.encounter || EncounterProto.create());

			if (settings.settings) {
				this.sim.fromProto(eventID, settings.settings);
			} else {
				const tankSpec = isTankSpec(this.player.spec);
				const healingSpec = isHealingSpec(this.player.spec);
				this.sim.applyDefaults(eventID, tankSpec, healingSpec);
			}

			// Needed because of new proto field addition. Can remove on 2022/11/14 (2 months).
			if (!isHealingSpec(this.player.spec)) {
				this.sim.setShowDamageMetrics(eventID, true);
			}
		});
	}

	splitRelevantOptions<T>(options: Array<StatOption<T> | null>): Array<T> {
		return options
			.filter(option => option != null)
			.filter(option =>
				this.individualConfig.includeBuffDebuffInputs.includes(option!.item) ||
				option!.stats.length == 0 ||
				option!.stats.some(stat => this.individualConfig.epStats.includes(stat)))
			.filter(option =>
				!this.individualConfig.excludeBuffDebuffInputs.includes(option!.item))
			.map(option => option!.item);
	}
}

export interface StatOption<T> {
	stats: Array<Stat>,
	item: T,
}
