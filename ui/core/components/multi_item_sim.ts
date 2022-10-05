import { getWowheadItemId } from '../proto_utils/equipped_item.js';
import { EquippedItem } from '../proto_utils/equipped_item.js';
import { getEmptyGemSocketIconUrl, gemMatchesSocket } from '../proto_utils/gems.js';
import { setGemSocketCssClass } from '../proto_utils/gems.js';
import { Stats } from '../proto_utils/stats.js';
import { enchantAppliesToItem } from '../proto_utils/utils.js';
import { Class, Enchant, Gem, GemColor } from '../proto/common.js';
import { HandType } from '../proto/common.js';
import { WeaponType } from '../proto/common.js';
import { Item } from '../proto/common.js';
import { ItemQuality } from '../proto/common.js';
import { ItemSlot } from '../proto/common.js';
import { ItemType } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import { getEnchantDescription } from '../proto_utils/enchants.js';
import { ActionId } from '../proto_utils/action_id.js';
import { slotNames } from '../proto_utils/names.js';
import { setItemQualityCssClass } from '../css_utils.js';
import { Player } from '../player.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { formatDeltaTextElem } from '../utils.js';
import { getEnumValues } from '../utils.js';

import { Component } from './component.js';
import { Popup } from './popup.js';
import { makePhaseSelector } from './other_inputs.js';
import { makeShow1hWeaponsSelector } from './other_inputs.js';
import { makeShow2hWeaponsSelector } from './other_inputs.js';
import { makeShowMatchingGemsSelector } from './other_inputs.js';

class MultiItemSim extends Popup {
	private player: Player<any>;
	private readonly contentElem: HTMLElement;

	constructor(parent: HTMLElement, player: Player<any>) {
		super(parent);
		this.player = player;

		this.rootElem.classList.add('multi-item-modal');
        this.contentElem = document.createElement("div");
        this.rootElem.appendChild(this.contentElem);

		this.addCloseButton();
	}
}