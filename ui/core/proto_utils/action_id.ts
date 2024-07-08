import { getWowheadLanguagePrefix } from '../constants/lang.js';
import { CHARACTER_LEVEL } from '../constants/mechanics.js';
import { ResourceType } from '../proto/api.js';
import { ActionID as ActionIdProto , OtherAction } from '../proto/common.js';
import { IconData ,
	UIItem as Item,
} from '../proto/ui.js';
import { Database } from './database.js';

// If true uses wotlkdb.com, else uses wowhead.com.
export const USE_WOTLK_DB = false;

// Uniquely identifies a specific item / spell / thing in WoW. This object is immutable.
export class ActionId {
	readonly itemId: number;
	readonly spellId: number;
	readonly otherId: OtherAction;
	readonly tag: number;

	readonly baseName: string; // The name without any tag additions.
	readonly name: string;
	readonly iconUrl: string;

	private constructor(itemId: number, spellId: number, otherId: OtherAction, tag: number, baseName: string, name: string, iconUrl: string) {
		this.itemId = itemId;
		this.spellId = spellId;
		this.otherId = otherId;
		this.tag = tag;

		switch (otherId) {
			case OtherAction.OtherActionNone:
				break;
			case OtherAction.OtherActionWait:
				baseName = '等待';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_pocketwatch_01.jpg';
				break;
			case OtherAction.OtherActionManaRegen:
				name = '法力回复';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				if (tag == 1) {
					name += ' (施法中)';
				} else if (tag == 2) {
					name += ' (未施法)';
				}
				break;
			case OtherAction.OtherActionEnergyRegen:
				baseName = '能量回复';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeEnergy];
				break;
			case OtherAction.OtherActionFocusRegen:
				baseName = '集中回复';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeFocus];
				break;
			case OtherAction.OtherActionManaGain:
				baseName = '法力获得';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				break;
			case OtherAction.OtherActionRageGain:
				baseName = '怒气获得';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeRage];
				break;
			case OtherAction.OtherActionAttack:
				name = '攻击';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_sword_04.jpg';
				if (tag == 1) {
					name += ' (主手)';
				} else if (tag == 2) {
					name += ' (副手)';
				}
				break;
			case OtherAction.OtherActionShoot:
				name = '射击';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/ability_marksmanship.jpg';
				break;
			case OtherAction.OtherActionPet:
				break;
			case OtherAction.OtherActionRefund:
				baseName = '退款';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_misc_coin_01.jpg';
				break;
			case OtherAction.OtherActionDamageTaken:
				baseName = '受到的伤害';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_sword_04.jpg';
				break;
			case OtherAction.OtherActionHealingModel:
				baseName = '治疗量';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_holy_renew.jpg';
				break;
			case OtherAction.OtherActionBloodRuneGain:
				baseName = '血符文获得';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_deathknight_deathstrike.jpg';
				break;
			case OtherAction.OtherActionFrostRuneGain:
				baseName = '冰霜符文获得';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_deathknight_deathstrike2.jpg';
				break;
			case OtherAction.OtherActionUnholyRuneGain:
				baseName = '邪恶符文获得';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_deathknight_empowerruneblade.jpg';
				break;
			case OtherAction.OtherActionDeathRuneGain:
				baseName = '死亡符文获得';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_deathknight_empowerruneblade.jpg';
				break;
			case OtherAction.OtherActionPotion:
				baseName = '药水';
				iconUrl = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/inv_alchemy_elixir_04.jpg';
				break;
		}

		this.baseName = baseName;
		this.name = name || baseName;
		this.iconUrl = iconUrl;
	}

	anyId(): number {
		return this.itemId || this.spellId || this.otherId;
	}

	equals(other: ActionId): boolean {
		return this.equalsIgnoringTag(other) && this.tag == other.tag;
	}

	equalsIgnoringTag(other: ActionId): boolean {
		return (
			this.itemId == other.itemId
			&& this.spellId == other.spellId
			&& this.otherId == other.otherId);
	}

	setBackground(elem: HTMLElement) {
		if (this.iconUrl) {
			elem.style.backgroundImage = `url('${this.iconUrl}')`;
		}
	}

	static makeItemUrl(id: number): string {
		return `https://www.wclbox.com/db/wotlk/item/${id}`;
	}
	static makeSpellUrl(id: number): string {
		return `https://www.wclbox.com/db/wotlk/spell/${id}`;
	}
	static makeQuestUrl(id: number): string {
		return `https://wowhead.com/wotlk/cn/quest=${id}`;
	}
	static makeNpcUrl(id: number): string {
		return `https://wowhead.com/wotlk/cn/npc=${id}`;
	}
	static makeZoneUrl(id: number): string {
		return `https://wowhead.com/wotlk/cn/zone=${id}`;
	}

	setWowheadHref(elem: HTMLAnchorElement) {
		if (this.itemId) {
			elem.href = ActionId.makeItemUrl(this.itemId);
		} else if (this.spellId) {
			elem.href = ActionId.makeSpellUrl(this.spellId);
		}
	}

	setBackgroundAndHref(elem: HTMLAnchorElement) {
		this.setBackground(elem);
		this.setWowheadHref(elem);
	}

	async fillAndSet(elem: HTMLAnchorElement, setHref: boolean, setBackground: boolean): Promise<ActionId> {
		const filled = await this.fill();
		if (setHref) {
			filled.setWowheadHref(elem);
		}
		if (setBackground) {
			filled.setBackground(elem);
		}
		return filled;
	}

	// Returns an ActionId with the name and iconUrl fields filled.
	// playerIndex is the optional index of the player to whom this ID corresponds.
	async fill(playerIndex?: number): Promise<ActionId> {
		if (this.name || this.iconUrl) {
			return this;
		}

		if (this.otherId) {
			return this;
		}

		const tooltipData = await ActionId.getTooltipData(this);

		const baseName = tooltipData['name'];
		let name = baseName;
		switch (baseName) {
			case 'Explosive Shot':
				if (this.spellId == 60053) {
					name += ' (R4)';
				} else if (this.spellId == 60052) {
					name += ' (R3)';
				}
				break;
			case 'Explosive Trap':
				if (this.tag == 1) {
					name += ' (编织)';
				}
				break;
			case 'Arcane Blast':
				if (this.tag == 1) {
					name += ' (无层数)';
				} else if (this.tag == 2) {
					name += ` (1层)`;
				} else if (this.tag > 2) {
					name += ` (${this.tag - 1}层)`;
				}
				break;
			case 'Hot Streak':
				if (this.tag) name += ' (暴击)';
				break;
			case 'Fireball':
			case 'Flamestrike':
				if (this.tag == 8) {
					name += ' (等级 8)';
				} else if (this.tag == 9) {
					name += ' (等级 9)';
				}
				break;
			case 'Pyroblast':
				if (this.tag) name += ' (DoT)';
				break;
			case 'Living Bomb':
				if (this.spellId == 55362) name += ' (爆炸)';
				break;
			case 'Evocation':
				if (this.tag == 1) {
					name += ' (1次)';
				} else if (this.tag == 2) {
					name += ' (2次)';
				} else if (this.tag == 3) {
					name += ' (3次)';
				} else if (this.tag == 4) {
					name += ' (4次)';
				} else if (this.tag == 5) {
					name += ' (5次)';
				}
				break;
			case 'Mind Flay':
				if (this.tag == 1) {
					name += ' (1次)';
				} else if (this.tag == 2) {
					name += ' (2次)';
				} else if (this.tag == 3) {
					name += ' (3次)';
				}
				break;
			case 'Mind Sear':
				if (this.tag == 1) {
					name += ' (1次)';
				} else if (this.tag == 2) {
					name += ' (2次)';
				} else if (this.tag == 3) {
					name += ' (3次)';
				}
				break;
			case 'Shattering Throw':
				if (this.tag === playerIndex) {
					name += ` (自身)`;
				}
				break;
			case 'Envenom':
			case 'Eviscerate':
			case 'Expose Armor':
			case 'Rupture':
			case 'Slice and Dice':
				if (this.tag) name += ` (${this.tag} 连击点)`;
				break;
			case 'Instant Poison IX':
			case 'Wound Poison VII':
				if (this.tag == 1) {
					name += ' (致命)';
				} else if (this.tag == 2) {
					name += ' (毒化)';
				}
				break;
			case 'Fan of Knives':
			case 'Killing Spree':
				if (this.tag == 1) {
					name += ' (主手)';
				} else if (this.tag == 2) {
					name += ' (副手)';
				}
				break;
			case 'Tricks of the Trade':
				if (this.tag == 1) {
					name += ' (非自身)';
				}
				break;
			case 'Chain Lightning':
			case 'Lightning Bolt':
				if (this.tag == 6) {
					name += ' (过载)';
				} else if (this.tag) {
					name += ` (${this.tag} 魔杖充能)`;
				}
				break;
			case 'Holy Shield':
				if (this.tag == 1) {
					name += ' (触发)';
				}
				break;
			case 'Righteous Vengeance':
				if (this.tag == 1) {
					name += ' (应用)';
				} else if (this.tag == 2) {
					name += ' (DoT)';
				}
				break;
			case 'Holy Vengeance':
				if (this.tag == 1) {
					name += ' (应用)';
				} else if (this.tag == 2) {
					name += ' (DoT)';
				}
				break;
			case 'Bloodlust':
			case 'Ferocious Inspiration':
			case 'Innervate':
			case 'Focus Magic':
			case 'Mana Tide Totem':
			case 'Power Infusion':
				if (this.tag != -1) {
					if (this.tag === playerIndex || playerIndex == undefined) {
						name += ` (自身)`;
					} else {
						name += ` (来自 #${this.tag + 1})`;
					}
				} else {
					name += " (团队)";
				}
				break;
			case 'Darkmoon Card: Crusade':
				if (this.tag == 1) {
					name += ' (近战)';
				} else if (this.tag == 2) {
					name += ' (法术)';
				}
				break;
			case 'Frozen Blows':
				if (this.tag == 1) {
					name += ' (物理)';
				} else if (this.tag == 2) {
					name += ' (冰霜)';
				}
				break;
			case 'Scourge Strike':
				if (this.tag == 1) {
					name += ' (物理)';
				} else if (this.tag == 2) {
					name += ' (暗影)';
				}
				break;
			case 'Heart Strike':
				if (this.tag == 2) {
					name += ' (次要目标)';
				}
				break;
			case 'Rune Strike':
				if (this.tag == 0) {
					name += ' (排队)';
				} else if (this.tag == 1) {
					name += ' (主手)';
				} else if (this.tag == 2) {
					name += ' (副手)';
				}
				break;
			case 'Frost Strike':
			case 'Plague Strike':
			case 'Blood Strike':
			case 'Death Strike':
			case 'Obliterate':
			case 'Blood-Caked Strike':
			case 'Lightning Speed':
			case 'Windfury Weapon':
			case 'Berserk':
				if (this.tag == 1) {
					name += ' (主手)';
				} else if (this.tag == 2) {
					name += ' (副手)';
				}
				break;
			case 'Battle Shout':
				if (this.tag == 1) {
					name += ' (快照)';
				}
				break;
			case 'Heroic Strike':
			case 'Cleave':
			case 'Maul':
				if (this.tag == 1) {
					name += ' (排队)';
				}
				break;
			case 'Whirlwind':
				if (this.tag == 1) {
					name += ' (副手)';
				}
				break;
			case 'Seed of Corruption':
				if (this.tag == 0) {
					name += ' (DoT)';
				} else if (this.tag == 1) {
					name += ' (爆炸)';
				}
				break;
			case 'Thunderfury':
				if (this.tag == 1) {
					name += ' (单目标)';
				} else if (this.tag == 2) {
					name += ' (多目标)';
				}
				break;
			default:
				if (this.tag) {
					name += ' (未知)';
				}
				break;
		}


		const idString = this.toProtoString();
		const iconOverrideId = idOverrides[idString] || null;

		let iconUrl = ActionId.makeIconUrl(tooltipData['icon']);
		if (iconOverrideId) {
			const overrideTooltipData = await ActionId.getTooltipData(iconOverrideId);
			iconUrl = ActionId.makeIconUrl(overrideTooltipData['icon']);
		}

		return new ActionId(this.itemId, this.spellId, this.otherId, this.tag, baseName, name, iconUrl);
	}

	toString(): string {
		return this.toStringIgnoringTag() + (this.tag ? ('-' + this.tag) : '');
	}

	toStringIgnoringTag(): string {
		if (this.itemId) {
			return 'item-' + this.itemId;
		} else if (this.spellId) {
			return 'spell-' + this.spellId;
		} else if (this.otherId) {
			return 'other-' + this.otherId;
		} else {
			throw new Error('Empty action id!');
		}
	}

	toProto(): ActionIdProto {
		const protoId = ActionIdProto.create({
			tag: this.tag,
		});

		if (this.itemId) {
			protoId.rawId = {
				oneofKind: 'itemId',
				itemId: this.itemId,
			};
		} else if (this.spellId) {
			protoId.rawId = {
				oneofKind: 'spellId',
				spellId: this.spellId,
			};
		} else if (this.otherId) {
			protoId.rawId = {
				oneofKind: 'otherId',
				otherId: this.otherId,
			};
		}

		return protoId;
	}

	toProtoString(): string {
		return ActionIdProto.toJsonString(this.toProto());
	}

	withoutTag(): ActionId {
		return new ActionId(this.itemId, this.spellId, this.otherId, 0, this.baseName, this.baseName, this.iconUrl);
	}

	static fromEmpty(): ActionId {
		return new ActionId(0, 0, OtherAction.OtherActionNone, 0, '', '', '');
	}

	static fromItemId(itemId: number, tag?: number): ActionId {
		return new ActionId(itemId, 0, OtherAction.OtherActionNone, tag || 0, '', '', '');
	}

	static fromSpellId(spellId: number, tag?: number): ActionId {
		return new ActionId(0, spellId, OtherAction.OtherActionNone, tag || 0, '', '', '');
	}

	static fromOtherId(otherId: OtherAction, tag?: number): ActionId {
		return new ActionId(0, 0, otherId, tag || 0, '', '', '');
	}

	static fromPetName(petName: string): ActionId {
		return petNameToActionId[petName] || new ActionId(0, 0, OtherAction.OtherActionPet, 0, petName, petName, petNameToIcon[petName] || '');
	}

	static fromItem(item: Item): ActionId {
		return ActionId.fromItemId(item.id);
	}

	static fromProto(protoId: ActionIdProto): ActionId {
		if (protoId.rawId.oneofKind == 'spellId') {
			return ActionId.fromSpellId(protoId.rawId.spellId, protoId.tag);
		} else if (protoId.rawId.oneofKind == 'itemId') {
			return ActionId.fromItemId(protoId.rawId.itemId, protoId.tag);
		} else if (protoId.rawId.oneofKind == 'otherId') {
			return ActionId.fromOtherId(protoId.rawId.otherId, protoId.tag);
		} else {
			return ActionId.fromEmpty();
		}
	}

	private static readonly logRegex = /{((SpellID)|(ItemID)|(OtherID)): (\d+)(, Tag: (-?\d+))?}/;
	private static readonly logRegexGlobal = new RegExp(ActionId.logRegex, 'g');
	private static fromMatch(match: RegExpMatchArray): ActionId {
		const idType = match[1];
		const id = parseInt(match[5]);
		return new ActionId(
			idType == 'ItemID' ? id : 0,
			idType == 'SpellID' ? id : 0,
			idType == 'OtherID' ? id : 0,
			match[7] ? parseInt(match[7]) : 0,
			'', '', '');
	}
	static fromLogString(str: string): ActionId {
		const match = str.match(ActionId.logRegex);
		if (match) {
			return ActionId.fromMatch(match);
		} else {
			console.warn('Failed to parse action id from log: ' + str);
			return ActionId.fromEmpty();
		}
	}

	static async replaceAllInString(str: string): Promise<string> {
		const matches = [...str.matchAll(ActionId.logRegexGlobal)];

		const replaceData = await Promise.all(matches.map(async match => {
			const actionId = ActionId.fromMatch(match);
			const filledId = await actionId.fill();
			return {
				firstIndex: match.index || 0,
				len: match[0].length,
				actionId: filledId,
			};
		}));

		// Loop in reverse order so we can greedily apply the string replacements.
		for (let i = replaceData.length - 1; i >= 0; i--) {
			const data = replaceData[i];
			str = str.substring(0, data.firstIndex) + data.actionId.name + str.substring(data.firstIndex + data.len);
		}

		return str;
	}

	private static makeIconUrl(iconLabel: string): string {
		if (USE_WOTLK_DB) {
			return `https://wotlkdb.com/static/images/wow/icons/large/${iconLabel}.jpg`;
		} else {
			return `https://db.newbeebox.com/wow/wz/images/wow/icons/large/${iconLabel}.jpg`;
		}
	}

	static async getTooltipData(actionId: ActionId): Promise<IconData> {
		if (actionId.itemId) {
			return await Database.getItemIconData(actionId.itemId);
		} else {
			return await Database.getSpellIconData(actionId.spellId);
		}
	}
}

