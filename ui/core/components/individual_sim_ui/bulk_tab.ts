import { IndividualSimUI } from '../../individual_sim_ui';
import { BulkComboResult, BulkSettings, ItemSpecWithSlot, ProgressMetrics, TalentLoadout } from '../../proto/api';
import { EquipmentSpec, GemColor, ItemSlot, ItemSpec, SimDatabase, SimEnchant, SimGem, SimItem, Spec } from '../../proto/common';
import { SavedTalents, UIEnchant, UIGem, UIItem, UIItem_FactionRestriction } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { Database } from '../../proto_utils/database';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { getEmptyGemSocketIconUrl } from '../../proto_utils/gems';
import { canEquipItem, getEligibleItemSlots } from '../../proto_utils/utils';
import { TypedEvent } from '../../typed_event';
import { EventID } from '../../typed_event.js';
import { BaseModal } from '../base_modal';
import { BooleanPicker } from '../boolean_picker';
import { Component } from '../component';
import { ContentBlock } from '../content_block';
import { ItemData, ItemList, ItemRenderer, SelectorModal, SelectorModalTabs } from '../gear_picker';
import { Importer } from '../importers';
import { ResultsViewer } from '../results_viewer';
import { SimTab } from '../sim_tab';

export class BulkGearJsonImporter<SpecType extends Spec> extends Importer {
	private readonly simUI: IndividualSimUI<SpecType>;
	private readonly bulkUI: BulkTab;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>, bulkUI: BulkTab) {
		super(parent, simUI, 'Bag Item Import', true);
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.descriptionElem.innerHTML = `
      <p>Import bag items from a JSON file, which can be created by the WowSimsExporter in-game AddOn.</p>
      <p>To import, upload the file or paste the text below, then click, 'Import'.</p>
    `;
	}

	async onImport(data: string) {
		try {
			const equipment = EquipmentSpec.fromJsonString(data, { ignoreUnknownFields: true });
			if (equipment?.items?.length > 0) {
				const db = await Database.loadLeftoversIfNecessary(equipment);
				const items = equipment.items.filter(spec => spec.id > 0 && db.lookupItemSpec(spec));
				if (items.length > 0) {
					this.bulkUI.addItems(items);
				}
			}
			this.close();
		} catch (e: any) {
			console.warn(e);
			alert(e.toString());
		}
	}
}

class BulkSimResultRenderer {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<Spec>, result: BulkComboResult, baseResult: BulkComboResult) {
		const dpsDivParent = document.createElement('div');
		dpsDivParent.classList.add('results-sim');

		const dpsDiv = document.createElement('div');
		dpsDiv.classList.add('bulk-result-body-dps', 'bulk-items-text-line', 'results-sim-dps', 'damage-metrics');
		dpsDivParent.appendChild(dpsDiv);

		const dpsNumber = document.createElement('span');
		dpsNumber.textContent = this.formatDps(result.unitMetrics?.dps?.avg!);
		dpsNumber.classList.add('topline-result-avg');
		dpsDiv.appendChild(dpsNumber);

		const dpsDelta = result.unitMetrics?.dps?.avg! - baseResult.unitMetrics?.dps?.avg!;
		const dpsDeltaSpan = document.createElement('span');
		dpsDeltaSpan.textContent = `${this.formatDpsDelta(dpsDelta)}`;
		dpsDeltaSpan.classList.add(dpsDelta >= 0 ? 'bulk-result-header-positive' : 'bulk-result-header-negative');
		dpsDiv.appendChild(dpsDeltaSpan);

		const itemsContainer = document.createElement('div');
		itemsContainer.classList.add('bulk-gear-combo');
		parent.appendChild(itemsContainer);
		parent.appendChild(dpsDivParent);

		const talentText = document.createElement('p');
		talentText.classList.add('talent-loadout-text');
		if (result.talentLoadout && typeof result.talentLoadout === 'object') {
			if (typeof result.talentLoadout.name === 'string') {
				talentText.textContent = 'Talent loadout used: ' + result.talentLoadout.name;
			}
		} else {
			talentText.textContent = 'Current talents';
		}

