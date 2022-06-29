import { ItemQuality } from './proto/common.js';

const itemQualityCssClasses: Record<ItemQuality, string> = {
	[ItemQuality.ItemQualityJunk]: 'item-quality-junk',
	[ItemQuality.ItemQualityCommon]: 'item-quality-common',
	[ItemQuality.ItemQualityUncommon]: 'item-quality-uncommon',
	[ItemQuality.ItemQualityRare]: 'item-quality-rare',
	[ItemQuality.ItemQualityEpic]: 'item-quality-epic',
	[ItemQuality.ItemQualityLegendary]: 'item-quality-legendary',
};
export function setItemQualityCssClass(elem: HTMLElement, quality: ItemQuality | null) {
	Object.values(itemQualityCssClasses).forEach(cssClass => elem.classList.remove(cssClass));

	if (quality) {
		elem.classList.add(itemQualityCssClasses[quality]);
	}
}
