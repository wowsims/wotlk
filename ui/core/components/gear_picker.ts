import { getWowheadItemId } from '/tbc/core/proto_utils/equipped_item.js';
import { EquippedItem } from '/tbc/core/proto_utils/equipped_item.js';
import { getEmptyGemSocketIconUrl, gemMatchesSocket } from '/tbc/core/proto_utils/gems.js';
import { setGemSocketCssClass } from '/tbc/core/proto_utils/gems.js';
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { enchantAppliesToItem } from '/tbc/core/proto_utils/utils.js';
import { Enchant, Gem, GemColor } from '/tbc/core/proto/common.js';
import { HandType } from '/tbc/core/proto/common.js';
import { WeaponType } from '/tbc/core/proto/common.js';
import { Item } from '/tbc/core/proto/common.js';
import { ItemQuality } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { enchantDescriptions } from '/tbc/core/constants/enchants.js';
import { ActionId } from '/tbc/core/proto_utils/action_id.js';
import { slotNames } from '/tbc/core/proto_utils/names.js';
import { setItemQualityCssClass } from '/tbc/core/css_utils.js';
import { Player } from '/tbc/core/player.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { formatDeltaTextElem } from '/tbc/core/utils.js';
import { getEnumValues } from '/tbc/core/utils.js';

import { Component } from './component.js';
import { Popup } from './popup.js';
import { makePhaseSelector } from './other_inputs.js';
import { makeShow1hWeaponsSelector } from './other_inputs.js';
import { makeShow2hWeaponsSelector } from './other_inputs.js';
import { makeShowMatchingGemsSelector } from './other_inputs.js';

declare var $: any;

export class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;

	constructor(parent: HTMLElement, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSide = document.createElement('div');
		leftSide.classList.add('gear-picker-left');
		this.rootElem.appendChild(leftSide);

		const rightSide = document.createElement('div');
		rightSide.classList.add('gear-picker-right');
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
		].map(slot => new ItemPicker(leftSide, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
			ItemSlot.ItemSlotRanged,
		].map(slot => new ItemPicker(rightSide, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);
	}
}

class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly player: Player<any>;
	private readonly iconElem: HTMLAnchorElement;
	private readonly nameElem: HTMLElement;
	private readonly enchantElem: HTMLElement;
	private readonly socketsContainerElem: HTMLElement;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];

	private _equippedItem: EquippedItem | null = null;


	constructor(parent: HTMLElement, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');
		this.slot = slot;
		this.player = player;

		this.rootElem.innerHTML = `
      <a class="item-picker-icon">
        <div class="item-picker-sockets-container">
        </div>
      </a>
      <div class="item-picker-labels-container">
        <span class="item-picker-name"></span>
        <span class="item-picker-enchant"></span>
      </div>
    `;

		this.iconElem = this.rootElem.getElementsByClassName('item-picker-icon')[0] as HTMLAnchorElement;
		this.nameElem = this.rootElem.getElementsByClassName('item-picker-name')[0] as HTMLElement;
		this.enchantElem = this.rootElem.getElementsByClassName('item-picker-enchant')[0] as HTMLElement;
		this.socketsContainerElem = this.rootElem.getElementsByClassName('item-picker-sockets-container')[0] as HTMLElement;

		this.item = player.getEquippedItem(slot);
		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(this.slot);
			this._enchants = this.player.getEnchants(this.slot);

			this.iconElem.addEventListener('click', event => {
				event.preventDefault();
				const selectorModal = new SelectorModal(this.rootElem.closest('.individual-sim-ui')!, this.player, this.slot, this._equippedItem, this._items, this._enchants);
			});
			this.iconElem.addEventListener('touchstart', event => {
				event.preventDefault();
				const selectorModal = new SelectorModal(this.rootElem.closest('.individual-sim-ui')!, this.player, this.slot, this._equippedItem, this._items, this._enchants);
			});
			this.iconElem.addEventListener('touchend', event => {
				event.preventDefault();
			});
		});
		player.gearChangeEmitter.on(() => {
			this.item = player.getEquippedItem(slot);
		});
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');

		this.nameElem.textContent = slotNames[this.slot];
		setItemQualityCssClass(this.nameElem, null);

		this.enchantElem.textContent = '';
		this.socketsContainerElem.innerHTML = '';

		if (newItem != null) {
			this.nameElem.textContent = newItem.item.name;
			setItemQualityCssClass(this.nameElem, newItem.item.quality);

			this.player.setWowheadData(newItem, this.iconElem);
			newItem.asActionId().fillAndSet(this.iconElem, true, true);

			if (newItem.enchant) {
				this.enchantElem.textContent = enchantDescriptions.get(newItem.enchant.id) || newItem.enchant.name;
			}

			newItem.item.gemSockets.forEach((socketColor, gemIdx) => {
				const gemIconElem = document.createElement('img');
				gemIconElem.classList.add('item-picker-gem-icon');
				setGemSocketCssClass(gemIconElem, socketColor);
				if (newItem.gems[gemIdx] == null) {
					gemIconElem.src = getEmptyGemSocketIconUrl(socketColor);
				} else {
					ActionId.fromItemId(newItem.gems[gemIdx]!.id).fill().then(filledId => {
						gemIconElem.src = filledId.iconUrl;
					});
				}
				this.socketsContainerElem.appendChild(gemIconElem);
			});
		}
		this._equippedItem = newItem;
	}
}

