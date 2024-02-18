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

export const raidSimStatus: LaunchStatus = LaunchStatus.Launched;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Launched,
	[Spec.SpecFeralDruid]: LaunchStatus.Launched,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Launched,
	[Spec.SpecRestorationDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecElementalShaman]: LaunchStatus.Launched,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Launched,
	[Spec.SpecRestorationShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecHunter]: LaunchStatus.Launched,
	[Spec.SpecMage]: LaunchStatus.Launched,
	[Spec.SpecRogue]: LaunchStatus.Launched,
	[Spec.SpecHolyPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Launched,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Launched,
	[Spec.SpecHealingPriest]: LaunchStatus.Alpha,
	[Spec.SpecShadowPriest]: LaunchStatus.Launched,
	[Spec.SpecSmitePriest]: LaunchStatus.Launched,
	[Spec.SpecWarlock]: LaunchStatus.Launched,
	[Spec.SpecWarrior]: LaunchStatus.Launched,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Launched,
	[Spec.SpecDeathknight]: LaunchStatus.Launched,
	[Spec.SpecTankDeathknight]: LaunchStatus.Launched,
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
