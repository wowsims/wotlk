import {
	Class,
	Race,
	RaidBuffs,
	Spec,
} from '../core/proto/common.js';
import { PaladinAura } from '../core/proto/paladin.js';
import { AirTotem, EarthTotem, FireTotem, WaterTotem } from '../core/proto/shaman.js';
import { Warlock_Options_Summon as WarlockSummon } from '../core/proto/warlock.js';
import { WarriorShout} from '../core/proto/warrior.js';
import { Component } from '../core/components/component.js';
import { Player } from "../core/player.js";
import { Raid } from "../core/raid.js";
import { ActionId } from '../core/proto_utils/action_id.js';
import {
	ClassSpecs,
	SpecTalents,
} from '../core/proto_utils/utils.js';
import { sum } from '../core/utils.js';

import { BuffBot } from "./buff_bot.js";
import { BuffBotId } from "./presets.js";
import { RaidSimUI } from './raid_sim_ui.js';

declare var tippy: any;

interface RaidStatsOptions {
	sections: Array<RaidStatsSectionOptions>,
}

interface RaidStatsSectionOptions {
	label: string,
	categories: Array<RaidStatsCategoryOptions>,
}

interface RaidStatsCategoryOptions {
	label: string,
	effects: Array<RaidStatsEffectOptions>,
}

type PlayerProvider = (player: Player<any>) => boolean;
type BuffBotProvider = (buffBot: BuffBot) => boolean;
type RaidProvider = (raid: Raid) => boolean;

interface RaidStatsEffectOptions {
	label: string,
	actionId?: ActionId,
	providedByPlayer?: PlayerProvider,
	providedByBuffBot?: BuffBotProvider,
	providedByRaid?: RaidProvider,
}

export class RaidStats extends Component {
	private readonly categories: Array<RaidStatsCategory>;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI) {
		super(parent, 'raid-stats');

		let categories: Array<RaidStatsCategory> = [];
		RAID_STATS_OPTIONS.sections.forEach(section => {
			const sectionElem = document.createElement('div');
			sectionElem.classList.add('raid-stats-section');
			this.rootElem.appendChild(sectionElem);
			sectionElem.innerHTML = `
				<span class="raid-stats-section-label"${section.label}></span>
				<div class="raid-stats-section-content"></div>
			`;
			const contentElem = sectionElem.getElementsByClassName('raid-stats-section-content')[0] as HTMLDivElement;

			section.categories.forEach(categoryOptions => {
				categories.push(new RaidStatsCategory(contentElem, raidSimUI, categoryOptions));
			});
		});
		this.categories = categories;

		raidSimUI.changeEmitter.on(eventID => this.categories.forEach(c => c.update()));
	}
}

class RaidStatsCategory extends Component {
	readonly raidSimUI: RaidSimUI;
	private readonly options: RaidStatsCategoryOptions;
	private readonly effects: Array<RaidStatsEffect>;
	private readonly counterElem: HTMLElement;
	private readonly tooltipElem: HTMLElement;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI, options: RaidStatsCategoryOptions) {
		super(parent, 'raid-stats-category');
		this.raidSimUI = raidSimUI;
		this.options = options;
		this.rootElem.innerHTML = `
			<span class="raid-stats-category-counter"></span>
			<span class="raid-stats-category-label">${options.label}</span>
		`;

		this.counterElem = this.rootElem.getElementsByClassName('raid-stats-category-counter')[0] as HTMLElement;
		this.tooltipElem = document.createElement('div');
		this.tooltipElem.classList.add('raid-stats-category-tooltip');

		this.effects = options.effects.map(opt => new RaidStatsEffect(this.tooltipElem, raidSimUI, opt));

		tippy(this.rootElem, {
			content: this.tooltipElem,
			allowHTML: true,
		});
	}

	update() {
		this.effects.forEach(effect => effect.update());

		const total = sum(this.effects.map(effect => effect.count));
		this.counterElem.textContent = String(total);
		if (total == 0) {
			this.rootElem.classList.remove('active');
		} else {
			this.rootElem.classList.add('active');
		}
	}
}

class RaidStatsEffect extends Component {
	readonly raidSimUI: RaidSimUI;
	private readonly options: RaidStatsEffectOptions;
	private readonly counterElem: HTMLElement;

