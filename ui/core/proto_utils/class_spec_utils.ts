import { Class, Spec } from '../proto/common.js';
import { classNames, titleIcons } from './utils.js'

// Class order is intentional and based off of the in-game order on the character creation screen.
// It also mirrors the ordering on the homepage.
export const classSpecMap: Map<Class, Array<Spec>> = new Map();

classSpecMap.set(Class.ClassUnknown, []);

classSpecMap.set(Class.ClassWarrior, [
	Spec.SpecWarrior,
	Spec.SpecProtectionWarrior
]);

classSpecMap.set(Class.ClassPaladin, [
	Spec.SpecProtectionPaladin,
	Spec.SpecRetributionPaladin,
]);

classSpecMap.set(Class.ClassHunter, [
	Spec.SpecHunter,
]);

classSpecMap.set(Class.ClassRogue, [
	Spec.SpecRogue,
]);

classSpecMap.set(Class.ClassPriest, [
	Spec.SpecHealingPriest,
	Spec.SpecShadowPriest,
	Spec.SpecSmitePriest,
]);

classSpecMap.set(Class.ClassDeathknight, [
	Spec.SpecDeathknight,
	Spec.SpecTankDeathknight,
]);

classSpecMap.set(Class.ClassShaman, [
	Spec.SpecElementalShaman,
	Spec.SpecEnhancementShaman,
]);

classSpecMap.set(Class.ClassMage, [
	Spec.SpecMage,
]);

classSpecMap.set(Class.ClassWarlock, [
	Spec.SpecWarlock,
]);

classSpecMap.set(Class.ClassDruid, [
	Spec.SpecBalanceDruid,
	Spec.SpecFeralDruid,
	Spec.SpecFeralTankDruid,
]);

export const classList: Class[] = Array.from(classSpecMap.keys()).filter( (klass) => klass != Class.ClassUnknown );

export function getClassIcon(classIndex: Class): string {
	let className = classNames[classIndex];

	return `/wotlk/assets/img/${className.toLowerCase().replace(/\s/g, '_')}_icon.png`
}

export function getSpecIcon(specIndex: Spec): string {
	return titleIcons[specIndex];
}
