import { ActionId } from './proto_utils/action_id.js';
import { BattleElixir, HandType } from './proto/common.js';
import { BooleanPicker, BooleanPickerConfig } from './components/boolean_picker.js';
import { CharacterStats, StatMods } from './components/character_stats.js';
import { Class } from './proto/common.js';
import { Conjured } from './proto/common.js';
import { Consumes } from './proto/common.js';
import { Cooldowns } from './proto/common.js';
import { CooldownsPicker } from './components/individual_sim_ui/cooldowns_picker.js';
import { Debuffs } from './proto/common.js';
import { DetailedResults } from './components/detailed_results.js';

import { CustomRotationPicker } from './components/individual_sim_ui/custom_rotation_picker.js';
import { Encounter as EncounterProto } from './proto/common.js';
import { Encounter } from './encounter.js';
import { EncounterPicker, EncounterPickerConfig } from './components/encounter_picker.js';
import { EnumPicker, EnumPickerConfig } from './components/enum_picker.js';
import { EquipmentSpec } from './proto/common.js';
import { EventID, TypedEvent } from './typed_event.js';
import { Flask } from './proto/common.js';
import { Food } from './proto/common.js';
import { Gear } from './proto_utils/gear.js';
import { GearPicker } from './components/gear_picker.js';
import { Glyphs } from './proto/common.js';
import { GuardianElixir } from './proto/common.js';
import { HealingModel } from './proto/common.js';
import { HunterPetTalentsPicker } from './talents/hunter_pet.js';
import { IconEnumPicker, IconEnumPickerConfig } from './components/icon_enum_picker.js';
import { IconPicker, IconPickerConfig } from './components/icon_picker.js';
import { ItemSlot } from './proto/common.js';
import { IndividualBuffs } from './proto/common.js';
import { IndividualSimSettings } from './proto/ui.js';
import { Input } from './components/input.js';
import { LogRunner } from './components/log_runner.js';
import { MobType } from './proto/common.js';
import { MultiIconPicker } from './components/multi_icon_picker.js';
import { NumberPicker, NumberPickerConfig } from './components/number_picker.js';
import { Party } from './party.js';
import { PartyBuffs } from './proto/common.js';
import { PetFood } from './proto/common.js';
import { Player as PlayerProto } from './proto/api.js';
import { Player } from './player.js';
import { Potions } from './proto/common.js';
import { Profession } from './proto/common.js';
import { Race } from './proto/common.js';
import { Raid } from './raid.js';
import { RaidBuffs } from './proto/common.js';
import { SavedDataConfig, SavedDataManager } from './components/saved_data_manager.js';
import { SavedEncounter } from './proto/ui.js';
import { SavedGearSet } from './proto/ui.js';
import { SavedSettings } from './proto/ui.js';
import { SavedTalents } from './proto/ui.js';
import { SettingsMenu } from './components/settings_menu.js';
import { Sim } from './sim.js';
import { SimOptions } from './proto/api.js';
import { SimSettings as SimSettingsProto } from './proto/ui.js';
import { SimUI, SimWarning } from './sim_ui.js';
import { Spec } from './proto/common.js';
import { SpecOptions } from './proto_utils/utils.js';
import { SpecRotation } from './proto_utils/utils.js';
import { Stat, PseudoStat } from './proto/common.js';
import { StatWeightsRequest, StatWeightsResult } from './proto/api.js';
import { Stats } from './proto_utils/stats.js';
import { Target } from './target.js';
import { Target as TargetProto } from './proto/common.js';
import { UnitStats } from './proto/common.js';
import { addRaidSimAction, RaidSimResultsManager } from './components/raid_sim_action.js';
import { addStatWeightsAction } from './components/stat_weights_action.js';
import { equalsOrBothNull, getEnumValues } from './utils.js';
import { getMetaGemConditionDescription } from './proto_utils/gems.js';
import { getTalentPoints } from './proto_utils/utils.js';
import { isDualWieldSpec } from './proto_utils/utils.js';
import { simLaunchStatuses } from './launched_sims.js';
import { makePetTypeInputConfig } from './talents/hunter_pet.js';
import { newGlyphsPicker } from './talents/factory.js';
import { newTalentsPicker } from './talents/factory.js';
import { professionNames, raceNames } from './proto_utils/names.js';
import { isHealingSpec, isTankSpec } from './proto_utils/utils.js';
import { specToEligibleRaces } from './proto_utils/utils.js';
import { specToLocalStorageKey } from './proto_utils/utils.js';

import { Tooltip } from 'bootstrap';

import * as Exporters from './components/exporters.js';
import * as Importers from './components/importers.js';
import * as IconInputs from './components/icon_inputs.js';
import * as InputHelpers from './components/input_helpers.js';
import * as Mechanics from './constants/mechanics.js';
import * as OtherConstants from './constants/other.js';
import * as Tooltips from './constants/tooltips.js';
import { SettingsTab } from './components/individual_sim_ui/settings_tab.js';
import { ContentBlock } from './components/content_block.js';

declare var tippy: any;
declare var pako: any;

const SAVED_GEAR_STORAGE_KEY = '__savedGear__';
const SAVED_ROTATION_STORAGE_KEY = '__savedRotation__';
const SAVED_SETTINGS_STORAGE_KEY = '__savedSettings__';
const SAVED_TALENTS_STORAGE_KEY = '__savedTalents__';

