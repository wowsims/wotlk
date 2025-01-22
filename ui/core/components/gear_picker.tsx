import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment, ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../css_utils';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Player } from '../player';
import { Class, GemColor, ItemQuality, ItemSlot, ItemSpec, ItemType } from '../proto/common';
import { DatabaseFilters, RepFaction, UIEnchant as Enchant, UIGem as Gem, UIItem as Item, UIItem_FactionRestriction } from '../proto/ui.js';
import { ActionId } from '../proto_utils/action_id';
import { getEnchantDescription, getUniqueEnchantString } from '../proto_utils/enchants';
import { EquippedItem } from '../proto_utils/equipped_item';
import { gemMatchesSocket, getEmptyGemSocketIconUrl } from '../proto_utils/gems';
import { difficultyNames, professionNames, REP_FACTION_NAMES, REP_LEVEL_NAMES, slotNames } from '../proto_utils/names.js';
import { Stats } from '../proto_utils/stats';
import { Sim } from '../sim.js';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { formatDeltaTextElem } from '../utils';
import { BaseModal } from './base_modal';
import { Component } from './component';
import { FiltersMenu } from './filters_menu';
import {
	makePhaseSelector,
	makeShow1hWeaponsSelector,
	makeShow2hWeaponsSelector,
	makeShowEPValuesSelector,
	makeShowMatchingGemsSelector,
} from './other_inputs';
import { Clusterize } from './virtual_scroll/clusterize.js';

const EP_TOOLTIP = `
	EP (Equivalence Points) is way of comparing items by multiplying the raw stats of an item with your current stat weights.
	More EP does not necessarily mean more DPS, as EP doesn't take into account stat caps and non-linear stat calculations.
`;

const createHeroicLabel = () => {
	return <span className="heroic-label">[H]</span>;
};

const createGemContainer = (socketColor: GemColor, gem: Gem | null) => {
	const gemIconElem = ref<HTMLImageElement>();

	const gemContainer = (
		<div className="gem-socket-container">
			<img ref={gemIconElem} className={`gem-icon ${gem == null ? 'hide' : ''}`} />
			<img className="socket-icon" src={getEmptyGemSocketIconUrl(socketColor)} />
		</div>
	);

	if (gem != null) {
		ActionId.fromItemId(gem.id)
			.fill()
			.then(filledId => {
				gemIconElem.value!.src = filledId.iconUrl;
			});
	}
	return gemContainer;
};

export class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSide = document.createElement('div');
		leftSide.classList.add('gear-picker-left', 'tab-panel-col');
		this.rootElem.appendChild(leftSide);

		const rightSide = document.createElement('div');
		rightSide.classList.add('gear-picker-right', 'tab-panel-col');
		this.rootElem.appendChild(rightSide);

		const leftItemPickers = [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotNeck,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotBack,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotMainHand,
			ItemSlot.ItemSlotOffHand,
			ItemSlot.ItemSlotRanged,
		].map(slot => new ItemPicker(leftSide, simUI, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
		].map(slot => new ItemPicker(rightSide, simUI, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);
	}
}

export class ItemRenderer extends Component {
	private readonly player: Player<any>;

	readonly iconElem: HTMLAnchorElement;
	readonly nameElem: HTMLAnchorElement;
	readonly ilvlElem: HTMLSpanElement;
	readonly enchantElem: HTMLAnchorElement;
	readonly socketsContainerElem: HTMLElement;

	constructor(parent: HTMLElement, root: HTMLElement, player: Player<any>) {
		super(parent, 'item-picker-root', root);
		this.player = player;

		const iconElem = ref<HTMLAnchorElement>();
		const nameElem = ref<HTMLAnchorElement>();
		const ilvlElem = ref<HTMLSpanElement>();
		const enchantElem = ref<HTMLAnchorElement>();
		const sce = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<>
				<div className="item-picker-icon-wrapper">
					<span className="item-picker-ilvl" ref={ilvlElem} />
					<a ref={iconElem} className="item-picker-icon" href="javascript:void(0)" attributes={{ role: 'button' }}>
						<div ref={sce} className="item-picker-sockets-container"></div>
					</a>
				</div>
				<div className="item-picker-labels-container">
					<a ref={nameElem} className="item-picker-name" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<a ref={enchantElem} className="item-picker-enchant" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
				</div>
			</>,
		);

