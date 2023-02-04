import { GemColor } from '../proto/common.js';
import { Profession } from '../proto/common.js';
import { getEnumValues } from '../utils.js';
import {
	UIGem as Gem,
} from '../proto/ui.js';

export const GEM_COLORS = (getEnumValues(GemColor) as Array<GemColor>).filter(color => color != GemColor.GemColorUnknown);
export const PRIMARY_COLORS = [GemColor.GemColorRed, GemColor.GemColorYellow, GemColor.GemColorBlue];
// Secondary is intentionally ordered so that it matches the inverse of PRIMARY_COLORS.
export const SECONDARY_COLORS = [GemColor.GemColorGreen, GemColor.GemColorPurple, GemColor.GemColorOrange];

export const socketToMatchingColors = new Map<GemColor, Array<GemColor>>();
socketToMatchingColors.set(GemColor.GemColorMeta, [GemColor.GemColorMeta]);
socketToMatchingColors.set(GemColor.GemColorBlue, [GemColor.GemColorBlue, GemColor.GemColorPurple, GemColor.GemColorGreen, GemColor.GemColorPrismatic]);
socketToMatchingColors.set(GemColor.GemColorRed, [GemColor.GemColorRed, GemColor.GemColorPurple, GemColor.GemColorOrange, GemColor.GemColorPrismatic]);
socketToMatchingColors.set(GemColor.GemColorYellow, [GemColor.GemColorYellow, GemColor.GemColorOrange, GemColor.GemColorGreen, GemColor.GemColorPrismatic]);
socketToMatchingColors.set(GemColor.GemColorPrismatic, [
	GemColor.GemColorRed,
	GemColor.GemColorOrange,
	GemColor.GemColorYellow,
	GemColor.GemColorGreen,
	GemColor.GemColorBlue,
	GemColor.GemColorPurple,
	GemColor.GemColorPrismatic,
]);

export function gemColorMatchesSocket(gemColor: GemColor, socketColor: GemColor) {
	return gemColor == socketColor || (socketToMatchingColors.has(socketColor) && socketToMatchingColors.get(socketColor)!.includes(gemColor));
}

// Whether the gem matches the given socket color, for the purposes of gaining the socket bonuses.
export function gemMatchesSocket(gem: Gem, socketColor: GemColor) {
	return gemColorMatchesSocket(gem.color, socketColor);
}

// Whether the gem is capable of slotting into a socket of the given color.
export function gemEligibleForSocket(gem: Gem, socketColor: GemColor) {
	return (gem.color == GemColor.GemColorMeta) == (socketColor == GemColor.GemColorMeta);
}

export function isUnrestrictedGem(gem: Gem, phase?: number): boolean {
	return !gem.unique &&
		gem.requiredProfession == Profession.ProfessionUnknown &&
		(phase == null || gem.phase <= phase);
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

	isCompareColorCondition(): boolean {
		return this.minRed == 0 && this.minYellow == 0 && this.minBlue == 0;
	}

	isOneOfEach(): boolean {
		return this.minRed == 1 && this.minYellow == 1 && this.minBlue == 1;
	}

	isTwoAndOne(): boolean {
		return [this.minRed, this.minYellow, this.minBlue].includes(2);
	}

	isThreeOfAColor(): boolean {
		return this.minRed == 3 || this.minYellow == 3 || this.minBlue == 3;
	}

	private static getNumInCategory(gemColor: GemColor, numRed: number, numYellow: number, numBlue: number): number {
		if (gemColor == GemColor.GemColorRed) {
			return numRed;
		} else if (gemColor == GemColor.GemColorYellow) {
			return numYellow;
		} else if (gemColor == GemColor.GemColorBlue) {
			return numBlue;
		} else {
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

// Keep these lists in order by item ID.
export const CHAOTIC_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41285, 'Requires at least 2 Blue Gems.', 0, 0, 2);
export const DESTRUCTIVE_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41307, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const EMBER_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41333, 'Requires at least 3 Red Gems.', 3, 0, 0);
export const ENIGMATIC_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41335, 'Requires at least 2 Red Gems and at least 1 Yellow Gem.', 2, 1, 0);
export const EFFULGENT_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41377, 'Requires at least 2 Blue Gems and at least 1 Red Gem.', 1, 0, 2);
export const SWIFT_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41339, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.', 1, 2, 0);
export const TIRELESS_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41375, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const REVITALIZING_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41376, 'Requires at least 2 Red Gems.', 2, 0, 0);
export const FORLORN_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41378, 'Requires at least 2 Yellow Gems and at least 1 Blue Gem.', 0, 2, 1);
export const IMPASSIVE_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41379, 'Requires at least 2 Red Gems and at least 1 Blue Gem.', 2, 0, 1);

