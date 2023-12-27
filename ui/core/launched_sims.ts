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

export const raidSimStatus: LaunchStatus = LaunchStatus.Beta;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Alpha,
	[Spec.SpecFeralDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecRestorationDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecElementalShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecRestorationShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecHunter]: LaunchStatus.Unlaunched,
	[Spec.SpecMage]: LaunchStatus.Alpha,
	[Spec.SpecRogue]: LaunchStatus.Unlaunched,
	[Spec.SpecHolyPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecHealingPriest]: LaunchStatus.Unlaunched,
	[Spec.SpecShadowPriest]: LaunchStatus.Alpha,
	[Spec.SpecWarlock]: LaunchStatus.Alpha,
	[Spec.SpecWarrior]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Unlaunched,
};

// Alpha and Beta show an info notice at the top of the page.
export const aplLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Beta,
	[Spec.SpecFeralDruid]: LaunchStatus.Launched,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Launched,
	[Spec.SpecRestorationDruid]: LaunchStatus.Launched,
	[Spec.SpecElementalShaman]: LaunchStatus.Beta,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Beta,
	[Spec.SpecRestorationShaman]: LaunchStatus.Launched,
	[Spec.SpecHunter]: LaunchStatus.Launched,
	[Spec.SpecMage]: LaunchStatus.Launched,
	[Spec.SpecRogue]: LaunchStatus.Beta,
	[Spec.SpecHolyPaladin]: LaunchStatus.Launched,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Launched,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Beta,
	[Spec.SpecHealingPriest]: LaunchStatus.Launched,
	[Spec.SpecShadowPriest]: LaunchStatus.Launched,
	[Spec.SpecWarlock]: LaunchStatus.Alpha,
	[Spec.SpecWarrior]: LaunchStatus.Alpha,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Launched,
};

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