		this.iconElem = iconElem.value!;
		this.nameElem = nameElem.value!;
		this.ilvlElem = ilvlElem.value!;
		this.enchantElem = enchantElem.value!;
		this.socketsContainerElem = sce.value!;
	}

	clear() {
		this.nameElem.removeAttribute('data-wowhead');
		this.nameElem.removeAttribute('href');
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');
		this.enchantElem.removeAttribute('data-wowhead');
		this.enchantElem.removeAttribute('href');
		this.iconElem.removeAttribute('href');

		this.iconElem.style.backgroundImage = '';
		this.enchantElem.innerText = '';
		this.socketsContainerElem.innerText = '';
		this.nameElem.textContent = '';
		this.ilvlElem.replaceChildren();
	}

	update(newItem: EquippedItem) {
		this.nameElem.textContent = newItem.item.name;
		if (newItem.item.heroic) {
			this.nameElem.insertAdjacentElement('beforeend', createHeroicLabel());
		} else {
			this.nameElem.querySelector('.heroic-label')?.remove();
		}
		this.ilvlElem.textContent = newItem.item.ilvl.toString();

		setItemQualityCssClass(this.nameElem, newItem.item.quality);

		this.player.setWowheadData(newItem, this.iconElem);
		this.player.setWowheadData(newItem, this.nameElem);
		newItem
			.asActionId()
			.fill()
			.then(filledId => {
				filledId.setBackgroundAndHref(this.iconElem);
				filledId.setWowheadHref(this.nameElem);
			});

		if (newItem.enchant) {
			getEnchantDescription(newItem.enchant).then(description => {
				this.enchantElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.enchant.spellId) {
				this.enchantElem.href = ActionId.makeSpellUrl(newItem.enchant.spellId);
				this.enchantElem.dataset.wowhead = `domain=wotlk&spell=${newItem.enchant.spellId}`;
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
				this.enchantElem.dataset.wowhead = `domain=wotlk&item=${newItem.enchant.itemId}`;
			}
			this.enchantElem.dataset.whtticon = 'false';
		}

		newItem.allSocketColors().forEach((socketColor, gemIdx) => {
			const gemContainer = createGemContainer(socketColor, newItem.gems[gemIdx]);

			if (gemIdx == newItem.numPossibleSockets - 1 && [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(newItem.item.type)) {
				const updateProfession = () => {
					if (this.player.isBlacksmithing()) {
						gemContainer.classList.remove('hide');
					} else {
						gemContainer.classList.add('hide');
					}
				};
				this.player.professionChangeEmitter.on(updateProfession);
				updateProfession();
			}
			this.socketsContainerElem.appendChild(gemContainer);
		});
	}
}

export class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly itemElem: ItemRenderer;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];
	private _equippedItem: EquippedItem | null = null;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');
		this.slot = slot;
		this.simUI = simUI;
		this.player = player;
		this.itemElem = new ItemRenderer(parent, this.rootElem, player);

		this.item = player.getEquippedItem(slot);
		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(this.slot);
			this._enchants = this.player.getEnchants(this.slot);

			const gearData = {
				equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
					this.player.equipItem(eventID, this.slot, equippedItem);
				},
				getEquippedItem: () => this.player.getEquippedItem(this.slot),
				changeEvent: player.gearChangeEmitter,
			};

			const openGearSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Items, gearData);
			};
			const openEnchantSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Enchants, gearData);
			};

			this.itemElem.iconElem.addEventListener('click', openGearSelector);
			this.itemElem.nameElem.addEventListener('click', openGearSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantSelector);
		});

		player.gearChangeEmitter.on(() => {
			this.item = player.getEquippedItem(slot);
		});
		player.professionChangeEmitter.on(() => {
			if (this._equippedItem != null) {
				this.player.setWowheadData(this._equippedItem, this.itemElem.iconElem);
			}
		});
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.itemElem.clear();
		this.itemElem.nameElem.textContent = slotNames.get(this.slot) ?? '';
		setItemQualityCssClass(this.itemElem.nameElem, null);

		if (newItem != null) {
			this.itemElem.update(newItem);
		} else {
			this.itemElem.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		}

		this._equippedItem = newItem;
	}

	private openSelectorModal(tab: SelectorModalTabs, gearData: GearData) {
		new SelectorModal(this.simUI.rootElem, this.simUI, this.player, {
			selectedTab: tab,
			slot: this.slot,
			equippedItem: this._equippedItem,
			eligibleItems: this._items,
			eligibleEnchants: this._enchants,
			gearData: gearData,
		});
	}
}

export class IconItemSwapPicker extends Component {
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly socketsContainerElem: HTMLElement;
	private readonly player: Player<any>;
	private readonly slot: ItemSlot;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'icon-picker-root');
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.slot = slot;

		this.iconAnchor = document.createElement('a');
		this.iconAnchor.classList.add('icon-picker-button');
		this.iconAnchor.target = '_blank';
		this.rootElem.prepend(this.iconAnchor);

		this.socketsContainerElem = document.createElement('div');
		this.socketsContainerElem.classList.add('item-picker-sockets-container');
		this.iconAnchor.appendChild(this.socketsContainerElem);

		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(slot);
			this._enchants = this.player.getEnchants(slot);
			const gearData = {
				equipItem: (eventID: EventID, newItem: EquippedItem | null) => {
					player.equipItemSwapitem(eventID, this.slot, newItem);
				},
				getEquippedItem: () => player.getItemSwapItem(this.slot),
				changeEvent: player.itemSwapChangeEmitter,
			};

			this.iconAnchor.addEventListener('click', (event: Event) => {
				event.preventDefault();
				new SelectorModal(simUI.rootElem, simUI, this.player, {
					selectedTab: SelectorModalTabs.Items,
					slot: this.slot,
					equippedItem: this.player.getItemSwapGear().getEquippedItem(slot),
					eligibleItems: this._items,
					eligibleEnchants: this._enchants,
					gearData: gearData,
				});
			});
		});

		player.itemSwapChangeEmitter.on(() => {
			this.update(player.getItemSwapGear().getEquippedItem(slot));
		});
	}

	update(newItem: EquippedItem | null) {
		this.iconAnchor.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconAnchor.removeAttribute('data-wowhead');
		this.iconAnchor.href = '#';
		this.socketsContainerElem.innerText = '';

		if (newItem) {
			this.iconAnchor.classList.add('active');

			newItem.asActionId().fillAndSet(this.iconAnchor, true, true);
			this.player.setWowheadData(newItem, this.iconAnchor);

			newItem.allSocketColors().forEach((socketColor, gemIdx) => {
				this.socketsContainerElem.appendChild(createGemContainer(socketColor, newItem.gems[gemIdx]));
			});
		} else {
			this.iconAnchor.classList.remove('active');
		}
	}
}

export interface GearData {
	equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => void;
	getEquippedItem: () => EquippedItem | null;
	changeEvent: TypedEvent<any>;
}