	curPlayers: Array<Player<any>>;
	curBuffBots: Array<BuffBot>;
	count: number;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI, options: RaidStatsEffectOptions) {
		super(parent, 'raid-stats-effect');
		this.raidSimUI = raidSimUI;
		this.options = options;

		this.curPlayers = [];
		this.curBuffBots = [];
		this.count = 0;

		this.rootElem.innerHTML = `
			<span class="raid-stats-effect-counter"></span>
			<a class="raid-stats-effect-icon"></a>
			<span class="raid-stats-effect-label">${options.label}</span>
		`;
		this.counterElem = this.rootElem.getElementsByClassName('raid-stats-effect-counter')[0] as HTMLElement;

		const iconElem = this.rootElem.getElementsByClassName('raid-stats-effect-icon')[0] as HTMLAnchorElement;
		if (options.actionId) {
			options.actionId.fillAndSet(iconElem, false, true);
		} else {
			iconElem.remove();
		}
	}

	update() {
		if (this.options.providedByPlayer) {
			this.curPlayers = this.raidSimUI.getPlayers().filter(p => p != null && this.options.providedByPlayer!(p)) as Array<Player<any>>;
		}
		if (this.options.providedByBuffBot) {
			this.curBuffBots = this.raidSimUI.getBuffBots().filter(b => this.options.providedByBuffBot!(b));
		}

		const providedByRaid = this.options.providedByRaid && this.options.providedByRaid(this.raidSimUI.sim.raid);

		this.count = this.curPlayers.length + this.curBuffBots.length + (providedByRaid ? 1 : 0);
		this.counterElem.textContent = String(this.count);
		if (this.count == 0) {
			this.rootElem.classList.remove('active');
		} else {
			this.rootElem.classList.add('active');
		}
	}
}

function negateIf(val: boolean, cond: boolean): boolean {
	return cond ? !val : val;
}

function playerClass(clazz: Class): PlayerProvider {
	return (player: Player<any>): boolean => {
		return player.getClass() == clazz;
	};
}
function playerClassAndTalentInternal<T extends Class>(clazz: T, talentName: keyof SpecTalents<ClassSpecs<T>>, negateTalent: boolean, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return (player: Player<any>): boolean => {
		return player.getClass() == clazz
			&& negateIf(Boolean((player as Player<ClassSpecs<T>>).getTalents()[talentName]), negateTalent)
			&& (!extraCondition || extraCondition(player as Player<ClassSpecs<T>>));
	};
}
function playerClassAndTalent<T extends Class>(clazz: T, talentName: keyof SpecTalents<ClassSpecs<T>>, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return playerClassAndTalentInternal(clazz, talentName, false, extraCondition);
}
function playerClassAndMissingTalent<T extends Class>(clazz: T, talentName: keyof SpecTalents<ClassSpecs<T>>, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return playerClassAndTalentInternal(clazz, talentName, true, extraCondition);
}
function playerSpecAndTalent<T extends Spec>(spec: T, talentName: keyof SpecTalents<T>): PlayerProvider {
	return (player: Player<any>): boolean => {
		return player.spec == spec && Boolean((player as Player<T>).getTalents()[talentName]);
	};
}

function buffBotAny(ids: Array<BuffBotId>): BuffBotProvider {
	return (buffBot: BuffBot): boolean => {
		return ids.includes(buffBot.id);
	};
}
function buffBotClass(clazz: Class): BuffBotProvider {
	return (buffBot: BuffBot): boolean => {
		return buffBot.getClass() == clazz;
	};
}

function raidBuff(buffName: keyof RaidBuffs): RaidProvider {
	return (raid: Raid): boolean => {
		return Boolean(raid.getBuffs()[buffName]);
	};
}

