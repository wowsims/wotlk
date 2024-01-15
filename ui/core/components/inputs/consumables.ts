// import { Potions } from "../../proto/common";
// import { ActionId } from "../../proto_utils/action_id";

// export const LesserManaPotion = { actionId: ActionId.fromItemId(3385), value: Potions.LesserManaPotion }
// export const ManaPotion = { actionId: ActionId.fromItemId(3385), value: Potions.ManaPotion }

// function makePotionInputFactory(consumesFieldName: keyof Consumes): (options: Array<Potions>, tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, Potions> {
// 	return makeConsumeInputFactory({
// 		consumesFieldName: consumesFieldName,
// 		allOptions: [
// 			{ actionId: ActionId.fromItemId(3385), value: Potions.LesserManaPotion },
// 			{ actionId: ActionId.fromItemId(3827), value: Potions.ManaPotion },
// 		] as Array<IconEnumValueConfig<Player<any>, Potions>>,
// 	});
// }
// export const makePotionsInput = makePotionInputFactory('defaultPotion');

// // TODO: Classic? 
// export const makeConjuredInput = makeConsumeInputFactory({
// 	consumesFieldName: 'defaultConjured',
// 	allOptions: [
// 		{ actionId: ActionId.fromItemId(4381), value: Conjured.ConjuredMinorRecombobulator, showWhen: (player: Player<any>) => player.getGear().hasTrinket(4381) },
// 		{ actionId: ActionId.fromItemId(12662), value: Conjured.ConjuredDemonicRune, showWhen: (p) => p.getLevel() >= 40 },
// 	] as Array<IconEnumValueConfig<Player<any>, Conjured>>
// });

// export const makeFlasksInput = makeConsumeInputFactory({
// 	consumesFieldName: 'flask',
// 	allOptions: [
// 		{ actionId: ActionId.fromItemId(13510), value: Flask.FlaskOfTheTitans },
// 		{ actionId: ActionId.fromItemId(13511), value: Flask.FlaskOfDistilledWisdom },
// 		{ actionId: ActionId.fromItemId(13512), value: Flask.FlaskOfSupremePower },
// 		{ actionId: ActionId.fromItemId(13513), value: Flask.FlaskOfChromaticResistance },
// 	] as Array<IconEnumValueConfig<Player<any>, Flask>>,
// });

// export const makeMainHandImbuesInput = makeConsumeInputFactory({
// 	consumesFieldName: 'mainHandImbue',
// 	allOptions: [
// 		// TODO: Classic hide when required level too high e.g. `showWhen: (p) =>  p.getLevel() >= 25`
// 		// Registering a `showWhen` is causing issues with event callback loops
// 		{ actionId: ActionId.fromItemId(20749), value: WeaponImbue.BrillianWizardOil},
// 		{ actionId: ActionId.fromItemId(20748), value: WeaponImbue.BrilliantManaOil},
// 		{ actionId: ActionId.fromItemId(12404), value: WeaponImbue.DenseSharpeningStone },
// 		{ actionId: ActionId.fromItemId(18262), value: WeaponImbue.ElementalSharpeningStone },
// 		{ actionId: ActionId.fromItemId(211848), value: WeaponImbue.BlackfathomManaOil},
// 		{ actionId: ActionId.fromItemId(211845), value: WeaponImbue.BlackfathomSharpeningStone},
// 		{ actionId: ActionId.fromSpellId(407975), value: WeaponImbue.WildStrikes},
// 	] as Array<IconEnumValueConfig<Player<any>, WeaponImbue>>,
// });

// export const makeOffHandImbuesInput = makeConsumeInputFactory({
// 	consumesFieldName: 'offHandImbue',
// 	allOptions: [
// 		// TODO: Classic hide when required level too high e.g. `showWhen: (p) =>  p.getLevel() >= 25`
// 		// Registering a `showWhen` is causing issues with event callback loops
// 		{ actionId: ActionId.fromItemId(20749), value: WeaponImbue.BrillianWizardOil},
// 		{ actionId: ActionId.fromItemId(20748), value: WeaponImbue.BrilliantManaOil},
// 		{ actionId: ActionId.fromItemId(12404), value: WeaponImbue.DenseSharpeningStone },
// 		{ actionId: ActionId.fromItemId(18262), value: WeaponImbue.ElementalSharpeningStone },
// 		{ actionId: ActionId.fromItemId(211848), value: WeaponImbue.BlackfathomManaOil},
// 		{ actionId: ActionId.fromItemId(211845), value: WeaponImbue.BlackfathomSharpeningStone},
		