export enum SelectorModalTabs {
	Items = 'Items',
	Enchants = 'Enchants',
	Gem1 = 'Gem1',
	Gem2 = 'Gem2',
	Gem3 = 'Gem3',
}

interface SelectorModalConfig {
	selectedTab: SelectorModalTabs;
	slot: ItemSlot;
	equippedItem: EquippedItem | null;
	eligibleItems: Array<Item>;
	eligibleEnchants: Array<Enchant>;
	gearData: GearData;
}

export class SelectorModal extends BaseModal {
	private readonly simUI: SimUI;
	private player: Player<any>;
	private config: SelectorModalConfig;
	private ilists: ItemList<any>[];

	private readonly tabsElem: HTMLElement;
	private readonly contentElem: HTMLElement;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, config: SelectorModalConfig) {
		super(parent, 'selector-modal');

		this.simUI = simUI;
		this.player = player;
		this.config = config;
		this.ilists = [];

		window.scrollTo({ top: 0 });

		this.header!.insertAdjacentElement('afterbegin', <ul className="nav nav-tabs selector-modal-tabs"></ul>);

		this.body.appendChild(<div className="tab-content selector-modal-tab-content"></div>);

		this.tabsElem = this.rootElem.querySelector('.selector-modal-tabs') as HTMLElement;
		this.contentElem = this.rootElem.querySelector('.selector-modal-tab-content') as HTMLElement;

		this.setData();

		this.body.appendChild(
			<div className="d-flex align-items-center form-text mt-3">
				<i className="fas fa-circle-exclamation fa-xl me-2"></i>
				<span>
					If gear is missing, check the selected phase and your gear filters.
					<br />
					If the problem persists, save any un-saved data, click the
					<i className="fas fa-cog mx-1"></i>
					to open your sim options, then click the "Restore Defaults".
				</span>
			</div>,
		);
	}

	// Could be 'Items' 'Enchants' or 'Gem1'-'Gem3'
	openTabName(name: string) {
		Array.from(this.tabsElem.getElementsByClassName('selector-modal-item-tab')).forEach(elem => {
			if (elem.getAttribute('data-content-id') == name + '-tab') {
				(elem as HTMLElement).click();
			}
		});
	}

	openTab(idx: number) {
		const elems = this.tabsElem.getElementsByClassName('selector-modal-item-tab');
		(elems[idx] as HTMLElement).click();
	}

	setData() {
		this.tabsElem.innerText = '';
		this.contentElem.innerText = '';

		const { slot, equippedItem, eligibleItems, eligibleEnchants, gearData } = this.config;

		this.addTab<Item>(
			'Items',
			eligibleItems.map(item => {
				return {
					item: item,
					id: item.id,
					actionId: ActionId.fromItem(item),
					name: item.name,
					quality: item.quality,
					heroic: item.heroic,
					phase: item.phase,
					baseEP: this.player.computeItemEP(item, slot),
					ignoreEPFilter: false,
					onEquip: (eventID, item: Item) => {
						const equippedItem = gearData.getEquippedItem();
						if (equippedItem) {
							gearData.equipItem(eventID, equippedItem.withItem(item));
						} else {
							gearData.equipItem(eventID, new EquippedItem(item));
						}
					},
				};
			}),
			item => this.player.computeItemEP(item, slot),
			equippedItem => equippedItem?.item,
			GemColor.GemColorUnknown,
			eventID => {
				gearData.equipItem(eventID, null);
				this.removeTabs('Gem');
			},
		);

		this.addTab<Enchant>(
			'Enchants',
			eligibleEnchants.map(enchant => {
				return {
					item: enchant,
					id: enchant.effectId,
					actionId: enchant.spellId ? ActionId.fromSpellId(enchant.spellId) : ActionId.fromItemId(enchant.itemId),
					name: enchant.name,
					quality: enchant.quality,
					phase: enchant.phase || 1,
					baseEP: this.player.computeStatsEP(new Stats(enchant.stats)),
					ignoreEPFilter: true,
					heroic: false,
					onEquip: (eventID, enchant: Enchant) => {
						const equippedItem = gearData.getEquippedItem();
						if (equippedItem) gearData.equipItem(eventID, equippedItem.withEnchant(enchant));
					},
				};
			}),
			enchant => this.player.computeEnchantEP(enchant),
			equippedItem => equippedItem?.enchant,
			GemColor.GemColorUnknown,
			eventID => {
				const equippedItem = gearData.getEquippedItem();
				if (equippedItem) gearData.equipItem(eventID, equippedItem.withEnchant(null));
			},
		);

		this.addGemTabs(slot, equippedItem, gearData);
	}

	protected override onShow(e: Event) {
		// Only refresh opened tab
		const t = e.target! as HTMLElement;
		const tab = t.querySelector<HTMLElement>('.active')!.dataset.contentId!;
		if (tab.includes('Item')) {
			this.ilists[0].sizeRefresh();
		} else if (tab.includes('Enchant')) {
			this.ilists[1].sizeRefresh();
		}
	}

	private addGemTabs(slot: ItemSlot, equippedItem: EquippedItem | null, gearData: GearData) {
		if (equippedItem == undefined) {
			return;
		}

		const socketBonusEP = this.player.computeStatsEP(new Stats(equippedItem.item.socketBonus)) / (equippedItem.item.gemSockets.length || 1);
		equippedItem.curSocketColors(this.player.isBlacksmithing()).forEach((socketColor, socketIdx) => {
			this.addTab<Gem>(
				'Gem ' + (socketIdx + 1),
				this.player.getGems(socketColor).map((gem: Gem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						heroic: false,
						baseEP: this.player.computeStatsEP(new Stats(gem.stats)),
						ignoreEPFilter: true,
						onEquip: (eventID, gem: Gem) => {
							const equippedItem = gearData.getEquippedItem();
							if (equippedItem) gearData.equipItem(eventID, equippedItem.withGem(gem, socketIdx));
						},
					};
				}),
				gem => {
					let gemEP = this.player.computeGemEP(gem);
					if (gemMatchesSocket(gem, socketColor)) {
						gemEP += socketBonusEP;
					}
					return gemEP;
				},
				equippedItem => equippedItem?.gems[socketIdx],
				socketColor,
				eventID => {
					const equippedItem = gearData.getEquippedItem();
					if (equippedItem) gearData.equipItem(eventID, equippedItem.withGem(null, socketIdx));
				},
				tabAnchor => {
					const gemContainer = createGemContainer(socketColor, null);
					tabAnchor.appendChild(gemContainer);
					tabAnchor.classList.add('selector-modal-tab-gem');

					const gemElem = tabAnchor.querySelector('.gem-icon') as HTMLElement;
					const emptySocketUrl = getEmptyGemSocketIconUrl(socketColor);

					const updateGemIcon = () => {
						const equippedItem = gearData.getEquippedItem();
						const gem = equippedItem?.gems[socketIdx];

						if (gem) {
							gemElem.classList.remove('hide');
							ActionId.fromItemId(gem.id)
								.fill()
								.then(filledId => {
									gemElem.setAttribute('src', filledId.iconUrl);
								});
						} else {
							gemElem.classList.add('hide');
							gemElem.setAttribute('src', emptySocketUrl);
						}
					};

					gearData.changeEvent.on(updateGemIcon);
					this.addOnDisposeCallback(() => gearData.changeEvent.off(updateGemIcon));
					updateGemIcon();
				},
			);
		});
	}

	/**
	 * Adds one of the tabs for the item selector menu.
	 *
	 * T is expected to be Item, Enchant, or Gem. Tab menus for all 3 looks extremely
	 * similar so this function uses extra functions to do it generically.
	 */
	private addTab<T>(
		label: string,
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => T | null | undefined,
		socketColor: GemColor,
		onRemove: (eventID: EventID) => void,
		setTabContent?: (tabElem: HTMLAnchorElement) => void,
	) {
		if (itemData.length == 0) {
			return;
		}

		const { slot, gearData } = this.config;
		const tabContentId = (label + '-tab').split(' ').join('');
		const selected = label === this.config.selectedTab;

		const tabAnchor = ref<HTMLAnchorElement>();
		this.tabsElem.appendChild(
			<li className="nav-item">
				<a
					ref={tabAnchor}
					className={`nav-link selector-modal-item-tab ${selected ? 'active' : ''}`}
					dataset={{
						label: label,
						contentId: tabContentId,
						bsToggle: 'tab',
						bsTarget: `#${tabContentId}`,
					}}
					attributes={{
						role: 'tab',
						'aria-selected': selected,
					}}
					type="button"></a>
			</li>,
		);

		if (setTabContent) {
			setTabContent(tabAnchor.value!);
		} else {
			tabAnchor.value!.textContent = label;
		}

		// TODO: do we need this check here?
		if (itemData.length == 0) {
			return;
		}

		const ilist = new ItemList<T>(
			this.contentElem,
			this.simUI,
			this.config,
			this.player,
			label,
			itemData,
			computeEP,
			equippedToItemFn,
			socketColor,
			onRemove,
			(itemData: ItemData<T>) => {
				const item = itemData.item;
				itemData.onEquip(TypedEvent.nextEventID(), item);

				// If the item changes, the gem slots might change, so remove and recreate the gem tabs
				if (Item.is(item)) {
					this.removeTabs('Gem');
					this.addGemTabs(slot, gearData.getEquippedItem(), gearData);
				}
			},
		);

		const invokeUpdate = () => {
			ilist.updateSelected();
		};
		const applyFilter = () => {
			ilist.applyFilters();
		};
		const hideOrShowEPValues = () => {
			ilist.hideOrShowEPValues();
		};
		// Add event handlers
		gearData.changeEvent.on(invokeUpdate);

		this.player.sim.phaseChangeEmitter.on(applyFilter);
		this.player.sim.filtersChangeEmitter.on(applyFilter);
		this.player.sim.showEPValuesChangeEmitter.on(hideOrShowEPValues);

		this.addOnDisposeCallback(() => {
			gearData.changeEvent.off(invokeUpdate);
			this.player.sim.phaseChangeEmitter.off(applyFilter);
			this.player.sim.filtersChangeEmitter.off(applyFilter);
			this.player.sim.showEPValuesChangeEmitter.off(hideOrShowEPValues);
			ilist.dispose();
		});

		tabAnchor.value!.addEventListener('shown.bs.tab', _event => {
			ilist.sizeRefresh();
		});

		this.ilists.push(ilist);
	}

	private removeTabs(labelSubstring: string) {
		const tabElems = Array.prototype.slice
			.call(this.tabsElem.getElementsByClassName('selector-modal-item-tab'))
			.filter(tab => tab.dataset.label.includes(labelSubstring));

		const contentElems = tabElems.map(tabElem => document.getElementById(tabElem.dataset.contentId!)).filter(tabElem => Boolean(tabElem));

		tabElems.forEach(elem => elem.parentElement.remove());
		contentElems.forEach(elem => elem!.remove());
	}
}

