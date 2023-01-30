import { Class, Spec } from './proto/common';
import { specToClass } from './proto_utils/utils';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
	Unlaunched,
	Alpha,
	Beta,
	Launched,
}

export const raidSimStatus: LaunchStatus = LaunchStatus.Alpha;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Beta,
	[Spec.SpecFeralDruid]: LaunchStatus.Beta,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Beta,
	[Spec.SpecRestorationDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecElementalShaman]: LaunchStatus.Beta,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Beta,
	[Spec.SpecRestorationShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecHunter]: LaunchStatus.Launched,
	[Spec.SpecMage]: LaunchStatus.Beta,
	[Spec.SpecRogue]: LaunchStatus.Beta,
	[Spec.SpecHolyPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Beta,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Beta,
	[Spec.SpecHealingPriest]: LaunchStatus.Alpha,
	[Spec.SpecShadowPriest]: LaunchStatus.Beta,
	[Spec.SpecSmitePriest]: LaunchStatus.Beta,
	[Spec.SpecWarlock]: LaunchStatus.Beta,
	[Spec.SpecWarrior]: LaunchStatus.Beta,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Beta,
	[Spec.SpecDeathknight]: LaunchStatus.Beta,
	[Spec.SpecTankDeathknight]: LaunchStatus.Alpha,
};

// Meme specs are excluded from title drop-down menu.
export const memeSpecs: Array<Spec> = [
	Spec.SpecSmitePriest,
];

export function getLaunchedSims(): Array<Spec> {
	return Object.keys(simLaunchStatuses)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => simLaunchStatuses[spec] > LaunchStatus.Unlaunched);
}

export function getLaunchedSimsForClass(klass: Class): Array<Spec> {
	return Object.keys(specToClass)
		.map(specStr => parseInt(specStr) as Spec)
		.filter(spec => specToClass[spec] == klass && isSimLaunched(spec));
}

export function isSimLaunched(specIndex: Spec): boolean {
	return simLaunchStatuses[specIndex] > LaunchStatus.Unlaunched;
}