class SelectorModal extends Popup {
	private player: Player<any>;
	private readonly tabsElem: HTMLElement;
	private readonly contentElem: HTMLElement;

	constructor(parent: HTMLElement, player: Player<any>, slot: ItemSlot, equippedItem: EquippedItem | null, eligibleItems: Array<Item>, eligibleEnchants: Array<Enchant>) {
		super(parent);
		this.player = player;

		this.rootElem.classList.add('selector-modal');
		this.rootElem.innerHTML = `
			<ul class="nav nav-tabs selector-modal-tabs">
			</ul>
			<div class="tab-content selector-modal-tab-content">
			</div>
		`;

		this.addCloseButton();
		this.tabsElem = this.rootElem.getElementsByClassName('selector-modal-tabs')[0] as HTMLElement;
		this.contentElem = this.rootElem.getElementsByClassName('selector-modal-tab-content')[0] as HTMLElement;

		this.setData(slot, equippedItem, eligibleItems, eligibleEnchants);
	}

	setData(slot: ItemSlot, equippedItem: EquippedItem | null, eligibleItems: Array<Item>, eligibleEnchants: Array<Enchant>) {
		this.tabsElem.innerHTML = '';
		this.contentElem.innerHTML = '';

		this.addTab(
			'Items',
			slot,
			equippedItem,
			eligibleItems.map(item => {
				return {
					item: item,
					id: item.id,
					actionId: ActionId.fromItem(item),
					name: item.name,
					quality: item.quality,
					phase: item.phase,
					baseEP: this.player.computeItemEP(item),
					ignoreEPFilter: false,
					onEquip: (eventID, item: Item) => {
						const equippedItem = this.player.getEquippedItem(slot);
						if (equippedItem) {
							this.player.equipItem(eventID, slot, equippedItem.withItem(item));
						} else {
							this.player.equipItem(eventID, slot, new EquippedItem(item));
						}
					},
				};
			}),
			item => this.player.computeItemEP(item),
			equippedItem => equippedItem?.item,
			GemColor.GemColorUnknown,
			eventID => {
				this.player.equipItem(eventID, slot, null);
			});

		this.addTab(
			'Enchants',
			slot,
			equippedItem,
			eligibleEnchants.map(enchant => {
				return {
					item: enchant,
					id: enchant.id,
					actionId: enchant.isSpellId ? ActionId.fromSpellId(enchant.id) : ActionId.fromItemId(enchant.id),
					name: enchant.name,
					quality: enchant.quality,
					phase: enchant.phase || 1,
					baseEP: this.player.computeStatsEP(new Stats(enchant.stats)),
					ignoreEPFilter: true,
					onEquip: (eventID, enchant: Enchant) => {
						const equippedItem = this.player.getEquippedItem(slot);
						if (equippedItem)
							this.player.equipItem(eventID, slot, equippedItem.withEnchant(enchant));
					},
				};
			}),
			enchant => this.player.computeEnchantEP(enchant),
			equippedItem => equippedItem?.enchant,
			GemColor.GemColorUnknown,
			eventID => {
				const equippedItem = this.player.getEquippedItem(slot);
				if (equippedItem)
					this.player.equipItem(eventID, slot, equippedItem.withEnchant(null));
			});

		this.addGemTabs(slot, equippedItem);
	}