// Some items/spells have weird icons, so use this to show a different icon instead.
const idOverrides: Record<string, ActionId> = {};
idOverrides[ActionId.fromSpellId(37212).toProtoString()] = ActionId.fromItemId(29035); // Improved Wrath of Air Totem
idOverrides[ActionId.fromSpellId(37223).toProtoString()] = ActionId.fromItemId(29040); // Improved Strength of Earth Totem
idOverrides[ActionId.fromSpellId(37447).toProtoString()] = ActionId.fromItemId(30720); // Serpent-Coil Braid
idOverrides[ActionId.fromSpellId(37443).toProtoString()] = ActionId.fromItemId(30196); // Robes of Tirisfal (4pc bonus)

export const defaultTargetIcon = 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_metamorphosis.jpg';

const petNameToActionId: Record<string, ActionId> = {
	'Gnomish Flame Turret': ActionId.fromItemId(23841),
	'Mirror Image': ActionId.fromSpellId(55342),
	'Water Elemental': ActionId.fromSpellId(31687),
	"Greater Fire Elemental": ActionId.fromSpellId(2894),
	'Shadowfiend': ActionId.fromSpellId(34433),
	'Spirit Wolf 1': ActionId.fromSpellId(51533),
	'Spirit Wolf 2': ActionId.fromSpellId(51533),
	'Rune Weapon': ActionId.fromSpellId(49028),
	'Bloodworm': ActionId.fromSpellId(50452),
	'Gargoyle': ActionId.fromSpellId(49206),
	'Ghoul': ActionId.fromSpellId(46584),
	'Army of the Dead': ActionId.fromSpellId(42650),
	'Valkyr': ActionId.fromSpellId(71844),
};