export const AUSTERE_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41380, 'Requires at least 2 Blue Gems and at least 1 Red Gem.', 1, 0, 2);
export const PERSISTENT_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41381, 'Requires at least 2 Yellow Gems and at least 1 Blue Gem.', 0, 2, 1);
export const TRENCHANT_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41382, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const INVIGORATING_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41385, 'Requires at least 2 Blue Gems and at least 1 Red Gem.', 1, 0, 2);
export const BEAMING_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41389, 'Requires at least 2 Red Gems and at least 1 Yellow Gem.', 2, 1, 0);
export const BRACING_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41395, 'Requires at least 2 Red Gems and at least 1 Blue Gem.', 2, 0, 1);
export const ETERNAL_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41396, 'Requires at least 2 Red Gems and at least 1 Blue Gem.', 2, 0, 1);
export const POWERFUL_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41397, 'Requires at least 3 Blue Gems.', 0, 0, 3);
export const RELENTLESS_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41398, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const THUNDERING_SKYFLARE_DIAMOND = MetaGemCondition.fromMinColors(41400, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const INSIGHTFUL_EARTHSIEGE_DIAMOND = MetaGemCondition.fromMinColors(41401, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const SWIFT_STARFLARE_DIAMOND = MetaGemCondition.fromMinColors(44076, 'Requires at least 2 Yellow Gems and at least 1 Red Gem.', 1, 2, 0);
export const TIRELESS_STARFLARE_DIAMOND = MetaGemCondition.fromMinColors(44078, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);
export const ENIGMATIC_STARFLARE_DIAMOND = MetaGemCondition.fromMinColors(44081, 'Requires at least 2 Red Gems and at least 1 Blue Gem.', 2, 0, 1);
export const IMPASSIVE_STARFLARE_DIAMOND = MetaGemCondition.fromMinColors(44082, 'Requires at least 2 Blue Gems and at least 1 Red Gem.', 1, 0, 2);
export const FORLORN_STARFLARE_DIAMOND = MetaGemCondition.fromMinColors(44084, 'Requires at least 2 Yellow Gems and at least 1 Blue Gem.', 0, 2, 1);
export const PERSISTENT_EARTHSHATTER_DIAMOND = MetaGemCondition.fromMinColors(44087, 'Requires at least 3 Blue Gems.', 0, 0, 3);
export const POWERFUL_EARTHSHATTER_DIAMOND = MetaGemCondition.fromMinColors(44088, 'Requires at least 2 Blue Gems and at least 1 Yellow Gem.', 0, 1, 2);
export const TRENCHANT_EARTHSHATTER_DIAMOND = MetaGemCondition.fromMinColors(44089, 'Requires at least 1 Red Gem, at least 1 Yellow Gem, and at least 1 Blue Gem.', 1, 1, 1);

// TBC GEMS
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

const emptyGemSocketIcons: Partial<Record<GemColor, string>> = {
	[GemColor.GemColorBlue]: 'https://wow.zamimg.com/images/icons/socket-blue.gif',
	[GemColor.GemColorMeta]: 'https://wow.zamimg.com/images/icons/socket-meta.gif',
	[GemColor.GemColorRed]: 'https://wow.zamimg.com/images/icons/socket-red.gif',
	[GemColor.GemColorYellow]: 'https://wow.zamimg.com/images/icons/socket-yellow.gif',
	[GemColor.GemColorPrismatic]: 'https://wow.zamimg.com/images/icons/socket-prismatic.gif',
};
export function getEmptyGemSocketIconUrl(color: GemColor): string {
	if (emptyGemSocketIcons[color])
		return emptyGemSocketIcons[color] as string;

	throw new Error('No empty socket url for gem socket color: ' + color);
}