		dpsDiv.appendChild(talentText);
		if (result.itemsAdded && result.itemsAdded.length > 0) {
			const equipBtn = document.createElement('button');
			equipBtn.textContent = 'Equip';
			equipBtn.classList.add('btn', 'btn-primary', 'bulk-equipit');
			equipBtn.onclick = () => {
				result.itemsAdded.forEach(itemAdded => {
					const item = simUI.sim.db.lookupItemSpec(itemAdded.item!);
					simUI.player.equipItem(TypedEvent.nextEventID(), itemAdded.slot, item);
					simUI.simHeader.activateTab('gear-tab');
				});
			};

			parent.appendChild(equipBtn);

			for (const is of result.itemsAdded) {
				const item = simUI.sim.db.lookupItemSpec(is.item!);
				const renderer = new ItemRenderer(parent, itemsContainer, simUI.player);
				renderer.update(item!);

				const p = document.createElement('a');
				p.classList.add('bulk-result-item-slot');
				p.textContent = this.itemSlotName(is);
				renderer.nameElem.appendChild(p);
			}
		} else if (!result.talentLoadout || typeof result.talentLoadout !== 'object') {
			const p = document.createElement('p');
			p.textContent = 'No changes - this is your currently equipped gear!';
			parent.appendChild(p);
			dpsDeltaSpan.textContent = '';
		}
	}

	private formatDps(dps: number): string {
		return (Math.round(dps * 100) / 100).toFixed(2);
	}

	private formatDpsDelta(delta: number): string {
		return (delta >= 0 ? '+' : '') + this.formatDps(delta);
	}

	private itemSlotName(is: ItemSpecWithSlot): string {
		return JSON.parse(ItemSpecWithSlot.toJsonString(is, { emitDefaultValues: true }))['slot'].replace('ItemSlot', '');
	}
}

export class BulkItemPicker extends Component {
	private readonly itemElem: ItemRenderer;
	readonly simUI: IndividualSimUI<Spec>;
	readonly bulkUI: BulkTab;
	readonly index: number;

	protected item: EquippedItem;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<Spec>, bulkUI: BulkTab, item: EquippedItem, index: number) {
		super(parent, 'bulk-item-picker');
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.index = index;
		this.item = item;
		this.itemElem = new ItemRenderer(parent, this.rootElem, simUI.player);

		this.simUI.sim.waitForInit().then(() => {
			this.setItem(item);
			const slot = getEligibleItemSlots(this.item.item)[0];
			const eligibleEnchants = this.simUI.sim.db.getEnchants(slot);
			const openEnchantGemSelector = (event: Event) => {
				event.preventDefault();
				const changeEvent = new TypedEvent<void>();
				const modal = new SelectorModal(this.bulkUI.rootElem, this.simUI, this.simUI.player, {
					selectedTab: SelectorModalTabs.Enchants,
					slot: slot,
					equippedItem: this.item,
					eligibleItems: new Array<UIItem>(),
					eligibleEnchants: eligibleEnchants,
					gearData: {
						equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
							if (equippedItem) {
								const allItems = this.bulkUI.getItems();
								allItems[this.index] = equippedItem.asSpec();
								this.item = equippedItem;
								this.bulkUI.setItems(allItems);
								changeEvent.emit(TypedEvent.nextEventID());
							}
						},
						getEquippedItem: () => this.item,
						changeEvent: changeEvent,
					},
				});

				if (eligibleEnchants.length > 0) {
					modal.openTabName('Enchants');
				} else if (this.item._gems.length > 0) {
					modal.openTabName('Gem1');
				}

				const destroyItemButton = document.createElement('button');
				destroyItemButton.textContent = 'Remove from Batch';
				destroyItemButton.classList.add('btn', 'btn-danger');
				destroyItemButton.onclick = () => {
					bulkUI.setItems(
						bulkUI.getItems().filter((item, idx) => {
							return idx != this.index;
						}),
					);
					modal.close();
				};
				const closeX = modal.header?.querySelector('.close-button');
				if (closeX != undefined) {
					modal.header?.insertBefore(destroyItemButton, closeX);
				}
			};

			this.itemElem.iconElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.nameElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantGemSelector);
		});
	}

	setItem(newItem: EquippedItem | null) {
		this.itemElem.clear();
		if (newItem != null) {
			this.itemElem.update(newItem);
			this.item = newItem;
		} else {
			this.itemElem.rootElem.style.opacity = '30%';
			this.itemElem.iconElem.style.backgroundImage = `url('/wotlk/assets/item_slots/empty.jpg')`;
			this.itemElem.nameElem.textContent = 'Add new item (not implemented)';
			this.itemElem.rootElem.style.alignItems = 'center';
		}
	}
}

