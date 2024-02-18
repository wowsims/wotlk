import {
	Class,
	RaidBuffs,
	Spec,
} from '../core/proto/common.js';
import { Component } from '../core/components/component.js';
import { Player } from "../core/player.js";
import { Raid } from "../core/raid.js";
import { ActionId } from '../core/proto_utils/action_id.js';
import {
	ClassSpecs,
	SpecTalents,
	specToClass,
	isTankSpec,
	isHealingSpec,
	isMeleeDpsSpec,
	isRangedDpsSpec,
	textCssClassForClass,
} from '../core/proto_utils/utils.js';
import { sum } from '../core/utils.js';

import { Hunter_Options_PetType as HunterPetType } from '../core/proto/hunter.js';
import { PaladinAura } from '../core/proto/paladin.js';
import { AirTotem, EarthTotem, FireTotem, WaterTotem } from '../core/proto/shaman.js';
import { Warlock_Options_Summon as WarlockSummon } from '../core/proto/warlock.js';
import { WarriorShout } from '../core/proto/warrior.js';

import { RaidSimUI } from './raid_sim_ui.js';
import { Tooltip } from 'bootstrap';

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

type PlayerProvider = { class?: Class, condition: (player: Player<any>) => boolean };
type RaidProvider = (raid: Raid) => boolean;

interface RaidStatsEffectOptions {
	label: string,
	actionId?: ActionId,
	playerData?: PlayerProvider,
	raidData?: RaidProvider,
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
				<div class="raid-stats-section-header">
					<label class="raid-stats-section-label form-label">${section.label}</label>
				</div>
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
		super(parent, 'raid-stats-category-root');
		this.raidSimUI = raidSimUI;
		this.options = options;

		this.rootElem.innerHTML = `
			<a href="javascript:void(0)" role="button" class="raid-stats-category">
				<span class="raid-stats-category-counter"></span>
				<span class="raid-stats-category-label">${options.label}</span>
			</a>
		`;

		this.counterElem = this.rootElem.querySelector('.raid-stats-category-counter') as HTMLElement;
		this.tooltipElem = document.createElement('div');
		this.tooltipElem.innerHTML = `
			<label class="raid-stats-category-label">${options.label}</label>
		`

		this.effects = options.effects.map(opt => new RaidStatsEffect(this.tooltipElem, raidSimUI, opt));

		if (options.effects.length != 1 || options.effects[0].playerData?.class) {
			const statsLink = this.rootElem.querySelector('.raid-stats-category') as HTMLElement;

			// Using the title option here because outerHTML sanitizes and filters out the img src options
			Tooltip.getOrCreateInstance(statsLink, {
				customClass: 'raid-stats-category-tooltip',
				html: true,
				placement: 'right',
				title: this.tooltipElem,
			})
		}
	}

	update() {
		this.effects.forEach(effect => effect.update());

		const total = sum(this.effects.map(effect => effect.count));
		this.counterElem.textContent = String(total);

		const statsLink = this.rootElem.querySelector('.raid-stats-category') as HTMLElement;

		if (total == 0) {
			statsLink?.classList.remove('active');
		} else {
			statsLink?.classList.add('active');
		}
	}
}

class RaidStatsEffect extends Component {
	readonly raidSimUI: RaidSimUI;
	private readonly options: RaidStatsEffectOptions;
	private readonly counterElem: HTMLElement;