// https://wowhead.com/wotlk/hunter-pets
const petNameToIcon: Record<string, string> = {
	'Bat': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_bat.jpg',
	'Bear': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_bear.jpg',
	'Bird of Prey': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_owl.jpg',
	'Boar': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_boar.jpg',
	'Carrion Bird': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_vulture.jpg',
	'Cat': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_cat.jpg',
	'Chimaera': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_chimera.jpg',
	'Core Hound': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_corehound.jpg',
	'Crab': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_crab.jpg',
	'Crocolisk': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_crocolisk.jpg',
	'Devilsaur': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_devilsaur.jpg',
	'Dragonhawk': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_dragonhawk.jpg',
	'Felguard': 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonfelguard.jpg',
	'Felhunter': 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonfelhunter.jpg',
	'Infernal': 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summoninfernal.jpg',
	'Gorilla': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_gorilla.jpg',
	'Hyena': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_hyena.jpg',
	'Imp': 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonimp.jpg',
	'Mirror Image': 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
	'Moth': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_moth.jpg',
	'Nether Ray': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_netherray.jpg',
	'Owl': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_owl.jpg',
	'Raptor': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_raptor.jpg',
	'Ravager': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_ravager.jpg',
	'Rhino': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_rhino.jpg',
	'Scorpid': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_scorpid.jpg',
	'Serpent': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_nature_guardianward.jpg',
	'Silithid': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_silithid.jpg',
	'Spider': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_spider.jpg',
	'Spirit Beast': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_druid_primalprecision.jpg',
	'Spore Bat': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_sporebat.jpg',
	'Succubus': 'https://db.newbeebox.com/wow/wz/images/wow/icons/large/spell_shadow_summonsuccubus.jpg',
	'Tallstrider': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_tallstrider.jpg',
	'Turtle': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_turtle.jpg',
	'Warp Stalker': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_warpstalker.jpg',
	'Wasp': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_wasp.jpg',
	'Wind Serpent': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_windserpent.jpg',
	'Wolf': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_wolf.jpg',
	'Worm': 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_pet_worm.jpg',
};