	private addGemTabs(slot: ItemSlot, equippedItem: EquippedItem | null) {
		if (equippedItem == undefined) {
			return;
		}

		const socketBonusEP = this.player.computeStatsEP(new Stats(equippedItem.item.socketBonus)) / equippedItem.item.gemSockets.length;
		equippedItem.item.gemSockets.forEach((socketColor, socketIdx) => {
			this.addTab(
				'Gem ' + (socketIdx + 1),
				slot,
				equippedItem,
				this.player.getGems(socketColor).map((gem: Gem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						baseEP: this.player.computeStatsEP(new Stats(gem.stats)),
						ignoreEPFilter: true,
						onEquip: (eventID, gem: Gem) => {
							const equippedItem = this.player.getEquippedItem(slot);
							if (equippedItem)
								this.player.equipItem(eventID, slot, equippedItem.withGem(gem, socketIdx));
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
					const equippedItem = this.player.getEquippedItem(slot);
					if (equippedItem)
						this.player.equipItem(eventID, slot, equippedItem.withGem(null, socketIdx));
				},
				tabAnchor => {
					tabAnchor.classList.add('selector-modal-tab-gem-icon');
					setGemSocketCssClass(tabAnchor, socketColor);

					const updateGemIcon = () => {
						const equippedItem = this.player.getEquippedItem(slot);
						const gem = equippedItem?.gems[socketIdx];

						if (gem) {
							ActionId.fromItemId(gem.id).fill().then(filledId => {
								tabAnchor.style.backgroundImage = `url('${filledId.iconUrl}')`;
							});
						} else {
							const url = getEmptyGemSocketIconUrl(socketColor);
							tabAnchor.style.backgroundImage = `url('${url}')`;
						}
					};

					this.player.gearChangeEmitter.on(updateGemIcon);
					updateGemIcon();
				});
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
		slot: ItemSlot,
		equippedItem: EquippedItem | null,
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => (T | null | undefined),
		socketColor: GemColor,
		onRemove: (eventID: EventID) => void,
		setTabContent?: (tabElem: HTMLAnchorElement) => void) {
		if (itemData.length == 0) {
			return;
		}

		if (slot == ItemSlot.ItemSlotTrinket1 || slot == ItemSlot.ItemSlotTrinket2) {
			// Trinket EP is weird so just sort by ilvl instead.
			itemData.sort((dataA, dataB) => (dataB.item as unknown as Item).ilvl - (dataA.item as unknown as Item).ilvl);
		} else {
			itemData.sort((dataA, dataB) => computeEP(dataB.item) - computeEP(dataA.item));
		}

		const tabElem = document.createElement('li');
		this.tabsElem.appendChild(tabElem);
		const tabContentId = (label + '-tab').split(' ').join('');
		tabElem.innerHTML = `<a class="selector-modal-item-tab" data-toggle="tab" href="#${tabContentId}"></a>`;

		const tabAnchor = tabElem.getElementsByClassName('selector-modal-item-tab')[0] as HTMLAnchorElement;
		tabAnchor.dataset.label = label;
		if (setTabContent) {
			setTabContent(tabAnchor);
		} else {
			tabAnchor.textContent = label;
		}

		const tabContent = document.createElement('div');
		tabContent.id = tabContentId;
		tabContent.classList.add('tab-pane', 'fade', 'selector-modal-tab-content');
		this.contentElem.appendChild(tabContent);
		tabContent.innerHTML = `
    <div class="selector-modal-tab-content-header">
      <button class="selector-modal-remove-button sim-button">Remove</button>
      <input class="selector-modal-search" type="text" placeholder="Search...">
      <div class="selector-modal-filter-bar-filler"></div>
      <div class="sim-input selector-modal-boolean-option selector-modal-show-1h-weapons"></div>
      <div class="sim-input selector-modal-boolean-option selector-modal-show-2h-weapons"></div>
      <div class="sim-input selector-modal-boolean-option selector-modal-show-matching-gems"></div>
      <div class="selector-modal-phase-selector"></div>
    </div>
    <ul class="selector-modal-list"></ul>
    `;

		const show1hWeaponsSelector = makeShow1hWeaponsSelector(tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement, this.player.sim);
		const show2hWeaponsSelector = makeShow2hWeaponsSelector(tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement, this.player.sim);
		if (label != 'Items' || slot != ItemSlot.ItemSlotMainHand) {
			(tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement).style.display = 'none';
			(tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement).style.display = 'none';
		}

		const showMatchingGemsSelector = makeShowMatchingGemsSelector(tabContent.getElementsByClassName('selector-modal-show-matching-gems')[0] as HTMLElement, this.player.sim);
		if (!label.startsWith('Gem')) {
			(tabContent.getElementsByClassName('selector-modal-show-matching-gems')[0] as HTMLElement).style.display = 'none';
		}

		const phaseSelector = makePhaseSelector(tabContent.getElementsByClassName('selector-modal-phase-selector')[0] as HTMLElement, this.player.sim);

		if (label == 'Items') {
			tabElem.classList.add('active', 'in');
			tabContent.classList.add('active', 'in');
		}

		const listElem = tabContent.getElementsByClassName('selector-modal-list')[0] as HTMLElement;

		const listItemElems = itemData.map((itemData, itemIdx) => {
			const item = itemData.item;
			const itemEP = computeEP(item);

			const listItemElem = document.createElement('li');
			listItemElem.classList.add('selector-modal-list-item');
			listElem.appendChild(listItemElem);

			listItemElem.dataset.idx = String(itemIdx);

			listItemElem.innerHTML = `
        <a class="selector-modal-list-item-icon"></a>
        <a class="selector-modal-list-item-name">${itemData.name}</a>
        <div class="selector-modal-list-item-padding"></div>
        <div class="selector-modal-list-item-ep">
					<span class="selector-modal-list-item-ep-value">${itemEP < 9.95 ? itemEP.toFixed(1) : Math.round(itemEP)}</span>
					<span class="selector-modal-list-item-ep-delta"></span>
				</div>
      `;

			if (slot == ItemSlot.ItemSlotTrinket1 || slot == ItemSlot.ItemSlotTrinket2) {
				const epElem = listItemElem.getElementsByClassName('selector-modal-list-item-ep')[0] as HTMLElement;
				epElem.style.display = 'none';
			}

			const iconElem = listItemElem.getElementsByClassName('selector-modal-list-item-icon')[0] as HTMLImageElement;
			itemData.actionId.fill().then(filledId => {
				filledId.setWowheadHref(listItemElem.children[0] as HTMLAnchorElement);
				filledId.setWowheadHref(listItemElem.children[1] as HTMLAnchorElement);
				iconElem.style.backgroundImage = `url('${filledId.iconUrl}')`;
			});

			const nameElem = listItemElem.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLImageElement;
			setItemQualityCssClass(nameElem, itemData.quality);

			const onclick = (event: Event) => {
				event.preventDefault();
				itemData.onEquip(TypedEvent.nextEventID(), item);

				// If the item changes, the gem slots might change, so remove and recreate the gem tabs
				if (Item.is(item)) {
					this.removeTabs('Gem');
					this.addGemTabs(slot, this.player.getEquippedItem(slot));
				}
			};
			nameElem.addEventListener('click', onclick);
			iconElem.addEventListener('click', onclick);

			return listItemElem;
		});

		const removeButton = tabContent.getElementsByClassName('selector-modal-remove-button')[0] as HTMLButtonElement;
		removeButton.addEventListener('click', event => {
			listItemElems.forEach(elem => elem.classList.remove('active'));
			onRemove(TypedEvent.nextEventID());
		});

		const updateSelected = () => {
			const newEquippedItem = this.player.getEquippedItem(slot);
			const newItem = equippedToItemFn(newEquippedItem);

			const newItemId = (newItem as any)?.id || null;
			const newEP = newItem ? computeEP(newItem) : 0;

			listItemElems.forEach(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = itemData[listItemIdx];
				const listItem = listItemData.item;

				elem.classList.remove('active');
				if (listItemData.id == newItemId) {
					elem.classList.add('active');
				}

				const epDeltaElem = elem.getElementsByClassName('selector-modal-list-item-ep-delta')[0] as HTMLSpanElement;
				epDeltaElem.textContent = '';
				if (listItem) {
					const listItemEP = computeEP(listItem);
					formatDeltaTextElem(epDeltaElem, newEP, listItemEP, 0);
				}
			});
		};
		updateSelected();
		this.player.gearChangeEmitter.on(updateSelected);

		const applyFilters = () => {
			let validItemElems = listItemElems;
			const currentEquippedItem = this.player.getEquippedItem(slot);

			validItemElems = validItemElems.filter(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = itemData[listItemIdx];

				if (label == 'Items') {
					const listItem = listItemData.item as unknown as Item;
					if (!this.player.sim.getShow1hWeapons() && listItem.weaponType != WeaponType.WeaponTypeUnknown && listItem.handType != HandType.HandTypeTwoHand) {
						return false;
					}
					if (!this.player.sim.getShow2hWeapons() && listItem.weaponType != WeaponType.WeaponTypeUnknown && listItem.handType == HandType.HandTypeTwoHand) {
						return false;
					}
				} else if (label == 'Enchants') {
					const listItem = listItemData.item as unknown as Enchant;
					if (currentEquippedItem && !enchantAppliesToItem(listItem, currentEquippedItem.item)) {
						return false;
					}
				} else if (label.startsWith('Gem')) {
					const listItem = listItemData.item as unknown as Gem;
					if (this.player.sim.getShowMatchingGems() && !gemMatchesSocket(listItem, socketColor)) {
						return false;
					}
				}

				if (listItemData.phase > this.player.sim.getPhase()) {
					return false;
				}

				if (searchInput.value.length > 0) {
					const searchQuery = searchInput.value.toLowerCase().split(" ");
					const name = listItemData.name.toLowerCase();

					var include = true;
					searchQuery.forEach(v => {
						if (!name.includes(v))
							include = false;
					});
					if (!include) {
						return false;
					}
				}

				return true;
			});

			let numShown = 0;
			listItemElems.forEach(elem => {
				if (validItemElems.includes(elem)) {
					elem.classList.remove('hidden');
					numShown++;
					if (numShown % 2 == 0) {
						elem.classList.remove('odd');
					} else {
						elem.classList.add('odd');
					}
				} else {
					elem.classList.add('hidden');
				}
			});
		};

		const searchInput = tabContent.getElementsByClassName('selector-modal-search')[0] as HTMLInputElement;
		searchInput.addEventListener('input', applyFilters);
		searchInput.addEventListener("keyup", ev => {
			if (ev.key == "Enter") {
				listItemElems.find(ele => {
					if (ele.classList.contains("hidden")) {
						return false;
					}
					const nameElem = ele.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLElement;
					nameElem.click();
					return true;
				});
			}
		});

		this.player.sim.phaseChangeEmitter.on(() => {
			applyFilters();
		});
		TypedEvent.onAny([
			this.player.sim.show1hWeaponsChangeEmitter,
			this.player.sim.show2hWeaponsChangeEmitter,
			this.player.sim.showMatchingGemsChangeEmitter,
		]).on(() => {
			applyFilters();
		});
		this.player.gearChangeEmitter.on(() => {
			applyFilters();
			updateSelected();
		});

		applyFilters();
	}

	private removeTabs(labelSubstring: string) {
		const tabElems = Array.prototype.slice.call(this.tabsElem.getElementsByClassName('selector-modal-item-tab'))
			.filter(tab => tab.dataset.label.includes(labelSubstring));

		const contentElems = tabElems
			.map(tabElem => document.getElementById(tabElem.getAttribute('href').substring(1)))
			.filter(tabElem => Boolean(tabElem));

		tabElems.forEach(elem => elem.parentElement.remove());
		contentElems.forEach(elem => elem!.remove());
	}
}

interface ItemData<T> {
	item: T,
	name: string,
	id: number,
	actionId: ActionId,
	quality: ItemQuality,
	phase: number,
	baseEP: number,
	ignoreEPFilter: boolean,
	onEquip: (eventID: EventID, item: T) => void,
}

const emptySlotIcons: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: 'https://cdn.seventyupgrades.com/item-slots/Head.jpg',
	[ItemSlot.ItemSlotNeck]: 'https://cdn.seventyupgrades.com/item-slots/Neck.jpg',
	[ItemSlot.ItemSlotShoulder]: 'https://cdn.seventyupgrades.com/item-slots/Shoulders.jpg',
	[ItemSlot.ItemSlotBack]: 'https://cdn.seventyupgrades.com/item-slots/Back.jpg',
	[ItemSlot.ItemSlotChest]: 'https://cdn.seventyupgrades.com/item-slots/Chest.jpg',
	[ItemSlot.ItemSlotWrist]: 'https://cdn.seventyupgrades.com/item-slots/Wrists.jpg',
	[ItemSlot.ItemSlotHands]: 'https://cdn.seventyupgrades.com/item-slots/Hands.jpg',
	[ItemSlot.ItemSlotWaist]: 'https://cdn.seventyupgrades.com/item-slots/Waist.jpg',
	[ItemSlot.ItemSlotLegs]: 'https://cdn.seventyupgrades.com/item-slots/Legs.jpg',
	[ItemSlot.ItemSlotFeet]: 'https://cdn.seventyupgrades.com/item-slots/Feet.jpg',
	[ItemSlot.ItemSlotFinger1]: 'https://cdn.seventyupgrades.com/item-slots/Finger.jpg',
	[ItemSlot.ItemSlotFinger2]: 'https://cdn.seventyupgrades.com/item-slots/Finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: 'https://cdn.seventyupgrades.com/item-slots/Trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: 'https://cdn.seventyupgrades.com/item-slots/Trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: 'https://cdn.seventyupgrades.com/item-slots/MainHand.jpg',
	[ItemSlot.ItemSlotOffHand]: 'https://cdn.seventyupgrades.com/item-slots/OffHand.jpg',
	[ItemSlot.ItemSlotRanged]: 'https://cdn.seventyupgrades.com/item-slots/Ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}