export interface ItemData<T> {
	item: T;
	name: string;
	id: number;
	actionId: ActionId;
	quality: ItemQuality;
	phase: number;
	baseEP: number;
	ignoreEPFilter: boolean;
	heroic: boolean;
	onEquip: (eventID: EventID, item: T) => void;
}

interface ItemDataWithIdx<T> {
	idx: number;
	data: ItemData<T>;
}

const emptySlotIcons: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: '/wotlk/assets/item_slots/head.jpg',
	[ItemSlot.ItemSlotNeck]: '/wotlk/assets/item_slots/neck.jpg',
	[ItemSlot.ItemSlotShoulder]: '/wotlk/assets/item_slots/shoulders.jpg',
	[ItemSlot.ItemSlotBack]: '/wotlk/assets/item_slots/shirt.jpg',
	[ItemSlot.ItemSlotChest]: '/wotlk/assets/item_slots/chest.jpg',
	[ItemSlot.ItemSlotWrist]: '/wotlk/assets/item_slots/wrists.jpg',
	[ItemSlot.ItemSlotHands]: '/wotlk/assets/item_slots/hands.jpg',
	[ItemSlot.ItemSlotWaist]: '/wotlk/assets/item_slots/waist.jpg',
	[ItemSlot.ItemSlotLegs]: '/wotlk/assets/item_slots/legs.jpg',
	[ItemSlot.ItemSlotFeet]: '/wotlk/assets/item_slots/feet.jpg',
	[ItemSlot.ItemSlotFinger1]: '/wotlk/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotFinger2]: '/wotlk/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: '/wotlk/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: '/wotlk/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: '/wotlk/assets/item_slots/mainhand.jpg',
	[ItemSlot.ItemSlotOffHand]: '/wotlk/assets/item_slots/offhand.jpg',
	[ItemSlot.ItemSlotRanged]: '/wotlk/assets/item_slots/ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}

