import { Gem } from '/tbc/core/proto/common.js';
import { GemColor } from '/tbc/core/proto/common.js';

const socketToMatchingColors = new Map<GemColor, Array<GemColor>>();
socketToMatchingColors.set(GemColor.GemColorMeta, [GemColor.GemColorMeta]);
socketToMatchingColors.set(GemColor.GemColorBlue, [GemColor.GemColorBlue, GemColor.GemColorPurple, GemColor.GemColorGreen, GemColor.GemColorPrismatic]);
socketToMatchingColors.set(GemColor.GemColorRed, [GemColor.GemColorRed, GemColor.GemColorPurple, GemColor.GemColorOrange, GemColor.GemColorPrismatic]);
socketToMatchingColors.set(GemColor.GemColorYellow, [GemColor.GemColorYellow, GemColor.GemColorOrange, GemColor.GemColorGreen, GemColor.GemColorPrismatic]);

// Whether the gem matches the given socket color, for the purposes of gaining the socket bonuses.
export function gemMatchesSocket(gem: Gem, socketColor: GemColor) {
	return gem.color == socketColor || (socketToMatchingColors.has(socketColor) && socketToMatchingColors.get(socketColor)!.includes(gem.color));
}

// Whether the gem is capable of slotting into a socket of the given color.
export function gemEligibleForSocket(gem: Gem, socketColor: GemColor) {
	return (gem.color == GemColor.GemColorMeta) == (socketColor == GemColor.GemColorMeta);
}


export class MetaGemCondition {
	readonly id: number;
	readonly description: string;

	readonly minRed: number;
	readonly minYellow: number;
	readonly minBlue: number;

	readonly compareColorGreater: GemColor;
	readonly compareColorLesser: GemColor;

	constructor(id: number, description: string, minRed: number, minYellow: number, minBlue: number, compareColorGreater: GemColor, compareColorLesser: GemColor) {
		this.id = id;
		this.description = description;
		this.minRed = minRed;
		this.minYellow = minYellow;
		this.minBlue = minBlue;
		this.compareColorGreater = compareColorGreater;
		this.compareColorLesser = compareColorLesser;

		metaGemConditions.set(this.id, this);
	}

	// Whether the condition is met, i.e. the meta gem is activated.
	isMet(numRed: number, numYellow: number, numBlue: number): boolean {
		if (!(numRed >= this.minRed && numYellow >= this.minYellow && numBlue >= this.minBlue)) {
			return false;
		}

		if (this.compareColorGreater == GemColor.GemColorUnknown) {
			return true;
		}

		const numGreater = MetaGemCondition.getNumInCategory(this.compareColorGreater, numRed, numYellow, numBlue);
		const numLesser = MetaGemCondition.getNumInCategory(this.compareColorLesser, numRed, numYellow, numBlue);
		return numGreater > numLesser;
	}

	private static getNumInCategory(gemColor: GemColor, numRed: number, numYellow: number, numBlue: number): number {
		if (gemColor == GemColor.GemColorRed) {
			return numRed;
		} else if (gemColor == GemColor.GemColorYellow) {
			return numYellow;
		} else if (gemColor == GemColor.GemColorBlue) {
			return numBlue;
		} else  {
			throw new Error('Invalid gem color for category check: ' + gemColor);
		}
	}

	static fromMinColors(id: number, description: string, minRed: number, minYellow: number, minBlue: number): MetaGemCondition {
		return new MetaGemCondition(id, description, minRed, minYellow, minBlue, GemColor.GemColorUnknown, GemColor.GemColorUnknown);
	}

	static fromCompareColors(id: number, description: string, compareColorGreater: GemColor, compareColorLesser: GemColor): MetaGemCondition {
		return new MetaGemCondition(id, description, 0, 0, 0, compareColorGreater, compareColorLesser);
	}
}

const metaGemConditions = new Map<number, MetaGemCondition>();

export function getMetaGemCondition(id: number): MetaGemCondition {
	if (!metaGemConditions.has(id)) {
		throw new Error('Missing meta gem condition for gem: ' + id);
	}

	return metaGemConditions.get(id)!;
}

export function isMetaGemActive(metaGem: Gem, numRed: number, numYellow: number, numBlue: number): boolean {
	return getMetaGemCondition(metaGem.id).isMet(numRed, numYellow, numBlue);
}

export function getMetaGemConditionDescription(metaGem: Gem): string {
	return getMetaGemCondition(metaGem.id).description;
}

// Keep these lists in alphabetical order, separated by condition type.