export class BulkTab extends SimTab {
	readonly simUI: IndividualSimUI<Spec>;

	readonly itemsChangedEmitter = new TypedEvent<void>();

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	readonly column1: HTMLElement = this.buildColumn(1, 'raid-settings-col');

	protected items: Array<ItemSpec> = new Array<ItemSpec>();

	private pendingResults: ResultsViewer;
	private pendingDiv: HTMLDivElement;

	// TODO: Make a real options probably
	private doCombos: boolean;
	private fastMode: boolean;
	private autoGem: boolean;
	private simTalents: boolean;
	private autoEnchant: boolean;
	private defaultGems: SimGem[];
	private savedTalents: TalentLoadout[];
	private gemIconElements: HTMLImageElement[];

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, simUI, { identifier: 'bulk-tab', title: 'Batch' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('bulk-tab-left', 'tab-panel-left');
		this.leftPanel.appendChild(this.column1);

		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('bulk-tab-right', 'tab-panel-right');

		this.pendingDiv = document.createElement('div');
		this.pendingDiv.classList.add('results-pending-overlay');
		this.pendingResults = new ResultsViewer(this.pendingDiv);
		this.pendingResults.hideAll();

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);
		this.contentContainer.appendChild(this.pendingDiv);

		this.doCombos = true;
		this.fastMode = true;
		this.autoGem = true;
		this.autoEnchant = true;
		this.savedTalents = [];
		this.simTalents = false;
		this.defaultGems = [UIGem.create(), UIGem.create(), UIGem.create(), UIGem.create()];
		this.gemIconElements = [];
		this.buildTabContent();

