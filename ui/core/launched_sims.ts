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
	[Spec.SpecFeralDruid]: LaunchStatus.Alpha,
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
	[Spec.SpecTankWarlock]: LaunchStatus.Alpha,
	[Spec.SpecWarrior]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Unlaunched,
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