export class ItemList<T> {
	private listElem: HTMLElement;
	private readonly player: Player<any>;
	private label: string;
	private slot: ItemSlot;
	private itemData: Array<ItemData<T>>;
	private itemsToDisplay: Array<number>;
	private currentFilters: DatabaseFilters;
	private searchInput: HTMLInputElement;
	private socketColor: GemColor;
	private computeEP: (item: T) => number;
	private equippedToItemFn: (equippedItem: EquippedItem | null) => T | null | undefined;
	private gearData: GearData;
	private tabContent: Element;
	private onItemClick: (itemData: ItemData<T>) => void;
	private scroller: Clusterize;

	constructor(
		parent: HTMLElement,
		simUI: SimUI,
		config: SelectorModalConfig,
		player: Player<any>,
		label: string,
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => T | null | undefined,
		socketColor: GemColor,
		onRemove: (eventID: EventID) => void,
		onItemClick: (itemData: ItemData<T>) => void,
	) {
		this.label = label;
		this.player = player;
		this.itemData = itemData;
		this.socketColor = socketColor;
		this.computeEP = computeEP;
		this.equippedToItemFn = equippedToItemFn;
		this.onItemClick = onItemClick;

		const { slot, gearData } = config;
		this.slot = slot;
		this.gearData = gearData;
		this.currentFilters = this.player.sim.getFilters();

		const tabContentId = (label + '-tab').split(' ').join('');
		const selected = label === config.selectedTab;

		const epButton = ref<HTMLButtonElement>();
		this.tabContent = (
			<div id={tabContentId} className={`selector-modal-tab-pane tab-pane fade ${selected ? 'active show' : ''}`}>
				<div className="selector-modal-filters">
					<input className="selector-modal-search form-control" type="text" placeholder="Search..." />
					{label == 'Items' && <button className="selector-modal-filters-button btn btn-primary">Filters</button>}
					<div className="selector-modal-phase-selector"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-1h-weapons"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-2h-weapons"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-matching-gems"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-ep-values"></div>
					<button className="selector-modal-simall-button btn btn-warning">Add to Batch Sim</button>
					<button className="selector-modal-remove-button btn btn-danger">Unequip Item</button>
				</div>
				<div className="selector-modal-list-labels">
					<label className="item-label">
						<small>Item</small>
					</label>
					<label className="source-label">
						<small>Source</small>
					</label>
					<label className="ilvl-label">
						<small>ILvl</small>
					</label>
					<label className="ep-label">
						<small>EP</small>
						<i className="fa-solid fa-plus-minus fa-2xs"></i>
						<button ref={epButton} className="btn btn-link p-0 ms-1">
							<i className="far fa-question-circle fa-lg"></i>
						</button>
					</label>
					<label className="favorite-label"></label>
				</div>
				<ul className="selector-modal-list"></ul>
			</div>
		);

		parent.appendChild(this.tabContent);

		new Tooltip(epButton.value!, {
			title: EP_TOOLTIP,
		});

		makeShow1hWeaponsSelector(this.tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement, player.sim);
		makeShow2hWeaponsSelector(this.tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement, player.sim);
		if (!(label == 'Items' && (slot == ItemSlot.ItemSlotMainHand || (slot == ItemSlot.ItemSlotOffHand && player.getClass() == Class.ClassWarrior)))) {
			(this.tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement).style.display = 'none';
			(this.tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement).style.display = 'none';
		}

		makeShowEPValuesSelector(this.tabContent.getElementsByClassName('selector-modal-show-ep-values')[0] as HTMLElement, player.sim);

		makeShowMatchingGemsSelector(this.tabContent.getElementsByClassName('selector-modal-show-matching-gems')[0] as HTMLElement, player.sim);
		if (!label.startsWith('Gem')) {
			(this.tabContent.getElementsByClassName('selector-modal-show-matching-gems')[0] as HTMLElement).style.display = 'none';
		}

		makePhaseSelector(this.tabContent.getElementsByClassName('selector-modal-phase-selector')[0] as HTMLElement, player.sim);

		if (label == 'Items') {
			const filtersButton = this.tabContent.getElementsByClassName('selector-modal-filters-button')[0] as HTMLElement;
			filtersButton.addEventListener('click', () => new FiltersMenu(parent, player, slot));
		}

		this.listElem = this.tabContent.getElementsByClassName('selector-modal-list')[0] as HTMLElement;

		this.itemsToDisplay = [];

		this.scroller = new Clusterize(
			{
				getNumberOfRows: () => {
					return this.itemsToDisplay.length;
				},
				generateRows: (startIdx, endIdx) => {
					const items = [];
					for (let i = startIdx; i < endIdx; ++i) {
						if (i >= this.itemsToDisplay.length) break;
						items.push(
							this.createItemElem({
								idx: this.itemsToDisplay[i],
								data: this.itemData[this.itemsToDisplay[i]],
							}),
						);
					}
					return items;
				},
			},
			{
				rows: [],
				scroll_elem: this.listElem,
				content_elem: this.listElem,
				item_height: 56,
				show_no_data_row: false,
				no_data_text: '',
				tag: 'li',
				rows_in_block: 16,
				blocks_in_cluster: 2,
			},
		);

		const removeButton = this.tabContent.getElementsByClassName('selector-modal-remove-button')[0] as HTMLButtonElement;
		removeButton.addEventListener('click', _event => {
			onRemove(TypedEvent.nextEventID());
		});

		if (label.startsWith('Enchants')) {
			removeButton.textContent = 'Remove Enchant';
		} else if (label.startsWith('Gem')) {
			removeButton.textContent = 'Remove Gem';
		}

		this.updateSelected();

		this.searchInput = this.tabContent.getElementsByClassName('selector-modal-search')[0] as HTMLInputElement;
		this.searchInput.addEventListener('input', () => this.applyFilters());

		const simAllButton = this.tabContent.getElementsByClassName('selector-modal-simall-button')[0] as HTMLButtonElement;
		if (label == 'Items') {
			simAllButton.hidden = !player.sim.getShowExperimental();
			player.sim.showExperimentalChangeEmitter.on(() => {
				simAllButton.hidden = !player.sim.getShowExperimental();
			});
			simAllButton.addEventListener('click', _event => {
				if (simUI instanceof IndividualSimUI) {
					const itemSpecs = Array<ItemSpec>();
					const isRangedOrTrinket =
						this.slot == ItemSlot.ItemSlotRanged || this.slot == ItemSlot.ItemSlotTrinket1 || this.slot == ItemSlot.ItemSlotTrinket2;

					const curItem = this.equippedToItemFn(this.player.getEquippedItem(this.slot));
					let curEP = 0;
					if (curItem != null) {
						curEP = this.computeEP(curItem);
					}

					for (const i of this.itemsToDisplay) {
						const idata = this.itemData[i];
						if (!isRangedOrTrinket && curEP > 0 && idata.baseEP < curEP / 2) {
							continue; // If we have EPs on current item, dont sim items with less than half the EP.
						}

						// Add any item that is either >0 EP or a trinket/ranged item.
						if (idata.baseEP > 0 || isRangedOrTrinket) {
							itemSpecs.push(ItemSpec.create({ id: idata.id }));
						}
					}
					simUI.bt.addItems(itemSpecs);
					// TODO: should we open the bulk sim UI or should we run in the background showing progress, and then sort the items in the picker?
				}
			});
		} else {
			// always hide non-items from being added to batch.
			simAllButton.hidden = true;
		}
	}

