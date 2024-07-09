import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player';
import { Spec , UnitReference, UnitReference_Type as UnitType } from '../core/proto/common.js';
import {
	DeathknightMajorGlyph,
} from '../core/proto/deathknight.js';
import { EventID, TypedEvent } from '../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfUnholyFrenzy = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'unholyFrenzyTarget',
	label: '自己使用邪恶狂热',
	labelTooltip: '对自己施放邪恶狂热。',
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecDeathknight>) => player.getSpecOptions().unholyFrenzyTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecDeathknight>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.unholyFrenzyTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().hysteria,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const StartingRunicPower = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathknight>({
	fieldName: 'startingRunicPower',
	label: '初始符文能量',
	labelTooltip: '每次迭代开始时的初始符文能量。',
});

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathknight>({
	fieldName: 'petUptime',
	label: '食尸鬼存活时间(%)',
	labelTooltip: '战斗期间你的食尸鬼将处于目标上的时间百分比。',
	percent: true,
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().masterOfGhouls,
});

export const DrwPestiApply = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'drwPestiApply',
	label: '使用符文刃舞瘟疫',
	labelTooltip: '目前有一个与符文刃舞和瘟疫有关的互动，你可以使用瘟疫来强制符文刃舞应用疾病。它仅适用于疾病雕文，并且必须有一个额外目标。',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalentTree() == 0 && (player.getGlyphs().major1 == DeathknightMajorGlyph.GlyphOfDisease || player.getGlyphs().major2 == DeathknightMajorGlyph.GlyphOfDisease || player.getGlyphs().major3 == DeathknightMajorGlyph.GlyphOfDisease),
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});

export const UseAMSInput = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecDeathknight>({
	fieldName: 'useAms',
	label: '使用反魔法护罩',
	labelTooltip: '在预测伤害时使用反魔法护罩以获得符文能量增益。',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});

export const AvgAMSSuccessRateInput = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathknight>({
	fieldName: 'avgAmsSuccessRate',
	label: '反魔法护罩成功率(%)',
	labelTooltip: '在反魔法护罩的5秒窗口期内受到伤害的几率。',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getSpecOptions().useAms == true && player.getTalents().howlingBlast,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});

export const AvgAMSHitInput = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecDeathknight>({
	fieldName: 'avgAmsHit',
	label: '平均反魔法护罩被击中',
	labelTooltip: '当AMS成功时，角色平均受到的伤害(+-10%)。',
	showWhen: (player: Player<Spec.SpecDeathknight>) => player.getSpecOptions().useAms == true && player.getTalents().howlingBlast,
	changeEmitter: (player: Player<Spec.SpecDeathknight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});
