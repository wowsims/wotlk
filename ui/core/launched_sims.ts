import { Spec } from './proto/common.js';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
    Unlaunched,
    Alpha,
    Beta,
    Launched,
}

export const raidSimLaunched = false;

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, LaunchStatus> = {
	[Spec.SpecBalanceDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecElementalShaman]: LaunchStatus.Unlaunched,
	[Spec.SpecEnhancementShaman]: LaunchStatus.Alpha,
	[Spec.SpecFeralDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecFeralTankDruid]: LaunchStatus.Unlaunched,
	[Spec.SpecHunter]: LaunchStatus.Alpha,
	[Spec.SpecMage]: LaunchStatus.Unlaunched,
	[Spec.SpecRogue]: LaunchStatus.Alpha,
	[Spec.SpecRetributionPaladin]: LaunchStatus.Alpha,
	[Spec.SpecProtectionPaladin]: LaunchStatus.Unlaunched,
	[Spec.SpecShadowPriest]: LaunchStatus.Alpha,
	[Spec.SpecWarlock]: LaunchStatus.Alpha,
	[Spec.SpecWarrior]: LaunchStatus.Unlaunched,
	[Spec.SpecProtectionWarrior]: LaunchStatus.Unlaunched,
	[Spec.SpecSmitePriest]: LaunchStatus.Unlaunched,
	[Spec.SpecDeathknight]: LaunchStatus.Alpha,
	[Spec.SpecTankDeathknight]: LaunchStatus.Unlaunched,
};

export function getLaunchedSims(): Array<Spec> {
    return Object.keys(simLaunchStatuses)
        .map(specStr => parseInt(specStr) as Spec)
        .filter(spec => simLaunchStatuses[spec] > LaunchStatus.Unlaunched);
}