	public sizeRefresh() {
		this.scroller.refresh(true);
		this.applyFilters();
	}

	public dispose() {
		this.scroller.dispose();
	}

	public updateSelected() {
		const newEquippedItem = this.gearData.getEquippedItem();
		const newItem = this.equippedToItemFn(newEquippedItem);

		const newItemId = newItem ? (this.label == 'Enchants' ? (newItem as unknown as Enchant).effectId : (newItem as unknown as Item | Gem).id) : 0;
		const newEP = newItem ? this.computeEP(newItem) : 0;

		this.scroller.elementUpdate(item => {
			const idx = (item as HTMLElement).dataset.idx!;
			const itemData = this.itemData[parseFloat(idx)];
			if (itemData.id == newItemId) item.classList.add('active');
			else item.classList.remove('active');

			const epDeltaElem = item.getElementsByClassName('selector-modal-list-item-ep-delta')[0] as HTMLSpanElement;
			if (epDeltaElem) {
				epDeltaElem.textContent = '';
				if (itemData.item) {
					const listItemEP = this.computeEP(itemData.item);
					if (newEP != listItemEP) {
						formatDeltaTextElem(epDeltaElem, newEP, listItemEP, 0);
					}
				}
			}
		});
	}

	public applyFilters() {
		this.currentFilters = this.player.sim.getFilters();
		let itemIdxs = new Array<number>(this.itemData.length);
		for (let i = 0; i < this.itemData.length; ++i) {
			itemIdxs[i] = i;
		}

		const currentEquippedItem = this.player.getEquippedItem(this.slot);

		if (this.label == 'Items') {
			itemIdxs = this.player.filterItemData(itemIdxs, i => this.itemData[i].item as unknown as Item, this.slot);
		} else if (this.label == 'Enchants') {
			itemIdxs = this.player.filterEnchantData(itemIdxs, i => this.itemData[i].item as unknown as Enchant, this.slot, currentEquippedItem);
		} else if (this.label.startsWith('Gem')) {
			itemIdxs = this.player.filterGemData(itemIdxs, i => this.itemData[i].item as unknown as Gem, this.slot, this.socketColor);
		}

		itemIdxs = itemIdxs.filter(i => {
			const listItemData = this.itemData[i];

			if (listItemData.phase > this.player.sim.getPhase()) {
				return false;
			}

			if (this.searchInput.value.length > 0) {
				const searchQuery = this.searchInput.value
					.toLowerCase()
					.replaceAll(/[^a-zA-Z0-9\s]/g, '')
					.split(' ');
				const name = listItemData.name.toLowerCase().replaceAll(/[^a-zA-Z0-9\s]/g, '');

				let include = true;
				searchQuery.forEach(v => {
					if (!name.includes(v)) include = false;
				});
				if (!include) {
					return false;
				}
			}

			return true;
		});

		let sortFn: (itemA: T, itemB: T) => number;
		if (this.slot == ItemSlot.ItemSlotTrinket1 || this.slot == ItemSlot.ItemSlotTrinket2) {
			// Trinket EP is weird so just sort by ilvl instead.
			sortFn = (itemA, itemB) => (itemB as unknown as Item).ilvl - (itemA as unknown as Item).ilvl;
		} else {
			sortFn = (itemA, itemB) => {
				const diff = this.computeEP(itemB) - this.computeEP(itemA);
				// if EP is same, sort by ilvl
				if (Math.abs(diff) < 0.01) return (itemB as unknown as Item).ilvl - (itemA as unknown as Item).ilvl;
				return diff;
			};
		}

		itemIdxs = itemIdxs.sort((dataA, dataB) => {
			const itemA = this.itemData[dataA];
			const itemB = this.itemData[dataB];
			if (this.isItemFavorited(itemA) && !this.isItemFavorited(itemB)) return -1;
			if (this.isItemFavorited(itemB) && !this.isItemFavorited(itemA)) return 1;

			return sortFn(itemA.item, itemB.item);
		});

		this.itemsToDisplay = itemIdxs;
		this.scroller.update();

		this.hideOrShowEPValues();
	}