// 	] as Array<IconEnumValueConfig<Player<any>, WeaponImbue>>,
// });

// export const makeFoodInput = makeConsumeInputFactory({
// 	consumesFieldName: 'food',
// 	allOptions: [
// 		{ actionId: ActionId.fromItemId(21072), value: Food.FoodSmokedSagefish, showWhen: (p) => p.getLevel() >= 10 },
// 		{ actionId: ActionId.fromItemId(13851), value: Food.FoodHotWolfRibs, showWhen: (p) => p.getLevel() >= 25 },
// 		{ actionId: ActionId.fromItemId(22480), value: Food.FoodTenderWolfSteak, showWhen: (p) => p.getLevel() >= 40 },
// 		{ actionId: ActionId.fromItemId(13931), value: Food.FoodNightfinSoup, showWhen: (p) => p.getLevel() >= 35 },
// 		{ actionId: ActionId.fromItemId(13928), value: Food.FoodGrilledSquid, showWhen: (p) => p.getLevel() >= 35 },
// 		{ actionId: ActionId.fromItemId(20452), value: Food.FoodSmokedDesertDumpling, showWhen: (p) => p.getLevel() >= 45 },
// 		{ actionId: ActionId.fromItemId(18254), value: Food.FoodRunnTumTuberSurprise, showWhen: (p) => p.getLevel() >= 45 },
// 		{ actionId: ActionId.fromItemId(13813), value: Food.FoodBlessedSunfruitJuice, showWhen: (p) => p.getLevel() >= 45 },
// 		{ actionId: ActionId.fromItemId(13810), value: Food.FoodBlessSunfruit, showWhen: (p) => p.getLevel() >= 45 },
// 		{ actionId: ActionId.fromItemId(21023), value: Food.FoodDirgesKickChimaerokChops, showWhen: (p) => p.getLevel() >= 55 },
// 	] as Array<IconEnumValueConfig<Player<any>, Food>>
// });

// export const AgilityBuffInput = makeConsumeInput('agilityElixir', [
// 	{ actionId: ActionId.fromItemId(13452), value: AgilityElixir.ElixirOfTheMongoose, showWhen: (p) => p.getLevel() >= 46 },
// 	{ actionId: ActionId.fromItemId(9187), value: AgilityElixir.ElixirOfGreaterAgility, showWhen: (p) => p.getLevel() >= 38},
//         { actionId: ActionId.fromItemId(3390), value: AgilityElixir.ElixirOfLesserAgility, showWhen: (p) => p.getLevel() >= 18},
// 	{ actionId: ActionId.fromItemId(10309), value: AgilityElixir.ScrollOfAgility},
// ] as Array<IconEnumValueConfig<Player<any>, AgilityElixir>>, (p) => p.getLevel() >= 18);

// export const StrengthBuffInput = makeConsumeInput('strengthBuff', [
// 	{ actionId: ActionId.fromItemId(12451), value: StrengthBuff.JujuPower, showWhen: (p) => p.getLevel() >= 46 },
// 	{ actionId: ActionId.fromItemId(9206), value: StrengthBuff.ElixirOfGiants, showWhen: (p) => p.getLevel() >= 46 },
//         { actionId: ActionId.fromItemId(3391), value: StrengthBuff.ElixirOfOgresStrength, showWhen: (p) => p.getLevel() >= 20},
// 	{ actionId: ActionId.fromItemId(10310), value: StrengthBuff.ScrollOfStrength },
// ] as Array<IconEnumValueConfig<Player<any>, StrengthBuff>>, (p) => p.getLevel() >= 20);