export const BRUTAL_EARTHSTORM_DIAMOND = MetaGemCondition.fromMinColors(25899, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems.', 2, 2, 2);
export const CHAOTIC_SKYFIRE_DIAMOND = MetaGemCondition.fromMinColors(34220, 'Requires at least 2 Blue Gems.', 0, 0, 2);
export const DESTRUCTIVE_SKYFIRE_DIAMOND = MetaGemCondition.fromMinColors(25890, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems.', 2, 2, 2);
export const EMBER_SKYFIRE_DIAMOND = MetaGemCondition.fromMinColors(35503, 'Requires at least 3 Red Gems.', 3, 0, 0);
export const ETERNAL_EARTHSTORM_DIAMOND = MetaGemCondition.fromMinColors(35501, 'Requires at least 2 Blue Gems and at least 1 Yellow Gem.', 0, 1, 2);
export const IMBUED_UNSTABLE_DIAMOND = MetaGemCondition.fromMinColors(32641, 'Requires at least 3 Yellow Gems.', 0, 3, 0);
export const INSIGHTFUL_EARTHSTORM_DIAMOND = MetaGemCondition.fromMinColors(25901, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems.', 2, 2, 2);
export const POWERFUL_EARTHSTORM_DIAMOND = MetaGemCondition.fromMinColors(25896, 'Requires at least 3 Blue Gems.', 0, 0, 3);
export const RELENTLESS_EARTHSTORM_DIAMOND = MetaGemCondition.fromMinColors(32409, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems.', 2, 2, 2);
export const SWIFT_SKYFIRE_DIAMOND = MetaGemCondition.fromMinColors(25894, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.', 1, 2, 0);
export const SWIFT_STARFIRE_DIAMOND = MetaGemCondition.fromMinColors(28557, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.', 1, 2, 0);
export const SWIFT_WINDFIRE_DIAMOND = MetaGemCondition.fromMinColors(28556, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.', 1, 2, 0);
export const TENACIOUS_EARTHSTORM_DIAMOND = MetaGemCondition.fromMinColors(25898, 'Requires at least 5 Blue Gems.', 0, 0, 5);
export const THUNDERING_SKYFIRE_DIAMOND = MetaGemCondition.fromMinColors(32410, 'Requires at least 2 Red Gems, at least 2 Yellow Gems, and at least 2 Blue Gems.', 2, 2, 2);

export const BRACING_EARTHSTORM_DIAMOND = MetaGemCondition.fromCompareColors(25897, 'Requires more Red Gems than Blue Gems.', GemColor.GemColorRed, GemColor.GemColorBlue);
export const ENIGMATIC_SKYFIRE_DIAMOND = MetaGemCondition.fromCompareColors(25895, 'Requires more Red Gems than Yellow Gems.', GemColor.GemColorRed, GemColor.GemColorYellow);
export const MYSTICAL_SKYFIRE_DIAMOND = MetaGemCondition.fromCompareColors(25893, 'Requires more Blue Gems than Yellow Gems.', GemColor.GemColorBlue, GemColor.GemColorYellow);
export const POTENT_UNSTABLE_DIAMOND = MetaGemCondition.fromCompareColors(32640, 'Requires more Blue Gems than Yellow Gems.', GemColor.GemColorBlue, GemColor.GemColorYellow);

const gemSocketCssClasses: Partial<Record<GemColor, string>> = {
	[GemColor.GemColorBlue]: 'socket-color-blue',
	[GemColor.GemColorMeta]: 'socket-color-meta',
	[GemColor.GemColorRed]: 'socket-color-red',
	[GemColor.GemColorYellow]: 'socket-color-yellow',
};
export function setGemSocketCssClass(elem: HTMLElement, color: GemColor) {
	Object.values(gemSocketCssClasses).forEach(cssClass => elem.classList.remove(cssClass));

	if (gemSocketCssClasses[color]) {
		elem.classList.add(gemSocketCssClasses[color] as string);
		return;
	}

	throw new Error('No css class for gem socket color: ' + color);
}

const emptyGemSocketIcons: Partial<Record<GemColor, string>> = {
	[GemColor.GemColorBlue]: 'https://wow.zamimg.com/images/icons/socket-blue.gif',
	[GemColor.GemColorMeta]: 'https://wow.zamimg.com/images/icons/socket-meta.gif',
	[GemColor.GemColorRed]: 'https://wow.zamimg.com/images/icons/socket-red.gif',
	[GemColor.GemColorYellow]: 'https://wow.zamimg.com/images/icons/socket-yellow.gif',
};
export function getEmptyGemSocketIconUrl(color: GemColor): string {
	if (emptyGemSocketIcons[color])
		return emptyGemSocketIcons[color] as string;

	throw new Error('No empty socket url for gem socket color: ' + color);
}
