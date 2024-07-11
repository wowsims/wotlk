import { Tooltip } from 'bootstrap';

import { Component } from '../core/components/component.js';
import { Player } from "../core/player.js";
import {
	Class,
	RaidBuffs,
	Spec,
} from '../core/proto/common.js';
import { Hunter_Options_PetType as HunterPetType } from '../core/proto/hunter.js';
import { PaladinAura } from '../core/proto/paladin.js';
import { AirTotem, EarthTotem, FireTotem, WaterTotem } from '../core/proto/shaman.js';
import { Warlock_Options_Summon as WarlockSummon } from '../core/proto/warlock.js';
import { WarriorShout } from '../core/proto/warrior.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import {
	ClassSpecs,
	isHealingSpec,
	isMeleeDpsSpec,
	isRangedDpsSpec,
	isTankSpec,
	SpecTalents,
	specToClass,
	textCssClassForClass,
} from '../core/proto_utils/utils.js';
import { Raid } from "../core/raid.js";
import { sum } from '../core/utils.js';
import { RaidSimUI } from './raid_sim_ui.js';

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

		const categories: Array<RaidStatsCategory> = [];
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
			label: '职责',
			categories: [
				{
					label: '坦克',
					effects: [
						{
							label: '坦克',
							playerData: { condition: player => isTankSpec(player.spec) },
						},
					],
				},
				{
					label: '治疗',
					effects: [
						{
							label: '治疗',
							playerData: { condition: player => isHealingSpec(player.spec) },
						},
					],
				},
				{
					label: '近战',
					effects: [
						{
							label: '近战',
							playerData: { condition: player => isMeleeDpsSpec(player.spec) },
						},
					],
				},
				{
					label: '远程',
					effects: [
						{
							label: '远程',
							playerData: { condition: player => isRangedDpsSpec(player.spec) },
						},
					],
				},
			],
		},
		{
			label: '增益',
			categories: [
				{
					label: '嗜血',
					effects: [
						{
							label: '嗜血',
							actionId: ActionId.fromSpellId(2825),
							playerData: playerClass(Class.ClassShaman),
						},
					],
				},
				{
					label: '属性',
					effects: [
						{
							label: '强化野性赐福',
							actionId: ActionId.fromSpellId(17051),
							playerData: playerClassAndTalent(Class.ClassDruid, 'improvedMarkOfTheWild'),
						},
						{
							label: '野性赐福',
							actionId: ActionId.fromSpellId(48470),
							playerData: playerClassAndMissingTalent(Class.ClassDruid, 'improvedMarkOfTheWild'),
						},
						{
							label: '狂野战鼓',
							actionId: ActionId.fromItemId(49634),
							raidData: raidBuff('drumsOfTheWild'),
						},
					],
				},
				{
					label: '属性 %',
					effects: [
						{
							label: '王者祝福',
							actionId: ActionId.fromSpellId(25898),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: '遗忘王者战鼓',
							actionId: ActionId.fromItemId(49633),
							raidData: raidBuff('drumsOfForgottenKings'),
						},
						{
							label: '庇护祝福',
							actionId: ActionId.fromSpellId(25899),
							playerData: playerClass(Class.ClassPaladin),
						},
					],
				},
				{
					label: '护甲',
					effects: [
						{
							label: '强化虔诚光环',
							actionId: ActionId.fromSpellId(20140),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'improvedDevotionAura', player => player.getSpecOptions().aura == PaladinAura.DevotionAura),
						},
						{
							label: '虔诚光环',
							actionId: ActionId.fromSpellId(48942),
							playerData: playerClassAndMissingTalent(Class.ClassPaladin, 'improvedDevotionAura', player => player.getSpecOptions().aura == PaladinAura.DevotionAura),
						},
						{
							label: '强化石肤图腾',
							actionId: ActionId.fromSpellId(16293),
							playerData: playerClassAndTalent(Class.ClassShaman, 'guardianTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StoneskinTotem),
						},
						{
							label: '石肤图腾',
							actionId: ActionId.fromSpellId(58753),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'guardianTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StoneskinTotem),
						},
						{
							label: '保护卷轴',
							actionId: ActionId.fromItemId(43468),
							raidData: raidBuff('scrollOfProtection'),
						},
					],
				},
				{
					label: '耐力',
					effects: [
						{
							label: '强化真言术：韧',
							actionId: ActionId.fromSpellId(14767),
							playerData: playerClassAndTalent(Class.ClassPriest, 'improvedPowerWordFortitude'),
						},
						{
							label: '真言术：韧',
							actionId: ActionId.fromSpellId(48161),
							playerData: playerClassAndMissingTalent(Class.ClassPriest, 'improvedPowerWordFortitude'),
						},
						{
							label: '耐力卷轴',
							actionId: ActionId.fromItemId(37094),
							raidData: raidBuff('scrollOfStamina'),
						},
					],
				},
				{
					label: '力量/敏捷',
					effects: [
						{
							label: '强化大地之力图腾',
							actionId: ActionId.fromSpellId(52456),
							playerData: playerClassAndTalent(Class.ClassShaman, 'enhancingTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StrengthOfEarthTotem),
						},
						{
							label: '大地之力图腾',
							actionId: ActionId.fromSpellId(58643),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'enhancingTotems', player => player.getSpecOptions().totems?.earth == EarthTotem.StrengthOfEarthTotem),
						},
						{
							label: '凛冬号角',
							actionId: ActionId.fromSpellId(57623),
							playerData: playerClass(Class.ClassDeathknight),
						},
						{
							label: '力量卷轴',
							actionId: ActionId.fromItemId(43466),
							raidData: raidBuff('scrollOfStrength'),
						},
						{
							label: '敏捷卷轴',
							actionId: ActionId.fromItemId(43464),
							raidData: raidBuff('scrollOfAgility'),
						},
					],
				},
				{
					label: '智力',
					effects: [
						{
							label: '奥术光辉',
							actionId: ActionId.fromSpellId(43002),
							playerData: playerClass(Class.ClassMage),
						},
						{
							label: '强化地狱猎犬',
							actionId: ActionId.fromSpellId(54038),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: '邪能智力',
							actionId: ActionId.fromSpellId(57567),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: '智力卷轴',
							actionId: ActionId.fromItemId(37092),
							raidData: raidBuff('scrollOfIntellect'),
						},
					],
				},
				{
					label: '精神',
					effects: [
						{
							label: '神圣之灵',
							actionId: ActionId.fromSpellId(48073),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: '强化地狱猎犬',
							actionId: ActionId.fromSpellId(54038),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: '邪能智力',
							actionId: ActionId.fromSpellId(57567),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedFelhunter', player => player.getSpecOptions().summon == WarlockSummon.Felhunter),
						},
						{
							label: '精神卷轴',
							actionId: ActionId.fromItemId(37098),
							raidData: raidBuff('scrollOfSpirit'),
						},
					],
				},
				{
					label: '攻击强度',
					effects: [
						{
							label: '强化力量祝福',
							actionId: ActionId.fromSpellId(20045),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: '力量祝福',
							actionId: ActionId.fromSpellId(48934),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: '强化战斗怒吼',
							actionId: ActionId.fromSpellId(12861),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle),
						},
						{
							label: '战斗怒吼',
							actionId: ActionId.fromSpellId(47436),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutBattle),
						},
					],
				},
				{
					label: '攻击强度 %',
					effects: [
						{
							label: '憎恶之力',
							actionId: ActionId.fromSpellId(53138),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'abominationsMight'),
						},
						{
							label: '狂暴',
							actionId: ActionId.fromSpellId(30809),
							playerData: playerClassAndTalent(Class.ClassShaman, 'unleashedRage'),
						},
						{
							label: '强击光环',
							actionId: ActionId.fromSpellId(19506),
							playerData: playerClassAndTalent(Class.ClassHunter, 'trueshotAura'),
						},
					],
				},
				{
					label: '伤害 %',
					effects: [
						{
							label: '圣洁惩戒',
							actionId: ActionId.fromSpellId(31869),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'sanctifiedRetribution'),
						},
						{
							label: '奥术增效',
							actionId: ActionId.fromSpellId(31583),
							playerData: playerClassAndTalent(Class.ClassMage, 'arcaneEmpowerment'),
						},
						{
							label: '凶猛灵感',
							actionId: ActionId.fromSpellId(34460),
							playerData: playerClassAndTalent(Class.ClassHunter, 'ferociousInspiration'),
						},
					],
				},
				{
					label: '减伤 %',
					effects: [
						{
							label: '新生希望',
							actionId: ActionId.fromSpellId(57472),
							playerData: playerClassAndTalent(Class.ClassPriest, 'renewedHope'),
						},
						{
							label: '庇护祝福',
							actionId: ActionId.fromSpellId(25899),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: '警戒',
							actionId: ActionId.fromSpellId(50720),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'vigilance'),
						},
					],
				},
				{
					label: '急速 %',
					effects: [
						{
							label: '迅捷惩戒',
							actionId: ActionId.fromSpellId(53648),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'swiftRetribution'),
						},
						{
							label: '强化枭兽形态',
							actionId: ActionId.fromSpellId(48396),
							playerData: playerClassAndTalent(Class.ClassDruid, 'improvedMoonkinForm'),
						},
					],
				},
				{
					label: '每5秒回复法力值',
					effects: [
						{
							label: '强化智慧祝福',
							actionId: ActionId.fromSpellId(20245),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'improvedBlessingOfWisdom'),
						},
						{
							label: '智慧祝福',
							actionId: ActionId.fromSpellId(48938),
							playerData: playerClassAndMissingTalent(Class.ClassPaladin, 'improvedBlessingOfWisdom'),
						},
						{
							label: '强化法力之泉图腾',
							actionId: ActionId.fromSpellId(16206),
							playerData: playerClassAndTalent(Class.ClassShaman, 'restorativeTotems', player => player.getSpecOptions().totems?.water == WaterTotem.ManaSpringTotem),
						},
						{
							label: '法力之泉图腾',
							actionId: ActionId.fromSpellId(58774),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'restorativeTotems', player => player.getSpecOptions().totems?.water == WaterTotem.ManaSpringTotem),
						},
					],
				},
				{
					label: '近战暴击',
					effects: [
						{
							label: '兽群领袖',
							actionId: ActionId.fromSpellId(17007),
							playerData: playerClassAndTalent(Class.ClassDruid, 'leaderOfThePack'),
						},
						{
							label: '暴怒',
							actionId: ActionId.fromSpellId(29801),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'rampage'),
						},
					],
				},
				{
					label: '近战急速',
					effects: [
						{
							label: '强化冰冷之爪',
							actionId: ActionId.fromSpellId(55610),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'improvedIcyTalons'),
						},
						{
							label: '强化风怒图腾',
							actionId: ActionId.fromSpellId(29193),
							playerData: playerClassAndTalent(Class.ClassShaman, 'improvedWindfuryTotem', player => player.getSpecOptions().totems?.air == AirTotem.WindfuryTotem),
						},
						{
							label: '风怒图腾',
							actionId: ActionId.fromSpellId(65990),
							playerData: playerClassAndMissingTalent(Class.ClassShaman, 'improvedWindfuryTotem', player => player.getSpecOptions().totems?.air == AirTotem.WindfuryTotem),
						},
					],
				},
				{
					label: '法术强度',
					effects: [
						{
							label: '恶魔契约',
							actionId: ActionId.fromSpellId(47240),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'demonicPact'),
						},
						{
							label: '天怒图腾',
							actionId: ActionId.fromSpellId(57722),
							playerData: playerClassAndTalent(Class.ClassShaman, 'totemOfWrath', player => player.getSpecOptions().totems?.fire == FireTotem.TotemOfWrath),
						},
						{
							label: '火舌图腾',
							actionId: ActionId.fromSpellId(58656),
							playerData: playerClass(Class.ClassShaman, player => player.getSpecOptions().totems?.fire == FireTotem.FlametongueTotem),
						},
					],
				},
				{
					label: '法术暴击',
					effects: [
						{
							label: '枭兽光环',
							actionId: ActionId.fromSpellId(24907),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'moonkinForm'),
						},
						{
							label: '元素之誓',
							actionId: ActionId.fromSpellId(51470),
							playerData: playerClassAndTalent(Class.ClassShaman, 'elementalOath'),
						},
					],
				},
				{
					label: '法术急速',
					effects: [
						{
							label: '空气之怒图腾',
							actionId: ActionId.fromSpellId(3738),
							playerData: playerClass(Class.ClassShaman, player => player.getSpecOptions().totems?.air == AirTotem.WrathOfAirTotem),
						},
					],
				},
				{
					label: '生命值',
					effects: [
						{
							label: '强化命令怒吼',
							actionId: ActionId.fromSpellId(12861),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutCommanding),
						},
						{
							label: '命令怒吼',
							actionId: ActionId.fromSpellId(47440),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'commandingPresence', player => player.getSpecOptions().shout == WarriorShout.WarriorShoutCommanding),
						},
						{
							label: '强化小鬼',
							actionId: ActionId.fromSpellId(18696),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedImp', player => player.getSpecOptions().summon == WarlockSummon.Imp),
						},
						{
							label: '血之契印',
							actionId: ActionId.fromSpellId(47982),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedImp', player => player.getSpecOptions().summon == WarlockSummon.Imp),
						},
					],
				},
				{
					label: '恢复',
					effects: [
						{
							label: '吸血鬼之触',
							actionId: ActionId.fromSpellId(48160),
							playerData: playerSpecAndTalent(Spec.SpecShadowPriest, 'vampiricTouch'),
						},
						{
							label: '智者审判',
							actionId: ActionId.fromSpellId(31878),
							playerData: playerSpecAndTalent(Spec.SpecRetributionPaladin, 'judgementsOfTheWise'),
						},
						{
							label: '狩猎小队',
							actionId: ActionId.fromSpellId(53292),
							playerData: playerSpecAndTalent(Spec.SpecHunter, 'huntingParty'),
						},
						{
							label: '强化灵魂榨取',
							actionId: ActionId.fromSpellId(54118),
							playerData: playerSpecAndTalent(Spec.SpecWarlock, 'improvedSoulLeech'),
						},
						{
							label: '漫长寒冬',
							actionId: ActionId.fromSpellId(44561),
							playerData: playerSpecAndTalent(Spec.SpecMage, 'enduringWinter'),
						},
					],
				},
			],
		},
		{
			label: '外部增益',
			categories: [
				{
					label: '激活',
					effects: [
						{
							label: '激活',
							actionId: ActionId.fromSpellId(29166),
							playerData: playerClass(Class.ClassDruid),
						},
					],
				},
				{
					label: '能量灌注',
					effects: [
						{
							label: '能量灌注',
							actionId: ActionId.fromSpellId(10060),
							playerData: playerClassAndTalent(Class.ClassPriest, 'powerInfusion'),
						},
					],
				},
				{
					label: '专注魔法',
					effects: [
						{
							label: '专注魔法',
							actionId: ActionId.fromSpellId(54648),
							playerData: playerClassAndTalent(Class.ClassMage, 'focusMagic'),
						},
					],
				},
				{
					label: '嫁祸诀窍',
					effects: [
						{
							label: '嫁祸诀窍',
							actionId: ActionId.fromSpellId(57933),
							playerData: playerClass(Class.ClassRogue),
						},
					],
				},
				{
					label: '邪恶狂热',
					effects: [
						{
							label: '邪恶狂热',
							actionId: ActionId.fromSpellId(49016),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'hysteria'),
						},
					],
				},
				{
					label: '痛苦压制',
					effects: [
						{
							label: '痛苦压制',
							actionId: ActionId.fromSpellId(33206),
							playerData: playerClassAndTalent(Class.ClassPriest, 'painSuppression'),
						},
					],
				},
				{
					label: '神圣守护者',
					effects: [
						{
							label: '神圣守护者',
							actionId: ActionId.fromSpellId(53530),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'divineGuardian'),
						},
					],
				},
			],
		},
		{
			label: 'DPS减益',
			categories: [
				{
					label: '主要破甲',
					effects: [
						{
							label: '破甲攻击',
							actionId: ActionId.fromSpellId(47467),
							playerData: playerClass(Class.ClassWarrior),
						},
						{
							label: '破甲',
							actionId: ActionId.fromSpellId(8647),
							playerData: playerClass(Class.ClassRogue),
						},
						{
							label: '酸液喷吐',
							actionId: ActionId.fromSpellId(55754),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Worm),
						},
					],
				},
				{
					label: '次要破甲',
					effects: [
						{
							label: '精灵之火',
							actionId: ActionId.fromSpellId(770),
							playerData: playerClass(Class.ClassDruid, player => player.spec != Spec.SpecRestorationDruid),
						},
						{
							label: '虚弱诅咒',
							actionId: ActionId.fromSpellId(50511),
							playerData: playerClass(Class.ClassWarlock),
						},
						{
							label: '蜂刺',
							actionId: ActionId.fromSpellId(56631),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Wasp),
						},
						{
							label: '孢子云雾',
							actionId: ActionId.fromSpellId(53598),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Bat),
						},
					],
				},
				{
					label: '物理易伤',
					effects: [
						{
							label: '血性狂乱',
							actionId: ActionId.fromSpellId(29859),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'bloodFrenzy'),
						},
						{
							label: '野蛮战斗',
							actionId: ActionId.fromSpellId(58413),
							playerData: playerClassAndTalent(Class.ClassRogue, 'savageCombat'),
						},
					],
				},
				{
					label: '流血',
					effects: [
						{
							label: '野性侵略',
							actionId: ActionId.fromSpellId(16862),
							playerData: playerClass(Class.ClassDruid, player => [Spec.SpecFeralDruid, Spec.SpecFeralTankDruid].includes(player.spec)),
						},
						{
							label: '创伤',
							actionId: ActionId.fromSpellId(46855),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'trauma'),
						},
						{
							label: '奔踏',
							actionId: ActionId.fromSpellId(57393),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.Rhino),
						},
					],
				},
				{
					label: '暴击',
					effects: [
						{
							label: '天怒图腾',
							actionId: ActionId.fromSpellId(30706),
							playerData: playerClassAndTalent(Class.ClassShaman, 'totemOfWrath', player => player.getSpecOptions().totems?.fire == FireTotem.TotemOfWrath),
						},
						{
							label: '十字军之心',
							actionId: ActionId.fromSpellId(20337),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'heartOfTheCrusader', player => [Spec.SpecRetributionPaladin, Spec.SpecProtectionPaladin].includes(player.spec)),
						},
						{
							label: '奇毒',
							actionId: ActionId.fromSpellId(58410),
							playerData: playerClassAndTalent(Class.ClassRogue, 'masterPoisoner'),
						},
					],
				},
				{
					label: '法术暴击',
					effects: [
						{
							label: '强化暗影箭',
							actionId: ActionId.fromSpellId(17803),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedShadowBolt'),
						},
						{
							label: '强化灼烧',
							actionId: ActionId.fromSpellId(12873),
							playerData: playerClassAndTalent(Class.ClassMage, 'improvedScorch'),
						},
						{
							label: '深冬之寒',
							actionId: ActionId.fromSpellId(28593),
							playerData: playerClassAndTalent(Class.ClassMage, 'wintersChill'),
						},
					],
				},
				{
					label: '法术命中',
					effects: [
						{
							label: '悲惨',
							actionId: ActionId.fromSpellId(33198),
							playerData: playerSpecAndTalent(Spec.SpecShadowPriest, 'misery'),
						},
						{
							label: '强化精灵之火',
							actionId: ActionId.fromSpellId(33602),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'improvedFaerieFire'),
						},
					],
				},
				{
					label: '法术伤害',
					effects: [
						{
							label: '黑色热疫使者',
							actionId: ActionId.fromSpellId(51161),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'ebonPlaguebringer'),
						},
						{
							label: '大地与月亮',
							actionId: ActionId.fromSpellId(48511),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'earthAndMoon'),
						},
						{
							label: '元素诅咒',
							actionId: ActionId.fromSpellId(47865),
							playerData: playerClass(Class.ClassWarlock),
						},
					],
				},
			],
		},
		{
			label: '减伤减益',
			categories: [
				{
					label: '攻击强度',
					effects: [
						 {
							label: '辩护',
							actionId: ActionId.fromSpellId(26016),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'vindication', player => [Spec.SpecRetributionPaladin, Spec.SpecProtectionPaladin].includes(player.spec)),
						},
						{
							label: '强化挫志怒吼',
							actionId: ActionId.fromSpellId(12879),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'improvedDemoralizingShout'),
						},
						{
							label: '挫志怒吼',
							actionId: ActionId.fromSpellId(47437),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'improvedDemoralizingShout'),
						},
						{
							label: '强化挫志咆哮',
							actionId: ActionId.fromSpellId(16862),
							playerData: playerSpecAndTalent(Spec.SpecFeralTankDruid, 'feralAggression'),
						},
						{
							label: '挫志咆哮',
							actionId: ActionId.fromSpellId(48560),
							playerData: playerSpecAndMissingTalent(Spec.SpecFeralTankDruid, 'feralAggression'),
						},
						{
							label: '强化虚弱诅咒',
							actionId: ActionId.fromSpellId(18180),
							playerData: playerClassAndTalent(Class.ClassWarlock, 'improvedCurseOfWeakness'),
						},
						{
							label: '虚弱诅咒',
							actionId: ActionId.fromSpellId(50511),
							playerData: playerClassAndMissingTalent(Class.ClassWarlock, 'improvedCurseOfWeakness'),
						},
						{
							label: '挫志尖啸',
							actionId: ActionId.fromSpellId(55487),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().petType == HunterPetType.CarrionBird),
						},
					],
				},
				{
					label: '攻击速度',
					effects: [
						{
							label: '强化雷霆一击',
							actionId: ActionId.fromSpellId(12666),
							playerData: playerClassAndTalent(Class.ClassWarrior, 'improvedThunderClap'),
						},
						{
							label: '雷霆一击',
							actionId: ActionId.fromSpellId(47502),
							playerData: playerClassAndMissingTalent(Class.ClassWarrior, 'improvedThunderClap'),
						},
						{
							label: '强化冰冷触摸',
							actionId: ActionId.fromSpellId(51456),
							playerData: playerClassAndTalent(Class.ClassDeathknight, 'improvedIcyTouch'),
						},
						{
							label: '冰冷触摸',
							actionId: ActionId.fromSpellId(51456),
							playerData: playerClassAndMissingTalent(Class.ClassDeathknight, 'improvedIcyTouch'),
						},
						{
							label: '正义审判',
							actionId: ActionId.fromSpellId(53696),
							playerData: playerClassAndTalent(Class.ClassPaladin, 'judgementsOfTheJust'),
						},
						{
							label: '感染伤口',
							actionId: ActionId.fromSpellId(48485),
							playerData: playerClassAndTalent(Class.ClassDruid, 'infectedWounds', player => [Spec.SpecFeralDruid, Spec.SpecFeralTankDruid].includes(player.spec)),
						},
					],
				},
				{
					label: '未命中',
					effects: [
						{
							label: '虫群',
							actionId: ActionId.fromSpellId(65855),
							playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'insectSwarm'),
						},
						{
							label: '毒蝎钉刺',
							actionId: ActionId.fromSpellId(3043),
							playerData: playerClass(Class.ClassHunter),
						},
					],
				},
			],
		},
	]
};