// export const SpellDamageBuff = makeConsumeInput('spellPowerBuff', [
// 	{ actionId: ActionId.fromItemId(9155), value: SpellPowerBuff.ArcaneElixir, showWhen: (p) => p.getLevel() >= 37 },
// 	{ actionId: ActionId.fromItemId(13454), value: SpellPowerBuff.GreaterArcaneElixir, showWhen: (p) => p.getLevel() >= 46 },
// ] as Array<IconEnumValueConfig<Player<any>, SpellPowerBuff>>, (p) => p.getLevel() >= 37);

// export const FireDamageBuff = makeConsumeInput('firePowerBuff', [
// 	{ actionId: ActionId.fromItemId(6373), value: FirePowerBuff.ElixirOfFirepower, showWhen: (p) => p.getLevel() >= 18 },
// 	{ actionId: ActionId.fromItemId(21546), value: FirePowerBuff.ElixirOfGreaterFirepower, showWhen: (p) => p.getLevel() >= 40 },
// ] as Array<IconEnumValueConfig<Player<any>, FirePowerBuff>>, (p) => p.getLevel() >= 18);

// export const ShadowDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(9264), fieldName: 'shadowPowerBuff', minLevel: 40});
// export const FrostDamageBuff = makeBooleanConsumeInput({id: ActionId.fromItemId(17708), fieldName: 'frostPowerBuff', minLevel: 40});

// export const FillerExplosiveInput = makeConsumeInput('fillerExplosive', [
// 	{ actionId: ActionId.fromItemId(18641), value: Explosive.ExplosiveDenseDynamite, showWhen: (p) => p.getLevel() >= 40 },
// 	{ actionId: ActionId.fromItemId(15993), value: Explosive.ExplosiveThoriumGrenade, showWhen: (p) => p.getLevel() >= 40 },
// ] as Array<IconEnumValueConfig<Player<any>, Explosive>>);

// export interface ConsumeInputFactoryArgs<T extends number> {
// 	consumesFieldName: keyof Consumes,
// 	allOptions: Array<IconEnumValueConfig<Player<any>, T>>,
// 	// Additional callback if logic besides syncing consumes is required
// 	onSet?: (eventID: EventID, p: Player<any>, newValue: T) => void
// 	showWhen?: (obj: Player<any>) => boolean,
// }
// function makeConsumeInputFactory<T extends number>(args: ConsumeInputFactoryArgs<T>): (options: Array<T>, tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
// 	return (options: Array<T>, tooltip?: string) => {
// 		return {
// 			type: 'iconEnum',
// 			tooltip: tooltip,
// 			numColumns: options.length > 5 ? 2 : 1,
// 			values: [
// 				{ value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>,
// 			].concat(options.map(option => args.allOptions.find(allOption => allOption.value == option)!)),
// 			equals: (a: T, b: T) => a == b,
// 			zeroValue: 0 as T,
// 			showWhen: args.showWhen,
// 			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.levelChangeEmitter, player.gearChangeEmitter]),
// 			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
// 			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
// 				const newConsumes = player.getConsumes();

// 				if (newConsumes[args.consumesFieldName] === newValue){
// 					return;
// 				}

// 				(newConsumes[args.consumesFieldName] as number) = newValue;
// 				TypedEvent.freezeAllAndDo(() => {
// 					player.setConsumes(eventID, newConsumes);
// 					if (args.onSet) {
// 						args.onSet(eventID, player, newValue as T);
// 					}
// 				});
// 			},
// 		};
// 	};
// }

// function makeConsumeInput<T extends number>(consumesFieldName: keyof Consumes, allOptions: Array<IconEnumValueConfig<Player<any>, T>>, showWhen?: (obj: Player<any>) => boolean, onSet?: (eventID: EventID, p: Player<any>, newValue: T) => void): InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
// 	const factory = makeConsumeInputFactory({
// 		consumesFieldName: consumesFieldName,
// 		allOptions: allOptions,
// 		onSet: onSet,
// 		showWhen: showWhen,
// 	});
// 	return factory(allOptions.map(option => option.value));
// }