const RAID_STATS_OPTIONS: RaidStatsOptions = {sections: [
	{
		label: 'Buffs',
		categories: [
			{
				label: 'Stats',
				effects: [
					{
						label: 'Improved Gift of the Wild',
						actionId: ActionId.fromSpellId(17051),
						providedByPlayer: playerClassAndTalent(Class.ClassDruid, 'improvedMarkOfTheWild'),
						providedByBuffBot: buffBotAny(['Resto Druid']),
					},
					{
						label: 'Gift of the Wild',
						actionId: ActionId.fromSpellId(48470),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassDruid, 'improvedMarkOfTheWild'),
					},
					{
						label: 'Drums of the Wild',
						actionId: ActionId.fromItemId(49634),
						providedByRaid: raidBuff('drumsOfTheWild'),
					},
				],
			},
			{
				label: 'Stats %',
				effects: [
					{
						label: 'Blessing of Kings',
						actionId: ActionId.fromSpellId(25898),
						providedByPlayer: playerClass(Class.ClassPaladin),
						providedByBuffBot: buffBotClass(Class.ClassPaladin),
					},
					{
						label: 'Drums of Forgotten Kings',
						actionId: ActionId.fromItemId(49633),
						providedByRaid: raidBuff('drumsOfForgottenKings'),
					},
					{
						label: 'Blessing of Sanctuary',
						actionId: ActionId.fromSpellId(25899),
						providedByPlayer: playerClass(Class.ClassPaladin),
						providedByBuffBot: buffBotClass(Class.ClassPaladin),
					},
				],
			},
			{
				label: 'Armor',
				effects: [
					{
						label: 'Improved Devotion Aura',
						actionId: ActionId.fromSpellId(20140),
						providedByPlayer: playerClassAndTalent(Class.ClassPaladin, 'improvedDevotionAura', player => player.getSpecOptions().aura == PaladinAura.DevotionAura),
					},
					{
						label: 'Devotion Aura',
						actionId: ActionId.fromSpellId(48942),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassPaladin, 'improvedDevotionAura', player => player.getSpecOptions().aura == PaladinAura.DevotionAura),
					},
					// TODO: Implement stoneskin totem
					//{
					//	label: 'Improved Stoneskin Totem',
					//	actionId: ActionId.fromSpellId(16293),
					//	providedByPlayer: playerClassAndTalent(Class.ClassShaman, 'guardianTotems', player => player.getRotation().totems?.earth == EarthTotem.StoneskinTotem),
					//},
					//{
					//	label: 'Stoneskin Totem',
					//	actionId: ActionId.fromSpellId(58753),
					//	providedByPlayer: playerClassAndMissingTalent(Class.ClassShaman, 'guardianTotems', player => player.getRotation().totems?.earth == EarthTotem.StoneskinTotem),
					//},
					{
						label: 'Scroll of Protection',
						actionId: ActionId.fromItemId(43468),
						providedByRaid: raidBuff('scrollOfProtection'),
					},
				],
			},
			{
				label: 'Stamina',
				effects: [
					{
						label: 'Improved Power Word Fortitude',
						actionId: ActionId.fromSpellId(14767),
						providedByPlayer: playerClassAndTalent(Class.ClassPriest, 'improvedPowerWordFortitude'),
						providedByBuffBot: buffBotClass(Class.ClassPriest),
					},
					{
						label: 'Power Word Fortitude',
						actionId: ActionId.fromSpellId(48161),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassPriest, 'improvedPowerWordFortitude'),
					},
					{
						label: 'Scroll of Stamina',
						actionId: ActionId.fromItemId(37094),
						providedByRaid: raidBuff('scrollOfStamina'),
					},
				],
			},
			{
				label: 'Strength',
				effects: [
					{
						label: 'Improved Strength of Earth Totem',
						actionId: ActionId.fromSpellId(52456),
						providedByPlayer: playerClassAndTalent(Class.ClassShaman, 'enhancingTotems', player => player.getRotation().totems?.earth == EarthTotem.StrengthOfEarthTotem),
					},
					{
						label: 'Strength of Earth Totem',
						actionId: ActionId.fromSpellId(58643),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassShaman, 'enhancingTotems', player => player.getRotation().totems?.earth == EarthTotem.StrengthOfEarthTotem),
					},
					{
						label: 'Horn of Winter',
						actionId: ActionId.fromSpellId(57643),
						providedByPlayer: playerClass(Class.ClassDeathknight),
						providedByBuffBot: buffBotClass(Class.ClassDeathknight),
					},
					{
						label: 'Scroll of Strength',
						actionId: ActionId.fromItemId(43466),
						providedByRaid: raidBuff('scrollOfStrength'),
					},
				],
			},
			{
				label: 'Agility',
				effects: [
					{
						label: 'Improved Strength of Earth Totem',
						actionId: ActionId.fromSpellId(52456),
						providedByPlayer: playerClassAndTalent(Class.ClassShaman, 'enhancingTotems', player => player.getRotation().totems?.earth == EarthTotem.StrengthOfEarthTotem),
					},
					{
						label: 'Strength of Earth Totem',
						actionId: ActionId.fromSpellId(58643),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassShaman, 'enhancingTotems', player => player.getRotation().totems?.earth == EarthTotem.StrengthOfEarthTotem),
					},
					{
						label: 'Horn of Winter',
						actionId: ActionId.fromSpellId(57643),
						providedByPlayer: playerClass(Class.ClassDeathknight),
						providedByBuffBot: buffBotClass(Class.ClassDeathknight),
					},
					{
						label: 'Scroll of Agility',
						actionId: ActionId.fromItemId(43464),
						providedByRaid: raidBuff('scrollOfAgility'),
					},
				],
			},
			{
				label: 'Intellect',
				effects: [
					{
						label: 'Arcane Brilliance',
						actionId: ActionId.fromSpellId(43002),
						providedByPlayer: playerClass(Class.ClassMage),
					},
					// TODO: Double-check the talent
					{
						label: 'Improved Fel Intelligence',
						actionId: ActionId.fromSpellId(54038),
						providedByPlayer: playerClassAndTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
					},
					{
						label: 'Fel Intelligence',
						actionId: ActionId.fromSpellId(58643),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
					},
					{
						label: 'Scroll of Intellect',
						actionId: ActionId.fromItemId(37092),
						providedByRaid: raidBuff('scrollOfIntellect'),
					},
				],
			},
			{
				label: 'Spirit',
				effects: [
					// TODO: Implement divine spirit
					{
						label: 'Divine Spirit',
						actionId: ActionId.fromSpellId(48073),
						providedByPlayer: playerClass(Class.ClassPriest),
					},
					// TODO: Double-check the talent
					{
						label: 'Improved Fel Intelligence',
						actionId: ActionId.fromSpellId(54038),
						providedByPlayer: playerClassAndTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
					},
					{
						label: 'Fel Intelligence',
						actionId: ActionId.fromSpellId(58643),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
					},
					{
						label: 'Scroll of Spirit',
						actionId: ActionId.fromItemId(37098),
						providedByRaid: raidBuff('scrollOfSpirit'),
					},
				],
			},
			{
				label: 'Atk Pwr',
				effects: [
					{
						label: 'Improved Blessing of Might',
						actionId: ActionId.fromSpellId(20045),
						providedByPlayer: playerClass(Class.ClassPaladin),
						providedByBuffBot: buffBotClass(Class.ClassPaladin),
					},
					{
						label: 'Blessing of Might',
						actionId: ActionId.fromSpellId(48934),
					},
					{
						label: 'Improved Battle Shout',
						actionId: ActionId.fromSpellId(12861),
						providedByPlayer: playerClassAndTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle),
					},
					{
						label: 'Battle Shout',
						actionId: ActionId.fromSpellId(47436),
						providedByPlayer: playerClassAndMissingTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle),
					},
				],
			},
			{
				label: 'Atk Pwr %',
				effects: [
					{
						label: 'Abomination\'s Might',
						actionId: ActionId.fromSpellId(53138),
						providedByPlayer: playerClassAndTalent(Class.ClassDeathknight, 'abominationsMight'),
					},
					{
						label: 'Unleashed Rage',
						actionId: ActionId.fromSpellId(30809),
						providedByPlayer: playerClassAndTalent(Class.ClassShaman, 'unleashedRage'),
					},
					{
						label: 'Trueshot Aura',
						actionId: ActionId.fromSpellId(19506),
						providedByPlayer: playerClassAndTalent(Class.ClassHunter, 'trueshotAura'),
					},
				],
			},
			{
				label: 'Damage %',
				effects: [
					{
						label: 'Sanctified Retribution',
						actionId: ActionId.fromSpellId(31869),
						providedByPlayer: playerClassAndTalent(Class.ClassPaladin, 'sanctifiedRetribution'),
					},
					{
						label: 'Arcane Empowerment',
						actionId: ActionId.fromSpellId(31583),
						providedByPlayer: playerClassAndTalent(Class.ClassMage, 'arcaneEmpowerment'),
					},
					{
						label: 'Ferocious Inspiration',
						actionId: ActionId.fromSpellId(34460),
						providedByPlayer: playerClassAndTalent(Class.ClassHunter, 'ferociousInspiration'),
					},
				],
			},
			{
				label: 'Mit %',
				effects: [
					{
						label: 'Renewed Hope',
						actionId: ActionId.fromSpellId(57472),
						providedByPlayer: playerClassAndTalent(Class.ClassPriest, 'renewedHope'),
					},
					{
						label: 'Blessing Of Sanctuary',
						actionId: ActionId.fromSpellId(25899),
						providedByPlayer: playerClass(Class.ClassPaladin),
					},
					{
						label: 'Vigilance',
						actionId: ActionId.fromSpellId(50720),
						providedByPlayer: playerClassAndTalent(Class.ClassWarrior, 'vigilance'),
					},
				],
			},
			{
				label: 'Replenishment',
				effects: [
					{
						label: 'Vampiric Touch',
						actionId: ActionId.fromSpellId(48160),
						providedByPlayer: playerSpecAndTalent(Spec.SpecShadowPriest, 'vampiricTouch'),
					},
					{
						label: 'Judgements of the Wise',
						actionId: ActionId.fromSpellId(31878),
						providedByPlayer: playerSpecAndTalent(Spec.SpecRetributionPaladin, 'judgementsOfTheWise'),
					},
					{
						label: 'Hunting Party',
						actionId: ActionId.fromSpellId(53292),
						providedByPlayer: playerSpecAndTalent(Spec.SpecHunter, 'huntingParty'),
					},
					{
						label: 'Improved Soul Leech',
						actionId: ActionId.fromSpellId(54118),
						providedByPlayer: playerSpecAndTalent(Spec.SpecWarlock, 'improvedSoulLeech'),
					},
					{
						label: 'Enduring Winter',
						actionId: ActionId.fromSpellId(44561),
						providedByPlayer: playerSpecAndTalent(Spec.SpecMage, 'enduringWinter'),
					},
				],
			},
		],
	},
	{
		label: 'Debuffs',
		categories: [
		],
	},
]};