	curPlayers: Array<Player<any>>;
	count: number;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI, options: RaidStatsEffectOptions) {
		super(parent, 'raid-stats-effect');
		this.raidSimUI = raidSimUI;
		this.options = options;

		this.curPlayers = [];
		this.count = 0;

		this.rootElem.innerHTML = `
			<span class="raid-stats-effect-counter"></span>
			<img class="raid-stats-effect-icon"></img>
			<span class="raid-stats-effect-label">${options.label}</span>
		`;
		this.counterElem = this.rootElem.querySelector('.raid-stats-effect-counter') as HTMLElement;

		if (this.options.playerData?.class) {
			const labelElem = this.rootElem.querySelector('.raid-stats-effect-label') as HTMLElement;
			const playerCssClass = textCssClassForClass(this.options.playerData.class);
			labelElem.classList.add(playerCssClass);
		}

		const iconElem = this.rootElem.querySelector('.raid-stats-effect-icon') as HTMLImageElement;
		if (options.actionId) {
			options.actionId.fill().then(actionId => iconElem.src = actionId.iconUrl);
		} else {
			iconElem.remove();
		}
	}

	update() {
		if (this.options.playerData) {
			this.curPlayers = this.raidSimUI.getActivePlayers().filter(p => this.options.playerData!.condition(p));
		}

		const raidData = this.options.raidData && this.options.raidData(this.raidSimUI.sim.raid);

		this.count = this.curPlayers.length + (raidData ? 1 : 0);
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

function playerClass<T extends Class>(clazz: T, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return {
		class: clazz,
		condition: (player: Player<any>): boolean => {
			return player.isClass(clazz)
				&& (!extraCondition || extraCondition(player));
		},
	};
}
function playerClassAndTalentInternal<T extends Class>(clazz: T, talentName: keyof SpecTalents<ClassSpecs<T>>, negateTalent: boolean, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return {
		class: clazz,
		condition: (player: Player<any>): boolean => {
			return player.isClass(clazz)
				&& negateIf(Boolean((player.getTalents() as any)[talentName]), negateTalent)
				&& (!extraCondition || extraCondition(player));
		},
	};
}
function playerClassAndTalent<T extends Class>(clazz: T, talentName: keyof SpecTalents<ClassSpecs<T>>, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return playerClassAndTalentInternal(clazz, talentName, false, extraCondition);
}
function playerClassAndMissingTalent<T extends Class>(clazz: T, talentName: keyof SpecTalents<ClassSpecs<T>>, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return playerClassAndTalentInternal(clazz, talentName, true, extraCondition);
}
function playerSpecAndTalentInternal<T extends Spec>(spec: T, talentName: keyof SpecTalents<T>, negateTalent: boolean, extraCondition?: (player: Player<T>) => boolean): PlayerProvider {
	return {
		class: specToClass[spec],
		condition: (player: Player<any>): boolean => {
			return player.isSpec(spec)
				&& negateIf(Boolean((player.getTalents() as any)[talentName]), negateTalent)
				&& (!extraCondition || extraCondition(player));
		},
	};
}
function playerSpecAndTalent<T extends Spec>(spec: T, talentName: keyof SpecTalents<T>, extraCondition?: (player: Player<T>) => boolean): PlayerProvider {
	return playerSpecAndTalentInternal(spec, talentName, false, extraCondition);
}
function playerSpecAndMissingTalent<T extends Spec>(spec: T, talentName: keyof SpecTalents<T>, extraCondition?: (player: Player<T>) => boolean): PlayerProvider {
	return playerSpecAndTalentInternal(spec, talentName, true, extraCondition);
}

function raidBuff(buffName: keyof RaidBuffs): RaidProvider {
	return (raid: Raid): boolean => {
		return Boolean(raid.getBuffs()[buffName]);
	};
}

const RAID_STATS_OPTIONS: RaidStatsOptions = {
	sections: [
		{
			label: 'Roles',
			categories: [
				{
					label: 'Tanks',
					effects: [
						{
							label: 'Tanks',
							playerData: { condition: player => isTankSpec(player.spec) },
						},
					],
				},
				{
					label: 'Healers',
					effects: [
						{
							label: 'Healers',
							playerData: { condition: player => isHealingSpec(player.spec) },
						},
					],
				},
				{
					label: 'Melee',
					effects: [
						{
							label: 'Melee',
							playerData: { condition: player => isMeleeDpsSpec(player.spec) },
						},
					],
				},
				{
					label: 'Ranged',
					effects: [
						{
							label: 'Ranged',
							playerData: { condition: player => isRangedDpsSpec(player.spec) },
						},
					],
				},
			],
		},
		{
			label: 'Buffs',
			categories: [
				{
					label: 'Bloodlust',
					effects: [
						{
							label: 'Bloodlust',
							actionId: ActionId.fromSpellId(2825),
							playerData: playerClass(Class.ClassShaman),
						},
					],
				},
				{
					label: 'Stats',
					effects: [
						{
							label: 'Improved Gift of the Wild',
							actionId: ActionId.fromSpellId(17051),
							playerData: playerClassAndTalent(Class.ClassDruid, 'improvedMarkOfTheWild'),
						},
						{
							label: 'Gift of the Wild',
							actionId: ActionId.fromSpellId(48470),
							playerData: playerClassAndMissingTalent(Class.ClassDruid, 'improvedMarkOfTheWild'),
						},
						{
							label: 'Drums of the Wild',
							actionId: ActionId.fromItemId(49634),
							raidData: raidBuff('drumsOfTheWild'),
						},
					],
				},
				{
					label: 'Stats %',
					effects: [
						{
							label: 'Blessing of Kings',
							actionId: ActionId.fromSpellId(25898),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Drums of Forgotten Kings',
							actionId: ActionId.fromItemId(49633),
							raidData: raidBuff('drumsOfForgottenKings'),
						},
						{
							label: 'Blessing of Sanctuary',
							actionId: ActionId.fromSpellId(25899),
							playerData: playerClass(Class.ClassPaladin),
						},
					],
				},
				{
					label: 'Armor',
					effects: [
						{
							label: 'Improved Devotion Aura',
							actionId: ActionId.fromSpellId(20140),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'improvedDevotionAura', player => player.getSpecOptions().aura == PaladinAura.DevotionAura),
						},
						{
							label: 'Devotion Aura',
							actionId: ActionId.fromSpellId(48942),
							playerData: playerClassAndMissingTalent(Class.ClassPaladin, 'improvedDevotionAura', player => player.getSpecOptions().aura == PaladinAura.DevotionAura),
						},
						{
							label: 'Improved Stoneskin Totem',
							actionId: ActionId.fromSpellId(16293),
							playerData: playerClassAndTalent(Class.ClassShaman, 'guardianTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StoneskinTotem),
						},
						{
							label: 'Stoneskin Totem',
							actionId: ActionId.fromSpellId(58753),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'guardianTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StoneskinTotem),
						},
						{
							label: 'Scroll of Protection',
							actionId: ActionId.fromItemId(43468),
							raidData: raidBuff('scrollOfProtection'),
						},
					],
				},
				{
					label: 'Stamina',
					effects: [
						{
							label: 'Improved Power Word Fortitude',
							actionId: ActionId.fromSpellId(14767),
							playerData: playerClassAndTalent(Class.ClassPriest, 'improvedPowerWordFortitude'),
						},
						{
							label: 'Power Word Fortitude',
							actionId: ActionId.fromSpellId(48161),
							playerData: playerClassAndMissingTalent(Class.ClassPriest, 'improvedPowerWordFortitude'),
						},
						{
							label: 'Scroll of Stamina',
							actionId: ActionId.fromItemId(37094),
							raidData: raidBuff('scrollOfStamina'),
						},
					],
				},
				{
					label: 'Str + Agi',
					effects: [
						{
							label: 'Improved Strength of Earth Totem',
							actionId: ActionId.fromSpellId(52456),
							playerData: playerClassAndTalent(Class.ClassShaman, 'enhancingTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StrengthOfEarthTotem),
						},
						{
							label: 'Strength of Earth Totem',
							actionId: ActionId.fromSpellId(58643),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'enhancingTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StrengthOfEarthTotem),
						},
						{
							label: 'Horn of Winter',
							actionId: ActionId.fromSpellId(57623),
							playerData: playerClass(Class.ClassDeathknight),
						},
						{
							label: 'Scroll of Strength',
							actionId: ActionId.fromItemId(43466),
							raidData: raidBuff('scrollOfStrength'),
						},
						{
							label: 'Scroll of Agility',
							actionId: ActionId.fromItemId(43464),
							raidData: raidBuff('scrollOfAgility'),
						},
					],
				},
				{
					label: 'Intellect',
					effects: [
						{
							label: 'Arcane Brilliance',
							actionId: ActionId.fromSpellId(43002),
							playerData: playerClass(Class.ClassMage),
						},
						{
							label: 'Improved Fel Intelligence',
							actionId: ActionId.fromSpellId(54038),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: 'Fel Intelligence',
							actionId: ActionId.fromSpellId(57567),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: 'Scroll of Intellect',
							actionId: ActionId.fromItemId(37092),
							raidData: raidBuff('scrollOfIntellect'),
						},
					],
				},
				{
					label: 'Spirit',
					effects: [
						{
							label: 'Divine Spirit',
							actionId: ActionId.fromSpellId(48073),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Improved Fel Intelligence',
							actionId: ActionId.fromSpellId(54038),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: 'Fel Intelligence',
							actionId: ActionId.fromSpellId(57567),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: 'Scroll of Spirit',
							actionId: ActionId.fromItemId(37098),
							raidData: raidBuff('scrollOfSpirit'),
						},
					],
				},
				{
					label: 'Atk Pwr',
					effects: [
						{
							label: 'Improved Blessing of Might',
							actionId: ActionId.fromSpellId(20045),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Blessing of Might',
							actionId: ActionId.fromSpellId(48934),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Improved Battle Shout',
							actionId: ActionId.fromSpellId(12861),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle),
						},
						{
							label: 'Battle Shout',
							actionId: ActionId.fromSpellId(47436),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle),
						},
					],
				},
				{
					label: 'Atk Pwr %',
					effects: [
						{
							label: 'Abomination\'s Might',
							actionId: ActionId.fromSpellId(53138),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'abominationsMight'),
						},
						{
							label: 'Unleashed Rage',
							actionId: ActionId.fromSpellId(30809),
							playerData: playerClassAndTalent(Class.ClassShaman, 'unleashedRage'),
						},
						{
							label: 'Trueshot Aura',
							actionId: ActionId.fromSpellId(19506),
							playerData: playerClassAndTalent(Class.ClassHunter, 'trueshotAura'),
						},
					],
				},
				{
					label: 'Damage %',
					effects: [
						{
							label: 'Sanctified Retribution',
							actionId: ActionId.fromSpellId(31869),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'sanctifiedRetribution'),
						},
						{
							label: 'Arcane Empowerment',
							actionId: ActionId.fromSpellId(31583),
							playerData: playerClassAndTalent(Class.ClassMage, 'arcaneEmpowerment'),
						},
						{
							label: 'Ferocious Inspiration',
							actionId: ActionId.fromSpellId(34460),
							playerData: playerClassAndTalent(Class.ClassHunter, 'ferociousInspiration'),
						},
					],
				},
				{
					label: 'Mit %',
					effects: [
						{
							label: 'Renewed Hope',
							actionId: ActionId.fromSpellId(57472),
							playerData: playerClassAndTalent(Class.ClassPriest, 'renewedHope'),
						},
						{
							label: 'Blessing Of Sanctuary',
							actionId: ActionId.fromSpellId(25899),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Vigilance',
							actionId: ActionId.fromSpellId(50720),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'vigilance'),
						},
					],
				},
				{
					label: 'Haste %',
					effects: [
						{
							label: 'Swift Retribution',
							actionId: ActionId.fromSpellId(53648),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'swiftRetribution'),
						},
						{
							label: 'Improved Moonkin Form',
							actionId: ActionId.fromSpellId(48396),
							playerData: playerClassAndTalent(Class.ClassDruid, 'improvedMoonkinForm'),
						},
					],
				},
				{
					label: 'MP5',
					effects: [
						{
							label: 'Improved Blessing of Wisdom',
							actionId: ActionId.fromSpellId(20245),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'improvedBlessingOfWisdom'),
						},
						{
							label: 'Blessing of Wisdom',
							actionId: ActionId.fromSpellId(48938),
							playerData: playerClassAndMissingTalent(Class.ClassPaladin, 'improvedBlessingOfWisdom'),
						},
						{
							label: 'Improved Mana Spring Totem',
							actionId: ActionId.fromSpellId(16206),
							playerData: playerClassAndTalent(Class.ClassShaman, 'restorativeTotems', player => player.getSpecOptions().totems?.water == WaterTotem.ManaSpringTotem),
						},
						{
							label: 'Mana Spring Totem',
							actionId: ActionId.fromSpellId(58774),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'restorativeTotems', player => player.getSpecOptions().totems?.water == WaterTotem.ManaSpringTotem),
						},
					],
				},
				{
					label: 'Melee Crit',
					effects: [
						{
							label: 'Leader of the Pack',
							actionId: ActionId.fromSpellId(17007),
							playerData: playerClassAndTalent(Class.ClassDruid, 'leaderOfThePack'),
						},
						{
							label: 'Rampage',
							actionId: ActionId.fromSpellId(29801),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'rampage'),
						},
					],
				},
				{
					label: 'Melee Haste',
					effects: [
						{
							label: 'Improved Icy Talons',
							actionId: ActionId.fromSpellId(55610),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'improvedIcyTalons'),
						},
						{
							label: 'Improved Windfury Totem',
							actionId: ActionId.fromSpellId(29193),
							playerData: playerClassAndTalent(Class.ClassShaman, 'improvedWindfuryTotem', player => player.getSpecOptions().totems?.air == AirTotem.WindfuryTotem),
						},
						{
							label: 'Windfury Totem',
							actionId: ActionId.fromSpellId(65990),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'improvedWindfuryTotem', player => player.getSpecOptions().totems?.air == AirTotem.WindfuryTotem),
						},
					],
				},
				{
					label: 'Spell Power',
					effects: [
						{
							label: 'Demonic Pact',
							actionId: ActionId.fromSpellId(47240),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'demonicPact'),
						},
						{
							label: 'Totem of Wrath',
							actionId: ActionId.fromSpellId(57722),
							playerData: playerClassAndTalent(Class.ClassShaman, 'totemOfWrath', player => player.getSpecOptions().totems?.fire == FireTotem.TotemOfWrath),
						},
						{
							label: 'Flametongue Totem',
							actionId: ActionId.fromSpellId(58656),
							playerData: playerClass(Class.ClassShaman, player => player.getSpecOptions().totems?.fire == FireTotem.FlametongueTotem),
						},
					],
				},
				{
					label: 'Spell Crit',
					effects: [
						{
							label: 'Moonkin Form',
							actionId: ActionId.fromSpellId(24907),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'moonkinForm'),
						},
						{
							label: 'Elemental Oath',
							actionId: ActionId.fromSpellId(51470),
							playerData: playerClassAndTalent(Class.ClassShaman, 'elementalOath'),
						},
					],
				},
				{
					label: 'Spell Haste',
					effects: [
						{
							label: 'Wrath of Air Totem',
							actionId: ActionId.fromSpellId(3738),
							playerData: playerClass(Class.ClassShaman, player => player.getSpecOptions().totems?.air == AirTotem.WrathOfAirTotem),
						},
					],
				},
				{
					label: 'Health',
					effects: [
						{
							label: 'Improved Commanding Shout',
							actionId: ActionId.fromSpellId(12861),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutCommanding),
						},
						{
							label: 'Commanding Shout',
							actionId: ActionId.fromSpellId(47440),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutCommanding),
						},
						{
							label: 'Improved Imp',
							actionId: ActionId.fromSpellId(18696),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedImp', player => player.getSpecOptions().summon == WarlockSummon.Imp),
						},
						{
							label: 'Blood Pact',
							actionId: ActionId.fromSpellId(47982),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedImp', player => player.getSpecOptions().summon == WarlockSummon.Imp),
						},
					],
				},
				{
					label: 'Replenishment',
					effects: [
						{
							label: 'Vampiric Touch',
							actionId: ActionId.fromSpellId(48160),
							playerData: playerSpecAndTalent(Spec.SpecShadowPriest, 'vampiricTouch'),
						},
						{
							label: 'Judgements of the Wise',
							actionId: ActionId.fromSpellId(31878),
							playerData: playerSpecAndTalent(Spec.SpecRetributionPaladin, 'judgementsOfTheWise'),
						},
						{
							label: 'Hunting Party',
							actionId: ActionId.fromSpellId(53292),
							playerData: playerSpecAndTalent(Spec.SpecHunter, 'huntingParty'),
						},
						{
							label: 'Improved Soul Leech',
							actionId: ActionId.fromSpellId(54118),
							playerData: playerSpecAndTalent(Spec.SpecWarlock, 'improvedSoulLeech'),
						},
						{
							label: 'Enduring Winter',
							actionId: ActionId.fromSpellId(44561),
							playerData: playerSpecAndTalent(Spec.SpecMage, 'enduringWinter'),
						},
					],
				},
			],
		},
		{
			label: 'External Buffs',
			categories: [
				{
					label: 'Innervate',
					effects: [
						{
							label: 'Innervate',
							actionId: ActionId.fromSpellId(29166),
							playerData: playerClass(Class.ClassDruid),
						},
					],
				},
				{
					label: 'Power Infusion',
					effects: [
						{
							label: 'Power Infusion',
							actionId: ActionId.fromSpellId(10060),
							playerData: playerClassAndTalent(Class.ClassPriest, 'powerInfusion'),
						},
					],
				},
				{
					label: 'Focus Magic',
					effects: [
						{
							label: 'Focus Magic',
							actionId: ActionId.fromSpellId(54648),
							playerData: playerClassAndTalent(Class.ClassMage, 'focusMagic'),
						},
					],
				},
				{
					label: 'Tricks of the Trade',
					effects: [
						{
							label: 'Tricks of the Trade',
							actionId: ActionId.fromSpellId(57933),
							playerData: playerClass(Class.ClassRogue),
						},
					],
				},
				{
					label: 'Unholy Frenzy',
					effects: [
						{
							label: 'Unholy Frenzy',
							actionId: ActionId.fromSpellId(49016),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'hysteria'),
						},
					],
				},
				{
					label: 'Pain Suppression',
					effects: [
						{
							label: 'Pain Suppression',
							actionId: ActionId.fromSpellId(33206),
							playerData: playerClassAndTalent(Class.ClassPriest, 'painSuppression'),
						},
					],
				},
				{
					label: 'Divine Guardian',
					effects: [
						{
							label: 'Divine Guardian',
							actionId: ActionId.fromSpellId(53530),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'divineGuardian'),
						},
					],
				},
			],
		},
		{
			label: 'DPS Debuffs',
			categories: [
				{
					label: 'Major ArP',
					effects: [
						{
							label: 'Sunder Armor',
							actionId: ActionId.fromSpellId(47467),
							playerData: playerClass(Class.ClassWarrior),
						},
						{
							label: 'Expose Armor',
							actionId: ActionId.fromSpellId(8647),
							playerData: playerClass(Class.ClassRogue),
						},
						{
							label: 'Acid Spit',
							actionId: ActionId.fromSpellId(55754),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Worm),
						},
					],
				},
				{
					label: 'Minor ArP',
					effects: [
						{
							label: 'Faerie Fire',
							actionId: ActionId.fromSpellId(770),
							playerData: playerClass(Class.ClassDruid, player => player.spec != Spec.SpecRestorationDruid),
						},
						{
							label: 'Curse of Weakness',
							actionId: ActionId.fromSpellId(50511),
							playerData: playerClass(Class.ClassWarlock),
						},
						{
							label: 'Sting',
							actionId: ActionId.fromSpellId(56631),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Wasp),
						},
						{
							label: 'Spore Cloud',
							actionId: ActionId.fromSpellId(53598),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Bat),
						},
					],
				},
				{
					label: 'Phys Vuln',
					effects: [
						{
							label: 'Blood Frenzy',
							actionId: ActionId.fromSpellId(29859),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'bloodFrenzy'),
						},
						{
							label: 'Savage Combat',
							actionId: ActionId.fromSpellId(58413),
							playerData: playerClassAndTalent(Class.ClassRogue, 'savageCombat'),
						},
					],
				},
				{
					label: 'Bleed',
					effects: [
						{
							label: 'Mangle',
							actionId: ActionId.fromSpellId(16862),
							playerData: playerClass(Class.ClassDruid, player => [Spec.SpecFeralDruid, Spec.SpecFeralTankDruid].includes(player.spec)),
						},
						{
							label: 'Trauma',
							actionId: ActionId.fromSpellId(46855),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'trauma'),
						},
						{
							label: 'Stampede',
							actionId: ActionId.fromSpellId(57393),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Rhino),
						},
					],
				},
				{
					label: 'Crit',
					effects: [
						{
							label: 'Totem of Wrath',
							actionId: ActionId.fromSpellId(30706),
							playerData: playerClassAndTalent(Class.ClassShaman, 'totemOfWrath', player => player.getSpecOptions().totems?.fire == FireTotem.TotemOfWrath),
						},
						{
							label: 'Heart of the Crusader',
							actionId: ActionId.fromSpellId(20337),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'heartOfTheCrusader', player => [Spec.SpecRetributionPaladin, Spec.SpecProtectionPaladin].includes(player.spec)),
						},
						{
							label: 'Master Poisoner',
							actionId: ActionId.fromSpellId(58410),
							playerData: playerClassAndTalent(Class.ClassRogue, 'masterPoisoner'),
						},
					],
				},
				{
					label: 'Spell Crit',
					effects: [
						{
							label: 'Improved Shadow Bolt',
							actionId: ActionId.fromSpellId(17803),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedShadowBolt'),
						},
						{
							label: 'Improved Scorch',
							actionId: ActionId.fromSpellId(12873),
							playerData: playerClassAndTalent(Class.ClassMage, 'improvedScorch'),
						},
						{
							label: 'Winter\'s Chill',
							actionId: ActionId.fromSpellId(28593),
							playerData: playerClassAndTalent(Class.ClassMage, 'wintersChill'),
						},
					],
				},
				{
					label: 'Spell Hit',
					effects: [
						{
							label: 'Misery',
							actionId: ActionId.fromSpellId(33198),
							playerData: playerSpecAndTalent(Spec.SpecShadowPriest, 'misery'),
						},
						{
							label: 'Improved Faerie Fire',
							actionId: ActionId.fromSpellId(33602),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'improvedFaerieFire'),
						},
					],
				},
				{
					label: 'Spell Dmg',
					effects: [
						{
							label: 'Ebon Plaguebringer',
							actionId: ActionId.fromSpellId(51161),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'ebonPlaguebringer'),
						},
						{
							label: 'Earth and Moon',
							actionId: ActionId.fromSpellId(48511),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'earthAndMoon'),
						},
						{
							label: 'Curse of Elements',
							actionId: ActionId.fromSpellId(47865),
							playerData: playerClass(Class.ClassWarlock),
						},
					],
				},
			],
		},
		{
			label: 'Mitigation Debuffs',
			categories: [
				{
					label: 'Atk Pwr',
					effects: [
						{
							label: 'Vindication',
							actionId: ActionId.fromSpellId(26016),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'vindication', player => [Spec.SpecRetributionPaladin, Spec.SpecProtectionPaladin].includes(player.spec)),
						},
						{
							label: 'Improved Demoralizing Shout',
							actionId: ActionId.fromSpellId(12879),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'improvedDemoralizingShout'),
						},
						{
							label: 'Demoralizing Shout',
							actionId: ActionId.fromSpellId(47437),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'improvedDemoralizingShout'),
						},
						{
							label: 'Improved Demoralizing Roar',
							actionId: ActionId.fromSpellId(16862),
							playerData: playerSpecAndTalent(Spec.SpecFeralTankDruid, 'feralAggression'),
						},
						{
							label: 'Demoralizing Roar',
							actionId: ActionId.fromSpellId(48560),
							playerData: playerSpecAndMissingTalent(Spec.SpecFeralTankDruid, 'feralAggression'),
						},
						{
							label: 'Improved Curse of Weakness',
							actionId: ActionId.fromSpellId(18180),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedCurseOfWeakness'),
						},
						{
							label: 'Curse of Weakness',
							actionId: ActionId.fromSpellId(50511),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedCurseOfWeakness'),
						},
						{
							label: 'Demoralizing Screech',
							actionId: ActionId.fromSpellId(55487),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.CarrionBird),
						},
					],
				},
				{
					label: 'Atk Speed',
					effects: [
						{
							label: 'Improved Thunder Clap',
							actionId: ActionId.fromSpellId(12666),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'improvedThunderClap'),
						},
						{
							label: 'Thunder Clap',
							actionId: ActionId.fromSpellId(47502),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'improvedThunderClap'),
						},
						{
							label: 'Improved Frost Fever',
							actionId: ActionId.fromSpellId(51456),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'improvedIcyTouch'),
						},
						{
							label: 'Frost Fever',
							actionId: ActionId.fromSpellId(51456),
							playerData: playerClassAndMissingTalent(Class.ClassDeathknight, 'improvedIcyTouch'),
						},
						{
							label: 'Judgements of the Just',
							actionId: ActionId.fromSpellId(53696),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'judgementsOfTheJust'),
						},
						{
							label: 'Infected Wounds',
							actionId: ActionId.fromSpellId(48485),
							playerData: playerClassAndTalent(Class.ClassDruid, 'infectedWounds', player => [Spec.SpecFeralDruid, Spec.SpecFeralTankDruid].includes(player.spec)),
						},
					],
				},
				{
					label: 'Miss',
					effects: [
						{
							label: 'Insect Swarm',
							actionId: ActionId.fromSpellId(65855),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'insectSwarm'),
						},
						{
							label: 'Scorpid Sting',
							actionId: ActionId.fromSpellId(3043),
							playerData: playerClass(Class.ClassHunter),
						},
					],
				},
			],
		},
	]
};