		this.simUI.sim.waitForInit().then(() => {
			this.loadSettings();
		});
	}

	private getSettingsKey(): string {
		return this.simUI.getStorageKey('bulk-settings.v1');
	}

	private loadSettings() {
		const storedSettings = window.localStorage.getItem(this.getSettingsKey());
		if (storedSettings != null) {
			const settings = BulkSettings.fromJsonString(storedSettings, {
				ignoreUnknownFields: true,
			});

			this.doCombos = settings.combinations;
			this.fastMode = settings.fastMode;
			this.autoEnchant = settings.autoEnchant;
			this.savedTalents = settings.talentsToSim;
			this.autoGem = settings.autoGem;
			this.simTalents = settings.simTalents;
			this.defaultGems = new Array<SimGem>(
				SimGem.create({ id: settings.defaultRedGem }),
				SimGem.create({ id: settings.defaultYellowGem }),
				SimGem.create({ id: settings.defaultBlueGem }),
				SimGem.create({ id: settings.defaultMetaGem }),
			);

			this.defaultGems.forEach((gem, idx) => {
				ActionId.fromItemId(gem.id)
					.fill()
					.then(filledId => {
						this.gemIconElements[idx].src = filledId.iconUrl;
					});
			});
		}
	}

	private storeSettings() {
		const settings = this.createBulkSettings();
		const setStr = BulkSettings.toJsonString(settings, { enumAsInteger: true });
		window.localStorage.setItem(this.getSettingsKey(), setStr);
	}

	protected createBulkSettings(): BulkSettings {
		return BulkSettings.create({
			items: this.items,
			// TODO(Riotdog-GehennasEU): Make all of these configurable.
			// For now, it's always constant iteration combinations mode for "sim my bags".
			combinations: this.doCombos,
			fastMode: this.fastMode,
			autoEnchant: this.autoEnchant,
			autoGem: this.autoGem,
			simTalents: this.simTalents,
			talentsToSim: this.savedTalents,
			defaultRedGem: this.defaultGems[0].id,
			defaultYellowGem: this.defaultGems[1].id,
			defaultBlueGem: this.defaultGems[2].id,
			defaultMetaGem: this.defaultGems[3].id,
			iterationsPerCombo: this.simUI.sim.getIterations(), // TODO(Riotdog-GehennasEU): Define a new UI element for the iteration setting.
		});
	}

	protected createBulkItemsDatabase(): SimDatabase {
		const itemsDb = SimDatabase.create();
		for (const is of this.items) {
			const item = this.simUI.sim.db.lookupItemSpec(is);
			if (!item) {
				throw new Error(`item with ID ${is.id} not found in database`);
			}
			itemsDb.items.push(SimItem.fromJson(UIItem.toJson(item.item), { ignoreUnknownFields: true }));
			if (item.enchant) {
				itemsDb.enchants.push(
					SimEnchant.fromJson(UIEnchant.toJson(item.enchant), {
						ignoreUnknownFields: true,
					}),
				);
			}
			for (const gem of item.gems) {
				if (gem) {
					itemsDb.gems.push(SimGem.fromJson(UIGem.toJson(gem), { ignoreUnknownFields: true }));
				}
			}
		}
		for (const gem of this.defaultGems) {
			if (gem.id > 0) {
				itemsDb.gems.push(gem);
			}
		}
		return itemsDb;
	}

	addItems(items: Array<ItemSpec>) {
		if (this.items.length == 0) {
			this.items = items;
		} else {
			this.items = this.items.concat(items);
		}
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setItems(items: Array<ItemSpec>) {
		this.items = items;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	clearItems() {
		this.items = new Array<ItemSpec>();
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	getItems(): Array<ItemSpec> {
		const result = new Array<ItemSpec>();
		this.items.forEach(spec => {
			result.push(ItemSpec.clone(spec));
		});
		return result;
	}

	setCombinations(doCombos: boolean) {
		this.doCombos = doCombos;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setFastMode(fastMode: boolean) {
		this.fastMode = fastMode;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	protected async runBulkSim(onProgress: (_?: any) => void) {
		this.pendingResults.setPending();

		try {
			await this.simUI.sim.runBulkSim(this.createBulkSettings(), this.createBulkItemsDatabase(), onProgress);
		} catch (e) {
			this.simUI.handleCrash(e);
		}
	}

	protected buildTabContent() {
		const itemsBlock = new ContentBlock(this.column1, 'bulk-items', {
			header: { title: 'Items' },
		});

		itemsBlock.bodyElement.classList.add('gear-picker-root');

		const noticeWorkInProgress = document.createElement('div');
		noticeWorkInProgress.classList.add('bulk-items-text-line');
		itemsBlock.bodyElement.appendChild(noticeWorkInProgress);
		noticeWorkInProgress.innerHTML =
			'<i>Notice: This is under very early but active development and experimental. You may also need to update your WoW AddOn if you want to import your bags.</i>';

		const itemTextIntro = document.createElement('div');
		itemTextIntro.classList.add('bulk-items-text-line');
		itemsBlock.bodyElement.appendChild(itemTextIntro);

		const itemList = document.createElement('div');

		itemList.classList.add('tab-panel-col', 'bulk-gear-combo');
		itemsBlock.bodyElement.appendChild(itemList);

		this.itemsChangedEmitter.on(() => {
			itemList.innerHTML = '';
			if (this.items.length > 0) {
				itemTextIntro.textContent = 'The following items will be simmed together with your equipped gear.';
				for (let i = 0; i < this.items.length; ++i) {
					const spec = this.items[i];
					const item = this.simUI.sim.db.lookupItemSpec(spec);
					const bulkItemPicker = new BulkItemPicker(itemList, this.simUI, this, item!, i);
				}
			}
		});

		this.clearItems();

		const resultsBlock = new ContentBlock(this.column1, 'bulk-results', {
			header: {
				title: 'Results',
				extraCssClasses: ['bulk-results-header'],
			},
		});

		resultsBlock.rootElem.hidden = true;
		resultsBlock.bodyElement.classList.add('gear-picker-root', 'tab-panel-col');

		this.simUI.sim.bulkSimStartEmitter.on(() => {
			resultsBlock.rootElem.hidden = true;
		});

		this.simUI.sim.bulkSimResultEmitter.on((_, bulkSimResult) => {
			resultsBlock.rootElem.hidden = bulkSimResult.results.length == 0;
			resultsBlock.bodyElement.innerHTML = '';

			for (const r of bulkSimResult.results) {
				const resultBlock = new ContentBlock(resultsBlock.bodyElement, 'bulk-result', {
					header: { title: '' },
					bodyClasses: ['bulk-results-body'],
				});
				new BulkSimResultRenderer(resultBlock.bodyElement, this.simUI, r, bulkSimResult.equippedGearResult!);
			}
		});

		const settingsBlock = new ContentBlock(this.rightPanel, 'bulk-settings', {
			header: { title: 'Setup' },
		});

		const bulkSimButton = document.createElement('button');
		bulkSimButton.classList.add('btn', 'btn-primary', 'w-100', 'bulk-settings-button');
		bulkSimButton.textContent = 'Simulate Batch';
		bulkSimButton.addEventListener('click', () => {
			this.pendingDiv.style.display = 'flex';
			this.leftPanel.classList.add('blurred');
			this.rightPanel.classList.add('blurred');

			const previousContents = bulkSimButton.innerHTML;
			bulkSimButton.disabled = true;
			bulkSimButton.classList.add('.disabled');
			bulkSimButton.innerHTML = `<i class="fa fa-spinner fa-spin"></i>&nbsp;Running`;

			let simStart = new Date().getTime();
			let lastTotal = 0;
			let rounds = 0;
			let currentRound = 0;
			let combinations = 0;

			this.runBulkSim((progressMetrics: ProgressMetrics) => {
				console.log(progressMetrics);

				const msSinceStart = new Date().getTime() - simStart;
				const iterPerSecond = progressMetrics.completedIterations / (msSinceStart / 1000);

				if (combinations == 0) {
					combinations = progressMetrics.totalSims;
				}
				if (this.fastMode) {
					if (rounds == 0 && progressMetrics.totalSims > 0) {
						rounds = Math.ceil(Math.log(progressMetrics.totalSims / 20) / Math.log(2)) + 1;
						currentRound = 1;
					}
					if (progressMetrics.totalSims < lastTotal) {
						currentRound += 1;
						simStart = new Date().getTime();
					}
				}

				this.setSimProgress(progressMetrics, iterPerSecond, currentRound, rounds, combinations);
				lastTotal = progressMetrics.totalSims;

				if (progressMetrics.finalBulkResult != null) {
					// reset state
					this.pendingDiv.style.display = 'none';
					this.leftPanel.classList.remove('blurred');
					this.rightPanel.classList.remove('blurred');

					this.pendingResults.hideAll();
					bulkSimButton.disabled = false;
					bulkSimButton.classList.remove('.disabled');
					bulkSimButton.innerHTML = previousContents;
				}
			});
		});

		settingsBlock.bodyElement.appendChild(bulkSimButton);

		const importButton = document.createElement('button');
		importButton.classList.add('btn', 'btn-secondary', 'w-100', 'bulk-settings-button');
		importButton.innerHTML = '<i class="fa fa-download"></i> Import From Bags';
		importButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI, this));
		settingsBlock.bodyElement.appendChild(importButton);

		const importFavsButton = document.createElement('button');
		importFavsButton.classList.add('btn', 'btn-secondary', 'w-100', 'bulk-settings-button');
		importFavsButton.innerHTML = '<i class="fa fa-download"></i> Import Favorites';
		importFavsButton.addEventListener('click', () => {
			const filters = this.simUI.player.sim.getFilters();
			const items = filters.favoriteItems.map(itemID => {
				return ItemSpec.create({ id: itemID });
			});
			this.addItems(items);
		});
		settingsBlock.bodyElement.appendChild(importFavsButton);

		const searchButton = document.createElement('button');
		const searchText = document.createElement('input');
		searchText.type = 'text';
		searchText.placeholder = 'search...';
		searchText.style.display = 'none';

		const searchResults = document.createElement('ul');
		searchResults.classList.add('batch-search-results');

		let allItems = Array<UIItem>();

		searchText.addEventListener('keyup', ev => {
			if (ev.key == 'Enter') {
				const toAdd = Array<ItemSpec>();
				searchResults.childNodes.forEach(node => {
					const strID = (node as HTMLElement).getAttribute('data-item-id');
					if (strID != null) {
						toAdd.push(ItemSpec.create({ id: Number.parseInt(strID) }));
					}
				});
				this.addItems(toAdd);
			}
		});

		searchText.addEventListener('input', e => {
			const searchString = searchText.value;
			searchResults.innerHTML = '';
			if (searchString.length == 0) {
				return;
			}
			const pieces = searchString.split(' ');

			let displayCount = 0;
			allItems.every(item => {
				let matched = true;
				const lcName = item.name.toLowerCase();
				const lcSetName = item.setName.toLowerCase();

				pieces.forEach(piece => {
					const lcPiece = piece.toLowerCase();
					if (!lcName.includes(lcPiece) && !lcSetName.includes(lcPiece)) {
						matched = false;
						return false;
					}
					return true;
				});

				if (matched) {
					const itemElement = document.createElement('li');
					itemElement.innerHTML = `<span>${item.name}</span>`;
					itemElement.setAttribute('data-item-id', item.id.toString());
					itemElement.addEventListener('click', ev => {
						this.addItems(Array<ItemSpec>(ItemSpec.create({ id: item.id })));
					});
					if (item.heroic) {
						const htxt = document.createElement('span');
						htxt.style.color = 'green';
						htxt.innerText = '[H]';
						itemElement.appendChild(htxt);
					}
					if (item.factionRestriction == UIItem_FactionRestriction.HORDE_ONLY) {
						const ftxt = document.createElement('span');
						ftxt.style.color = 'red';
						ftxt.innerText = '(H)';
						itemElement.appendChild(ftxt);
					}
					if (item.factionRestriction == UIItem_FactionRestriction.ALLIANCE_ONLY) {
						const ftxt = document.createElement('span');
						ftxt.style.color = 'blue';
						ftxt.innerText = '(A)';
						itemElement.appendChild(ftxt);
					}
					searchResults.append(itemElement);
					displayCount++;
				}

				return displayCount < 10;
			});
		});

		searchButton.classList.add('btn', 'btn-secondary', 'w-100', 'bulk-settings-button');
		const baseSearchHTML = '<i class="fa fa-search"></i> Add Item';
		searchButton.innerHTML = baseSearchHTML;
		searchButton.addEventListener('click', () => {
			if (searchText.style.display == 'none') {
				searchButton.innerHTML = 'Close Search Results';
				allItems = this.simUI.sim.db.getAllItems().filter(item => {
					return canEquipItem(item, this.simUI.player.spec, undefined);
				});
				searchText.style.display = 'block';
				searchText.focus();
			} else {
				searchButton.innerHTML = baseSearchHTML;
				searchText.style.display = 'none';
				searchResults.innerHTML = '';
			}
		});
		settingsBlock.bodyElement.appendChild(searchButton);
		settingsBlock.bodyElement.appendChild(searchText);
		settingsBlock.bodyElement.appendChild(searchResults);

		const clearButton = document.createElement('button');
		clearButton.classList.add('btn', 'btn-secondary', 'w-100', 'bulk-settings-button');
		clearButton.textContent = 'Clear All';
		clearButton.addEventListener('click', () => {
			this.clearItems();
			resultsBlock.rootElem.hidden = true;
			resultsBlock.bodyElement.innerHTML = '';
		});
		settingsBlock.bodyElement.appendChild(clearButton);

		// Talents to sim
		const talentsToSimDiv = document.createElement('div');
		if (this.simTalents) {
			talentsToSimDiv.style.display = 'flex';
		} else {
			talentsToSimDiv.style.display = 'none';
		}
		talentsToSimDiv.classList.add('talents-picker-container');
		const talentsLabel = document.createElement('label');
		talentsLabel.innerText = 'Pick talents to sim (will increase time to sim)';
		talentsToSimDiv.appendChild(talentsLabel);
		const talentsContainerDiv = document.createElement('div');
		talentsContainerDiv.classList.add('talents-container');

		const dataStr = window.localStorage.getItem(this.simUI.getSavedTalentsStorageKey());

		let jsonData;
		try {
			if (dataStr !== null) {
				jsonData = JSON.parse(dataStr);
			}
		} catch (e) {
			console.warn('Invalid json for local storage value: ' + dataStr);
		}
		const handleToggle = (frag: HTMLElement, load: TalentLoadout) => {
			const chipDiv = frag.querySelector('.saved-data-set-chip');
			const exists = this.savedTalents.some(talent => talent.name === load.name); // Replace 'id' with your unique identifier

			console.log('Exists:', exists);
			console.log('Load Object:', load);
			console.log('Saved Talents Before Update:', this.savedTalents);

			if (exists) {
				// If the object exists, find its index and remove it
				const indexToRemove = this.savedTalents.findIndex(talent => talent.name === load.name);
				this.savedTalents.splice(indexToRemove, 1);
				chipDiv?.classList.remove('active');
			} else {
				// If the object does not exist, add it
				this.savedTalents.push(load);
				chipDiv?.classList.add('active');
			}

			console.log('Updated savedTalents:', this.savedTalents);
		};
		for (const name in jsonData) {
			try {
				console.log(name, jsonData[name]);
				const savedTalentLoadout = SavedTalents.fromJson(jsonData[name]);
				var loadout = {
					talentsString: savedTalentLoadout.talentsString,
					glyphs: savedTalentLoadout.glyphs,
					name: name,
				};

				const index = this.savedTalents.findIndex(talent => JSON.stringify(talent) === JSON.stringify(loadout));
				const talentFragment = document.createElement('fragment');
				talentFragment.innerHTML = `
					<div class="saved-data-set-chip badge rounded-pill ${index !== -1 ? 'active' : ''}">
						<a href="javascript:void(0)" class="saved-data-set-name" role="button">${name}</a>
					</div>`;

				console.log('Adding event for loadout', loadout);
				// Wrap the event listener addition in an IIFE
				(function (talentFragment, loadout) {
					talentFragment.addEventListener('click', () => handleToggle(talentFragment, loadout));
				})(talentFragment, loadout);

				talentsContainerDiv.appendChild(talentFragment);
			} catch (e) {
				console.log(e);
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}

		talentsToSimDiv.append(talentsContainerDiv);
		//////////////////////
		////////////////////////////////////

		// Default Gem Options
		const defaultGemDiv = document.createElement('div');
		if (this.autoGem) {
			defaultGemDiv.style.display = 'flex';
		} else {
			defaultGemDiv.style.display = 'none';
		}

		defaultGemDiv.classList.add('default-gem-container');
		const gemLabel = document.createElement('label');
		gemLabel.innerText = 'Defaults for Auto Gem';
		defaultGemDiv.appendChild(gemLabel);

		const gemSocketsDiv = document.createElement('div');
		gemSocketsDiv.classList.add('sockets-container');

		Array<GemColor>(GemColor.GemColorRed, GemColor.GemColorYellow, GemColor.GemColorBlue, GemColor.GemColorMeta).forEach((socketColor, socketIndex) => {
			const gemFragment = document.createElement('fragment');
			gemFragment.innerHTML = `
          <div class="gem-socket-container">
            <img class="gem-icon" />
            <img class="socket-icon" />
          </div>
        `;

			const gemContainer = gemFragment.children[0] as HTMLElement;
			this.gemIconElements.push(gemContainer.querySelector('.gem-icon') as HTMLImageElement);
			const socketIconElem = gemContainer.querySelector('.socket-icon') as HTMLImageElement;
			socketIconElem.src = getEmptyGemSocketIconUrl(socketColor);

			let selector: GemSelectorModal;

			const handleChoose = (itemData: ItemData<UIGem>) => {
				this.defaultGems[socketIndex] = itemData.item;
				this.storeSettings();
				ActionId.fromItemId(itemData.id)
					.fill()
					.then(filledId => {
						this.gemIconElements[socketIndex].src = filledId.iconUrl;
					});
				selector.close();
			};

			const openGemSelector = (color: GemColor, socketIndex: number) => {
				return (event: Event) => {
					if (selector == null) {
						selector = new GemSelectorModal(this.simUI.rootElem, this.simUI, socketColor, handleChoose);
					}
					selector.show();
				};
			};

			this.gemIconElements[socketIndex].addEventListener('click', openGemSelector(socketColor, socketIndex));
			gemContainer.addEventListener('click', openGemSelector(socketColor, socketIndex));
			gemSocketsDiv.appendChild(gemContainer);
		});
		defaultGemDiv.appendChild(gemSocketsDiv);

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			label: 'Fast Mode',
			labelTooltip: 'Fast mode reduces accuracy but will run faster.',
			changedEvent: (obj: BulkTab) => this.itemsChangedEmitter,
			getValue: obj => this.fastMode,
			setValue: (id: EventID, obj: BulkTab, value: boolean) => {
				obj.fastMode = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			label: 'Combinations',
			labelTooltip:
				'When checked bulk simulator will create all possible combinations of the items. When disabled trinkets and rings will still run all combinations becausee they have two slots to fill each.',
			changedEvent: (obj: BulkTab) => this.itemsChangedEmitter,
			getValue: obj => this.doCombos,
			setValue: (id: EventID, obj: BulkTab, value: boolean) => {
				obj.doCombos = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			label: 'Auto Enchant',
			labelTooltip: 'When checked bulk simulator apply the current enchant for a slot to each replacement item it can.',
			changedEvent: (obj: BulkTab) => this.itemsChangedEmitter,
			getValue: obj => this.autoEnchant,
			setValue: (id: EventID, obj: BulkTab, value: boolean) => {
				obj.autoEnchant = value;
				if (value) {
					defaultGemDiv.style.display = 'flex';
				} else {
					defaultGemDiv.style.display = 'none';
				}
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			label: 'Auto Gem',
			labelTooltip: 'When checked bulk simulator will fill any un-filled gem sockets with default gems.',
			changedEvent: (obj: BulkTab) => this.itemsChangedEmitter,
			getValue: obj => this.autoGem,
			setValue: (id: EventID, obj: BulkTab, value: boolean) => {
				obj.autoGem = value;
				if (value) {
					defaultGemDiv.style.display = 'flex';
				} else {
					defaultGemDiv.style.display = 'none';
				}
			},
		});

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			label: 'Sim Talents',
			labelTooltip: 'When checked bulk simulator will sim chosen talent setups. Warning, it might cause the bulk sim to run for a lot longer',
			changedEvent: (obj: BulkTab) => this.itemsChangedEmitter,
			getValue: obj => this.simTalents,
			setValue: (id: EventID, obj: BulkTab, value: boolean) => {
				obj.simTalents = value;
				if (value) {
					talentsToSimDiv.style.display = 'flex';
				} else {
					talentsToSimDiv.style.display = 'none';
				}
			},
		});

		settingsBlock.bodyElement.appendChild(defaultGemDiv);
		settingsBlock.bodyElement.appendChild(talentsToSimDiv);
	}

	private setSimProgress(progress: ProgressMetrics, iterPerSecond: number, currentRound: number, rounds: number, combinations: number) {
		const secondsRemain = ((progress.totalIterations - progress.completedIterations) / iterPerSecond).toFixed();

		let roundsText = '';
		if (rounds > 0) {
			roundsText = `${currentRound} / ${rounds} refining rounds`;
		}

		this.pendingResults.setContent(`
      <div class="results-sim">
        <div class="">${combinations} total combinations.</div>
        <div class="">${roundsText}</div>
        <div class=""> ${progress.completedSims} / ${progress.totalSims}<br>simulations complete</div>
        <div class="">
          ${progress.completedIterations} / ${progress.totalIterations}<br>iterations complete
        </div>
        <div class="">
          ${secondsRemain} seconds remaining.
        </div>
      </div>
    `);
	}
}

class GemSelectorModal extends BaseModal {
	private readonly simUI: IndividualSimUI<Spec>;

	private readonly contentElem: HTMLElement;
	private ilist: ItemList<UIGem> | null;
	private socketColor: GemColor;
	private onSelect: (itemData: ItemData<UIGem>) => void;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<Spec>, socketColor: GemColor, onSelect: (itemData: ItemData<UIGem>) => void) {
		super(parent, 'selector-modal');

		this.simUI = simUI;
		this.onSelect = onSelect;
		this.socketColor = socketColor;
		this.ilist = null;

		window.scrollTo({ top: 0 });

		this.header!.insertAdjacentHTML('afterbegin', `<span>Choose Default Gem</span>`);
		this.body.innerHTML = `<div class="tab-content selector-modal-tab-content"></div>`;
		this.contentElem = this.rootElem.querySelector('.selector-modal-tab-content') as HTMLElement;
	}

	show() {
		// construct item list the first time its opened.
		// This makes startup faster and also means we are sure to have item database loaded.
		if (this.ilist == null) {
			this.ilist = new ItemList<UIGem>(
				this.body,
				this.simUI,
				{
					selectedTab: SelectorModalTabs.Gem1,
					slot: ItemSlot.ItemSlotHead,
					equippedItem: null,
					eligibleItems: new Array<UIItem>(),
					eligibleEnchants: new Array<UIEnchant>(),
					gearData: {
						equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {},
						getEquippedItem: () => null,
						changeEvent: new TypedEvent(), // FIXME
					},
				},
				this.simUI.player,
				'Gem1',
				this.simUI.player.getGems(this.socketColor).map((gem: UIGem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						heroic: false,
						baseEP: 0,
						ignoreEPFilter: true,
						onEquip: (eventID, gem: UIGem) => {},
					};
				}),
				gem => {
					return this.simUI.player.computeGemEP(gem);
				},
				() => {
					return null;
				},
				this.socketColor,
				() => {},
				this.onSelect,
			);

			// let invokeUpdate = () => {this.ilist?.updateSelected()}
			const applyFilter = () => {
				this.ilist?.applyFilters();
			};
			// Add event handlers
			// this.itemsChangedEmitter.on(invokeUpdate);

			this.addOnDisposeCallback(() => this.ilist?.dispose());

			this.simUI.sim.phaseChangeEmitter.on(applyFilter);
			this.simUI.sim.filtersChangeEmitter.on(applyFilter);
			// gearData.changeEvent.on(applyFilter);
		}

		this.open();
	}
}
