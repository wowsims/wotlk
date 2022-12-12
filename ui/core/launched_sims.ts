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

export const raidSimLaunched = false;
export const raidSimStatus: LaunchStatus = LaunchStatus.Alpha;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Alpha,
	[Spec.SpecElementalShaman]: LaunchStatus.Alpha,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Alpha,
	[Spec.SpecFeralDruid]: LaunchStatus.Alpha,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecHunter]: LaunchStatus.Beta,
	[Spec.SpecMage]: LaunchStatus.Alpha,
	[Spec.SpecRogue]: LaunchStatus.Alpha,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Alpha,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Alpha,
	[Spec.SpecHealingPriest]: LaunchStatus.Alpha,
	[Spec.SpecShadowPriest]: LaunchStatus.Alpha,
	[Spec.SpecSmitePriest]: LaunchStatus.Alpha,
	[Spec.SpecWarlock]: LaunchStatus.Alpha,
	[Spec.SpecWarrior]: LaunchStatus.Alpha,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Alpha,
	[Spec.SpecDeathknight]: LaunchStatus.Alpha,
	[Spec.SpecTankDeathknight]: LaunchStatus.Unlaunched,
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