export function getPetIconFromName(name: string): string|ActionId|undefined {
	return petNameToActionId[name] || petNameToIcon[name];
}

export const resourceTypeToIcon: Record<ResourceType, string> = {
	[ResourceType.ResourceTypeNone]: '',
	[ResourceType.ResourceTypeHealth]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/inv_elemental_mote_life01.jpg',
	[ResourceType.ResourceTypeMana]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/inv_elemental_mote_mana.jpg',
	[ResourceType.ResourceTypeEnergy]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_shadow_shadowworddominate.jpg',
	[ResourceType.ResourceTypeRage]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/spell_misc_emotionangry.jpg',
	[ResourceType.ResourceTypeComboPoints]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/inv_mace_2h_pvp410_c_01.jpg',
	[ResourceType.ResourceTypeFocus]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/ability_hunter_focusfire.jpg',
	[ResourceType.ResourceTypeRunicPower]: 'https://db.newbeebox.com/wow/wz/images/wow/icons/medium/inv_sword_62.jpg',
	[ResourceType.ResourceTypeBloodRune]: '/wotlk/assets/img/blood_rune.png',
	[ResourceType.ResourceTypeFrostRune]: '/wotlk/assets/img/frost_rune.png',
	[ResourceType.ResourceTypeUnholyRune]: '/wotlk/assets/img/unholy_rune.png',
	[ResourceType.ResourceTypeDeathRune]: '/wotlk/assets/img/death_rune.png',
};