export type InputConfig<ModObject> = (
	InputHelpers.TypedBooleanPickerConfig<ModObject> |
	InputHelpers.TypedNumberPickerConfig<ModObject> |
	InputHelpers.TypedEnumPickerConfig<ModObject> |
	InputHelpers.TypedCustomRotationPickerConfig<any, any>
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
		rotation?: Array<SavedDataConfig<Player<any>, string>>,
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

	constructor(parentElem: HTMLElement, player: Player<SpecType>, config: IndividualSimUIConfig<SpecType>) {
		super(parentElem, player.sim, {
			spec: player.spec,
			knownIssues: config.knownIssues,
			launchStatus: simLaunchStatuses[player.spec],
		});
		this.rootElem.classList.add('individual-sim-ui', config.cssClass);
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
		this.addSettingsTab();
		this.addTalentsTab();

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
		this.addTab('Gear', 'gear-tab', `
			<div class="gear-tab-columns">
				<div class="left-gear-panel">
					<div class="gear-picker"></div>
				</div>
				<div class="right-gear-panel">
					<div class="saved-gear-manager"></div>
				</div>
			</div>
		`);

		const gearPicker = new GearPicker(this.rootElem.getElementsByClassName('gear-picker')[0] as HTMLElement, this.player);

		const savedGearManager = new SavedDataManager<Player<any>, SavedGearSet>(this.rootElem.getElementsByClassName('saved-gear-manager')[0] as HTMLElement, this.player, {
			header: {title: "Gear Sets"},
			label: 'Gear Set',
			storageKey: this.getSavedGearStorageKey(),
			getData: (player: Player<any>) => {
				return SavedGearSet.create({
					gear: player.getGear().asSpec(),
					bonusStatsStats: player.getBonusStats().toProto(),
				});
			},
			setData: (eventID: EventID, player: Player<any>, newSavedGear: SavedGearSet) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setGear(eventID, this.sim.db.lookupEquipmentSpec(newSavedGear.gear || EquipmentSpec.create()));
					if (newSavedGear.bonusStats && newSavedGear.bonusStats.some(s => s != 0)) {
						player.setBonusStats(eventID, new Stats(newSavedGear.bonusStats));
					} else {
						player.setBonusStats(eventID, Stats.fromProto(newSavedGear.bonusStatsStats || UnitStats.create()));
					}
				});
			},
			changeEmitters: [this.player.changeEmitter],
			equals: (a: SavedGearSet, b: SavedGearSet) => SavedGearSet.equals(a, b),
			toJson: (a: SavedGearSet) => SavedGearSet.toJson(a),
			fromJson: (obj: any) => SavedGearSet.fromJson(obj),
		});

		this.sim.waitForInit().then(() => {
			savedGearManager.loadUserData();
			this.individualConfig.presets.gear.forEach(presetGear => {
				savedGearManager.addSavedData({
					name: presetGear.name,
					tooltip: presetGear.tooltip,
					isPreset: true,
					data: SavedGearSet.create({
						// Convert to gear and back so order is always the same.
						gear: this.sim.db.lookupEquipmentSpec(presetGear.gear).asSpec(),
						bonusStatsStats: new Stats().toProto(),
					}),
					enableWhen: presetGear.enableWhen,
				});
			});
		});
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

		const savedTalentsManager = new SavedDataManager<Player<any>, SavedTalents>(this.rootElem.getElementsByClassName('saved-talents-manager')[0] as HTMLElement, this.player, {
			label: 'Talents',
			header: {title: 'Saved Talents'},
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
		});

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
				const petTalentsPicker = new HunterPetTalentsPicker(this.rootElem.getElementsByClassName('pet-talents-picker')[0] as HTMLElement, this.player as Player<Spec.SpecHunter>);

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

	private addDetailedResultsTab() {
		this.addTab('Results', 'detailed-results-tab', `
			<div class="detailed-results">
			</div>
		`);

		const detailedResults = new DetailedResults(this.rootElem.getElementsByClassName('detailed-results')[0] as HTMLElement, this, this.raidSimResultsManager!);
	}

	private addLogTab() {
		this.addTab('Log', 'log-tab', `
			<div class="log-runner">
			</div>
		`);

		const logRunner = new LogRunner(this.rootElem.getElementsByClassName('log-runner')[0] as HTMLElement, this);
	}

	private addTopbarComponents() {
		this.simHeader.addImportLink('JSON', parent => new Importers.IndividualJsonImporter(parent, this), true);
		this.simHeader.addImportLink('80U', parent => new Importers.Individual80UImporter(parent, this), true);
		this.simHeader.addImportLink('Addon', parent => new Importers.IndividualAddonImporter(parent, this), true);

		this.simHeader.addExportLink('Link', parent => new Exporters.IndividualLinkExporter(parent, this), false);
		this.simHeader.addExportLink('JSON', parent => new Exporters.IndividualJsonExporter(parent, this), true);
		this.simHeader.addExportLink('80U EP', parent => new Exporters.Individual80UEPExporter(parent, this), false);
		this.simHeader.addExportLink('Pawn EP', parent => new Exporters.IndividualPawnEPExporter(parent, this), false);
	}

	applyDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			const tankSpec = isTankSpec(this.player.spec);
			const healingSpec = isHealingSpec(this.player.spec);

			//Special case for Totem of Wrath keeps buff and debuff sync'd
			const towEnabled =  this.individualConfig.defaults.raidBuffs.totemOfWrath || this.individualConfig.defaults.debuffs.totemOfWrath
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
					this.sim.raid.setTanks(eventID, [this.player.makeRaidTarget()]);
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