	public hideOrShowEPValues() {
		const labels = this.tabContent.getElementsByClassName('ep-label');
		const container = this.tabContent.getElementsByClassName('selector-modal-list');
		const show = this.player.sim.getShowEPValues();
		const display = show ? '' : 'none';

		for (const label of labels) {
			(label as HTMLElement).style.display = display;
		}

		for (const c of container) {
			if (show) c.classList.remove('hide-ep');
			else c.classList.add('hide-ep');
		}
	}

	private createItemElem(item: ItemDataWithIdx<T>): JSX.Element {
		const itemData = item.data;
		const itemEP = this.computeEP(itemData.item);

		const equippedItem = this.equippedToItemFn(this.gearData.getEquippedItem());
		const equippedItemID = equippedItem
			? this.label == 'Enchants'
				? (equippedItem as unknown as Enchant).effectId
				: (equippedItem as unknown as Item).id
			: 0;
		const equippedItemEP = equippedItem ? this.computeEP(equippedItem) : 0;

		const nameElem = ref<HTMLLabelElement>();
		const anchorElem = ref<HTMLAnchorElement>();
		const iconElem = ref<HTMLImageElement>();
		const listItemElem = (
			<li className={`selector-modal-list-item ${equippedItemID == itemData.id ? 'active' : ''}`} dataset={{ idx: item.idx.toString() }}>
				<div className="selector-modal-list-label-cell">
					<a className="selector-modal-list-item-link" ref={anchorElem} dataset={{ whtticon: 'false' }}>
						<img className="selector-modal-list-item-icon" ref={iconElem}></img>
						<label className="selector-modal-list-item-name" ref={nameElem}>
							{itemData.name}
							{itemData.heroic && createHeroicLabel()}
						</label>
					</a>
				</div>
			</li>
		);

		if (this.label == 'Items') {
			listItemElem.appendChild(
				<div className="selector-modal-list-item-source-container">{this.getSourceInfo(itemData.item as unknown as Item, this.player.sim)}</div>,
			);
			listItemElem.appendChild(
				<div className="selector-modal-list-item-ilvl-container">{(itemData.item as unknown as Item).ilvl}</div>,
			);
		}

		if (this.slot != ItemSlot.ItemSlotTrinket1 && this.slot != ItemSlot.ItemSlotTrinket2) {
			listItemElem.appendChild(
				<div className="selector-modal-list-item-ep">
					<span className="selector-modal-list-item-ep-value">{itemEP < 9.95 ? itemEP.toFixed(1).toString() : Math.round(itemEP).toString()}</span>
					<span
						className="selector-modal-list-item-ep-delta"
						ref={e => itemData.item && equippedItemEP != itemEP && formatDeltaTextElem(e, equippedItemEP, itemEP, 0)}></span>
				</div>,
			);
		}

		const favoriteElem = ref<HTMLButtonElement>();
		listItemElem.appendChild(
			<div>
				<button
					className="selector-modal-list-item-favorite btn btn-link p-0"
					ref={favoriteElem}
					onclick={() => setFavorite(listItemElem.dataset.fav == 'false')}>
					<i className="fa-star fa-xl"></i>
				</button>
			</div>,
		);

		anchorElem.value!.addEventListener('click', (event: Event) => {
			event.preventDefault();
			if (event.target === favoriteElem.value) return false;
			this.onItemClick(itemData);
		});

		itemData.actionId.fill().then(filledId => {
			filledId.setWowheadHref(anchorElem.value!);
			iconElem.value!.src = filledId.iconUrl;
		});

		setItemQualityCssClass(nameElem.value!, itemData.quality);

		new Tooltip(favoriteElem.value!, {
			title: 'Add to favorites',
		});
		const setFavorite = (isFavorite: boolean) => {
			const filters = this.player.sim.getFilters();
			if (this.label == 'Items') {
				const favId = itemData.id;
				if (isFavorite) {
					filters.favoriteItems.push(favId);
				} else {
					const favIdx = filters.favoriteItems.indexOf(favId);
					if (favIdx != -1) {
						filters.favoriteItems.splice(favIdx, 1);
					}
				}
			} else if (this.label == 'Enchants') {
				const favId = getUniqueEnchantString(itemData.item as unknown as Enchant);
				if (isFavorite) {
					filters.favoriteEnchants.push(favId);
				} else {
					const favIdx = filters.favoriteEnchants.indexOf(favId);
					if (favIdx != -1) {
						filters.favoriteEnchants.splice(favIdx, 1);
					}
				}
			} else if (this.label.startsWith('Gem')) {
				const favId = itemData.id;
				if (isFavorite) {
					filters.favoriteGems.push(favId);
				} else {
					const favIdx = filters.favoriteGems.indexOf(favId);
					if (favIdx != -1) {
						filters.favoriteGems.splice(favIdx, 1);
					}
				}
			}
			favoriteElem.value!.children[0].classList.toggle('fas');
			favoriteElem.value!.children[0].classList.toggle('far');
			listItemElem.dataset.fav = isFavorite.toString();

			this.player.sim.setFilters(TypedEvent.nextEventID(), filters);
		};

		const isFavorite = this.isItemFavorited(itemData);

		if (isFavorite) {
			favoriteElem.value!.children[0].classList.add('fas');
			listItemElem.dataset.fav = 'true';
		} else {
			favoriteElem.value!.children[0].classList.add('far');
			listItemElem.dataset.fav = 'false';
		}

		return listItemElem;
	}

	private isItemFavorited(itemData: ItemData<T>): boolean {
		if (this.label == 'Items') {
			return this.currentFilters.favoriteItems.includes(itemData.id);
		} else if (this.label == 'Enchants') {
			return this.currentFilters.favoriteEnchants.includes(getUniqueEnchantString(itemData.item as unknown as Enchant));
		} else if (this.label.startsWith('Gem')) {
			return this.currentFilters.favoriteGems.includes(itemData.id);
		}
		return false;
	}

	private getSourceInfo(item: Item, sim: Sim): JSX.Element {
		const makeAnchor = (href: string, inner: string | JSX.Element) => {
			return (
				<a href={href} target="_blank" dataset={{ whtticon: 'false' }}>
					<small>{inner}</small>
				</a>
			);
		};

		if (!item.sources || item.sources.length == 0) {
			return <></>;
		}

		let source = item.sources[0];
		if (source.source.oneofKind == 'crafted') {
			const src = source.source.crafted;

			if (src.spellId) {
				return makeAnchor(ActionId.makeSpellUrl(src.spellId), professionNames.get(src.profession) ?? 'Unknown');
			}
			return makeAnchor(ActionId.makeItemUrl(item.id), professionNames.get(src.profession) ?? 'Unknown');
		} else if (source.source.oneofKind == 'drop') {
			const src = source.source.drop;
			const zone = sim.db.getZone(src.zoneId);
			const npc = sim.db.getNpc(src.npcId);
			if (!zone) {
				throw new Error('No zone found for item: ' + item);
			}

			const category = src.category ? ` - ${src.category}` : '';
			if (npc) {
				return makeAnchor(
					ActionId.makeNpcUrl(npc.id),
					<span>
						{zone.name} ({difficultyNames.get(src.difficulty) ?? 'Unknown'})
						<br />
						{npc.name + category}
					</span>,
				);
			} else if (src.otherName) {
				return makeAnchor(
					ActionId.makeZoneUrl(zone.id),
					<span>
						{zone.name}
						<br />
						{src.otherName}
					</span>,
				);
			}
			return makeAnchor(ActionId.makeZoneUrl(zone.id), zone.name);
		} else if (source.source.oneofKind == 'quest' && source.source.quest.name) {
			const src = source.source.quest;
			return makeAnchor(
				ActionId.makeQuestUrl(src.id),
				<span>
					Quest
					{item.factionRestriction == UIItem_FactionRestriction.ALLIANCE_ONLY && (
						<img src="/wotlk/assets/img/alliance.png" className="ms-1" width="15" height="15" />
					)}
					{item.factionRestriction == UIItem_FactionRestriction.HORDE_ONLY && (
						<img src="/wotlk/assets/img/horde.png" className="ms-1" width="15" height="15" />
					)}
					<br />
					{src.name}
				</span>,
			);
		} else if ((source = item.sources.find(source => source.source.oneofKind == 'rep') ?? source).source.oneofKind == 'rep') {
			const factionNames = item.sources
				.filter(source => source.source.oneofKind == 'rep')
				.map(source =>
					source.source.oneofKind == 'rep' ? REP_FACTION_NAMES[source.source.rep.repFactionId] : REP_FACTION_NAMES[RepFaction.RepFactionUnknown],
				);
			const src = source.source.rep;
			return makeAnchor(
				ActionId.makeItemUrl(item.id),
				<>
					{factionNames.map(name => (
						<span>
							{name}
							{item.factionRestriction == UIItem_FactionRestriction.ALLIANCE_ONLY && (
								<img src="/wotlk/assets/img/alliance.png" className="ms-1" width="15" height="15" />
							)}
							{item.factionRestriction == UIItem_FactionRestriction.HORDE_ONLY && (
								<img src="/wotlk/assets/img/horde.png" className="ms-1" width="15" height="15" />
							)}
							<br />
						</span>
					))}
					<span>{REP_LEVEL_NAMES[src.repLevel]}</span>
				</>,
			);
		} else if (source.source.oneofKind == 'soldBy') {
			const src = source.source.soldBy;
			return makeAnchor(
				ActionId.makeNpcUrl(src.npcId),
				<span>
					Sold by
					<br />
					{src.npcName}
				</span>,
			);
		}
		return <></>;
	}
}
